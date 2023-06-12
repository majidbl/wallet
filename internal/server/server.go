package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcOpentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/majidbl/wallet/config"
	"github.com/majidbl/wallet/internal/interceptors"
	"github.com/majidbl/wallet/internal/middlewares"
	transactionV1 "github.com/majidbl/wallet/internal/transaction/delivery/http/v1"
	transactionNats "github.com/majidbl/wallet/internal/transaction/delivery/nats"
	transactionRepository "github.com/majidbl/wallet/internal/transaction/repository"
	transactionUsecase "github.com/majidbl/wallet/internal/transaction/usecase"
	walletGrpc "github.com/majidbl/wallet/internal/wallet/delivery/grpc"
	walletV1 "github.com/majidbl/wallet/internal/wallet/delivery/http/v1"
	"github.com/majidbl/wallet/internal/wallet/repository"
	"github.com/majidbl/wallet/internal/wallet/usecase"
	"github.com/majidbl/wallet/pkg/logger"
	walletService "github.com/majidbl/wallet/proto/wallet"
	"github.com/nats-io/stan.go"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	maxHeaderBytes  = 1 << 20
	gzipLevel       = 5
	stackSize       = 1 << 10 // 1 KB
	csrfTokenHeader = "X-CSRF-Token"
	bodyLimit       = "2M"
)

type server struct {
	log      logger.Logger
	cfg      *config.Config
	natsConn stan.Conn
	dbx      *pgxpool.Pool
	tracer   opentracing.Tracer
	echo     *echo.Echo
	redis    *redis.Client
}

// NewServer constructor
func NewServer(
	log logger.Logger,
	cfg *config.Config,
	natsConn stan.Conn,
	db *pgxpool.Pool,
	tracer opentracing.Tracer,
	redis *redis.Client,
) *server {
	return &server{
		log:      log,
		cfg:      cfg,
		natsConn: natsConn,
		dbx:      db,
		tracer:   tracer,
		redis:    redis,
		echo:     echo.New(),
	}
}

// Run start application
func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	metricsServer := echo.New()

	go func() {
		metricsServer.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		s.log.Infof("Metrics server is running on port: %s", s.cfg.Metrics.Port)
		if err := metricsServer.Start(s.cfg.Metrics.Port); err != nil {
			s.log.Error(err)
			cancel()
		}
	}()

	trxPublisher := transactionNats.NewPublisher(s.natsConn)
	transactionPgRepo := transactionRepository.NewTransactionPGRepository(s.dbx)
	transactionUC := transactionUsecase.NewTransactionUseCase(s.log, transactionPgRepo, trxPublisher)

	walletPgRepo := repository.NewWalletPGRepository(s.dbx)
	walletRedisRepo := repository.NewWalletRedisRepository(s.redis)
	walletUC := usecase.NewWalletUseCase(s.log, walletPgRepo, walletRedisRepo, transactionPgRepo)

	im := interceptors.NewInterceptorManager(s.log, s.cfg)
	mw := middlewares.NewMiddlewareManager(s.log, s.cfg)

	validate := validator.New()

	v1 := s.echo.Group("/api/v1")
	v1.Use(mw.Metrics)

	v1.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	walletHandlers := walletV1.NewWalletHandlers(
		v1.Group("/wallet"),
		walletUC,
		s.log,
		validate,
	)
	walletHandlers.MapRoutes()

	transactionHandlers := transactionV1.NewTransactionHandlers(
		v1.Group("/transactions"),
		transactionUC,
		s.log,
		validate,
	)
	transactionHandlers.MapRoutes()

	go func() {
		s.log.Infof("Server is listening on PORT: %s", s.cfg.HTTP.Port)
		s.runHttpServer()
	}()

	l, err := net.Listen("tcp", s.cfg.GRPC.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	defer l.Close()

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: s.cfg.GRPC.MaxConnectionIdle * time.Minute,
			Timeout:           s.cfg.GRPC.Timeout * time.Second,
			MaxConnectionAge:  s.cfg.GRPC.MaxConnectionAge * time.Minute,
			Time:              s.cfg.GRPC.Timeout * time.Minute,
		}),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			otgrpc.OpenTracingServerInterceptor(s.tracer, otgrpc.LogPayloads()),
			grpcCtxTags.UnaryServerInterceptor(),
			grpcOpentracing.UnaryServerInterceptor(),
			grpcPrometheus.UnaryServerInterceptor,
			grpcRecovery.UnaryServerInterceptor(),
			im.Logger,
		),
		),
	)

	walletGRPCService := walletGrpc.NewWalletGRPCService(walletUC, s.log, validate)
	walletService.RegisterWalletServiceServer(grpcServer, walletGRPCService)
	grpcPrometheus.Register(grpcServer)

	s.log.Infof("GRPC Server is listening on port: %s", s.cfg.GRPC.Port)
	s.log.Fatal(grpcServer.Serve(l))

	if s.cfg.HTTP.Development {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.log.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.log.Errorf("ctx.Done: %v", done)
	}

	if err = s.echo.Server.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "echo.Server.Shutdown")
	}

	if err = metricsServer.Shutdown(ctx); err != nil {
		s.log.Errorf("metricsServer.Shutdown: %v", err)
	}

	grpcServer.GracefulStop()
	s.log.Info("Server Exited Properly")

	return nil
}

package main

import (
	"log"

	"github.com/majidbl/wallet/config"
	"github.com/majidbl/wallet/internal/server"
	"github.com/majidbl/wallet/pkg/jaeger"
	"github.com/majidbl/wallet/pkg/logger"
	"github.com/majidbl/wallet/pkg/nats"
	"github.com/majidbl/wallet/pkg/postgresql"
	"github.com/majidbl/wallet/pkg/redis"
	"github.com/opentracing/opentracing-go"
)

// @title Wallet microservice
// @version 1.0
// @description Wallet microservice
// @termsOfService http://swagger.io/terms/

// @host localhost:5000
// @BasePath /api/v1
func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting wallet microservice")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, DevelopmentMode: %s",
		cfg.AppVersion,
		cfg.Logger.Level,
		cfg.HTTP.Development,
	)
	appLogger.Infof("Success loaded config: %+v", cfg.AppVersion)

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		appLogger.Fatalf("NewRedisClient: %+v", err)
	}

	appLogger.Infof("Redis connected: %+v", redisClient.PoolStats())

	natsConn, err := nats.NewNatsConnect(cfg, appLogger)
	if err != nil {
		appLogger.Fatalf("NewNatsConnect: %+v", err)
	}
	appLogger.Infof(
		"Nats Connected: Status: %+v IsConnected: %v ConnectedUrl: %v ConnectedServerId: %v",
		natsConn.NatsConn().Status(),
		natsConn.NatsConn().IsConnected(),
		natsConn.NatsConn().ConnectedUrl(),
		natsConn.NatsConn().ConnectedServerId(),
	)

	pgxPool, err := postgresql.NewPgxConn(cfg)
	if err != nil {
		appLogger.Fatalf("NewPgxConn: %+v", err)
	}
	appLogger.Infof("PostgreSQL connected: %+v", pgxPool.Stat().TotalConns())

	s := server.NewServer(appLogger, cfg, natsConn, pgxPool, tracer, redisClient)

	appLogger.Fatal(s.Run())
}

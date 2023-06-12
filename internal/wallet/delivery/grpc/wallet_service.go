package grpc

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"

	"github.com/majidbl/wallet/internal/models"
	"github.com/majidbl/wallet/internal/wallet"
	grpcErrors "github.com/majidbl/wallet/pkg/grpc_errors"
	"github.com/majidbl/wallet/pkg/logger"
	walletService "github.com/majidbl/wallet/proto/wallet"
)

type walletGRPCService struct {
	walletUC  wallet.UseCase
	log       logger.Logger
	validator *validator.Validate
	walletService.UnimplementedWalletServiceServer
}

// NewWalletGRPCService wallet gRPC service constructor
func NewWalletGRPCService(
	walletUC wallet.UseCase,
	log logger.Logger,
	validator *validator.Validate,
) *walletGRPCService {
	return &walletGRPCService{walletUC: walletUC, log: log, validator: validator}
}

// Create wallet
func (e *walletGRPCService) Create(ctx context.Context, req *walletService.CreateReq) (*walletService.CreateRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chargeService.Create")
	defer span.Finish()
	createRequests.Inc()

	CreateRequest := &models.Wallet{
		ID:          uuid.New(),
		Name:        req.Name,
		Mobile:      req.Mobile,
		Balance:     req.Balance,
		Avatar:      &req.Avatar,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	if err := e.validator.StructCtx(ctx, CreateRequest); err != nil {
		errorRequests.Inc()
		e.log.Errorf("validator.StructCtx: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	if err := e.walletUC.Create(ctx, CreateRequest); err != nil {
		errorRequests.Inc()
		e.log.Errorf("walletUC.Create: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	successRequests.Inc()
	return &walletService.CreateRes{Status: "Ok"}, nil
}

// Charge wallet
func (e *walletGRPCService) Charge(ctx context.Context, req *walletService.ChargeReq) (*walletService.ChargeRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chargeService.Charge")
	defer span.Finish()
	chargeRequests.Inc()

	chargeRequest := &models.ChargeWalletReq{
		Mobile:    req.Mobile,
		Amount:    req.Amount,
		UpdatedAt: time.Now(),
	}

	if err := e.validator.StructCtx(ctx, chargeRequest); err != nil {
		errorRequests.Inc()
		e.log.Errorf("validator.StructCtx: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	if err := e.walletUC.ChargeWallet(ctx, chargeRequest); err != nil {
		errorRequests.Inc()
		e.log.Errorf("walletUC.ChargeWallet: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	successRequests.Inc()
	return &walletService.ChargeRes{Status: "Ok"}, nil
}

// GetByID find single wallet by id
func (e *walletGRPCService) GetByID(ctx context.Context, req *walletService.GetByIDReq) (*walletService.GetByIDRes, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chargeService.GetByID")
	defer span.Finish()

	getByIdRequests.Inc()

	walletUUID, err := uuid.Parse(req.GetWalletID())
	if err != nil {
		errorRequests.Inc()
		e.log.Errorf("uuid.parse: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	res, err := e.walletUC.GetByID(ctx, walletUUID)
	if err != nil {
		errorRequests.Inc()
		e.log.Errorf("walletUC.GetByID: %v", err)
		return nil, grpcErrors.ErrorResponse(err, err.Error())
	}

	successRequests.Inc()
	return &walletService.GetByIDRes{Wallet: res.ToProto()}, nil
}

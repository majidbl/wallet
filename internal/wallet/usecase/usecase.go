package usecase

import (
	"context"
	"github.com/majidbl/wallet/pkg/sql_errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/wallet/internal/models"
	"github.com/majidbl/wallet/internal/transaction"
	"github.com/majidbl/wallet/internal/wallet"
	"github.com/majidbl/wallet/pkg/logger"
)

type walletUseCase struct {
	log               logger.Logger
	walletPGRepo      wallet.PGRepository
	transactionPGRepo transaction.PGRepository
	redisRepo         wallet.RedisRepository
}

// NewWalletUseCase wallet usecase constructor
func NewWalletUseCase(
	log logger.Logger,
	walletPGRepo wallet.PGRepository,
	redisRepo wallet.RedisRepository,
	transactionPGRepo transaction.PGRepository,
) *walletUseCase {
	return &walletUseCase{
		log:               log,
		walletPGRepo:      walletPGRepo,
		transactionPGRepo: transactionPGRepo,
		redisRepo:         redisRepo,
	}
}

var _ wallet.UseCase = &walletUseCase{}

// Create creates a new wallet and saves it in the database
func (wu *walletUseCase) Create(ctx context.Context, wallet *models.Wallet) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletUseCase.Create")
	defer span.Finish()

	_, err := wu.walletPGRepo.Create(ctx, wallet)
	if err != nil {
		return errors.Wrap(err, "walletPGRepo.Create")
	}

	return nil
}

// GetByID fnd wallet by id
func (wu *walletUseCase) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletUseCase.GetByID")
	defer span.Finish()

	cached, err := wu.redisRepo.GetWalletByID(ctx, walletID)
	if err != nil && err != redis.Nil {
		wu.log.Errorf("redisRepo.GetWalletByID: %v", err)
	}

	if cached != nil {
		return cached, nil
	}

	w, err := wu.walletPGRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, errors.Wrap(err, "walletPGRepo.GetByID")
	}

	if err = wu.redisRepo.SetWalletByID(ctx, w); err != nil {
		wu.log.Errorf("redisRepo.SetWallet: %v", err)
	}

	return w, nil
}

// GetByMobile retrieves a wallet by its mobile number
func (wu *walletUseCase) GetByMobile(ctx context.Context, mobile string) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletUseCase.GetByMobile")
	defer span.Finish()

	cached, err := wu.redisRepo.GetWalletByMobile(ctx, mobile)
	if err != nil && err != redis.Nil {
		wu.log.Errorf("redisRepo.GetWalletByMobile: %v", err)
	}

	if cached != nil {
		return cached, nil
	}

	w, err := wu.walletPGRepo.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, errors.Wrap(err, "walletPGRepo.GetByMobile")
	}

	if err = wu.redisRepo.SetWalletByMobile(ctx, w); err != nil {
		wu.log.Errorf("redisRepo.SetWallet: %v", err)
	}

	return w, nil
}

// ChargeWallet charge wallet using mobile
func (wu *walletUseCase) ChargeWallet(ctx context.Context, req *models.ChargeWalletReq) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "walletUseCase.ChargeWallet")
	defer span.Finish()

	var walletByMobile *models.Wallet
	var getWalletByMobileErr error

	walletByMobile, getWalletByMobileErr = wu.GetByMobile(ctx, req.Mobile)
	if getWalletByMobileErr != nil {
		if !errors.As(getWalletByMobileErr, &sql_errors.SqlNotFound) {
			return errors.Wrap(getWalletByMobileErr, "walletPGRepo.GetByMobile")
		}

		walletByMobile = &models.Wallet{
			ID:        uuid.New(),
			Mobile:    req.Mobile,
			Balance:   req.Amount,
			CreatedAt: time.Now(),
		}

		createWalletErr := wu.Create(ctx, walletByMobile)
		if createWalletErr != nil {
			return errors.Wrap(createWalletErr, "walletPGRepo.Create")
		}
	}

	tx, updateErr := wu.walletPGRepo.UpdateBalanceX(ctx, &models.UpdateWalletBalanceReq{
		WalletID: walletByMobile.ID,
		Amount:   walletByMobile.Balance + req.Amount,
	}, nil)

	if updateErr != nil {
		return errors.Wrap(updateErr, "walletPGRepo.UpdateBalanceX")
	}

	_, _, trxErr := wu.transactionPGRepo.CreateX(
		ctx,
		&models.CreateTransactionReq{
			WalletID: walletByMobile.ID,
			Amount:   req.Amount,
			Type:     models.Charge,
		},
		tx)

	if trxErr != nil {
		if rollBackErr := wu.walletPGRepo.RollBack(ctx, tx); rollBackErr != nil {
			return errors.Wrap(rollBackErr, "tx.RollBack")
		}

		return trxErr
	}

	if commitErr := wu.walletPGRepo.Commit(ctx, tx); commitErr != nil {
		return errors.Wrap(commitErr, "tx.Commit")
	}

	walletByMobile.Balance = walletByMobile.Balance + req.Amount

	if err := wu.redisRepo.SetWalletByMobile(ctx, walletByMobile); err != nil {
		wu.log.Errorf("redisRepo.SetWallet: %v", err)
	}

	walletChargeCounter.Inc()
	return nil
}

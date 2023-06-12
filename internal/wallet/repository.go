package wallet

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/majidbl/wallet/internal/models"
)

// PGRepository Wallet postgresql repository interface
type PGRepository interface {
	Create(ctx context.Context, wallet *models.Wallet) (*models.Wallet, error)
	GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
	GetByMobile(ctx context.Context, mobile string) (*models.Wallet, error)
	UpdateBalance(ctx context.Context, wallet *models.UpdateWalletBalanceReq) error
	UpdateBalanceX(ctx context.Context, wallet *models.UpdateWalletBalanceReq, tx pgx.Tx) (pgx.Tx, error)
	RollBack(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}

// RedisRepository redis wallet repository interface
type RedisRepository interface {
	SetWalletByMobile(ctx context.Context, wallet *models.Wallet) error
	SetWalletByID(ctx context.Context, wallet *models.Wallet) error
	GetWalletByMobile(ctx context.Context, mobile string) (*models.Wallet, error)
	GetWalletByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
}

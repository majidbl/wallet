package wallet

import (
	"context"

	"github.com/google/uuid"

	"github.com/majidbl/wallet/internal/models"
)

// UseCase Wallet usecase interface
type UseCase interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
	GetByMobile(ctx context.Context, mobile string) (*models.Wallet, error)
	ChargeWallet(ctx context.Context, req *models.ChargeWalletReq) error
}

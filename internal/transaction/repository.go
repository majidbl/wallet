package transaction

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/majidbl/wallet/internal/models"
)

// PGRepository Transaction postgresql repository interface
type PGRepository interface {
	Create(ctx context.Context, transaction *models.CreateTransactionReq) (*models.Transaction, error)
	CreateX(ctx context.Context, transaction *models.CreateTransactionReq, tx pgx.Tx) (*models.Transaction, pgx.Tx, error)
	GetByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error)
	GetByWalletID(ctx context.Context, walletId uuid.UUID) ([]*models.Transaction, error)
	RollBack(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}

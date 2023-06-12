package transaction

import (
	"context"

	"github.com/google/uuid"

	"github.com/majidbl/wallet/internal/models"
)

// UseCase Transaction usecase interface
type UseCase interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	PublishCreate(ctx context.Context, transaction *models.Transaction) error
	GetByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error)
	GetList(ctx context.Context, transactionID uuid.UUID) ([]*models.Transaction, error)
}

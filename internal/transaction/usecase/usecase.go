package usecase

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/wallet/internal/models"
	"github.com/majidbl/wallet/internal/transaction"
	"github.com/majidbl/wallet/internal/transaction/delivery/nats"
	"github.com/majidbl/wallet/pkg/logger"
)

const (
	createTransactionSubject = "transaction:create"
)

type transactionUseCase struct {
	log               logger.Logger
	transactionPGRepo transaction.PGRepository
	publisher         nats.Publisher
}

// NewTransactionUseCase transaction usecase constructor
func NewTransactionUseCase(
	log logger.Logger,
	transactionPGRepo transaction.PGRepository,
	publisher nats.Publisher,
) *transactionUseCase {
	return &transactionUseCase{
		log:               log,
		transactionPGRepo: transactionPGRepo,
		publisher:         publisher,
	}
}

var _ transaction.UseCase = &transactionUseCase{}

// Create new transaction saves in db
func (e *transactionUseCase) Create(ctx context.Context, req *models.Transaction) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "transactionUseCase.Create")
	defer span.Finish()

	tType, err := models.GetTransactionType(req.Type)
	if err != nil {
		span.SetTag("error:", err)
		return errors.Wrap(err, "getTransactionType")
	}

	created, err := e.transactionPGRepo.Create(ctx, &models.CreateTransactionReq{
		WalletID: req.WalletID,
		Amount:   req.Amount,
		Type:     tType,
	})
	if err != nil {
		span.SetTag("error:", err)
		return errors.Wrap(err, "transactionPGRepo.Create")
	}

	transactionBytes, err := json.Marshal(created)
	if err != nil {
		span.SetTag("error:", err)
		return errors.Wrap(err, "json.Marshal")
	}

	return e.publisher.Publish(createTransactionSubject, transactionBytes)
}

// GetByID fnd transaction by id
func (e *transactionUseCase) GetByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "transactionUseCase.GetByID")
	defer span.Finish()

	trx, err := e.transactionPGRepo.GetByID(ctx, transactionID)
	if err != nil {
		return nil, errors.Wrap(err, "transactionPGRepo.GetByID")
	}

	return trx, nil
}

// PublishCreate publish create transaction event to message broker
func (e *transactionUseCase) PublishCreate(ctx context.Context, transaction *models.Transaction) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "transactionUseCase.PublishCreate")
	defer span.Finish()

	transactionBytes, err := json.Marshal(transaction)
	if err != nil {
		return errors.Wrap(err, "json.Marshal")
	}

	return e.publisher.Publish(createTransactionSubject, transactionBytes)
}

func (e *transactionUseCase) GetList(ctx context.Context, transactionID uuid.UUID) ([]*models.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/wallet/internal/models"
	"github.com/majidbl/wallet/internal/transaction"
)

type transactionPGRepository struct {
	db *pgxpool.Pool
}

// NewTransactionPGRepository Transaction postgresql repository constructor
func NewTransactionPGRepository(db *pgxpool.Pool) *transactionPGRepository {
	return &transactionPGRepository{db: db}
}

var _ transaction.PGRepository = &transactionPGRepository{}

// Create  new transaction
func (tr *transactionPGRepository) Create(ctx context.Context, createRequest *models.CreateTransactionReq) (*models.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "transactionPGRepository.Create")
	defer span.Finish()

	var t models.Transaction
	if err := tr.db.QueryRow(
		ctx,
		createTransactionQuery,
		createRequest.WalletID,
		createRequest.Amount,
		createRequest.Type.String(),
		time.Now(),
	).Scan(
		&t.ID,
		&t.WalletID,
		&t.Amount,
		&t.Type,
		&t.CreatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &t, nil
}

// CreateX new transactional create
func (tr *transactionPGRepository) CreateX(ctx context.Context, createRequest *models.CreateTransactionReq, tx pgx.Tx) (*models.Transaction, pgx.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "transactionPGRepository.CreateX")
	defer span.Finish()

	if tx == nil {
		var beginTxErr error
		tx, beginTxErr = tr.db.Begin(ctx)
		if beginTxErr != nil {
			return nil, nil, beginTxErr
		}
	}

	var t models.Transaction
	if err := tr.db.QueryRow(
		ctx,
		createTransactionQuery,
		createRequest.WalletID,
		createRequest.Amount,
		createRequest.Type.String(),
		time.Now(),
	).Scan(
		&t.ID,
		&t.WalletID,
		&t.Amount,
		&t.Type,
		&t.CreatedAt,
	); err != nil {
		return nil, nil, errors.Wrap(err, "Scan")
	}

	return &t, tx, nil
}

// GetByID get single transaction by id
func (tr *transactionPGRepository) GetByID(ctx context.Context, transactionID uuid.UUID) (*models.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "transactionPGRepository.GetByID")
	defer span.Finish()

	var t models.Transaction
	if err := tr.db.QueryRow(ctx, getByIDQuery, transactionID).
		Scan(
			&t.ID,
			&t.WalletID,
			&t.Amount,
			&t.Type,
			&t.CreatedAt,
		); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return &t, nil
}

// GetByWalletID get all transactions of a wallet by wallet_id
func (tr *transactionPGRepository) GetByWalletID(ctx context.Context, walletId uuid.UUID) ([]*models.Transaction, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "transactionPGRepository.GetByWalletID")
	defer span.Finish()

	rows, err := tr.db.Query(ctx, getByWalletIDQuery, walletId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		var t models.Transaction
		err = rows.Scan(&t.ID, &t.WalletID, &t.Amount, &t.Type, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil

}

func (tr *transactionPGRepository) RollBack(ctx context.Context, tx pgx.Tx) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "transactionPGRepository.RollBack")
	defer span.Finish()

	if tx == nil {
		return fmt.Errorf("transactions not begin")
	}

	return tx.Rollback(ctx)
}

func (tr *transactionPGRepository) Commit(ctx context.Context, tx pgx.Tx) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "transactionPGRepository.Commit")
	defer span.Finish()

	if tx == nil {
		return fmt.Errorf("transactions not begin")
	}

	return tx.Commit(ctx)
}

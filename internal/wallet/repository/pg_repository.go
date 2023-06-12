package repository

import (
	"context"
	"fmt"
	"github.com/majidbl/wallet/pkg/sql_errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/majidbl/wallet/internal/models"
	"github.com/majidbl/wallet/internal/wallet"
)

type walletPGRepository struct {
	db *pgxpool.Pool
}

// NewWalletPGRepository Wallet postgresql repository constructor
func NewWalletPGRepository(db *pgxpool.Pool) *walletPGRepository {
	return &walletPGRepository{db: db}
}

var _ wallet.PGRepository = &walletPGRepository{}

// Create  new wallet
func (wr *walletPGRepository) Create(ctx context.Context, createReq *models.Wallet) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletPGRepository.Create")
	defer span.Finish()

	var w models.WalletModel
	if err := wr.db.QueryRow(
		ctx,
		createWalletQuery,
		&createReq.Name,
		&createReq.Mobile,
		&createReq.Balance,
		&createReq.Avatar,
		&createReq.Description,
	).Scan(
		&w.ID,
		&w.Name,
		&w.Mobile,
		&w.Balance,
		&w.Avatar,
		&w.Description,
		&w.CreatedAt,
		&w.UpdatedAt,
	); err != nil {
		return nil, errors.Wrap(err, "Scan")
	}

	return w.Entity(), nil
}

// GetByID get single wallet by id
func (wr *walletPGRepository) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletPGRepository.GetByID")
	defer span.Finish()

	var w models.WalletModel
	if err := wr.db.QueryRow(ctx, getByIDQuery, walletID).
		Scan(
			&w.ID,
			&w.Name,
			&w.Mobile,
			&w.Balance,
			&w.Avatar,
			&w.Description,
			&w.CreatedAt,
			&w.UpdatedAt,
		); err != nil {
		return nil, sql_errors.ParseSqlErrors(errors.Wrap(err, "Scan"))
	}

	return w.Entity(), nil
}

// GetByMobile get single wallet by wallet
func (wr *walletPGRepository) GetByMobile(ctx context.Context, mobile string) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletPGRepository.GetByMobile")
	defer span.Finish()

	var w models.WalletModel
	if err := wr.db.QueryRow(ctx, getByMobileQuery, mobile).
		Scan(
			&w.ID,
			&w.Name,
			&w.Mobile,
			&w.Balance,
			&w.Avatar,
			&w.Description,
			&w.CreatedAt,
			&w.UpdatedAt,
		); err != nil {
		return nil, sql_errors.ParseSqlErrors(errors.Wrap(err, "Scan"))
	}

	return w.Entity(), nil
}

// UpdateBalance updates the balance of a wallet
func (wr *walletPGRepository) UpdateBalance(ctx context.Context, request *models.UpdateWalletBalanceReq) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletPGRepository.UpdateBalance")
	defer span.Finish()

	_, err := wr.db.Exec(ctx, updateBalanceQuery, request.Amount, request.WalletID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateBalanceX updates the balance of a wallet within a transaction
func (wr *walletPGRepository) UpdateBalanceX(ctx context.Context, request *models.UpdateWalletBalanceReq, tx pgx.Tx) (pgx.Tx, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletPGRepository.UpdateBalanceX")
	defer span.Finish()

	if tx == nil {
		var beginTxErr error
		tx, beginTxErr = wr.db.Begin(ctx)
		if beginTxErr != nil {
			return nil, beginTxErr
		}
	}

	_, err := tx.Exec(ctx, updateBalanceQuery, request.Amount, request.WalletID)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// RollBack rolls back a transaction
func (wr *walletPGRepository) RollBack(ctx context.Context, tx pgx.Tx) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletPGRepository.RollBack")
	defer span.Finish()

	if tx == nil {
		return fmt.Errorf("transactions not begin")
	}

	return tx.Rollback(ctx)
}

// Commit commits a transaction
func (wr *walletPGRepository) Commit(ctx context.Context, tx pgx.Tx) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletPGRepository.Commit")
	defer span.Finish()

	if tx == nil {
		return fmt.Errorf("transactions not begin")
	}

	return tx.Commit(ctx)
}

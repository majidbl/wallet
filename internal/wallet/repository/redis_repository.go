package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/majidbl/wallet/internal/wallet"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/google/uuid"
	"github.com/majidbl/wallet/internal/models"
)

const (
	prefix       = "wallets"
	mobilePrefix = "mobile"
	idPrefix     = "id"
	expiration   = time.Second * 3600
)

type walletRedisRepository struct {
	redis *redis.Client
}

// NewWalletRedisRepository wallets redis repository constructor
func NewWalletRedisRepository(redis *redis.Client) *walletRedisRepository {
	return &walletRedisRepository{redis: redis}
}

var _ wallet.RedisRepository = &walletRedisRepository{}

func (e *walletRedisRepository) SetWalletByID(ctx context.Context, wallet *models.Wallet) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletRedisRepository.SetWalletByID")
	defer span.Finish()

	walletBytes, err := json.Marshal(wallet)
	if err != nil {
		return errors.Wrap(err, "walletRedisRepository.Marshal.SetWalletByID")
	}

	return e.redis.HSet(ctx, e.createKeyById(wallet.ID), prefix, string(walletBytes), expiration).Err()
}

func (e *walletRedisRepository) SetWalletByMobile(ctx context.Context, wallet *models.Wallet) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletRedisRepository.SetWallet")
	defer span.Finish()

	walletBytes, err := json.Marshal(wallet)
	if err != nil {
		return errors.Wrap(err, "walletRedisRepository.Marshal.SetWalletByID")
	}

	return e.redis.HSet(ctx, e.createKeyById(wallet.ID), prefix, string(walletBytes), expiration).Err()
}

func (e *walletRedisRepository) GetWalletByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletRedisRepository.GetWalletByID")
	defer span.Finish()

	result, err := e.redis.HGet(ctx, idPrefix, e.createKeyById(walletID)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "walletRedisRepository.redis.GetWalletByID")
	}

	var res models.Wallet
	if err := json.Unmarshal(result, &res); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return &res, nil
}

func (e *walletRedisRepository) GetWalletByMobile(ctx context.Context, mobile string) (*models.Wallet, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "walletRedisRepository.GetWalletByMobile")
	defer span.Finish()

	result, err := e.redis.HGet(ctx, mobilePrefix, e.createKeyByMobile(mobile)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "walletRedisRepository.redis.GetWalletByMobile")
	}

	var res models.Wallet
	if err = json.Unmarshal(result, &res); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal")
	}
	return &res, nil
}

func (e *walletRedisRepository) createKeyById(walletID uuid.UUID) string {
	return fmt.Sprintf("%s: %s", idPrefix, walletID.String())
}

func (e *walletRedisRepository) createKeyByMobile(mobile string) string {
	return fmt.Sprintf("%s: %s", mobilePrefix, mobile)
}

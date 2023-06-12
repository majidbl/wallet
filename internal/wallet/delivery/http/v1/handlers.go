package v1

import (
	"fmt"
	"github.com/majidbl/wallet/pkg/validation"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/majidbl/wallet/internal/models"
	"github.com/majidbl/wallet/internal/wallet"
	httpErrors "github.com/majidbl/wallet/pkg/http_errors"
	"github.com/majidbl/wallet/pkg/logger"
)

type walletHandlers struct {
	group    *echo.Group
	walletUC wallet.UseCase
	log      logger.Logger
	validate *validator.Validate
}

// NewWalletHandlers walletHandlers constructor
func NewWalletHandlers(
	group *echo.Group,
	walletUC wallet.UseCase,
	log logger.Logger,
	validate *validator.Validate,
) *walletHandlers {
	return &walletHandlers{group: group, walletUC: walletUC, log: log, validate: validate}
}

// Create New Wallet
// @Tags Wallet
// @Summary Create new wallet
// @Description Create new wallet and send it
// @Accept json
// @Produce json
// @Success 201 {object} models.Wallet
// @Router /wallet [post]
func (h *walletHandlers) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "walletHandlers.Create")
		defer span.Finish()
		createRequests.Inc()

		var createWalletReq models.Wallet

		if err := c.Bind(&createWalletReq); err != nil {
			errorRequests.Inc()
			h.log.Errorf("c.Bind: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if err := h.validate.StructCtx(ctx, &createWalletReq); err != nil {
			errorRequests.Inc()
			h.log.Errorf("validate.StructCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		if valid := validation.ValidatePhoneNumber(createWalletReq.Mobile); !valid {
			errorRequests.Inc()
			validationErr := fmt.Errorf("invalid phone")
			h.log.Errorf("validation.ValidatePhoneNumber: %v", validationErr)
			return httpErrors.ErrorCtxResponse(c, validationErr)
		}

		if err := h.walletUC.Create(ctx, &createWalletReq); err != nil {
			errorRequests.Inc()
			h.log.Errorf("walletUC.Create: %v", err)
			span.LogKV("err", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.NoContent(http.StatusCreated)
	}
}

// GetByID Get Wallet by ID
// @Tags Wallet
// @Summary Get wallet by id
// @Description Get wallet by wallet uuid
// @Accept json
// @Produce json
// @Param wallet_id path string true "wallet_id"
// @Success 200 {object} models.Wallet
// @Router /wallet/{wallet_id} [get]
func (h *walletHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "walletHandlers.GetByID")
		defer span.Finish()
		getByIdRequests.Inc()

		walletUUID, err := uuid.Parse(c.Param("wallet_id"))
		if err != nil {
			errorRequests.Inc()
			h.log.Errorf("uuid.FromString: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		m, err := h.walletUC.GetByID(ctx, walletUUID)
		if err != nil {
			errorRequests.Inc()
			h.log.Errorf("walletUC.GetByID: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, m)
	}
}

// GetBalance Get Wallet Balance by ID
// @Tags Wallet
// @Summary Get wallet balance by id
// @Description Get wallet by wallet uuid
// @Accept json
// @Produce json
// @Param wallet_id path string true "wallet_id"
// @Success 200 {object} models.Wallet
// @Router /wallet/balance/{wallet_id} [get]
func (h *walletHandlers) GetBalance() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "walletHandlers.GetBalance")
		defer span.Finish()
		getBalanceRequests.Inc()

		walletUUID, err := uuid.Parse(c.Param("wallet_id"))
		if err != nil {
			errorRequests.Inc()
			h.log.Errorf("uuid.FromString: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		m, err := h.walletUC.GetByID(ctx, walletUUID)
		if err != nil {
			errorRequests.Inc()
			h.log.Errorf("walletUC.GetByID: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, map[string]int64{
			"balance": m.Balance,
		})
	}
}

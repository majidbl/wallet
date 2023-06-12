package v1

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"

	"github.com/majidbl/wallet/internal/transaction"
	httpErrors "github.com/majidbl/wallet/pkg/http_errors"
	"github.com/majidbl/wallet/pkg/logger"
)

type transactionHandlers struct {
	group         *echo.Group
	transactionUC transaction.UseCase
	log           logger.Logger
	validate      *validator.Validate
}

// NewTransactionHandlers transactionHandlers constructor
func NewTransactionHandlers(group *echo.Group, transactionUC transaction.UseCase, log logger.Logger, validate *validator.Validate) *transactionHandlers {
	return &transactionHandlers{group: group, transactionUC: transactionUC, log: log, validate: validate}
}

// GetByID get trx with id
// @Tags Transactions
// @Summary Get transaction by id
// @Description Get transaction by transaction uuid
// @Accept json
// @Produce json
// @Param transaction_id path string true "transaction_id"
// @Success 200 {object} models.Transaction
// @Router /transaction/{transaction_id} [get]
func (h *transactionHandlers) GetByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "transactionHandlers.GetByID")
		defer span.Finish()
		getByIdRequests.Inc()

		transactionUUID, err := uuid.Parse(c.Param("transaction_id"))
		if err != nil {
			errorRequests.Inc()
			h.log.Errorf("uuid.FromString: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		m, err := h.transactionUC.GetByID(ctx, transactionUUID)
		if err != nil {
			errorRequests.Inc()
			h.log.Errorf("transactionUC.GetByID: %v", err)
			return httpErrors.ErrorCtxResponse(c, err)
		}

		successRequests.Inc()
		return c.JSON(http.StatusOK, m)
	}
}

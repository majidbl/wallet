package v1

// MapRoutes transactions REST API routes
func (h *transactionHandlers) MapRoutes() {
	h.group.GET("/:transaction_id", h.GetByID())
}

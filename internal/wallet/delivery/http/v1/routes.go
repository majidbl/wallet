package v1

// MapRoutes wallets REST API routes
func (h *walletHandlers) MapRoutes() {
	h.group.POST("", h.Create())
	h.group.GET("/:wallet_id", h.GetByID())
	h.group.GET("/balance/:wallet_id", h.GetBalance())
}

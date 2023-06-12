package v1

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_wallet_success_incoming_messages_total",
		Help: "The total number of success incoming wallet HTTP requests",
	})
	errorRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_wallet_error_incoming_message_total",
		Help: "The total number of error incoming wallet HTTP requests",
	})
	createRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_wallet_create_incoming_requests_total",
		Help: "The total number of incoming create wallet HTTP requests",
	})
	getByIdRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_wallet_get_by_id_incoming_requests_total",
		Help: "The total number of incoming get by id wallet HTTP requests",
	})
	getBalanceRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_wallet_get_balance_incoming_requests_total",
		Help: "The total number of incoming get balance wallet HTTP requests",
	})
)

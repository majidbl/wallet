package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	successRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_wallet_success_incoming_messages_total",
		Help: "The total number of success incoming wallet GRPC requests",
	})
	errorRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_wallet_error_incoming_message_total",
		Help: "The total number of error incoming wallet GRPC requests",
	})
	createRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_wallet_create_incoming_requests_total",
		Help: "The total number of incoming create wallet GRPC requests",
	})
	chargeRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_wallet_charge_incoming_requests_total",
		Help: "The total number of incoming create wallet GRPC requests",
	})
	getByIdRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_wallet_get_by_id_incoming_requests_total",
		Help: "The total number of incoming get by id wallet GRPC requests",
	})
)

package nats

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalSubscribeMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nats_transaction_incoming_messages_total",
		Help: "The total number of incoming transaction NATS messages",
	})
	successSubscribeMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nats_transaction_success_incoming_messages_total",
		Help: "The total number of success transaction NATS messages",
	})
	errorSubscribeMessages = promauto.NewCounter(prometheus.CounterOpts{
		Name: "nats_transaction_error_incoming_messages_total",
		Help: "The total number of error transaction NATS messages",
	})
)

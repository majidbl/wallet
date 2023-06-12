package usecase

import "github.com/prometheus/client_golang/prometheus"

var (
	walletChargeCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "wallet_charge_total",
		Help: "Total number of wallet charges",
	})
)

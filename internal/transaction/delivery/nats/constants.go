package nats

import "time"

const (
	ackWait     = 60 * time.Second
	durableName = "transaction-dur"
	maxInflight = 25

	createTransactionWorkers = 6

	createTransactionSubject = "transaction:create"
	transactionGroupName     = "transaction_service"

	deadLetterQueueSubject = "transaction:errors"
	maxRedeliveryCount     = 3
)

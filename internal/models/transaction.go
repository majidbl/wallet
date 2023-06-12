package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID        uuid.UUID `json:"id"`
	WalletID  uuid.UUID `json:"wallet_id"`
	Amount    int64     `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTransactionReq struct {
	WalletID uuid.UUID       `json:"wallet_id"`
	Amount   int64           `json:"amount"`
	Type     TransactionType `json:"type"`
}

type TransactionErrorMsg struct {
	Subject   string    `json:"subject"`
	Sequence  uint64    `json:"sequence"`
	Data      []byte    `json:"data"`
	Timestamp int64     `json:"topic"`
	Error     string    `json:"error"`
	Time      time.Time `json:"time"`
}

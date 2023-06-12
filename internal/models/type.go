package models

import (
	"fmt"
	"strings"
)

type TransactionType int

const (
	Income TransactionType = iota
	Expense
	Transfer
	Deposit
	Withdrawal
	Refund
	Payment
	Conversion
	Interest
	Adjustment
	Charge
)

var transactionTypeStrings = [...]string{
	"Income",
	"Expense",
	"Transfer",
	"Deposit",
	"Withdrawal",
	"Refund",
	"Payment",
	"Conversion",
	"Interest",
	"Adjustment",
	"Charge",
}

// String returns the string representation of the TransactionType
func (t TransactionType) String() string {
	if t < Income || t > Charge {
		return fmt.Sprintf("Unknown TransactionType: %d", t)
	}
	return transactionTypeStrings[t]
}

var transactionTypeName = map[TransactionType]string{
	Income:     "Income",
	Expense:    "Expense",
	Transfer:   "Transfer",
	Deposit:    "Deposit",
	Withdrawal: "Withdrawal",
	Refund:     "Refund",
	Payment:    "Payment",
	Conversion: "Conversion",
	Interest:   "Interest",
	Adjustment: "Adjustment",
	Charge:     "Charge",
}

// StringMap returns the string representation of the TransactionType
func (t TransactionType) StringMap() (string, bool) {
	if s, ok := transactionTypeName[t]; ok {
		return s, true
	}
	return fmt.Sprintf("unknown TransactionType: %d", t), false
}

// GetTransactionType returns the TransactionType based on its string representation
func GetTransactionType(s string) (TransactionType, error) {
	s = strings.ToLower(s)
	for t, str := range transactionTypeName {
		if strings.ToLower(str) == s {
			return t, nil
		}
	}
	return 0, fmt.Errorf("unknown TransactionType: %s", s)
}

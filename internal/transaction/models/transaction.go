package models

import "time"

const (
	MaxTransactionCount = 10
	MaxAttemptsCount    = 5
)

var (
	ProcessingErr TransactionStatus = 0
	InProcessing  TransactionStatus = 1
	Processed     TransactionStatus = 2

	InputType  TransactionsType = 1
	OutputType TransactionsType = 2
)

type TransactionStatus int
type TransactionsType int

type Transaction struct {
	ID       int       `json:"ID" db:"id"`
	UserID   int       `json:"user_id" db:"user_id"`
	Attempts int       `json:"attempts" db:"attempts"`
	Status   int       `json:"status" db:"status"`
	Type     int       `json:"type" db:"type"`
	Amount   float64   `json:"amount" db:"amount"`
	Date     time.Time `json:"date" db:"date"`
}

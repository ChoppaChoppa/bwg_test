package models

import (
	"sync"
	"time"
)

var (
	ProcessingErr TransactionStatus = 0
	InProcessing  TransactionStatus = 1
	Processed     TransactionStatus = 2

	InputType  TransactionsType = 1
	OutputType TransactionsType = 2
)

type User struct {
	ID       int    `json:"ID" db:"id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

type UserTransaction struct {
	Active       bool
	Ch           chan struct{}
	mutex        *sync.Mutex
	Transactions []Transaction
}

type Balance struct {
	ID      int     `json:"ID" db:"id"`
	UserID  int     `json:"userID" db:"user_id"`
	Balance float64 `json:"balance" db:"balance"`
}

type Transaction struct {
	ID     int       `json:"ID" db:"id"`
	UserID int       `json:"user_id" db:"user_id"`
	Status int       `json:"status" db:"status"`
	Type   int       `json:"type" db:"type"`
	Amount float64   `json:"amount" db:"amount"`
	Date   time.Time `json:"date" db:"date"`
}

type TransactionStatus int
type TransactionsType int

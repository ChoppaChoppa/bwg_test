package models

import "errors"

var (
	ErrWrongAmount    error = errors.New("wrong amount")
	ErrPositiveAmount error = errors.New("pq: new row for relation \\\"balances\\\" violates check constraint \\\"positive_balance\\\"")
)

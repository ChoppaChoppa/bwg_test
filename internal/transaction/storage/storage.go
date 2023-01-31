package storage

import (
	"bwg_test/internal/transaction/models"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type storage struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) *storage {
	return &storage{
		conn: conn,
	}
}

func (s *storage) OutputTransaction(ctx context.Context, transaction *models.Transaction) (float64, error) {
	transaction.Amount = transaction.Amount * -1
	balance, err := s.updateBalance(ctx, transaction)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func (s *storage) InputTransaction(ctx context.Context, transaction *models.Transaction) (float64, error) {
	balance, err := s.updateBalance(ctx, transaction)
	if err != nil {
		return 0, err
	}

	return balance, err
}

func (s *storage) NewTransaction(ctx context.Context, transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions(user_id, status, type, amount, date)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.conn.ExecContext(ctx, query,
		transaction.UserID,
		transaction.Status,
		transaction.Type,
		transaction.Amount,
		transaction.Date,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) UnhandledTransactions(ctx context.Context) ([]*models.Transaction, error) {
	query := `
		SELECT id, user_id, status, type, amount, date
		FROM transactions
		WHERE status = $1
	`

	var transactions []*models.Transaction
	if err := s.conn.SelectContext(ctx, &transactions, query, int(models.InProcessing)); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s storage) updateBalance(ctx context.Context, transaction *models.Transaction) (float64, error) {
	queryUpdateTransaction := `
		UPDATE transactions
		SET status = $1
		WHERE user_id = $2
	`

	queryUpdateBalance := `
		UPDATE balances
		SET balance = balance + $1
		WHERE user_id = $2
		RETURNING balance
	`

	tx, err := s.conn.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, queryUpdateTransaction, models.Processed, transaction.UserID)
	if err != nil {
		return 0, err
	}

	var balance float64
	err = tx.QueryRowContext(ctx, queryUpdateBalance, transaction.Amount, transaction.UserID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return balance, nil
}

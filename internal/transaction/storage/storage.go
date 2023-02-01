package storage

import (
	"bwg_test/internal/transaction/models"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"strings"
)

type storage struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) *storage {
	return &storage{
		conn: conn,
	}
}

func (s *storage) OutputTransaction(ctx context.Context, transaction *models.Transaction) error {
	amount := transaction.Amount * -1
	_, err := s.updateBalance(ctx, transaction, amount)
	if err != nil {
		//
		if strings.EqualFold(err.Error(), models.ErrPositiveAmount.Error()) {
			if err = s.addAttempt(ctx, transaction); err != nil {
				return err
			}
		}
		return err
	}

	return nil
}

func (s *storage) InputTransaction(ctx context.Context, transaction *models.Transaction) error {
	_, err := s.updateBalance(ctx, transaction, transaction.Amount)
	if err != nil {
		if strings.EqualFold(err.Error(), models.ErrPositiveAmount.Error()) {
			if err = s.addAttempt(ctx, transaction); err != nil {
				return err
			}
		}
		return err
	}

	return err
}

func (s *storage) addAttempt(ctx context.Context, transaction *models.Transaction) error {
	query := `
		UPDATE transactions
		SET attempts = attempts + 1
		WHERE transactions.id = $1
		RETURNING attempts
	`

	if err := s.conn.QueryRowContext(ctx, query, transaction.ID).Scan(&transaction.Attempts); err != nil {
		return err
	}

	if transaction.Attempts == models.MaxAttemptsCount {
		if err := s.DeleteTransaction(ctx, transaction); err != nil {
			return err
		}
	}

	return nil
}

func (s *storage) DeleteTransaction(ctx context.Context, transaction *models.Transaction) error {
	query := `
		UPDATE transactions
		SET status = $1
		WHERE id = $2
	`

	if _, err := s.conn.ExecContext(ctx, query, models.ProcessingErr, transaction.ID); err != nil {
		return err
	}

	return nil
}

func (s *storage) NewTransaction(ctx context.Context, transaction *models.Transaction) error {
	query := `
		INSERT INTO transactions(user_id, attempts, status, type, amount)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.conn.ExecContext(ctx, query,
		transaction.UserID,
		transaction.Attempts,
		transaction.Status,
		transaction.Type,
		transaction.Amount,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) UnhandledTransactions(ctx context.Context) ([]*models.Transaction, error) {
	query := `
		SELECT id, user_id, attempts, status, type, amount, date
		FROM transactions
		WHERE status = $1
		ORDER BY date
	`

	var transactions []*models.Transaction
	if err := s.conn.SelectContext(ctx, &transactions, query, int(models.InProcessing)); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *storage) GetTransactions(ctx context.Context, userID int) ([]*models.Transaction, error) {
	query := `
		SELECT user_id, attempts, status, type, amount, date
		FROM transactions
		WHERE user_id = $1
	`

	var transactions []*models.Transaction
	if err := s.conn.SelectContext(ctx, &transactions, query, userID); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *storage) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	query := `
		SELECT user_id, balance
		FROM balances
		WHERE user_id = $1
	`

	var balance models.Balance
	if err := s.conn.GetContext(ctx, &balance, query, userID); err != nil {
		return nil, err
	}
	return &balance, nil
}

func (s *storage) updateBalance(ctx context.Context, transaction *models.Transaction, amount float64) (float64, error) {
	queryUpdateTransaction := `
		UPDATE transactions
		SET status = $1,
		    attempts = attempts + 1
		WHERE id = $2
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

	_, err = tx.ExecContext(ctx, queryUpdateTransaction, models.Processed, transaction.ID)
	if err != nil {
		return 0, err
	}

	var balance float64
	err = tx.QueryRowContext(ctx, queryUpdateBalance, amount, transaction.UserID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return balance, nil
}

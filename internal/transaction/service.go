package transaction

import (
	"bwg_test/internal/transaction/models"
	"context"
	"github.com/rs/zerolog"
	"sync"
)

type IStorage interface {
	InputTransaction(ctx context.Context, transaction *models.Transaction) error
	OutputTransaction(ctx context.Context, transaction *models.Transaction) error
	NewTransaction(ctx context.Context, transaction *models.Transaction) error
	DeleteTransaction(ctx context.Context, transaction *models.Transaction) error
	UnhandledTransactions(ctx context.Context) ([]*models.Transaction, error)
	GetTransactions(ctx context.Context, userID int) ([]*models.Transaction, error)
	GetBalance(ctx context.Context, userID int) (*models.Balance, error)
}

type service struct {
	logger  zerolog.Logger
	storage IStorage
	wg      *sync.WaitGroup
}

func New(ctx context.Context, logger zerolog.Logger, storage IStorage) *service {
	svc := &service{
		logger:  logger,
		storage: storage,
		wg:      &sync.WaitGroup{},
	}

	go svc.GetUnhandledTransactions(ctx)

	return svc
}

func (s *service) GetUnhandledTransactions(ctx context.Context) error {
	for {
		transactions, err := s.storage.UnhandledTransactions(ctx)
		if err != nil {
			s.logger.Error().Err(err).Msg("failed to get unhandled tx")
			return err
		}

		s.wg.Add(1)
		go func() {
			if err = s.transactionHandler(ctx, transactions); err != nil {
				s.logger.Error().Err(err).Msg("transaction handler err")
			}
		}()

		s.wg.Wait()
	}
}

func (s *service) transactionHandler(ctx context.Context, transactions []*models.Transaction) error {
	if len(transactions) > models.MaxTransactionCount {
		s.wg.Add(1)
		go func() {
			if err := s.transactionHandler(ctx, transactions[models.MaxTransactionCount:]); err != nil {
				s.logger.Error().Err(err).Msg("transaction handler")
			}
		}()
		transactions = transactions[:models.MaxTransactionCount]
	}

	for _, v := range transactions {
		if v.Attempts >= models.MaxAttemptsCount {
			if err := s.storage.DeleteTransaction(ctx, v); err != nil {
				s.logger.Error().Err(err).Msg("failed to delete transaction")
			}
			continue
		}

		switch v.Type {
		case int(models.InputType):
			if err := s.storage.InputTransaction(ctx, v); err != nil {
				s.logger.Error().Err(err).Msg("filed to do input transaction")
			}
			break
		case int(models.OutputType):
			if err := s.storage.OutputTransaction(ctx, v); err != nil {
				s.logger.Error().Err(err).Msg("filed to do output transaction")
			}
			break
		}
	}

	s.wg.Done()
	return nil
}

func (s *service) Input(ctx context.Context, transaction *models.Transaction) error {
	transaction.Attempts = 0
	transaction.Status = int(models.InProcessing)
	transaction.Type = int(models.InputType)

	if err := s.validate(transaction); err != nil {
		s.logger.Error().Err(err).Msg("failed to validate")
		return err
	}

	if err := s.storage.NewTransaction(ctx, transaction); err != nil {
		s.logger.Error().Err(err).Msg("failed to create new transaction")
		return err
	}

	return nil
}

func (s *service) Output(ctx context.Context, transaction *models.Transaction) error {
	transaction.Attempts = 0
	transaction.Status = int(models.InProcessing)
	transaction.Type = int(models.OutputType)

	if err := s.validate(transaction); err != nil {
		s.logger.Error().Err(err).Msg("failed to validate")
		return err
	}

	if err := s.storage.NewTransaction(ctx, transaction); err != nil {
		s.logger.Error().Err(err).Msg("failed to create new transaction")
		return err
	}

	return nil
}

func (s *service) GetTransactions(ctx context.Context, userID int) ([]*models.Transaction, error) {
	transactions, err := s.storage.GetTransactions(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get transactions")
		return nil, err
	}

	return transactions, nil
}

func (s *service) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	balance, err := s.storage.GetBalance(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get balance")
		return nil, err
	}

	return balance, nil
}

func (s *service) validate(transaction *models.Transaction) error {
	if transaction.Amount <= 0 {
		return models.ErrWrongAmount
	}

	return nil
}

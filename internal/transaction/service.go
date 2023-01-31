package transaction

import (
	"bwg_test/internal/transaction/models"
	"context"
	"github.com/rs/zerolog"
	"sync"
)

const maxTransactionCount = 10

type IStorage interface {
	InputTransaction(ctx context.Context, transaction *models.Transaction) (float64, error)
	OutputTransaction(ctx context.Context, transaction *models.Transaction) (float64, error)
	NewTransaction(ctx context.Context, transaction *models.Transaction) error
	UnhandledTransactions(ctx context.Context) ([]*models.Transaction, error)
}

type service struct {
	logger  zerolog.Logger
	storage IStorage
	users   map[int]*models.UserTransaction
	wg      *sync.WaitGroup
}

func New(ctx context.Context, logger zerolog.Logger, storage IStorage) *service {
	svc := &service{
		logger:  logger,
		storage: storage,
		users:   make(map[int]*models.UserTransaction),
		wg:      &sync.WaitGroup{},
	}

	go svc.GetUnhandledTransactions(ctx)

	return svc
}

func (s *service) Input(ctx context.Context, transaction *models.Transaction) error {
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
	if len(transactions) > maxTransactionCount {
		s.wg.Add(1)
		//TODO: сделать обработчик ошибки чтобы подсчитывал кол-во неудачных попыток
		go s.transactionHandler(ctx, transactions[maxTransactionCount:])
		transactions = transactions[:maxTransactionCount]
	}

	for _, v := range transactions {
		switch v.Type {
		case int(models.InputType):
			if _, err := s.storage.InputTransaction(ctx, v); err != nil {
				s.logger.Error().Err(err).Msg("filed to do input transaction")
				return err
			}
			break
		case int(models.OutputType):
			if _, err := s.storage.OutputTransaction(ctx, v); err != nil {
				s.logger.Error().Err(err).Msg("filed to do output transaction")
				return err
			}
			break
		}
	}

	s.wg.Done()
	return nil
}

func (s *service) validate(transaction *models.Transaction) error {
	if transaction.Amount <= 0 {
		return models.ErrWrongAmount
	}

	return nil
}

//transactions, ok := s.users[transaction.UserID]
//if ok {
//if transactions.Active {
//<-transactions.Ch
//balance, err := s.storage.InputTransaction(ctx, transaction)
//if err != nil {
//return 0, err
//}
//
////amount := transaction.Amount
////balance := amount + 10
////<-time.After(time.Duration(amount) * time.Second)
//return balance, nil
//}
//}

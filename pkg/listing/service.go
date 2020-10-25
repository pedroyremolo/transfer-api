package listing

import "context"

type Service interface {
	GetAccountBalanceByID(ctx context.Context, id string) (float64, error)
	GetAccountByCPF(ctx context.Context, cpf string) (Account, error)
	GetAccounts(ctx context.Context) ([]Account, error)
	GetTransfersByAccountID(ctx context.Context, id string) (AccountTransfers, error)
}

type Repository interface {
	GetAccountByID(ctx context.Context, id string) (Account, error)
	GetAccountByCPF(ctx context.Context, cpf string) (Account, error)
	GetAccounts(ctx context.Context) ([]Account, error)
	GetTransfersByKey(ctx context.Context, key string, value string) ([]Transfer, error)
}

type service struct {
	r Repository
}

func NewService(repository Repository) Service {
	return &service{repository}
}

func (s *service) GetAccountBalanceByID(ctx context.Context, id string) (float64, error) {
	account, err := s.r.GetAccountByID(ctx, id)
	if err != nil {
		return 0, err
	}
	return account.Balance, nil
}

func (s *service) GetAccounts(ctx context.Context) ([]Account, error) {
	return s.r.GetAccounts(ctx)
}

func (s *service) GetAccountByCPF(ctx context.Context, cpf string) (Account, error) {
	return s.r.GetAccountByCPF(ctx, cpf)
}

func (s *service) GetTransfersByAccountID(ctx context.Context, id string) (AccountTransfers, error) {
	var sentTransfers, receivedTransfers []Transfer
	var err error

	sentTransfers, err = s.r.GetTransfersByKey(ctx, "account_origin_id", id)
	if err != nil {
		return AccountTransfers{}, err
	}
	receivedTransfers, err = s.r.GetTransfersByKey(ctx, "account_destination_id", id)
	if err != nil {
		return AccountTransfers{}, err
	}

	return AccountTransfers{Sent: sentTransfers, Received: receivedTransfers}, nil
}

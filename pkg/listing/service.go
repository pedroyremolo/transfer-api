package listing

import "context"

type Service interface {
	GetAccountBalanceByID(ctx context.Context, id string) (float64, error)
	GetAccounts(ctx context.Context) ([]Account, error)
}

type Repository interface {
	GetAccountByID(ctx context.Context, id string) (Account, error)
	GetAccounts(ctx context.Context) ([]Account, error)
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

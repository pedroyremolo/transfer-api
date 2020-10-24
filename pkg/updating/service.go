package updating

import "context"

type Service interface {
	UpdateAccounts(ctx context.Context, accounts ...Account) error
}

type Repository interface {
	UpdateAccounts(ctx context.Context, accounts []Account) error
}

type service struct {
	r Repository
}

func (s *service) UpdateAccounts(ctx context.Context, accounts ...Account) error {
	return s.r.UpdateAccounts(ctx, accounts)
}

func NewService(repository Repository) Service {
	return &service{r: repository}
}

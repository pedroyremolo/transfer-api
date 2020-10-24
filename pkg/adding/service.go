package adding

import (
	"context"
	"time"
)

type Service interface {
	AddAccount(ctx context.Context, account Account) (string, error)
	AddTransfer(ctx context.Context, transfer Transfer) (string, error)
}

type Repository interface {
	AddAccount(ctx context.Context, account Account) (string, error)
	AddTransfer(ctx context.Context, transfer Transfer) (string, error)
}

type service struct {
	r Repository
}

func (s *service) AddAccount(ctx context.Context, account Account) (string, error) {
	account.CreatedAt = time.Now().UTC()
	id, err := s.r.AddAccount(ctx, account)
	return id, err
}

func (s *service) AddTransfer(ctx context.Context, transfer Transfer) (string, error) {
	transfer.CreatedAt = time.Now().UTC()
	id, err := s.r.AddTransfer(ctx, transfer)
	return id, err
}

func NewService(r Repository) Service {
	return &service{r}
}

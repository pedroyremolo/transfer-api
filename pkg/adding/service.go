package adding

import (
	"context"
	"time"

	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/sirupsen/logrus"
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
	r   Repository
	log *logrus.Logger
}

func (s *service) AddAccount(ctx context.Context, account Account) (string, error) {
	s.log.Infof("adding account %v", account)
	account.CreatedAt = time.Now().UTC()
	id, err := s.r.AddAccount(ctx, account)
	if err != nil {
		s.log.Errorf("err %v when adding account to repository", err)
	}
	s.log.Infof("account %v added with success", account)
	return id, err
}

func (s *service) AddTransfer(ctx context.Context, transfer Transfer) (string, error) {
	s.log.Infof("adding transfer %v", transfer)
	transfer.CreatedAt = time.Now().UTC()
	id, err := s.r.AddTransfer(ctx, transfer)
	if err != nil {
		s.log.Errorf("err %v when adding transfer to repository", err)
	}
	s.log.Infof("transfer %v added with success", transfer)
	return id, err
}

func NewService(r Repository) Service {
	return &service{r, lgr.NewDefaultLogger()}
}

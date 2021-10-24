package updating

import (
	"context"

	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/sirupsen/logrus"
)

type Service interface {
	UpdateAccounts(ctx context.Context, accounts ...Account) error
}

type Repository interface {
	UpdateAccounts(ctx context.Context, accounts []Account) error
}

type service struct {
	r   Repository
	log *logrus.Logger
}

func (s *service) UpdateAccounts(ctx context.Context, accounts ...Account) error {
	if err := s.r.UpdateAccounts(ctx, accounts); err != nil {
		s.log.Errorf("Err %v when updating %v", err, accounts)
		return err
	}
	return nil
}

func NewService(repository Repository) Service {
	return &service{r: repository, log: lgr.NewDefaultLogger()}
}

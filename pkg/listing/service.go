package listing

import (
	"context"

	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/sirupsen/logrus"
)

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
	r   Repository
	log *logrus.Logger
}

func NewService(repository Repository) Service {
	return &service{repository, lgr.NewDefaultLogger()}
}

func (s *service) GetAccountBalanceByID(ctx context.Context, id string) (float64, error) {
	s.log.Infof("Getting account balance by id %s", id)
	account, err := s.r.GetAccountByID(ctx, id)
	if err != nil {
		s.log.Errorf("Err %v occurred when getting account", err)
		return 0, err
	}
	s.log.Infof("Getting account balance by id %s", id)
	return account.Balance, nil
}

func (s *service) GetAccounts(ctx context.Context) ([]Account, error) {
	s.log.Info("Retrieving all accounts")
	accounts, err := s.r.GetAccounts(ctx)
	if err != nil {
		s.log.Errorf("Err %v occurred when retrieving all accounts", err)
		return nil, err
	}
	return accounts, nil
}

func (s *service) GetAccountByCPF(ctx context.Context, cpf string) (Account, error) {
	s.log.Infof("Retrieving account by CPF %s", cpf)
	account, err := s.r.GetAccountByCPF(ctx, cpf)
	if err != nil {
		s.log.Errorf("Err %v occurred when retrieving account by CPF %s", err, cpf)
		return Account{}, err
	}
	return account, nil
}

func (s *service) GetTransfersByAccountID(ctx context.Context, id string) (AccountTransfers, error) {
	s.log.Infof("Retrieving transfers of account %s", id)
	var sentTransfers, receivedTransfers []Transfer
	var err error

	sentTransfers, err = s.r.GetTransfersByKey(ctx, "account_origin_id", id)
	if err != nil {
		s.log.Errorf("Err %v when retrieving sent transfers from acc %s", err, id)
		return AccountTransfers{}, err
	}
	receivedTransfers, err = s.r.GetTransfersByKey(ctx, "account_destination_id", id)
	if err != nil {
		s.log.Errorf("Err %v when retrieving received transfers from acc %s", err, id)
		return AccountTransfers{}, err
	}

	return AccountTransfers{Sent: sentTransfers, Received: receivedTransfers}, nil
}

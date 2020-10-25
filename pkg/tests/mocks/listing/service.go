package listing

import (
	"context"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
)

type MockService struct {
	Balance          float64
	Accounts         []listing.Account
	Account          listing.Account
	AccountTransfers listing.AccountTransfers
	CallsToFail      int
	Err              error
}

func (s *MockService) GetAccountBalanceByID(_ context.Context, _ string) (float64, error) {
	var err error
	s.CallsToFail--
	if s.CallsToFail <= 0 {
		err = s.Err
	}
	return s.Balance, err
}

func (s *MockService) GetAccounts(_ context.Context) ([]listing.Account, error) {
	return s.Accounts, s.Err
}

func (s *MockService) GetAccountByCPF(_ context.Context, _ string) (listing.Account, error) {
	return s.Account, s.Err
}

func (s *MockService) GetTransfersByAccountID(ctx context.Context, id string) (listing.AccountTransfers, error) {
	return s.AccountTransfers, s.Err
}

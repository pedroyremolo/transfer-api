package adding

import (
	"context"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
)

type MockService struct {
	ID  string
	Err error
}

func (s *MockService) AddAccount(_ context.Context, _ adding.Account) (string, error) {
	return s.ID, s.Err
}

func (s *MockService) AddTransfer(_ context.Context, _ adding.Transfer) (string, error) {
	return s.ID, s.Err
}

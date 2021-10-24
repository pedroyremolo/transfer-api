package updating

import (
	"context"

	"github.com/pedroyremolo/transfer-api/pkg/updating"
)

type MockService struct {
	Err error
}

func (m *MockService) UpdateAccounts(_ context.Context, _ ...updating.Account) error {
	return m.Err
}

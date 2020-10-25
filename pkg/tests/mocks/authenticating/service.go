package authenticating

import (
	"context"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
)

type MockService struct {
	Token authenticating.Token
	Err   error
}

func (m *MockService) Sign(_ context.Context, _ authenticating.Login, _ string, _ string) (authenticating.Token, error) {
	return m.Token, m.Err
}

func (m *MockService) Verify(_ context.Context, _ string) (authenticating.Token, error) {
	return m.Token, m.Err
}

package rest

import (
	"context"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"testing"
)

func TestHandler(t *testing.T) {
	s := &mockService{}
	handler := Handler(s)

	if handler == nil {
		t.Errorf("Expected an implementation of http.Handler, got %s", handler)
	}
}

type mockService struct {
	callCount int
}

func (m *mockService) AddAccount(_ context.Context, _ adding.Account) (string, error) {
	m.callCount++
	return "", nil
}

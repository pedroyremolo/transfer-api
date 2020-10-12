package adding

import (
	"testing"
	"time"
)

func TestAddAccount(t *testing.T) {
	zeroTime := time.Time{}
	a := Account{
		Name:    "Gopher",
		CPF:     "11111111030",
		Secret:  "g0rul&zz",
		Balance: 8000.00,
	}

	mockRepo := new(mockStorage)
	s := NewService(mockRepo)

	err := s.AddAccount(a)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if mockRepo.a.CreatedAt == zeroTime {
		t.Errorf("Expected account with CreatedAt near %s, got %s", time.Now().UTC(), mockRepo.a.CreatedAt)
	}
}

type mockStorage struct {
	a Account
}

func (m *mockStorage) AddAccount(account Account) error {
	m.a = account

	return nil
}

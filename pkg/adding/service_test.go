package adding

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

var zeroTime = time.Time{}

func TestAddAccount(t *testing.T) {

	a := Account{
		Name:      "Gopher",
		CPF:       "11111111030",
		Secret:    []byte("g0rul&zz"),
		Balance:   8000.00,
		CreatedAt: time.Time{},
	}

	mockRepo := new(mockStorage)
	s := NewService(mockRepo)
	ctx := context.TODO()
	id, err := s.AddAccount(ctx, a)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if mockRepo.a.CreatedAt == zeroTime {
		t.Errorf("Expected account with CreatedAt near %s, got %s", time.Now().UTC(), mockRepo.a.CreatedAt)
	}

	if id != mockRepo.oid.Hex() {
		t.Errorf("Expected id %s, got %s", id, mockRepo.oid.Hex())
	}
}

type mockStorage struct {
	a   Account
	oid primitive.ObjectID
}

func (m *mockStorage) AddAccount(_ context.Context, account Account) (string, error) {
	m.a = account
	m.oid = primitive.NewObjectID()
	return m.oid.Hex(), nil
}

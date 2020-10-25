package adding

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
	"time"
)

var zeroTime = time.Time{}

func TestService_AddAccount(t *testing.T) {
	tt := []struct {
		name        string
		account     Account
		expectedErr error
	}{
		{
			name: "When successfully adds",
			account: Account{
				Name:      "Gopher",
				CPF:       "11111111030",
				Secret:    []byte("g0rul&zz"),
				Balance:   8000.00,
				CreatedAt: time.Time{},
			},
		},
		{
			name: "When an error occurs",
			account: Account{
				Name:      "Gopher",
				CPF:       "11111111030",
				Secret:    []byte("g0rul&zz"),
				Balance:   8000.00,
				CreatedAt: time.Time{},
			},
			expectedErr: errors.New("foo"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockStorage)
			mockRepo.expectedErr = tc.expectedErr
			s := NewService(mockRepo)
			ctx := context.TODO()
			id, err := s.AddAccount(ctx, tc.account)

			if err != nil && tc.expectedErr == nil {
				t.Errorf("Expected nil, got %s", err)
			}

			if mockRepo.a.CreatedAt == zeroTime && tc.expectedErr == nil {
				t.Errorf("Expected account with CreatedAt near %s, got %s", time.Now().UTC(), mockRepo.a.CreatedAt)
			}

			if id != mockRepo.oid.Hex() && tc.expectedErr == nil {
				t.Errorf("Expected id %s, got %s", id, mockRepo.oid.Hex())
			}
		})
	}
}

func TestService_AddTransfer(t *testing.T) {
	tt := []struct {
		name        string
		transfer    Transfer
		expectedErr error
	}{
		{
			name: "When successfully adds",
			transfer: Transfer{
				OriginAccountID:      "4f89a4fs9864a",
				DestinationAccountID: "fas64fa684fa9",
				Amount:               50.00,
				CreatedAt:            time.Time{},
			},
		},
		{
			name: "When an error occurs",
			transfer: Transfer{
				OriginAccountID:      "4f89a4fs9864a",
				DestinationAccountID: "fas64fa684fa9",
				Amount:               50.00,
				CreatedAt:            time.Time{},
			},
			expectedErr: errors.New("foo"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mockStorage)
			mockRepo.expectedErr = tc.expectedErr
			s := NewService(mockRepo)
			ctx := context.TODO()
			id, err := s.AddTransfer(ctx, tc.transfer)

			if err != nil && tc.expectedErr == nil {
				t.Errorf("Expected nil, got %s", err)
			}

			if mockRepo.t.CreatedAt == zeroTime && tc.expectedErr == nil {
				t.Errorf("Expected account with CreatedAt near %s, got %s", time.Now().UTC(), mockRepo.t.CreatedAt)
			}

			if id != mockRepo.oid.Hex() && tc.expectedErr == nil {
				t.Errorf("Expected id %s, got %s", id, mockRepo.oid.Hex())
			}
		})
	}
}

type mockStorage struct {
	a           Account
	t           Transfer
	expectedErr error
	oid         primitive.ObjectID
}

func (m *mockStorage) AddAccount(_ context.Context, account Account) (string, error) {
	m.a = account
	return m.idOrErr()
}

func (m *mockStorage) AddTransfer(_ context.Context, transfer Transfer) (string, error) {
	m.t = transfer
	return m.idOrErr()
}

func (m *mockStorage) idOrErr() (string, error) {
	m.oid = primitive.NewObjectID()

	if m.expectedErr != nil {
		return "", m.expectedErr
	}

	return m.oid.Hex(), nil
}

package updating

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestService_UpdateAccounts(t *testing.T) {
	tt := []struct {
		name        string
		accounts    []Account
		expectedErr error
	}{
		{
			name:        "When updating occurs successfully",
			accounts:    []Account{{ID: "4896a4fs98a", Balance: 56.98}, {ID: "4f98a49f8", Balance: 89.63}},
			expectedErr: nil,
		},
		{
			name:        "When updating occurs successfully",
			accounts:    []Account{{ID: "4896a4fs98a", Balance: 56.98}, {ID: "4f98a49f8", Balance: 89.63}},
			expectedErr: errors.New("foo"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := new(mockRepository)
			r.expectedErr = tc.expectedErr

			s := NewService(r)

			err := s.UpdateAccounts(context.TODO(), tc.accounts...)

			if err != tc.expectedErr {
				t.Errorf("UpdateAccounts() err = %v, expected %v", err, tc.expectedErr)
			}

			if !reflect.DeepEqual(r.accounts, tc.accounts) && tc.expectedErr == nil {
				t.Errorf("Expected to update %v, got %v", tc.accounts, r.accounts)
			}
		})
	}
}

type mockRepository struct {
	accounts    []Account
	expectedErr error
}

func (m *mockRepository) UpdateAccounts(_ context.Context, accounts []Account) error {
	if m.expectedErr != nil {
		return m.expectedErr
	}
	m.accounts = append(m.accounts, accounts...)
	return nil
}

package listing

import (
	"context"
	"errors"
	"testing"
)

func TestService_GetAccountBalanceByID(t *testing.T) {
	tt := []struct {
		name       string
		id         string
		repository *mockListingRepository
	}{
		{
			name: "When runs smoothly",
			id:   "4d6as4d6a84d6as4wq4",
			repository: &mockListingRepository{
				expectedAccount: Account{Balance: 42.42},
				expectedError:   nil,
			},
		},
		{
			name: "When can't find an account with the given ID",
			repository: &mockListingRepository{
				expectedAccount: Account{},
				expectedError:   errors.New("couldn't find the informed account"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := NewService(tc.repository)
			balance, err := s.GetAccountBalanceByID(context.TODO(), tc.id)
			if err != tc.repository.expectedError {
				t.Errorf("Expected err %s; got %s", tc.repository.expectedError, err)
			}
			if balance != tc.repository.expectedAccount.Balance {
				t.Errorf("Expected balance %.2f; got %.2f", tc.repository.expectedAccount.Balance, balance)
			}
		})
	}
}

type mockListingRepository struct {
	expectedAccount Account
	expectedError   error
}

func (m *mockListingRepository) GetAccountByID(_ context.Context, _ string) (Account, error) {
	return m.expectedAccount, m.expectedError
}

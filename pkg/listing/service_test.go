package listing

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"
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

func TestService_GetAccounts(t *testing.T) {
	currentTime := time.Now().UTC()
	tt := []struct {
		name       string
		repository *mockListingRepository
	}{
		{
			name: "When runs smoothly",
			repository: &mockListingRepository{
				expectedAccounts: []Account{
					{
						ID:        "g4a68vf6a4g96ws84g",
						Name:      "Monkey D. Luffy",
						CPF:       "11111111030",
						Secret:    "onepiece42",
						Balance:   100000.00,
						CreatedAt: &currentTime,
					},
					{
						ID:        "8h964dsa6gb51qa98w1698",
						Name:      "Harry Potter",
						CPF:       "95360976055",
						Secret:    "rh934h@",
						Balance:   40000.00,
						CreatedAt: &currentTime,
					},
				},
			},
		},
		{
			name: "When no account was found",
			repository: &mockListingRepository{
				expectedAccounts: []Account{},
			},
		},
		{
			name: "When an error is emitted by repository",
			repository: &mockListingRepository{
				expectedAccounts: []Account{},
				expectedError:    errors.New("foo"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := NewService(tc.repository)
			accounts, err := s.GetAccounts(context.TODO())
			if err != tc.repository.expectedError {
				t.Errorf("Expected err %s; got %s", tc.repository.expectedError, err)
			}
			if len(accounts) != len(tc.repository.expectedAccounts) {
				t.Errorf("Expected accounts %v; got %v", tc.repository.expectedAccounts, accounts)
			}
		})
	}
}

func TestService_GetAccountByCPF(t *testing.T) {
	tt := []struct {
		name       string
		cpf        string
		repository *mockListingRepository
	}{
		{
			name: "When runs smoothly",
			cpf:  "4d6as4d6a84d6as4wq4",
			repository: &mockListingRepository{
				expectedAccount: Account{
					ID:      "g4a68vf6a4g96ws84g",
					Name:    "Monkey D. Luffy",
					CPF:     "11111111030",
					Secret:  "onepiece42",
					Balance: 100000.00,
				},
			},
		},
		{
			name: "When can't find an account with the given cpf",
			repository: &mockListingRepository{
				expectedError: errors.New("couldn't find the informed account"),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := NewService(tc.repository)
			account, err := s.GetAccountByCPF(context.TODO(), tc.cpf)
			if err != tc.repository.expectedError {
				t.Errorf("Expected err %s; got %s", tc.repository.expectedError, err)
			}
			if account.CPF != tc.repository.expectedAccount.CPF {
				t.Errorf("Expected cpf %s; got %s", tc.repository.expectedAccount.CPF, account.CPF)
			}
		})
	}
}

func TestService_GetAccountTransfersByID(t *testing.T) {
	accId := "wr896q4c3ar46"
	transfers := []Transfer{{
		ID:                   "6f5a4f56a",
		OriginAccountID:      accId,
		DestinationAccountID: "r4wq861a65f8qr6",
		Amount:               56.23,
		CreatedAt:            time.Time{},
	}, {
		ID:                   "9w8qe74981q",
		OriginAccountID:      accId,
		DestinationAccountID: "r9849a8c96a8w6",
		Amount:               23.23,
		CreatedAt:            time.Time{},
	}}
	tt := []struct {
		name       string
		id         string
		repository *mockListingRepository
		want       AccountTransfers
		wantErr    bool
	}{
		{
			name: "When successfully retrieve account transfers",
			id:   accId,
			repository: &mockListingRepository{
				expectedTransfers: transfers,
			},
			want: AccountTransfers{
				Sent:     transfers,
				Received: transfers,
			},
		},
		{
			name: "When an error occurs when retrieving sent transfers",
			id:   accId,
			repository: &mockListingRepository{
				expectedError: errors.New("foo"),
				callsToFail:   1,
			},
			want:    AccountTransfers{},
			wantErr: true,
		},
		{
			name: "When an error occurs when retrieving received transfers",
			id:   accId,
			repository: &mockListingRepository{
				expectedError: errors.New("foo"),
				callsToFail:   2,
			},
			want:    AccountTransfers{},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := NewService(tc.repository)
			got, err := s.GetTransfersByAccountID(context.TODO(), tc.id)

			if err != nil && !tc.wantErr {
				t.Errorf("GetTransfersByAccountID() err = %v; want err %v", err, tc.wantErr)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected account tranfers %v, got %v", tc.want, got)
			}
		})
	}
}

type mockListingRepository struct {
	expectedAccounts  []Account
	expectedAccount   Account
	expectedTransfers []Transfer
	callsToFail       int
	expectedError     error
}

func (m *mockListingRepository) GetAccounts(_ context.Context) ([]Account, error) {
	return m.expectedAccounts, m.expectedError
}

func (m *mockListingRepository) GetAccountByID(_ context.Context, _ string) (Account, error) {
	return m.expectedAccount, m.expectedError
}

func (m *mockListingRepository) GetAccountByCPF(_ context.Context, _ string) (Account, error) {
	return m.expectedAccount, m.expectedError
}

func (m *mockListingRepository) GetTransfersByKey(_ context.Context, _ string, _ string) ([]Transfer, error) {
	var err error
	m.callsToFail--
	if m.callsToFail <= 0 {
		err = m.expectedError
	}
	return m.expectedTransfers, err
}

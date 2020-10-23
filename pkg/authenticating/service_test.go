package authenticating

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
)

func TestService_Sign(t *testing.T) {
	type args struct {
		login        Login
		secretDigest string
		clientID     string
	}
	login := Login{
		CPF:    "11111111030",
		Secret: "65416949",
	}
	rightSecretDigestBytes, _ := bcrypt.GenerateFromPassword(
		[]byte(fmt.Sprintf(`"%s"`, login.Secret)),
		bcrypt.DefaultCost,
	)
	rightSecretDigest := string(rightSecretDigestBytes)
	oid := primitive.NewObjectID()
	tt := []struct {
		name       string
		args       args
		repository mockRepository
		gatekeeper mockGatekeeper
		wantErr    bool
	}{
		{
			name: "When everything runs smoothly",
			args: args{
				login:        login,
				secretDigest: rightSecretDigest,
				clientID:     "sa1685fd4w1a489f49asf",
			},
			repository: mockRepository{},
			gatekeeper: mockGatekeeper{
				expectedToken: Token{
					ID:       &oid,
					ClientID: "sa1685fd4w1a489f49asf",
					Digest:   "sa1685fd4w1a489f49asf.fasofapogkapog.gasjkgpoaskgpoa",
				},
			},
		},
		{
			name: "When an error occurs at Login.Password validation",
			args: args{
				login:        login,
				secretDigest: "fooo",
				clientID:     "sa1685fd4w1a489f49asf",
			},
			repository: mockRepository{},
			gatekeeper: mockGatekeeper{
				expectedErr: errors.New("foo"),
			},
			wantErr: true,
		},
		{
			name: "When an error occurs at token sign flow",
			args: args{
				login:        login,
				secretDigest: rightSecretDigest,
				clientID:     "sa1685fd4w1a489f49asf",
			},
			repository: mockRepository{},
			gatekeeper: mockGatekeeper{
				expectedErr: errors.New("foo"),
			},
			wantErr: true,
		},
		{
			name: "When an error occurs at repository insertion",
			args: args{
				login:        login,
				secretDigest: "4f89sa4fd6sa4c1984f",
				clientID:     "sa1685fd4w1a489f49asf",
			},
			repository: mockRepository{
				expectedErr: errors.New("repo err"),
			},
			gatekeeper: mockGatekeeper{
				expectedToken: Token{
					ID:       &oid,
					ClientID: "sa1685fd4w1a489f49asf",
					Digest:   "sa1685fd4w1a489f49asf.fasofapogkapog.gasjkgpoaskgpoa",
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := NewService(&tc.repository, &tc.gatekeeper)
			token, err := s.Sign(nil, tc.args.login, tc.args.secretDigest, tc.args.clientID)
			if (err != nil) != tc.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(tc.gatekeeper.expectedToken, token) && !tc.wantErr {
				t.Errorf("Expected token %v, got %v", tc.gatekeeper.expectedToken, token)
			}
		})
	}
}

func TestService_Verify(t *testing.T) {
	oid := primitive.NewObjectID()
	defaultToken := Token{
		ID:       &oid,
		ClientID: "4a89f4a1s98fa",
		Digest:   "a9ifa09sfamfk90asf.fafajrqr9qkf0mas09f.fqj09fj0ajf0a",
	}
	tt := []struct {
		name        string
		tokenDigest string
		gatekeeper  mockGatekeeper
		repository  mockRepository
		wantErr     bool
	}{
		{
			name:        "When token is ok was issued by this app",
			tokenDigest: defaultToken.Digest,
			gatekeeper: mockGatekeeper{
				expectedToken: defaultToken,
			},
			repository: mockRepository{
				expectedToken: defaultToken,
			},
		},
		{
			name:        "When token is not valid",
			tokenDigest: defaultToken.Digest,
			gatekeeper: mockGatekeeper{
				expectedErr: errors.New("token is invalid"),
			},
			wantErr: true,
		},
		{
			name:        "When token is valid but was not issued here",
			tokenDigest: defaultToken.Digest,
			gatekeeper: mockGatekeeper{
				expectedToken: defaultToken,
			},
			repository: mockRepository{
				expectedErr: errors.New("token not found"),
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			s := NewService(&tc.repository, &tc.gatekeeper)
			token, err := s.Verify(context.TODO(), tc.tokenDigest)

			if err != nil && !tc.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && !reflect.DeepEqual(tc.gatekeeper.expectedToken, token) {
				t.Errorf("Expected token")
			}
		})
	}
}

type mockGatekeeper struct {
	expectedToken Token
	expectedErr   error
}

func (m *mockGatekeeper) Sign(_ string) (Token, error) {
	return m.expectedToken, m.expectedErr
}

func (m *mockGatekeeper) Verify(_ string) (Token, error) {
	return m.expectedToken, m.expectedErr
}

type mockRepository struct {
	expectedToken Token
	expectedErr   error
}

func (m *mockRepository) AddToken(_ context.Context, _ Token) error {
	return m.expectedErr
}

func (m *mockRepository) GetTokenByID(_ context.Context, _ primitive.ObjectID) (Token, error) {
	return m.expectedToken, m.expectedErr
}

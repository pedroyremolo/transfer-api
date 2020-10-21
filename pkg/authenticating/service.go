package authenticating

import (
	"context"
	"errors"
)

var InvalidLoginErr = errors.New("it seems your login credentials are invalid, verify them and try again")

type Service interface {
	Sign(ctx context.Context, login Login, secretDigest string, clientID string) (Token, error)
	Verify(ctx context.Context, inToken Token) bool
}

type Repository interface {
	AddToken(ctx context.Context, token Token) error
}

type Gatekeeper interface {
	Sign(login Login, secretDigest string, clientID string) (Token, error)
	Verify(tokenDigest string) Token
}

type service struct {
	r Repository
	g Gatekeeper
}

func NewService(repository Repository, gatekeeper Gatekeeper) Service {
	return &service{
		repository,
		gatekeeper,
	}
}

func (s *service) Sign(ctx context.Context, login Login, secretDigest string, clientID string) (Token, error) {
	token, err := s.g.Sign(login, secretDigest, clientID)
	if err != nil {
		return Token{}, InvalidLoginErr
	}

	err = s.r.AddToken(ctx, token)
	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func (s *service) Verify(_ context.Context, _ Token) bool {
	panic("implement me")
}

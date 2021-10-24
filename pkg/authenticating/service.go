package authenticating

import (
	"context"
	"errors"
	"fmt"

	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var InvalidLoginErr = errors.New("it seems your login credentials are invalid, verify them and try again")
var ProtectedRouteErr = errors.New("it seems you don't have or didn't pass valid credentials to this route")

type Service interface {
	Sign(ctx context.Context, login Login, secretDigest string, clientID string) (Token, error)
	Verify(ctx context.Context, tokenDigest string) (Token, error)
}

type Repository interface {
	AddToken(ctx context.Context, token Token) error
	GetTokenByID(ctx context.Context, id primitive.ObjectID) (Token, error)
}

type Gatekeeper interface {
	Sign(clientID string) (Token, error)
	Verify(tokenDigest string) (Token, error)
}

type service struct {
	r   Repository
	g   Gatekeeper
	log *logrus.Logger
}

func NewService(repository Repository, gatekeeper Gatekeeper) Service {
	return &service{
		repository,
		gatekeeper,
		lgr.NewDefaultLogger(),
	}
}

func (s *service) Sign(ctx context.Context, login Login, secretDigest string, clientID string) (Token, error) {
	s.log.Infof("Signing token to clientID %s", clientID)
	secret := []byte(fmt.Sprintf(`"%s"`, login.Secret))
	if err := bcrypt.CompareHashAndPassword([]byte(secretDigest), secret); err != nil {
		s.log.Errorf("Err %v occurred when validating login secret", err)
		return Token{}, InvalidLoginErr
	}

	token, err := s.g.Sign(clientID)
	if err != nil {
		s.log.Errorf("Err %v occurred when gatekeeper signs token", err)
		return Token{}, InvalidLoginErr
	}

	err = s.r.AddToken(ctx, token)
	if err != nil {
		s.log.Errorf("Err %v occurred when repo tried to add token", err)
		return Token{}, err
	}

	s.log.Infof("Successfully signed token %v", token)
	return token, nil
}

func (s *service) Verify(ctx context.Context, tokenDigest string) (Token, error) {
	s.log.Infof("Verifying tokenDigest %v", tokenDigest)
	token, err := s.g.Verify(tokenDigest)
	if err != nil {
		s.log.Errorf("Err %v occurred when gatekeeper verified tokenDigest %v", err, tokenDigest)
		return Token{}, err
	}

	_, err = s.r.GetTokenByID(ctx, *token.ID)
	if err != nil {
		s.log.Errorf("Err %v when retrieving token %v from repository", err, token)
		return Token{}, err
	}
	s.log.Infof("Token %v successfully verified", token)
	return token, nil
}

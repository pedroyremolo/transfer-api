package jwt

import (
	"fmt"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

var (
	secret = os.Getenv("APP_JWT_GATEKEEPER_SECRET")
	issuer = os.Getenv("APP_JWT_GATEKEEPER_ISSUER")
)

type Gatekeeper struct {
	hs  *jwt.HMACSHA
	iss string
}

func NewGatekeeperFromEnv() *Gatekeeper {
	return &Gatekeeper{
		hs:  jwt.NewHS256([]byte(secret)),
		iss: issuer,
	}
}

func NewGatekeeper(tokenSecret string, iss string) *Gatekeeper {
	return &Gatekeeper{
		hs:  jwt.NewHS256([]byte(tokenSecret)),
		iss: iss,
	}
}

func (g *Gatekeeper) Sign(login authenticating.Login, secretDigest string, clientID string) (authenticating.Token, error) {
	secret := []byte(fmt.Sprintf(`"%s"`, login.Secret))
	if err := bcrypt.CompareHashAndPassword([]byte(secretDigest), secret); err != nil {
		return authenticating.Token{}, err
	}

	currentTime := time.Now().UTC()
	id := primitive.NewObjectID()
	token, err := jwt.Sign(Token{
		Payload: jwt.Payload{
			Issuer:         g.iss,
			ExpirationTime: jwt.NumericDate(currentTime.Add(time.Minute * 30)),
			IssuedAt:       jwt.NumericDate(currentTime),
			JWTID:          id.Hex(),
		},
		ClientID: clientID,
	}, g.hs)
	if err != nil {
		return authenticating.Token{}, err
	}

	return authenticating.Token{ID: &id, ClientID: clientID, Digest: string(token)}, nil
}

func (g *Gatekeeper) Verify(_ string) authenticating.Token {
	panic("implement me")
}

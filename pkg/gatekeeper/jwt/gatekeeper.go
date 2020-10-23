package jwt

import (
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"time"
)

type Gatekeeper struct {
	hs  *jwt.HMACSHA
	iss string
}

func NewGatekeeperFromEnv() *Gatekeeper {
	secret := os.Getenv("APP_JWT_GATEKEEPER_SECRET")
	issuer := os.Getenv("APP_JWT_GATEKEEPER_ISSUER")

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

func (g *Gatekeeper) Sign(clientID string) (authenticating.Token, error) {
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

func (g *Gatekeeper) Verify(tokenDigest string) (authenticating.Token, error) {
	var jwtToken Token
	var oid primitive.ObjectID
	now := time.Now().UTC()
	expValidator := jwt.ExpirationTimeValidator(now)
	issValidator := jwt.IssuerValidator(g.iss)
	validatePayload := jwt.ValidatePayload(&jwtToken.Payload, issValidator, expValidator)
	_, err := jwt.Verify([]byte(tokenDigest), g.hs, &jwtToken, validatePayload)
	if err != nil {
		return authenticating.Token{}, err
	}

	oid, err = primitive.ObjectIDFromHex(jwtToken.JWTID)
	if err != nil {
		return authenticating.Token{}, err
	}

	token := authenticating.Token{
		ID:       &oid,
		ClientID: jwtToken.ClientID,
		Digest:   tokenDigest,
	}
	return token, nil
}

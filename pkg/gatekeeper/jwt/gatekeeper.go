package jwt

import (
	"os"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Gatekeeper struct {
	hs  *jwt.HMACSHA
	iss string
	log *logrus.Logger
}

func NewGatekeeperFromEnv() *Gatekeeper {
	secret := os.Getenv("APP_JWT_GATEKEEPER_SECRET")
	issuer := os.Getenv("APP_JWT_GATEKEEPER_ISSUER")

	return &Gatekeeper{
		hs:  jwt.NewHS256([]byte(secret)),
		iss: issuer,
		log: lgr.NewDefaultLogger(),
	}
}

func NewGatekeeper(tokenSecret string, iss string) *Gatekeeper {
	return &Gatekeeper{
		hs:  jwt.NewHS256([]byte(tokenSecret)),
		iss: iss,
		log: lgr.NewDefaultLogger(),
	}
}

func (g *Gatekeeper) Sign(clientID string) (authenticating.Token, error) {
	g.log.Infof("Trying to emit a token for clientID %s", clientID)
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
		g.log.Errorf("Error %v when signing token", err)
		return authenticating.Token{}, err
	}

	return authenticating.Token{ID: &id, ClientID: clientID, Digest: string(token)}, nil
}

func (g *Gatekeeper) Verify(tokenDigest string) (authenticating.Token, error) {
	g.log.Infof("Trying to verify tokenDigest %s", tokenDigest)

	var jwtToken Token
	var oid primitive.ObjectID

	now := time.Now().UTC()
	expValidator := jwt.ExpirationTimeValidator(now)
	issValidator := jwt.IssuerValidator(g.iss)
	validatePayload := jwt.ValidatePayload(&jwtToken.Payload, issValidator, expValidator)

	_, err := jwt.Verify([]byte(tokenDigest), g.hs, &jwtToken, validatePayload)
	if err != nil {
		g.log.Errorf("Error %v when verifying tokenDigest %s", err, tokenDigest)
		return authenticating.Token{}, err
	}

	oid, err = primitive.ObjectIDFromHex(jwtToken.JWTID)
	if err != nil {
		g.log.Errorf("Error %v when converting jwti %s to oid", err, jwtToken.JWTID)
		return authenticating.Token{}, err
	}

	token := authenticating.Token{
		ID:       &oid,
		ClientID: jwtToken.ClientID,
		Digest:   tokenDigest,
	}

	g.log.Infof("tokenDigest %s successfully verified", tokenDigest)
	return token, nil
}

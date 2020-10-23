package jwt

import "github.com/gbrlsnchs/jwt/v3"

type Token struct {
	jwt.Payload
	ClientID string `json:"client_id"`
}

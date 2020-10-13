package adding

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrPswIsNotString = errors.New("the informed password is not a string")

// Account is the representation of an account to be added
type Account struct {
	Name      string  `json:"name"`
	CPF       string  `json:"cpf"`
	Secret    secret  `json:"secret"`
	Balance   float64 `json:"balance"`
	CreatedAt time.Time
}

// secret represents the user authentication password
type secret []byte

// UnmarshalJSON takes the secret sent as string and encrypt it
func (s *secret) UnmarshalJSON(b []byte) error {
	var pswStr string
	err := json.Unmarshal(b, &pswStr)

	if err != nil {
		return ErrPswIsNotString
	}

	password, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	*s = password

	return nil
}

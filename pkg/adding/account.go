package adding

import (
	"encoding/json"
	"fmt"
	"github.com/Nhanderu/brdoc"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Account is the representation of an account to be added
type Account struct {
	Name      string  `json:"name"`
	CPF       cpf     `json:"cpf"`
	Secret    secret  `json:"secret"`
	Balance   float64 `json:"balance"`
	CreatedAt time.Time
}

// cpf represents the brazilian id
type cpf string

// secret represents the user authentication password
type secret []byte

// UnmarshalJSON Unmarshaler implementation that takes the string and verify it's a valid CPF
func (c *cpf) UnmarshalJSON(b []byte) error {
	var document string
	err := json.Unmarshal(b, &document)
	if err != nil {
		// TODO Err logging
		return &ErrInvalidAccountField{
			field:   "cpf",
			message: "the informed cpf is not a string",
		}
	}
	if !brdoc.IsCPF(document) {
		// TODO Err logging
		return &ErrInvalidAccountField{
			field:   "cpf",
			message: fmt.Sprintf("%s is not a valid cpf", document),
		}
	}

	*c = cpf(document)

	return nil
}

// UnmarshalJSON Unmarshaler implementation that takes the sent secret, as string, and encrypt it
func (s *secret) UnmarshalJSON(b []byte) error {
	var pswStr string
	err := json.Unmarshal(b, &pswStr)

	if err != nil {
		//TODO Err logging
		return &ErrInvalidAccountField{
			field:   "password",
			message: "the informed password is not a string",
		}
	}

	password, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	*s = password

	return nil
}

type ErrInvalidAccountField struct {
	field   string
	message string
}

func (e *ErrInvalidAccountField) Error() string {
	return fmt.Sprintf("Field %s contains an invalid value: %s", e.field, e.message)
}

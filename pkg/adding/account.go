package adding

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Nhanderu/brdoc"
	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"golang.org/x/crypto/bcrypt"
)

// Account is the representation of an account to be added
type Account struct {
	Name      name    `json:"name"`
	CPF       cpf     `json:"cpf"`
	Secret    secret  `json:"secret"`
	Balance   balance `json:"balance"`
	CreatedAt time.Time
}

type (
	// name represents the user Full Name
	name string
	// cpf represents the brazilian id
	cpf string
	// secret represents the user authentication Password
	secret []byte
	// balance represents the initial account Balance
	balance float64
)

func (n *name) UnmarshalJSON(b []byte) error {
	var log = lgr.NewDefaultLogger()
	var fullName string

	if err := json.Unmarshal(b, &fullName); err != nil {
		log.Errorf("Err %v happened when unmarshal name", err)
		return &ErrInvalidAccountField{
			field:   "name",
			message: "name is not of string type",
		}
	}

	if len(fullName) == 0 {
		log.Error("Name cannot be empty")
		return &ErrInvalidAccountField{
			field:   "name",
			message: "name must be informed, therefore should not be empty",
		}
	}

	*n = name(fullName)
	return nil
}

// UnmarshalJSON Unmarshaler implementation that takes the string and verify it's a valid CPF
func (c *cpf) UnmarshalJSON(b []byte) error {
	var log = lgr.NewDefaultLogger()
	var document string
	err := json.Unmarshal(b, &document)
	if err != nil {
		log.Errorf("Err %v happened when unmarshal cpf", err)
		return &ErrInvalidAccountField{
			field:   "cpf",
			message: "the informed cpf is not a string",
		}
	}
	if !brdoc.IsCPF(document) {
		log.Errorf("Document %v is not a valid cpf", document)
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
	var log = lgr.NewDefaultLogger()
	var pswStr string
	err := json.Unmarshal(b, &pswStr)

	if err != nil {
		log.Errorf("Err %v happened when unmarshal secret", err)
		return &ErrInvalidAccountField{
			field:   "password",
			message: "the informed password is not a string",
		}
	}

	password, err := bcrypt.GenerateFromPassword(b, bcrypt.DefaultCost)

	if err != nil {
		log.Error("An err occurred when encrypting secret")
		return err
	}

	*s = password

	return nil
}

func (bc *balance) UnmarshalJSON(b []byte) error {
	var log = lgr.NewDefaultLogger()
	var incomingBalance float64
	if err := json.Unmarshal(b, &incomingBalance); err != nil {
		log.Errorf("Err %v happened when unmarshal balance", err)
		return &ErrInvalidAccountField{
			field:   "balance",
			message: "the informed balance is not a number",
		}
	}
	if incomingBalance < 0 {
		log.Error("Balance cannot be negative")
		return &ErrInvalidAccountField{
			field:   "balance",
			message: "can't start an account with negative balance",
		}
	}
	*bc = balance(incomingBalance)
	return nil
}

type ErrInvalidAccountField struct {
	field   string
	message string
}

func (e *ErrInvalidAccountField) Error() string {
	if e.field == "" {
		panic("ErrInvalidAccountField.Error usage without field")
	}
	return fmt.Sprintf("Field %s contains an invalid value: %s", e.field, e.message)
}

package adding

import "time"

// Account is the representation of an account to be added
type Account struct {
	Name      string  `json:"name"`
	CPF       string  `json:"cpf"`
	Secret    string  `json:"secret"`
	Balance   float64 `json:"balance"`
	CreatedAt time.Time
}

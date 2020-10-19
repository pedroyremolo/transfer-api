package listing

type Account struct {
	Name    string  `json:"name,omitempty"`
	CPF     string  `json:"cpf,omitempty"`
	Secret  string  `json:"-"`
	Balance float64 `json:"balance"`
}

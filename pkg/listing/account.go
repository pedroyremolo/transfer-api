package listing

type Account struct {
	Name    string  `json:"name"`
	CPF     string  `json:"cpf"`
	Secret  string  `json:"-"`
	Balance float64 `json:"balance"`
}

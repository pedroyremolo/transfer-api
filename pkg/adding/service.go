package adding

import "time"

type Service interface {
	AddAccount(account Account) error
}

type Repository interface {
	AddAccount(account Account) error
}

type service struct {
	r Repository
}

func (s *service) AddAccount(account Account) error {
	account.CreatedAt = time.Now().UTC()
	err := s.r.AddAccount(account)
	return err
}

func NewService(r Repository) Service {
	return &service{r}
}

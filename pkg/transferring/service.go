package transferring

import (
	"errors"
	"math"
	"math/big"
)

var ErrNotEnoughBalance = errors.New("not enough balance to execute this operation")

type Service interface {
	BalanceBetweenAccounts(originBalance float64, destinationBalance float64, amount float64) (newOriBalance float64, newDstBalance float64, err error)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) BalanceBetweenAccounts(oBalance float64, dBalance float64, amount float64) (newOBalance float64, newDBalance float64, err error) {
	preciseAmount := big.NewFloat(amount)
	preciseOBalance := big.NewFloat(oBalance)
	if preciseAmount.Cmp(preciseOBalance) > 0 {
		err = ErrNotEnoughBalance
		return
	}
	newOBalance = math.Round((oBalance-amount)*100) / 100
	newDBalance = math.Round((dBalance+amount)*100) / 100
	return
}

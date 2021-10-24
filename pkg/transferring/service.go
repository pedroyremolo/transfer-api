package transferring

import (
	"errors"
	"math"
	"math/big"

	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/sirupsen/logrus"
)

var ErrNotEnoughBalance = errors.New("not enough balance to execute this operation")

type Service interface {
	BalanceBetweenAccounts(originBalance float64, destinationBalance float64, amount float64) (newOriBalance float64, newDstBalance float64, err error)
}

type service struct {
	log *logrus.Logger
}

func NewService() Service {
	return &service{
		log: lgr.NewDefaultLogger(),
	}
}

func (s *service) BalanceBetweenAccounts(oBalance float64, dBalance float64, amount float64) (newOBalance float64, newDBalance float64, err error) {
	s.log.Infof("Transferring amount %.2f from balance %.2f to balance %.2f", amount, oBalance, dBalance)
	preciseAmount := big.NewFloat(amount)
	preciseOBalance := big.NewFloat(oBalance)
	if preciseAmount.Cmp(preciseOBalance) > 0 {
		s.log.Errorf("Balance %.2f is lower than amount %.2f", oBalance, amount)
		err = ErrNotEnoughBalance
		return
	}
	newOBalance = math.Round((oBalance-amount)*100) / 100
	newDBalance = math.Round((dBalance+amount)*100) / 100
	s.log.Infof(
		"New origin balance %.2f and destination balance %.2f after transferring amount %.2f",
		newOBalance,
		newDBalance,
		amount,
	)
	return
}

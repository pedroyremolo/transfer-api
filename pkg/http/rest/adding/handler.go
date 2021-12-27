package adding

import (
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger  *logrus.Entry
	service adding.Service
}

func NewHandler(logger *logrus.Entry, service adding.Service) Handler {
	return Handler{
		logger:  logger,
		service: service,
	}
}

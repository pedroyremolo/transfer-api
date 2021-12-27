package listing

import (
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger *logrus.Entry

	service listing.Service
}

func NewHandler(logger *logrus.Entry, service listing.Service) Handler {
	return Handler{
		logger:  logger,
		service: service,
	}
}

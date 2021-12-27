package authenticating

import (
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger         *logrus.Entry
	service        authenticating.Service
	listingService listing.Service
}

func NewHandler(logger *logrus.Entry, service authenticating.Service, listingService listing.Service) Handler {
	return Handler{
		logger:         logger,
		service:        service,
		listingService: listingService,
	}
}

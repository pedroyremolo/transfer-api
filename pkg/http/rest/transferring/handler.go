package transferring

import (
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/transferring"
	"github.com/pedroyremolo/transfer-api/pkg/updating"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	logger          *logrus.Entry
	service         transferring.Service
	listingService  listing.Service
	addingService   adding.Service
	updatingService updating.Service
}

func NewHandler(
	logger *logrus.Entry,
	service transferring.Service,
	listingService listing.Service,
	addingService adding.Service,
	updatingService updating.Service,
) Handler {
	return Handler{
		logger:          logger,
		service:         service,
		listingService:  listingService,
		addingService:   addingService,
		updatingService: updatingService,
	}
}

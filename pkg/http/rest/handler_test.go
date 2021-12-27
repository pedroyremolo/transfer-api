package rest

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	am "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/adding"
	aum "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/authenticating"
	lm "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/listing"
	tm "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/transferring"
)

func TestHandler(t *testing.T) {
	log := lgr.NewDefaultLogger()
	logger := logrus.NewEntry(log)

	addingHandlerMock := am.HandlerMock{}
	transferringHandlerMock := tm.HandlerMock{}
	listingHandlerMock := &lm.HandlerMock{}
	authHandlerMock := &aum.HandlerMock{}

	handler := Handler(logger, addingHandlerMock, transferringHandlerMock, authHandlerMock, listingHandlerMock)

	if handler == nil {
		t.Errorf("Expected an implementation of http.Handler, got %s", handler)
	}
}

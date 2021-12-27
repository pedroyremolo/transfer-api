package listing

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"github.com/pedroyremolo/transfer-api/pkg/tests/helpers"
	"github.com/pedroyremolo/transfer-api/pkg/tests/mocks/listing"
	"github.com/sirupsen/logrus"
)

func TestGetAccountBalanceByID(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())
	tt := []struct {
		name             string
		id               string
		service          *listing.MockService
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "When successfully returns",
			id:               "a6sf46af6af",
			service:          &listing.MockService{Balance: 42.42},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"balance":42.42}`,
		},
		{
			name:             "When no account was found with the given id",
			id:               "a6sf46af6af",
			service:          &listing.MockService{Err: mongodb.ErrNoAccountWasFound},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: `{"status_code":404,"message":"no account was found with the given filter parameters"}`,
		},
		{
			name:             "When unexpected errors inside the service occurs",
			id:               "a6sf46af6af",
			service:          &listing.MockService{Err: errors.New("foo")},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"status_code":500,"message":"foo"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHandler(logger, tc.service)
			w := httptest.NewRecorder()
			target := fmt.Sprintf("/accounts/%s/balance", tc.id)
			r := httptest.NewRequest(http.MethodGet, target, nil)

			handler.GetBalanceByID(w, r, httprouter.Params{{
				Key:   "id",
				Value: tc.id,
			}})

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}
			helpers.AssertResponseJSON(t, w, tc.expectedResponse)
		})
	}
}

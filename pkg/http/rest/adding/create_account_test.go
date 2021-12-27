package adding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"github.com/pedroyremolo/transfer-api/pkg/tests/helpers"
	am "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/adding"
	"github.com/sirupsen/logrus"
)

func TestCreateAccount(t *testing.T) {
	log := lgr.NewDefaultLogger()
	logger := logrus.NewEntry(log)

	tt := []struct {
		name                string
		service             *am.MockService
		reqBodyJSON         string
		expectedStatus      int
		expectedErrResponse string
	}{
		{
			name:           "When successfully returns",
			service:        &am.MockService{ID: "a6sf46af6af"},
			reqBodyJSON:    `{"name": "Jane Doe","cpf": "11111111030","secret": "254855","balance": 50.00}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:                "When the sent json is not valid",
			service:             &am.MockService{ID: ""},
			reqBodyJSON:         `{"name": "Jane Doe","cpf": "111111110301","secret": "254855","balance": 50.00}`,
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: `{"status_code":400,"message":"Field cpf contains an invalid value: 111111110301 is not a valid cpf"}`,
		},
		{
			name:                "When the sent cpf is already into DB",
			service:             &am.MockService{Err: mongodb.ErrCPFAlreadyExists},
			reqBodyJSON:         `{"name": "Jane Doe","cpf": "11111111030","secret": "254855","balance": 50.00}`,
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: `{"status_code":400,"message":"this cpf could not be inserted in our DB"}`,
		},
		{
			name:                "When unexpected errors inside the service occurs",
			service:             &am.MockService{Err: errors.New("foo")},
			reqBodyJSON:         `{"name": "Jane Doe","cpf": "11111111030","secret": "254855","balance": 50.00}`,
			expectedStatus:      http.StatusInternalServerError,
			expectedErrResponse: `{"status_code":500,"message":"foo"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHandler(logger, tc.service)

			var reqBody string
			jsonBuffer := bytes.NewBuffer([]byte(tc.reqBodyJSON))
			if err := json.NewEncoder(jsonBuffer).Encode(reqBody); err != nil {
				t.Fatalf("Could not encode %s as json", tc.reqBodyJSON)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/accounts", jsonBuffer)

			handler.CreateAccount(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}
			if len(tc.expectedErrResponse) > 0 {
				helpers.AssertResponseJSON(t, w, tc.expectedErrResponse)
				return
			}
			if w.Header().Get("Location") != fmt.Sprintf("/%s", tc.service.ID) {
				t.Errorf("Expected Location header /%s; got %v", tc.service.ID, w.Header().Get("Location"))
			}
		})
	}
}

package authenticating

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"github.com/pedroyremolo/transfer-api/pkg/tests/helpers"
	aum "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/authenticating"
	lm "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/listing"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestLogin(t *testing.T) {
	tokenOID := primitive.NewObjectID()
	log := lgr.NewDefaultLogger()
	logger := logrus.NewEntry(log)

	account := listing.Account{
		ID:      "hg94gs8a41v685s4g89",
		Name:    "John Doe",
		CPF:     "11111111030",
		Secret:  "894d9a8",
		Balance: 0,
	}
	token := authenticating.Token{
		ID:       &tokenOID,
		ClientID: account.ID,
		Digest:   "e4af98as986a96f84af.d8a694f6a5f1sa86f1a98g.4da89s4fda98f498ga",
	}

	tt := []struct {
		name             string
		reqBodyJSON      string
		listingService   *lm.MockService
		authService      *aum.MockService
		expectedResponse string
		expectedStatus   int
	}{
		{
			name:        "When login credentials are okay and token is successfully generated",
			reqBodyJSON: fmt.Sprintf(`{"cpf":"%s","secret":"%s"}`, account.CPF, account.Secret),
			listingService: &lm.MockService{
				Account: account,
			},
			authService: &aum.MockService{
				Token: token,
			},
			expectedResponse: fmt.Sprintf(`{"token":"%s"}`, token.Digest),
			expectedStatus:   http.StatusOK,
		},
		{
			name:        "When cpf credential is not in our repository",
			reqBodyJSON: fmt.Sprintf(`{"cpf":"%s","secret":"%s"}`, account.CPF, account.Secret),
			listingService: &lm.MockService{
				Err: mongodb.ErrNoAccountWasFound,
			},
			expectedResponse: `{"status_code":403,"message":"it seems your login credentials are invalid, verify them and try again"}`,
			expectedStatus:   http.StatusForbidden,
		},
		{
			name:        "When cpf credential is in our repository, but password is invalid",
			reqBodyJSON: fmt.Sprintf(`{"cpf":"%s","secret":"%s"}`, account.CPF, "deuruim"),
			listingService: &lm.MockService{
				Account: account,
			},
			authService: &aum.MockService{
				Err: authenticating.InvalidLoginErr,
			},
			expectedResponse: `{"status_code":403,"message":"it seems your login credentials are invalid, verify them and try again"}`,
			expectedStatus:   http.StatusForbidden,
		},
		{
			name:             "When payload is invalid",
			reqBodyJSON:      fmt.Sprintf(`{"cpf":"%s","secret":123498}`, account.CPF),
			listingService:   &lm.MockService{},
			authService:      &aum.MockService{},
			expectedResponse: `{"status_code":400,"message":"Invalid Login entity: expected type string, got number at field secret"}`,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:        "When unexpected errors occurs at repo operations",
			reqBodyJSON: fmt.Sprintf(`{"cpf":"%s","secret":"%s"}`, account.CPF, account.Secret),
			listingService: &lm.MockService{
				Err: errors.New("foo unexpected"),
			},
			authService:      &aum.MockService{},
			expectedResponse: `{"status_code":500,"message":"foo unexpected"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name:           "When unexpected errors occurs at gatekeeper operations",
			reqBodyJSON:    fmt.Sprintf(`{"cpf":"%s","secret":"%s"}`, account.CPF, account.Secret),
			listingService: &lm.MockService{},
			authService: &aum.MockService{
				Err: errors.New("foo unexpected"),
			},
			expectedResponse: `{"status_code":500,"message":"foo unexpected"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHandler(logger, tc.authService, tc.listingService)

			var reqBody string
			jsonBuffer := bytes.NewBuffer([]byte(tc.reqBodyJSON))
			if err := json.NewEncoder(jsonBuffer).Encode(reqBody); err != nil {
				t.Fatalf("Could not encode %s as json", tc.reqBodyJSON)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/login", jsonBuffer)

			handler.Login(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}

			helpers.AssertResponseJSON(t, w, tc.expectedResponse)
		})
	}
}

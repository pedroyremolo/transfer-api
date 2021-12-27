package transferring

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pedroyremolo/transfer-api/pkg"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/pedroyremolo/transfer-api/pkg/tests/helpers"
	am "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/adding"
	aum "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/authenticating"
	lm "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/listing"
	tm "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/transferring"
	um "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/updating"
	"github.com/sirupsen/logrus"
)

func TestMakeTransfer(t *testing.T) {
	log := lgr.NewDefaultLogger()
	logger := logrus.NewEntry(log)

	tt := []struct {
		name                string
		reqBodyJSON         string
		reqHeader           http.Header
		authService         *aum.MockService
		listingService      *lm.MockService
		transferringService *tm.MockService
		updatingService     *um.MockService
		addingService       *am.MockService
		expectedResponse    string
		expectedStatus      int
	}{
		{
			name:        "When transfer successfully occurs",
			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
				"Content-Type":  []string{"application/json"},
			},
			authService: &aum.MockService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
				Err: nil,
			},
			listingService: &lm.MockService{
				Balance: 22.22,
			},
			transferringService: &tm.MockService{},
			updatingService:     &um.MockService{},
			addingService: &am.MockService{
				ID: "f1869a4f9a84f89sa",
			},
			expectedStatus: http.StatusOK,
		},
		/* 		{
		   			name:        "When no auth header is sent",
		   			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
		   			reqHeader: http.Header{
		   				"Content-Type": []string{"application/json"},
		   			},
		   			expectedStatus:   http.StatusUnauthorized,
		   			expectedResponse: `{"status_code":401,"message":"it seems you don't have or didn't pass valid credentials to this route"}`,
		   		},
		   		{
		   			name:        "When malformed auth header is sent",
		   			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
		   			reqHeader: http.Header{
		   				"Authorization": []string{"Bearerea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
		   				"Content-Type":  []string{"application/json"},
		   			},
		   			expectedStatus:   http.StatusUnauthorized,
		   			expectedResponse: `{"status_code":401,"message":"it seems you don't have or didn't pass valid credentials to this route"}`,
		   		},
		   		{
		   			name:        "When auth header is not of Bearer type",
		   			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
		   			reqHeader: http.Header{
		   				"Authorization": []string{"Basic ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
		   				"Content-Type":  []string{"application/json"},
		   			},
		   			expectedStatus:   http.StatusUnauthorized,
		   			expectedResponse: `{"status_code":401,"message":"it seems you don't have or didn't pass valid credentials to this route"}`,
		   		},
		   		{
		   			name:        "When auth header contains an invalid token",
		   			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
		   			reqHeader: http.Header{
		   				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
		   				"Content-Type":  []string{"application/json"},
		   			},
		   			authService: &aum.MockService{
		   				Err: errors.New("foo"),
		   			},
		   			expectedStatus:   http.StatusUnauthorized,
		   			expectedResponse: `{"status_code":401,"message":"it seems you don't have or didn't pass valid credentials to this route"}`,
		   		}, */ // to be excluded after middleware refactoring
		{
			name:        "When req body cannot be deserialized as transfer",
			reqBodyJSON: `{"account_destination_id":123,"amount":11.11}`,
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
				"Content-Type":  []string{"application/json"},
			},
			authService: &aum.MockService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
				Err: nil,
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"status_code":400,"message":"Invalid Transfer entity: expected type string, got number at field account_destination_id"}`,
		},
		{
			name:        "When fails to retrieve origin account balance",
			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
				"Content-Type":  []string{"application/json"},
			},
			authService: &aum.MockService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
				Err: nil,
			},
			listingService: &lm.MockService{
				CallsToFail: 1,
				Err:         errors.New("foo"),
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"status_code":500,"message":"foo"}`,
		},
		{
			name:        "When fails to retrieve destination account balance",
			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
				"Content-Type":  []string{"application/json"},
			},
			authService: &aum.MockService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
				Err: nil,
			},
			listingService: &lm.MockService{
				CallsToFail: 2,
				Err:         errors.New("foo"),
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"status_code":500,"message":"foo"}`,
		},
		{
			name:        "When there's no destination account with the informed id",
			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
				"Content-Type":  []string{"application/json"},
			},
			authService: &aum.MockService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
				Err: nil,
			},
			listingService: &lm.MockService{
				CallsToFail: 2,
				Err:         errors.New("no account was found with the given filter parameters"),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"status_code":400,"message":"no account was found with the given filter parameters"}`,
		},
		{
			name:        "When origin account has not enough balance to transfer",
			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
				"Content-Type":  []string{"application/json"},
			},
			authService: &aum.MockService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
			},
			listingService: &lm.MockService{
				Balance: 22.22,
			},
			transferringService: &tm.MockService{
				Err: errors.New("not enough origin balance"),
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"status_code":400,"message":"not enough origin balance"}`,
		},
		{
			name:        "When fails to update accounts balances",
			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
				"Content-Type":  []string{"application/json"},
			},
			authService: &aum.MockService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
			},
			listingService: &lm.MockService{
				Balance: 22.22,
			},
			transferringService: &tm.MockService{},
			updatingService: &um.MockService{
				Err: errors.New("foo"),
			},
			expectedResponse: `{"status_code":500,"message":"foo"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name:        "When fails to add transfers to db",
			reqBodyJSON: `{"account_destination_id":"5f8f8ccb30a1cd7511c5cb70","amount":11.11}`,
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
				"Content-Type":  []string{"application/json"},
			},
			authService: &aum.MockService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
			},
			listingService: &lm.MockService{
				Balance: 22.22,
			},
			transferringService: &tm.MockService{},
			updatingService:     &um.MockService{},
			addingService: &am.MockService{
				Err: errors.New("foo"),
			},
			expectedResponse: `{"status_code":500,"message":"foo"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHandler(logger, tc.transferringService, tc.listingService, tc.addingService, tc.updatingService)

			var reqBody string
			jsonBuffer := bytes.NewBuffer([]byte(tc.reqBodyJSON))
			if err := json.NewEncoder(jsonBuffer).Encode(reqBody); err != nil {
				t.Fatalf("Could not encode %s as json", tc.reqBodyJSON)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/transfers", jsonBuffer)
			r = r.WithContext(context.WithValue(r.Context(), pkg.AccountID, "4a6sgf4as6g"))
			r.Header = tc.reqHeader

			handler.MakeTransfer(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}

			helpers.AssertResponseJSON(t, w, tc.expectedResponse)
		})
	}
}

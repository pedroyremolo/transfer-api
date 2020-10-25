package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	am "github.com/pedroyremolo/transfer-api/pkg/mocks/adding"
	lm "github.com/pedroyremolo/transfer-api/pkg/mocks/listing"
	um "github.com/pedroyremolo/transfer-api/pkg/mocks/updating"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	a := &am.MockService{}
	l := &lm.MockService{}
	auth := &mockAuthenticatingService{}
	tf := &mockTransferringService{}
	u := &um.MockService{}

	handler := Handler(a, l, auth, tf, u)

	if handler == nil {
		t.Errorf("Expected an implementation of http.Handler, got %s", handler)
	}
}

func TestAddAccount(t *testing.T) {
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
			var reqBody string
			jsonBuffer := bytes.NewBuffer([]byte(tc.reqBodyJSON))
			if err := json.NewEncoder(jsonBuffer).Encode(reqBody); err != nil {
				t.Fatalf("Could not encode %s as json", tc.reqBodyJSON)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/accounts", jsonBuffer)
			s := tc.service
			h := addAccount(s)

			h(w, r, nil)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}
			if len(tc.expectedErrResponse) > 0 {
				assertResponseJSON(t, w, tc.expectedErrResponse)
				return
			}
			if w.Header().Get("Location") != fmt.Sprintf("/%s", tc.service.ID) {
				t.Errorf("Expected Location header /%s; got %v", tc.service.ID, w.Header().Get("Location"))
			}
		})
	}
}

func TestGetAccountBalanceByID(t *testing.T) {
	tt := []struct {
		name             string
		id               string
		service          *lm.MockService
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "When successfully returns",
			id:               "a6sf46af6af",
			service:          &lm.MockService{Balance: 42.42},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"balance":42.42}`,
		},
		{
			name:             "When no account was found with the given id",
			id:               "a6sf46af6af",
			service:          &lm.MockService{Err: mongodb.ErrNoAccountWasFound},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: `{"status_code":404,"message":"no account was found with the given filter parameters"}`,
		},
		{
			name:             "When unexpected errors inside the service occurs",
			id:               "a6sf46af6af",
			service:          &lm.MockService{Err: errors.New("foo")},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"status_code":500,"message":"foo"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			target := fmt.Sprintf("/accounts/%s/balance", tc.id)
			r := httptest.NewRequest(http.MethodGet, target, nil)
			s := tc.service
			h := getAccountBalanceByID(s)

			h(w, r, httprouter.Params{{
				Key:   "id",
				Value: tc.id,
			}})

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}
			assertResponseJSON(t, w, tc.expectedResponse)
		})
	}
}

func TestGetAccounts(t *testing.T) {
	currentTime := time.Now().UTC()
	tt := []struct {
		name             string
		service          *lm.MockService
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "When successfully returns",
			service: &lm.MockService{
				Accounts: []listing.Account{
					{
						ID:        "g4a68vf6a4g96ws84g",
						Name:      "Monkey D. Luffy",
						CPF:       "11111111030",
						Secret:    "t89awsg4189a1f9a8s1d",
						Balance:   100000.42,
						CreatedAt: &currentTime,
					},
					{
						ID:        "8h964dsa6gb51qa98w1698",
						Name:      "Harry Potter",
						CPF:       "95360976055",
						Secret:    "4wq89fa6s19q8etg498a",
						Balance:   40000.42,
						CreatedAt: &currentTime,
					},
				},
			},
			expectedStatus: http.StatusOK,
			expectedResponse: fmt.Sprintf(
				`[{"id":"g4a68vf6a4g96ws84g","name":"Monkey D. Luffy","cpf":"11111111030","balance":100000.42,"created_at":"%s"},{"id":"8h964dsa6gb51qa98w1698","name":"Harry Potter","cpf":"95360976055","balance":40000.42,"created_at":"%s"}]`,
				currentTime.Format(time.RFC3339Nano), currentTime.Format(time.RFC3339Nano),
			),
		},
		{
			name:             "When no account was found",
			service:          &lm.MockService{Accounts: []listing.Account{}},
			expectedStatus:   http.StatusOK,
			expectedResponse: `[]`,
		},
		{
			name:             "When unexpected errors inside the service occurs",
			service:          &lm.MockService{Err: errors.New("foo")},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"status_code":500,"message":"foo"}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			target := "/accounts"
			r := httptest.NewRequest(http.MethodGet, target, nil)
			s := tc.service
			h := getAccounts(s)

			h(w, r, nil)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}
			assertResponseJSON(t, w, tc.expectedResponse)
		})
	}
}

func TestLogin(t *testing.T) {
	tokenOID := primitive.NewObjectID()
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
		authService      *mockAuthenticatingService
		expectedResponse string
		expectedStatus   int
	}{
		{
			name:        "When login credentials are okay and token is successfully generated",
			reqBodyJSON: fmt.Sprintf(`{"cpf":"%s","secret":"%s"}`, account.CPF, account.Secret),
			listingService: &lm.MockService{
				Account: account,
			},
			authService: &mockAuthenticatingService{
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
			expectedResponse: `{"status_code":401,"message":"it seems your login credentials are invalid, verify them and try again"}`,
			expectedStatus:   http.StatusUnauthorized,
		},
		{
			name:        "When cpf credential is in our repository, but password is invalid",
			reqBodyJSON: fmt.Sprintf(`{"cpf":"%s","secret":"%s"}`, account.CPF, "deuruim"),
			listingService: &lm.MockService{
				Account: account,
			},
			authService: &mockAuthenticatingService{
				Err: authenticating.InvalidLoginErr,
			},
			expectedResponse: `{"status_code":401,"message":"it seems your login credentials are invalid, verify them and try again"}`,
			expectedStatus:   http.StatusUnauthorized,
		},
		{
			name:             "When payload is invalid",
			reqBodyJSON:      fmt.Sprintf(`{"cpf":"%s","secret":123498}`, account.CPF),
			listingService:   &lm.MockService{},
			authService:      &mockAuthenticatingService{},
			expectedResponse: `{"status_code":400,"message":"Invalid Login entity: expected type string, got number at field secret"}`,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:        "When unexpected errors occurs at repo operations",
			reqBodyJSON: fmt.Sprintf(`{"cpf":"%s","secret":"%s"}`, account.CPF, account.Secret),
			listingService: &lm.MockService{
				Err: errors.New("foo unexpected"),
			},
			authService:      &mockAuthenticatingService{},
			expectedResponse: `{"status_code":500,"message":"foo unexpected"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
		{
			name:           "When unexpected errors occurs at gatekeeper operations",
			reqBodyJSON:    fmt.Sprintf(`{"cpf":"%s","secret":"%s"}`, account.CPF, account.Secret),
			listingService: &lm.MockService{},
			authService: &mockAuthenticatingService{
				Err: errors.New("foo unexpected"),
			},
			expectedResponse: `{"status_code":500,"message":"foo unexpected"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var reqBody string
			jsonBuffer := bytes.NewBuffer([]byte(tc.reqBodyJSON))
			if err := json.NewEncoder(jsonBuffer).Encode(reqBody); err != nil {
				t.Fatalf("Could not encode %s as json", tc.reqBodyJSON)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/accounts", jsonBuffer)
			h := login(tc.authService, tc.listingService)

			h(w, r, nil)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}

			assertResponseJSON(t, w, tc.expectedResponse)
		})
	}
}

func TestTransfer(t *testing.T) {
	tt := []struct {
		name                string
		reqBodyJSON         string
		reqHeader           http.Header
		authService         *mockAuthenticatingService
		listingService      *lm.MockService
		transferringService *mockTransferringService
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
			authService: &mockAuthenticatingService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
				Err: nil,
			},
			listingService: &lm.MockService{
				Balance: 22.22,
			},
			transferringService: &mockTransferringService{},
			updatingService:     &um.MockService{},
			addingService: &am.MockService{
				ID: "f1869a4f9a84f89sa",
			},
			expectedStatus: http.StatusOK,
		},
		{
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
			authService: &mockAuthenticatingService{
				Err: errors.New("foo"),
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: `{"status_code":401,"message":"it seems you don't have or didn't pass valid credentials to this route"}`,
		},
		{
			name:        "When req body cannot be deserialized as transfer",
			reqBodyJSON: `{"account_destination_id":123,"amount":11.11}`,
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
				"Content-Type":  []string{"application/json"},
			},
			authService: &mockAuthenticatingService{
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
			authService: &mockAuthenticatingService{
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
			authService: &mockAuthenticatingService{
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
			authService: &mockAuthenticatingService{
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
			authService: &mockAuthenticatingService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
			},
			listingService: &lm.MockService{
				Balance: 22.22,
			},
			transferringService: &mockTransferringService{
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
			authService: &mockAuthenticatingService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
			},
			listingService: &lm.MockService{
				Balance: 22.22,
			},
			transferringService: &mockTransferringService{},
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
			authService: &mockAuthenticatingService{
				Token: authenticating.Token{
					ClientID: "4f98as4f98sa496a1f",
				},
			},
			listingService: &lm.MockService{
				Balance: 22.22,
			},
			transferringService: &mockTransferringService{},
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
			var reqBody string
			jsonBuffer := bytes.NewBuffer([]byte(tc.reqBodyJSON))
			if err := json.NewEncoder(jsonBuffer).Encode(reqBody); err != nil {
				t.Fatalf("Could not encode %s as json", tc.reqBodyJSON)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/transfers", jsonBuffer)
			r.Header = tc.reqHeader
			h := transfer(tc.addingService, tc.authService, tc.listingService, tc.transferringService, tc.updatingService)

			h(w, r, nil)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}

			assertResponseJSON(t, w, tc.expectedResponse)
		})
	}
}

type mockTransferringService struct {
	Err error
}

func (m *mockTransferringService) BalanceBetweenAccounts(originBalance float64, destinationBalance float64, _ float64) (_ float64, _ float64, _ error) {
	return originBalance, destinationBalance, m.Err
}

type mockAuthenticatingService struct {
	Token authenticating.Token
	Err   error
}

func (m *mockAuthenticatingService) Sign(_ context.Context, _ authenticating.Login, _ string, _ string) (authenticating.Token, error) {
	return m.Token, m.Err
}

func (m *mockAuthenticatingService) Verify(_ context.Context, _ string) (authenticating.Token, error) {
	return m.Token, m.Err
}

func assertResponseJSON(t *testing.T, w *httptest.ResponseRecorder, expectedResponseJSON string) {
	t.Helper()
	respBodyBytes, err := ioutil.ReadAll(w.Body)
	respBody := string(bytes.TrimSpace(respBodyBytes))
	if err != nil {
		t.Fatal("Unable to read response from Recorder")
	}
	if respBody != expectedResponseJSON {
		t.Errorf("Expected response %s; got %s", expectedResponseJSON, respBody)
	}
}

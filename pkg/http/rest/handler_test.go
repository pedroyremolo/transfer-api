package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler(t *testing.T) {
	a := &mockAddingService{}
	l := &mockListingService{}
	handler := Handler(a, l)

	if handler == nil {
		t.Errorf("Expected an implementation of http.Handler, got %s", handler)
	}
}

func TestAddAccount(t *testing.T) {
	tt := []struct {
		name                string
		service             *mockAddingService
		reqBodyJSON         string
		expectedStatus      int
		expectedErrResponse string
	}{
		{
			name:           "When successfully returns",
			service:        &mockAddingService{ID: "a6sf46af6af"},
			reqBodyJSON:    `{"name": "Jane Doe","cpf": "11111111030","secret": "254855","balance": 50.00}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:                "When the sent json is not valid",
			service:             &mockAddingService{ID: ""},
			reqBodyJSON:         `{"name": "Jane Doe","cpf": "111111110301","secret": "254855","balance": 50.00}`,
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: `{"status_code":400,"message":"Field cpf contains an invalid value: 111111110301 is not a valid cpf"}`,
		},
		{
			name:                "When the sent cpf is already into DB",
			service:             &mockAddingService{Err: mongodb.ErrCPFAlreadyExists},
			reqBodyJSON:         `{"name": "Jane Doe","cpf": "11111111030","secret": "254855","balance": 50.00}`,
			expectedStatus:      http.StatusBadRequest,
			expectedErrResponse: `{"status_code":400,"message":"this cpf could not be inserted in our DB"}`,
		},
		{
			name:                "When unexpected errors inside the service occurs",
			service:             &mockAddingService{Err: errors.New("foo")},
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
		service          *mockListingService
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "When successfully returns",
			id:               "a6sf46af6af",
			service:          &mockListingService{Balance: 42.42},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"balance":42.42}`,
		},
		{
			name:             "When no account was found with the given id",
			id:               "a6sf46af6af",
			service:          &mockListingService{Err: mongodb.ErrNoAccountWasFound},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: `{"status_code":404,"message":"no account was found with the given filter parameters"}`,
		},
		{
			name:             "When unexpected errors inside the service occurs",
			id:               "a6sf46af6af",
			service:          &mockListingService{Err: errors.New("foo")},
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
		service          *mockListingService
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "When successfully returns",
			service: &mockListingService{
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
			service:          &mockListingService{Accounts: []listing.Account{}},
			expectedStatus:   http.StatusOK,
			expectedResponse: `[]`,
		},
		{
			name:             "When unexpected errors inside the service occurs",
			service:          &mockListingService{Err: errors.New("foo")},
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

type mockAddingService struct {
	ID  string
	Err error
}

func (m *mockAddingService) AddAccount(_ context.Context, _ adding.Account) (string, error) {
	return m.ID, m.Err
}

type mockListingService struct {
	Balance  float64
	Accounts []listing.Account
	Err      error
}

func (m *mockListingService) GetAccountBalanceByID(_ context.Context, _ string) (float64, error) {
	return m.Balance, m.Err
}

func (m *mockListingService) GetAccounts(_ context.Context) ([]listing.Account, error) {
	return m.Accounts, m.Err
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

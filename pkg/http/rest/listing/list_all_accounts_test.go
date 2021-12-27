package listing

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/tests/helpers"
	lm "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/listing"
	"github.com/sirupsen/logrus"
)

func TestGetAccounts(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())
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
			handler := NewHandler(logger, tc.service)
			w := httptest.NewRecorder()
			target := "/accounts"
			r := httptest.NewRequest(http.MethodGet, target, nil)

			handler.ListAllAccounts(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}
			helpers.AssertResponseJSON(t, w, tc.expectedResponse)
		})
	}
}

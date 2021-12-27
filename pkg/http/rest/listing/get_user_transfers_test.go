package listing

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pedroyremolo/transfer-api/pkg"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/tests/helpers"
	lm "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/listing"
	"github.com/sirupsen/logrus"
)

func TestGetUserTransfers(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())
	defaultClientID := "jff46as84dcsa365418"
	defaultSentTransfer := listing.Transfer{
		ID:                   "4as6g84as68gf4as",
		OriginAccountID:      defaultClientID,
		DestinationAccountID: "4896as4rfa689tqwrtg",
		Amount:               23.32,
		CreatedAt:            time.Time{},
	}
	defaultReceivedTransfer := listing.Transfer{
		ID:                   "t4a8g496ag49ga",
		OriginAccountID:      "4896as4rfa689tqwrtg",
		DestinationAccountID: defaultClientID,
		Amount:               23.32,
		CreatedAt:            time.Time{},
	}
	tt := []struct {
		name             string
		reqHeader        http.Header
		listingService   *lm.MockService
		expectedResponse string
		expectedStatus   int
	}{
		{
			name: "When successfully retrieves account transfers",
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
			},
			listingService: &lm.MockService{
				AccountTransfers: listing.AccountTransfers{
					Sent:     []listing.Transfer{defaultSentTransfer},
					Received: []listing.Transfer{defaultReceivedTransfer},
				},
				Err: nil,
			},
			expectedResponse: `{"sent":[{"id":"4as6g84as68gf4as","account_origin_id":"jff46as84dcsa365418","account_destination_id":"4896as4rfa689tqwrtg","amount":23.32,"created_at":"0001-01-01T00:00:00Z"}],"received":[{"id":"t4a8g496ag49ga","account_origin_id":"4896as4rfa689tqwrtg","account_destination_id":"jff46as84dcsa365418","amount":23.32,"created_at":"0001-01-01T00:00:00Z"}]}`,
			expectedStatus:   http.StatusOK,
		},
		{
			name: "When fails to retrieve account transfers",
			reqHeader: http.Header{
				"Authorization": []string{"Bearer ea4984da84fa8e.ae498f4a9e8f.af84a9f64a9"},
			},
			listingService: &lm.MockService{
				Err: errors.New("db error"),
			},
			expectedResponse: `{"status_code":500,"message":"db error"}`,
			expectedStatus:   http.StatusInternalServerError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHandler(logger, tc.listingService)
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/transfers", nil)
			r = r.WithContext(context.WithValue(r.Context(), pkg.AccountID, defaultClientID))
			r.Header = tc.reqHeader

			handler.GetUserTransfers(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected response status %v; got %v", tc.expectedStatus, w.Code)
			}

			helpers.AssertResponseJSON(t, w, tc.expectedResponse)
		})
	}
}

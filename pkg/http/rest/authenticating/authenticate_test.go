package authenticating

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pedroyremolo/transfer-api/pkg"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
	mocks "github.com/pedroyremolo/transfer-api/pkg/tests/mocks/authenticating"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestHandler_Authenticate(t *testing.T) {
	log := logrus.New()
	logger := logrus.NewEntry(log)

	// success case
	t.Run("should successfully authenticate and return account id sent by ctx", func(t *testing.T) {
		oid := primitive.NewObjectID()
		service := &mocks.MockService{
			Token: authenticating.Token{
				ID:       &oid,
				ClientID: "101",
				Digest:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			},
		}
		expectedStatus := http.StatusOK

		h := Handler{
			logger:         logger,
			service:        service,
			listingService: nil,
		}

		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/foo", nil)
		request.Header.Add("Authorization", fmt.Sprint(bearerAuthType, " ", service.Token.Digest))

		h.Authenticate(fakeHandlerFunc(expectedStatus, nil)).ServeHTTP(recorder, request)

		if strings.TrimSpace(recorder.Body.String()) != service.Token.ClientID {
			t.Errorf("expected account id %s, but got %s", service.Token.ClientID, strings.TrimSpace(recorder.Body.String()))
		}

		if recorder.Code != expectedStatus {
			t.Errorf("expected status code %d, but got %d", expectedStatus, recorder.Code)
		}
	})

	// error cases
	type fields struct {
		logger         *logrus.Entry
		tokenDigest    string
		wantServiceErr error
		authHeader     string
	}
	tests := []struct {
		name       string
		fields     fields
		wantStatus int
		wantErr    error
	}{
		{
			name: "should fail to authenticate when header is not of type bearer",
			fields: fields{
				logger:         logger,
				tokenDigest:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				authHeader:     "Basic eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				wantServiceErr: nil,
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    authenticating.ProtectedRouteErr,
		},
		{
			name: "should fail to authenticate when header is malformed",
			fields: fields{
				logger:         logger,
				tokenDigest:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				authHeader:     "BasiceyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				wantServiceErr: nil,
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    authenticating.ProtectedRouteErr,
		},
		{
			name: "should fail to authenticate when header is empty",
			fields: fields{
				logger:         logger,
				tokenDigest:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				authHeader:     "",
				wantServiceErr: authenticating.ProtectedRouteErr,
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    authenticating.ProtectedRouteErr,
		},
		{
			name: "should fail to authenticate when verify fails",
			fields: fields{
				logger:         logger,
				tokenDigest:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				authHeader:     "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				wantServiceErr: authenticating.ProtectedRouteErr,
			},
			wantStatus: http.StatusUnauthorized,
			wantErr:    authenticating.ProtectedRouteErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &mocks.MockService{
				Err: tt.fields.wantServiceErr,
			}
			expectedResponse, _ := json.Marshal(rest.ErrorResponse{StatusCode: tt.wantStatus, Message: tt.wantErr.Error()})

			h := Handler{
				logger:         tt.fields.logger,
				service:        service,
				listingService: nil,
			}

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/foo", nil)
			request.Header.Add("Authorization", tt.fields.authHeader)

			h.Authenticate(fakeHandlerFunc(tt.wantStatus, nil)).ServeHTTP(recorder, request)

			if strings.TrimSpace(recorder.Body.String()) != string(expectedResponse) {
				t.Errorf("expected %s, but got %s", expectedResponse, strings.TrimSpace(recorder.Body.String()))
			}

			if recorder.Code != tt.wantStatus {
				t.Errorf("expected status code %d, but got %d", tt.wantStatus, recorder.Code)
			}
		})
	}
}

func fakeHandlerFunc(statusCode int, err error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			w.WriteHeader(statusCode)
			fmt.Fprintln(w, err.Error())
		}
		accountID := r.Context().Value(pkg.AccountID).(string)
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, accountID)
	}
}

package jwt

import (
	"reflect"
	"testing"

	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
)

func TestGatekeeper_Sign(t *testing.T) {
	gk := NewGatekeeper("testSecret", "test")
	type args struct {
		clientID string
	}
	tt := []struct {
		name          string
		args          args
		wantErr       bool
		expectedToken authenticating.Token
	}{
		{
			name: "When successfully return the token",
			args: args{
				clientID: "fa684sf896asf49a8",
			},
			expectedToken: authenticating.Token{
				ClientID: "fa684sf896asf49a8",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gk.Sign(tc.args.clientID)
			if (err != nil) != tc.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if tc.expectedToken.ClientID != got.ClientID {
				t.Errorf("Expected id %s, got %s", tc.args.clientID, got.ClientID)
			}
		})
	}
}

func TestGatekeeper_Verify(t *testing.T) {
	gk := NewGatekeeper("testSecret", "test")
	clientID := "4sfa9684fsa698"
	token, _ := gk.Sign(clientID)
	tt := []struct {
		name        string
		tokenDigest string
		want        authenticating.Token
		wantErr     bool
	}{
		{
			name:        "When there's a valid token",
			tokenDigest: token.Digest,
			want:        token,
			wantErr:     false,
		},
		{
			name:        "When there's an invalid token",
			tokenDigest: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0RXJyIiwiZXhwIjoxNjAzNDY2OTEyLCJpYXQiOjE2MDM0NjUxMTIsImp0aSI6IjVmOTJlZjk4MGRlOWZmMGY0N2MzNjc2YiIsImNsaWVudF9pZCI6IjRzZmE5Njg0ZnNhNjk4In0.yIJbqOSFDjZ7gTlLRHK-Wm_WO_Ghh_d1SwYdMjiZKps",
			wantErr:     true,
		},
		{
			name:        "When algorithm is invalid",
			tokenDigest: "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ0ZXN0IiwiZXhwIjoxNjAzNDY2OTEyLCJpYXQiOjE2MDM0NjUxMTIsImp0aSI6IjVmOTJlZjk4MGRlOWZmMGY0N2MzNjc2YiIsImNsaWVudF9pZCI6IjRzZmE5Njg0ZnNhNjk4In0.Z5llFB6oUgt-KMshrFL7R7EN3FzkBWyalcDs4XZuGQ5r1HGXFnVXdHaQOOluuVtF",
			wantErr:     true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gk.Verify(tc.tokenDigest)
			if !tc.wantErr && (err != nil) {
				t.Errorf("Verify() error = %v, wantErr %v", err, tc.wantErr)
			}

			if !tc.wantErr && !reflect.DeepEqual(token, got) {
				t.Errorf("Expected token %v, got %v", token, got)
			}
		})
	}
}

package jwt

import (
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"reflect"
	"testing"
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
	gk := &Gatekeeper{hs: jwt.NewHS256([]byte("test"))}
	tt := []struct {
		name        string
		tokenDigest string
		want        authenticating.Token
	}{
		// TODO: Add test cases.
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := gk.Verify(tc.tokenDigest); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Verify() = %v, want %v", got, tc.want)
			}
		})
	}
}

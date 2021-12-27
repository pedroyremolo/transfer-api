package authenticating

import (
	"context"
	"net/http"
	"strings"

	"github.com/pedroyremolo/transfer-api/pkg"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
)

const bearerAuthType = "Bearer"

func (h Handler) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHd := r.Header.Get("Authorization")
		if authHd == "" {
			h.logger.Error("Empty Authorization header")
			rest.SetJSONError(h.logger, authenticating.ProtectedRouteErr, http.StatusUnauthorized, w)
			return
		}

		hdSpaceIndex := strings.Index(authHd, " ")
		if hdSpaceIndex == -1 {
			h.logger.Error("Authorization header is malformed")
			rest.SetJSONError(h.logger, authenticating.ProtectedRouteErr, http.StatusUnauthorized, w)
			return
		}

		authType, tokenDigest := authHd[:hdSpaceIndex], authHd[strings.LastIndex(authHd, " ")+1:]
		if authType != bearerAuthType {
			h.logger.Error("Authorization header is not of type Bearer")
			rest.SetJSONError(h.logger, authenticating.ProtectedRouteErr, http.StatusUnauthorized, w)
			return
		}

		token, err := h.service.Verify(ctx, tokenDigest)
		if err != nil {
			h.logger.Errorf("Err %s when trying to verify token", err)
			rest.SetJSONError(h.logger, authenticating.ProtectedRouteErr, http.StatusUnauthorized, w)
			return
		}

		r = r.WithContext(context.WithValue(ctx, pkg.AccountID, token.ClientID))

		next(w, r)
	}
}

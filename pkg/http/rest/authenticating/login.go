package authenticating

import (
	"encoding/json"
	"net/http"

	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
)

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	ctx := r.Context()

	var login authenticating.Login
	if err := decoder.Decode(&login); err != nil {
		rest.SetJSONError(h.logger, err, http.StatusBadRequest, w)
		return
	}

	account, err := h.listingService.GetAccountByCPF(ctx, login.CPF)
	if err != nil {
		if err.Error() == mongodb.ErrNoAccountWasFound.Error() {
			rest.SetJSONError(h.logger, authenticating.InvalidLoginErr, http.StatusForbidden, w)
			return
		}
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
		return
	}

	var token authenticating.Token
	token, err = h.service.Sign(ctx, login, account.Secret, account.ID)
	if err != nil {
		if err.Error() == authenticating.InvalidLoginErr.Error() {
			rest.SetJSONError(h.logger, authenticating.InvalidLoginErr, http.StatusForbidden, w)
			return
		}
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", rest.DefaultContentType)
	_ = json.NewEncoder(w).Encode(authenticating.Token{Digest: token.Digest})
}

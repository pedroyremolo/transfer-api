package adding

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
)

func (h Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	ctx := r.Context()
	var account adding.Account
	if err := decoder.Decode(&account); err != nil {
		rest.SetJSONError(h.logger, err, http.StatusBadRequest, w)
		return
	}

	id, err := h.service.AddAccount(ctx, account)
	if err != nil {
		if err.Error() == mongodb.ErrCPFAlreadyExists.Error() {
			rest.SetJSONError(h.logger, err, http.StatusBadRequest, w)
			return
		}
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", rest.DefaultContentType)
	w.Header().Set("Location", fmt.Sprintf("/%s", id))
	w.WriteHeader(http.StatusCreated)
}

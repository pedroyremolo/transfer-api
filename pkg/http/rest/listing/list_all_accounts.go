package listing

import (
	"encoding/json"
	"net/http"

	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
)

func (h Handler) ListAllAccounts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accounts, err := h.service.GetAccounts(ctx)
	if err != nil {
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
		return
	}
	w.Header().Set("Content-Type", rest.DefaultContentType)
	_ = json.NewEncoder(w).Encode(accounts)
}

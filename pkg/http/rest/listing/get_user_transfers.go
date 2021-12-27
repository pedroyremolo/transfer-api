package listing

import (
	"encoding/json"
	"net/http"

	"github.com/pedroyremolo/transfer-api/pkg"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
)

func (h Handler) GetUserTransfers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountID := ctx.Value(pkg.AccountID).(string)

	accountTransfers, err := h.service.GetTransfersByAccountID(ctx, accountID)
	if err != nil {
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", rest.DefaultContentType)
	_ = json.NewEncoder(w).Encode(accountTransfers)
}

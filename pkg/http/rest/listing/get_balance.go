package listing

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
)

func (h Handler) GetBalanceByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	ctx := r.Context()

	balance, err := h.service.GetAccountBalanceByID(ctx, id)
	if err != nil {
		if err.Error() == mongodb.ErrNoAccountWasFound.Error() {
			rest.SetJSONError(h.logger, err, http.StatusNotFound, w)
			return
		}
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", rest.DefaultContentType)
	_ = json.NewEncoder(w).Encode(listing.Account{Balance: balance})
}

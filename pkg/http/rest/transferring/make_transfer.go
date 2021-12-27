package transferring

import (
	"encoding/json"
	"net/http"

	"github.com/pedroyremolo/transfer-api/pkg"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/http/rest"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"github.com/pedroyremolo/transfer-api/pkg/updating"
)

func (h Handler) MakeTransfer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	originAccountID := ctx.Value(pkg.AccountID).(string)
	/* token, err := verifyTokenFromAuthHeader(r, auth, ctx)
	if err != nil {
		SetJSONError(err, http.StatusUnauthorized, w)
		return
	} */ // will be replaced by middleware
	decoder := json.NewDecoder(r.Body)
	var transfer adding.Transfer
	if err := decoder.Decode(&transfer); err != nil {
		rest.SetJSONError(h.logger, err, http.StatusBadRequest, w)
		return
	}

	originBalance, err := h.listingService.GetAccountBalanceByID(ctx, originAccountID) //take account id from ctx
	if err != nil {
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
		return
	}
	destinationBalance, err := h.listingService.GetAccountBalanceByID(ctx, transfer.DestinationAccountID)
	if err != nil {
		if err.Error() == mongodb.ErrNoAccountWasFound.Error() {
			rest.SetJSONError(h.logger, err, http.StatusBadRequest, w)
			return
		}
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
		return
	}

	var origin, destination updating.Account
	origin.ID = originAccountID
	destination.ID = transfer.DestinationAccountID
	transfer.OriginAccountID = originAccountID

	origin.Balance, destination.Balance, err = h.service.BalanceBetweenAccounts(originBalance, destinationBalance, transfer.Amount)
	if err != nil {
		rest.SetJSONError(h.logger, err, http.StatusBadRequest, w)
		return
	}

	if err = h.updatingService.UpdateAccounts(ctx, origin, destination); err != nil {
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
		return
	}

	if _, err = h.addingService.AddTransfer(ctx, transfer); err != nil {
		rest.SetJSONError(h.logger, err, http.StatusInternalServerError, w)
	}
}

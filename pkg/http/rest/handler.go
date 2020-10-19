package rest

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"log"
	"net/http"
)

const (
	defaultContentType = "application/json"
)

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func setJSONError(err error, status int, w http.ResponseWriter) {
	log.Println(err.Error())
	response := ErrorResponse{
		StatusCode: status,
		Message:    err.Error(),
	}
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func Handler(a adding.Service) http.Handler {
	router := httprouter.New()

	router.POST("/accounts", addAccount(a))

	return router
}

func addAccount(a adding.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)
		w.Header().Set("Content-Type", defaultContentType)
		ctx := r.Context()
		var account adding.Account
		if err := decoder.Decode(&account); err != nil {
			setJSONError(err, http.StatusBadRequest, w)
			return
		}

		id, err := a.AddAccount(ctx, account)
		if err != nil {
			if err.Error() == "this cpf could not be inserted in our DB" {
				setJSONError(err, http.StatusBadRequest, w)
				return
			}
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/%s", id))
		w.WriteHeader(http.StatusCreated)
	}
}

func getAccountBalanceByID(l listing.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Header().Set("Content-Type", defaultContentType)
		id := p.ByName("id")
		ctx := r.Context()

		balance, err := l.GetAccountBalanceByID(ctx, id)
		if err != nil {
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}

		_ = json.NewEncoder(w).Encode(listing.Account{Balance: balance})
		w.WriteHeader(http.StatusOK)
	}
}

package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"log"
	"net/http"
)

const (
	defaultContentType = "application/json"
)

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func setJSONError(err error, status int, w http.ResponseWriter) {
	log.Println(err.Error())
	var jsonDecodeErr *json.UnmarshalTypeError
	var message string
	if errors.As(err, &jsonDecodeErr) {
		message = fmt.Sprintf(
			"Invalid %s entity: expected type %s, got %s at field %s",
			jsonDecodeErr.Struct,
			jsonDecodeErr.Type.Name(),
			jsonDecodeErr.Value,
			jsonDecodeErr.Field)
	} else {
		message = err.Error()
	}
	response := ErrorResponse{
		StatusCode: status,
		Message:    message,
	}
	w.Header().Set("Content-Type", defaultContentType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func Handler(a adding.Service, l listing.Service) http.Handler {
	router := httprouter.New()

	router.POST("/accounts", addAccount(a))
	router.GET("/accounts", getAccounts(l))
	router.GET("/accounts/:id/balance", getAccountBalanceByID(l))

	return router
}

func addAccount(a adding.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		decoder := json.NewDecoder(r.Body)
		ctx := r.Context()
		var account adding.Account
		if err := decoder.Decode(&account); err != nil {
			setJSONError(err, http.StatusBadRequest, w)
			return
		}

		id, err := a.AddAccount(ctx, account)
		if err != nil {
			if err.Error() == mongodb.ErrCPFAlreadyExists.Error() {
				setJSONError(err, http.StatusBadRequest, w)
				return
			}
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Set("Content-Type", defaultContentType)
		w.Header().Set("Location", fmt.Sprintf("/%s", id))
		w.WriteHeader(http.StatusCreated)
	}
}

func getAccountBalanceByID(l listing.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id := p.ByName("id")
		ctx := r.Context()

		balance, err := l.GetAccountBalanceByID(ctx, id)
		if err != nil {
			if err.Error() == mongodb.ErrNoAccountWasFound.Error() {
				setJSONError(err, http.StatusNotFound, w)
				return
			}
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Set("Content-Type", defaultContentType)
		_ = json.NewEncoder(w).Encode(listing.Account{Balance: balance})
	}
}

func getAccounts(l listing.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := r.Context()
		accounts, err := l.GetAccounts(ctx)
		if err != nil {
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}
		w.Header().Set("Content-Type", defaultContentType)
		_ = json.NewEncoder(w).Encode(accounts)
	}
}

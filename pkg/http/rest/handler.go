package rest

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"io"
	"log"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func setJSONError(err error, status int, w io.Writer) {
	log.Println(err.Error())
	response := ErrorResponse{
		StatusCode: status,
		Message:    err.Error(),
	}

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
		w.Header().Set("Content-Type", "application/json")
		ctx := r.Context()
		var account adding.Account
		if err := decoder.Decode(&account); err != nil {
			setJSONError(err, http.StatusBadRequest, w)
			return
		}

		id, err := a.AddAccount(ctx, account)
		if err != nil {
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Set("Location", fmt.Sprintf("/%s", id))
		w.WriteHeader(http.StatusCreated)
	}
}

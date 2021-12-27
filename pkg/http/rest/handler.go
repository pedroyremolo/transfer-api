package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/sirupsen/logrus"
)

const (
	DefaultContentType = "application/json"
)

type AddingHandler interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
}

type TransferringHandler interface {
	MakeTransfer(w http.ResponseWriter, r *http.Request)
}

type AuthenticatingHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Authenticate(next http.HandlerFunc) http.HandlerFunc
}

type ListingHandler interface {
	GetBalanceByID(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	ListAllAccounts(w http.ResponseWriter, r *http.Request)
	GetUserTransfers(w http.ResponseWriter, r *http.Request)
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

var log *logrus.Logger

func Handler(logger *logrus.Entry, addingHandler AddingHandler, transferringHandler TransferringHandler, authenticatingHandler AuthenticatingHandler, listingHandler ListingHandler) http.Handler {
	router := httprouter.New()
	log = lgr.NewDefaultLogger()
	router.HandlerFunc(http.MethodPost, "/accounts", addingHandler.CreateAccount)
	router.HandlerFunc(http.MethodGet, "/accounts", listingHandler.ListAllAccounts)
	router.GET("/accounts/:id/balance", listingHandler.GetBalanceByID)

	router.HandlerFunc(http.MethodPost, "/login", authenticatingHandler.Login)
	router.HandlerFunc(http.MethodPost, "/transfers", authenticatingHandler.Authenticate(transferringHandler.MakeTransfer))
	router.HandlerFunc(http.MethodGet, "/transfers", listingHandler.GetUserTransfers)

	return router
}

func SetJSONError(logger *logrus.Entry, err error, status int, w http.ResponseWriter) {
	if logger == nil {
		logger = logrus.NewEntry(log)
	}
	logger.Errorf("Set err %v as JSON", err)
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
	w.Header().Set("Content-Type", DefaultContentType)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/pedroyremolo/transfer-api/pkg/adding"
	"github.com/pedroyremolo/transfer-api/pkg/authenticating"
	"github.com/pedroyremolo/transfer-api/pkg/listing"
	"github.com/pedroyremolo/transfer-api/pkg/log/lgr"
	"github.com/pedroyremolo/transfer-api/pkg/storage/mongodb"
	"github.com/pedroyremolo/transfer-api/pkg/transferring"
	"github.com/pedroyremolo/transfer-api/pkg/updating"
	"github.com/sirupsen/logrus"
)

const (
	defaultContentType = "application/json"
)

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

var log *logrus.Logger

func Handler(a adding.Service, l listing.Service, auth authenticating.Service, t transferring.Service, u updating.Service) http.Handler {
	router := httprouter.New()
	log = lgr.NewDefaultLogger()
	router.POST("/accounts", addAccount(a))
	router.GET("/accounts", getAccounts(l))
	router.GET("/accounts/:id/balance", getAccountBalanceByID(l))

	router.POST("/login", login(auth, l))

	router.POST("/transfers", transfer(a, auth, l, t, u))
	router.GET("/transfers", getAccountTransfers(auth, l))

	return router
}

func setJSONError(err error, status int, w http.ResponseWriter) {
	log.Errorf("Set err %v as JSON", err)
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

func decodeJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(v); err != nil {
		return err
	}
	return nil
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

func login(auth authenticating.Service, l listing.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		decoder := json.NewDecoder(r.Body)
		ctx := r.Context()
		var login authenticating.Login
		if err := decoder.Decode(&login); err != nil {
			setJSONError(err, http.StatusBadRequest, w)
			return
		}

		account, err := l.GetAccountByCPF(ctx, login.CPF)
		if err != nil {
			if err.Error() == mongodb.ErrNoAccountWasFound.Error() {
				setJSONError(authenticating.InvalidLoginErr, http.StatusForbidden, w)
				return
			}
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}
		var token authenticating.Token
		token, err = auth.Sign(ctx, login, account.Secret, account.ID)
		if err != nil {
			if err.Error() == authenticating.InvalidLoginErr.Error() {
				setJSONError(authenticating.InvalidLoginErr, http.StatusForbidden, w)
				return
			}
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Set("Content-Type", defaultContentType)
		_ = json.NewEncoder(w).Encode(authenticating.Token{Digest: token.Digest})
	}
}

func transfer(a adding.Service, auth authenticating.Service, l listing.Service, t transferring.Service, u updating.Service) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := r.Context()
		token, err := verifyTokenFromAuthHeader(r, auth, ctx)
		if err != nil {
			setJSONError(err, http.StatusUnauthorized, w)
			return
		}

		var transfer adding.Transfer
		if err = decodeJSON(r, &transfer); err != nil {
			setJSONError(err, http.StatusBadRequest, w)
			return
		}

		originBalance, err := l.GetAccountBalanceByID(ctx, token.ClientID)
		if err != nil {
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}
		destinationBalance, err := l.GetAccountBalanceByID(ctx, token.ClientID)
		if err != nil {
			if err.Error() == mongodb.ErrNoAccountWasFound.Error() {
				setJSONError(err, http.StatusBadRequest, w)
				return
			}
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}

		var origin, destination updating.Account
		origin.ID = token.ClientID
		destination.ID = transfer.DestinationAccountID
		transfer.OriginAccountID = token.ClientID

		origin.Balance, destination.Balance, err = t.BalanceBetweenAccounts(originBalance, destinationBalance, transfer.Amount)
		if err != nil {
			setJSONError(err, http.StatusBadRequest, w)
			return
		}
		if err = u.UpdateAccounts(ctx, origin, destination); err != nil {
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}
		if _, err = a.AddTransfer(ctx, transfer); err != nil {
			setJSONError(err, http.StatusInternalServerError, w)
		}
	}
}

func getAccountTransfers(auth authenticating.Service, l listing.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		ctx := r.Context()
		token, err := verifyTokenFromAuthHeader(r, auth, ctx)
		if err != nil {
			setJSONError(err, http.StatusUnauthorized, w)
			return
		}

		accountTransfers, err := l.GetTransfersByAccountID(ctx, token.ClientID)
		if err != nil {
			setJSONError(err, http.StatusInternalServerError, w)
			return
		}

		w.Header().Set("Content-Type", defaultContentType)
		_ = json.NewEncoder(w).Encode(accountTransfers)
	}
}

func verifyTokenFromAuthHeader(r *http.Request, auth authenticating.Service, ctx context.Context) (authenticating.Token, error) {
	authHd := r.Header.Get("Authorization")
	if authHd == "" {
		log.Error("Empty Authorization header")
		return authenticating.Token{}, authenticating.ProtectedRouteErr
	}
	hdSpaceIndex := strings.Index(authHd, " ")
	if hdSpaceIndex == -1 {
		log.Error("Authorization header is malformed")
		return authenticating.Token{}, authenticating.ProtectedRouteErr
	}
	authType, tokenDigest := authHd[:hdSpaceIndex], authHd[strings.LastIndex(authHd, " ")+1:]
	if authType != "Bearer" {
		log.Error("Authorization header is not of type Bearer")
		return authenticating.Token{}, authenticating.ProtectedRouteErr
	}

	token, err := auth.Verify(ctx, tokenDigest)
	if err != nil {
		log.Errorf("Err %v when trying to verify tokenDigest %s", err, tokenDigest)
		return authenticating.Token{}, authenticating.ProtectedRouteErr
	}
	return token, err
}

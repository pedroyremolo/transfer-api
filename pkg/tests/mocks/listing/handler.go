package listing

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HandlerMock struct {
}

func (h HandlerMock) GetBalanceByID(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	panic("not implemented") // TODO: Implement
}

func (h HandlerMock) ListAllAccounts(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h HandlerMock) GetUserTransfers(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

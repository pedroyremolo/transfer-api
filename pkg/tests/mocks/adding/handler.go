package adding

import "net/http"

type HandlerMock struct {
}

func (f HandlerMock) CreateAccount(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

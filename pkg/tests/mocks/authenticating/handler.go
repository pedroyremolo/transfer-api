package authenticating

import "net/http"

type HandlerMock struct {
}

func (h HandlerMock) Login(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

func (h HandlerMock) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return next
}

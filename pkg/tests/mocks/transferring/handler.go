package transferring

import "net/http"

type HandlerMock struct {
}

func (h HandlerMock) MakeTransfer(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
}

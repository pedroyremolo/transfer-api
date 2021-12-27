package helpers

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func AssertResponseJSON(t *testing.T, w *httptest.ResponseRecorder, expectedResponseJSON string) {
	t.Helper()
	respBodyBytes, err := ioutil.ReadAll(w.Body)
	respBody := string(bytes.TrimSpace(respBodyBytes))
	if err != nil {
		t.Fatal("Unable to read response from Recorder")
	}
	if respBody != expectedResponseJSON {
		t.Errorf("Expected response %s; got %s", expectedResponseJSON, respBody)
	}
}

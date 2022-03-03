package glockify

import (
	"net/http"
)

const (
	dummyAPIKey = "dummy"
)

func checkAuthHeader(r *http.Request) bool {
	api := r.Header.Get("X-Api-Key")
	return api == dummyAPIKey
}

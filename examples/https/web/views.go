package web

import (
	"bitbucket.org/aukbit/pluto/reply"
	"net/http"
)

func GetHandler (w http.ResponseWriter, r *http.Request){
	// TODO decorate this header
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	reply.Json(w, r, http.StatusOK, `{"message":"Hello Gopher"}`)
}

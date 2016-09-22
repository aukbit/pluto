package web

import (
	"bitbucket.org/aukbit/pluto/reply"
	"net/http"
	"log"
)

func GetHandler (w http.ResponseWriter, r *http.Request){
	// get context
	ctx := r.Context()
	// get service from context by service name
	s := ctx.Value("pluto_web")
	log.Printf("teste %s", s)
	reply.Json(w, r, http.StatusOK, `{"message":"Hello Gopher"}`)
}

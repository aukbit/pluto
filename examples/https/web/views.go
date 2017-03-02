package web

import (
	"github.com/aukbit/pluto/reply"
	"net/http"
)

type Message struct {
	Message    string	`json:"message"`
}

func GetHandler (w http.ResponseWriter, r *http.Request){
	m := &Message{"Hello Gopher"}
	reply.Json(w, r, http.StatusOK, m)
}

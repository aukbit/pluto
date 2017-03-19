package router

import "net/http"

// Mux interface to expose router struct
type Mux interface {
	GET(string, HandlerFunc)
	POST(string, HandlerFunc)
	PUT(string, HandlerFunc)
	DELETE(string, HandlerFunc)
	HandleFunc(string, string, HandlerFunc)
	ServeHTTP(http.ResponseWriter, *http.Request)
	WrapperMiddleware(...Middleware)
}

// New creates a new router
func New() *Router {
	return NewRouter()
}

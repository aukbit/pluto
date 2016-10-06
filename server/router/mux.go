package router

import "net/http"

// Mux interface to expose router struct
type Mux interface {
	GET(string, Handler)
	POST(string, Handler)
	PUT(string, Handler)
	DELETE(string, Handler)
	Handle(string, string, Handler)
	ServeHTTP(http.ResponseWriter, *http.Request)
	AddMiddleware(...Middleware) // deprecated
	WrapperMiddleware(...Middleware)
}

// NewMux creates a new router
func NewMux() Mux {
	return newRouter()
}

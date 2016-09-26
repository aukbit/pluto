package server

import (
	"net"
	"crypto/tls"
	"net/http"
	"bitbucket.org/aukbit/pluto/server/router"
)

type httpsServer struct {
	defaultServer
}

// listenTLS based on http.ListenAndServeTLS
// listens on the TCP network address srv.Addr
// If srv.Addr is blank, ":https" is used.
// returns nil or new listener
func (s *httpsServer) listen() (net.Listener, error) {

	addr := s.cfg.Addr
	if addr == "" {
		addr = ":https"
	}

	ln, err := tls.Listen("tcp", addr, s.cfg.TLSConfig)
	if err != nil {
		return nil, err
	}
	return ln, nil
}

// middlewareStrictSecurityHeader Middleware to wrap all handlers with
// Strict-Transport-Security header
func middlewareStrictSecurityHeader() router.Middleware {
    return func(h router.Handler) router.Handler {
        return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		h.ServeHTTP(w, r)
	}
    }
}


package server

import (
	"net"
	"crypto/tls"
	"net/http"
	"bitbucket.org/aukbit/pluto/server/router"
	"log"
	"syscall"
	"os/signal"
	"os"
)

type httpsServer struct {
	defaultServer
}

// Run Server
func (s *httpsServer) Run() error {
	if err := s.start(); err != nil {
		return err
	}
	// parse address for host, port
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	sig := <-ch
	log.Printf("----- %s signal %v received ", s.cfg.Name, sig)
	return s.Stop()
}

func (s *httpsServer) start() error {

	ln, err := s.listen()
	if err != nil {
		return err
	}

	if err := s.serve(ln); err != nil {
		return err
	}
	go s.waitSignal(ln)
	return nil
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


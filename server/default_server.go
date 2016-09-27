package server

import (
	"log"
	"net"
	"net/http"
	"time"
	"syscall"
	"os/signal"
	"os"
	"crypto/tls"
	"bitbucket.org/aukbit/pluto/server/router"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type defaultServer struct {
	cfg 			*Config
	close 			chan bool
}

// NewServer will instantiate a new defaultServer with the given config
func newServer(cfgs ...ConfigFunc) *defaultServer {
	c := newConfig(cfgs...)
	return &defaultServer{c, make(chan bool)}
}

// Run Server
func (s *defaultServer) Run() error {
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

// Stop stops server by sending a message to close the listener via channel
func (s *defaultServer) Stop() error {
	s.close <-true
	return nil
}

func (s *defaultServer) Config() *Config {
	cfg := s.cfg
	return cfg
}

func (s *defaultServer) start() (err error) {

	var ln net.Listener

	switch s.cfg.Format {
	case "https":
		ln, err = s.listenTLS()
		if err != nil {
			return err
		}
	default:
		ln, err = s.listen()
		if err != nil {
			return err
		}
	}

	switch s.cfg.Format {
	case "grpc":
		if err := s.serveGRPC(ln); err != nil {
			return err
		}
	default:
		if err := s.serve(ln); err != nil {
			return err
		}
	}

	go s.waitSignal(ln)
	return nil
}

// listen based on http.ListenAndServe
// listens on the TCP network address srv.Addr
// If srv.Addr is blank, ":http" is used.
// returns nil or new listener
func (s *defaultServer) listen() (net.Listener, error) {

	addr := s.cfg.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	ln = net.Listener(TcpKeepAliveListener{ln.(*net.TCPListener)})

	return ln, nil
}

// listenTLS based on http.ListenAndServeTLS
// listens on the TCP network address srv.Addr
// If srv.Addr is blank, ":https" is used.
// returns nil or new listener
func (s *defaultServer) listenTLS() (net.Listener, error) {

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

// serve based on http.ListenAndServe
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// serve always returns a non-nil error.
func (s *defaultServer) serve(ln net.Listener) error {

	srv := &http.Server{
		// handler to invoke, http.DefaultServeMux if nil
		Handler: s.cfg.Mux,

		// ReadTimeout is used by the http server to set a maximum duration before
		// timing out read of the request. The default timeout is 10 seconds.
		ReadTimeout: 10 * time.Second,

		// WriteTimeout is used by the http server to set a maximum duration before
		// timing out write of the response. The default timeout is 10 seconds.
		WriteTimeout: 10 * time.Second,

		TLSConfig:    s.cfg.TLSConfig,
	}
	go func() {
		if err := srv.Serve(ln); err != nil {
			log.Fatalf("ERROR %s srv.Serve(ln) %v", s.cfg.Name, err)
		}
	}()

	log.Printf("----- %s %s listening on %s", s.cfg.Format, s.cfg.Name, ln.Addr().String())
	return nil
}

// serve serves *grpc.Server
func (s *defaultServer) serveGRPC(ln net.Listener) (err error) {

	srv := s.cfg.GRPCServer

	go func() {
		if err := srv.Serve(ln); err != nil {
			log.Fatalf("ERROR %s g.Serve(lis) %v", s.cfg.Name, err)
		}
	}()

	log.Printf("----- %s %s listening on %s", s.cfg.Format, s.cfg.Name, ln.Addr().String())
	return nil
}

// waitSignal to be used as go routine waiting for a signal to stop the service
func (s *defaultServer) waitSignal (ln net.Listener) {
	// Waits for call to stop
	<-s.close
	log.Printf("CLOSE %s %s received", s.cfg.Format, s.cfg.Name)
	// close listener
	if err := ln.Close(); err != nil {
		log.Fatalf("ERROR %s %s ln.Close() %v", s.cfg.Format, s.cfg.Name, err)
	}
	log.Printf("----- %s %s listener closed", s.cfg.Format, s.cfg.Name)
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
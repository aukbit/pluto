package server

import (
	"log"
	"net"
	"net/http"
	"time"
	"syscall"
	"os/signal"
	"os"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type defaultServer struct {
	cfg 			*Config
	close 			chan bool
}

// NewServer will instantiate a new defaultServer with the given config
func newServer(cfgs ...ConfigFunc) Server {
	c := newConfig(cfgs...)
	switch c.Format {
	case "grpc":
		ds := defaultServer{c, make(chan bool)}
		return &grpcServer{ds}
	case "https":
		c.Mux.AddMiddleware(middlewareStrictSecurityHeader())
		ds := defaultServer{c, make(chan bool)}
		return &httpsServer{ds}
	default:
		return &defaultServer{c, make(chan bool)}
	}
	return nil
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

func (s *defaultServer) start() error {

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

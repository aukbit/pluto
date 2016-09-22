package server

import (
	"log"
	"net"
	"net/http"
	"crypto/tls"
	"time"
	"syscall"
	"os/signal"
	"os"
	"errors"
	"bitbucket.org/aukbit/pluto/server/router"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type defaultServer struct {
	cfg 			*Config
	mux			router.Mux
	// close chan for graceful shutdown
	close 			chan bool
}

// NewServer will instantiate a new Server with the given config
func newDefaultServer(cfgs ...ConfigFunc) Server {
	c := newConfig(cfgs...)
	return &defaultServer{cfg: c, mux: c.Mux, close: make(chan bool)}
}

func (s *defaultServer) Init(cfgs ...ConfigFunc) error {
	for _, c := range cfgs {
		c(s.cfg)
	}
	s.mux = s.cfg.Mux
	return nil
}

func (s *defaultServer) Config() *Config {
	cfg := s.cfg
	return cfg
}

// Run
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

// Stop sends message to close the listener via channel
func (s *defaultServer) Stop() error {
	s.close <-true
	return nil
}

// start start the Server
func (s *defaultServer) start() (err error) {
	log.Printf("START %s %s \t%s", s.cfg.Format, s.cfg.Name, s.cfg.Id)
	if s.mux == nil{
		return errors.New("Handlers not set up. Server will not start.")
	}
	var ln net.Listener
	if s.cfg.Format == "https" {
		ln, err = s.listenTLS()
		if err != nil {
			return err
		}
	} else {
		ln, err = s.listen()
		if err != nil {
			return err
		}
	}
	err = s.serve(ln)
	if err != nil {
		return err
	}
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
		Handler: s.mux,

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
	//
	log.Printf("----- %s %s listening on %s", s.cfg.Format, s.cfg.Name, ln.Addr().String())
	go func() {
		// Waits for call to stop
		<-s.close
		log.Printf("CLOSE %s received", s.cfg.Name)
		// close listener
		if err := ln.Close(); err != nil {
			log.Fatalf("ERROR %s ln.Close() %v", s.cfg.Name, err)
		}
		log.Printf("----- %s listener closed", s.cfg.Name)
	}()

	return nil
}

// listenAndServe based on http.ListenAndServe
// listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// If srv.Addr is blank, ":http" is used.
// ListenAndServe always returns a non-nil error.
func (s *defaultServer) listenAndServe() error {
	addr := s.cfg.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	ln = net.Listener(TcpKeepAliveListener{ln.(*net.TCPListener)})

	srv := http.Server{
		// handler to invoke, http.DefaultServeMux if nil
		Handler: s.mux,

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
	//
	log.Printf("----- %s %s listening on %s", s.cfg.Format, s.cfg.Name, ln.Addr().String())
	//
	go func() {
		// Waits for call to stop
		<-s.close
		log.Printf("CLOSE %s received", s.cfg.Name)
		// close listener
		if err := ln.Close(); err != nil {
			log.Fatalf("ERROR %s ln.Close() %v", s.cfg.Name, err)
		}
		log.Printf("----- %s listener closed", s.cfg.Name)
	}()

	return nil
}
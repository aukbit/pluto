package server

import (
	"log"
	"net"
	"net/http"
	"pluto/server/router"
	"time"
	"syscall"
	"os/signal"
	"os"
	"errors"
	"fmt"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type defaultServer struct {
	cfg 			*Config
	mux 			*router.Router
	// close chan for graceful shutdown
	close 			chan bool
}

// NewServer will instantiate a new Server with the given config
func newDefaultServer(cfgs ...ConfigFunc) Server {
	c := newConfig(cfgs...)
	return &defaultServer{cfg: c, close: make(chan bool)}
}

func (s *defaultServer) Init(cfgs ...ConfigFunc) error {
	for _, c := range cfgs {
		c(s.cfg)
	}
	return nil
}

func (s *defaultServer) Router(mux *router.Router) error {
	if mux == nil {
		s.mux = router.NewRouter()
	} else {
		s.mux = mux
	}
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

// start start the Server
func (s *defaultServer) start() error {
	log.Printf("START %s %s", s.cfg.Name, s.cfg.Id)
	if s.mux == nil{
		return errors.New("Handlers not set up. Server will not start.")
	}
	// start go routine
	go func(){
		if err := s.listenAndServe(); err != nil{
			log.Fatal(fmt.Sprintf("ERROR s.listenAndServe() %v", err))
		}
	}()
	return nil
}

// Stop sends message to close the listener via channel
func (s *defaultServer) Stop() error {
	s.close <-true
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

	// TODO config option for TLS

	httpServer := http.Server{
		// handler to invoke, http.DefaultServeMux if nil
		Handler: s.mux,

		// ReadTimeout is used by the http server to set a maximum duration before
		// timing out read of the request. The default timeout is 10 seconds.
		ReadTimeout: 10 * time.Second,

		// WriteTimeout is used by the http server to set a maximum duration before
		// timing out write of the response. The default timeout is 10 seconds.
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		if err := httpServer.Serve(ln); err != nil {
			log.Fatal(fmt.Sprintf("ERROR httpServer.Serve(ln) %v", err))
		}
	}()
	//
	log.Printf("----- %s listening on %s", s.cfg.Name, ln.Addr().String())
	//
	go func() {
		// Waits for call to stop
		<-s.close
		log.Printf("CLOSE %s received", s.cfg.Name)
		// close listener
		if err := ln.Close(); err != nil {
			log.Fatal(fmt.Sprintf("ERROR ln.Close() %v", err))
		}
		log.Printf("----- %s listener closed", s.cfg.Name)
	}()

	return nil
}
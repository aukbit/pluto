package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/uber-go/zap"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type defaultServer struct {
	cfg   *Config
	close chan bool
	wg    *sync.WaitGroup
}

// NewServer will instantiate a new defaultServer with the given config
func newServer(cfgs ...ConfigFunc) *defaultServer {
	c := newConfig(cfgs...)
	return &defaultServer{cfg: c, close: make(chan bool), wg: &sync.WaitGroup{}}
}

// Run Server
func (s *defaultServer) Run() error {
	if err := s.start(); err != nil {
		return err
	}
	// wait for go routines to finish
	s.wg.Wait()
	logger.Info("exit",
		zap.String("server", s.cfg.Name),
		zap.String("id", s.cfg.ID),
		zap.String("format", s.cfg.Format),
	)
	return nil
}

// Stop stops server by sending a message to close the listener via channel
func (s *defaultServer) Stop() {
	logger.Info("stop",
		zap.String("server", s.cfg.Name),
		zap.String("id", s.cfg.ID),
		zap.String("format", s.cfg.Format),
	)
	s.close <- true
}

func (s *defaultServer) Config() *Config {
	cfg := s.cfg
	return cfg
}

func (s *defaultServer) start() (err error) {
	logger.Info("start",
		zap.String("server", s.cfg.Name),
		zap.String("id", s.cfg.ID),
		zap.String("format", s.cfg.Format))
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

	// add go routine to WaitGroup
	s.wg.Add(1)
	go s.waitUntilStop(ln)
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
	ln = net.Listener(TCPKeepAliveListener{ln.(*net.TCPListener)})

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

		TLSConfig: s.cfg.TLSConfig,
	}
	// add go routine to WaitGroup
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := srv.Serve(ln); err != nil {
			if err.Error() == errClosing(ln).Error() {
				return
			}
			logger.Error("Serve(ln)",
				zap.String("server", s.cfg.Name),
				zap.String("id", s.cfg.ID),
				zap.String("format", s.cfg.Format),
				zap.String("port", ln.Addr().String()),
				zap.String("err", err.Error()))
			return
		}
	}()
	return nil
}

// errClosing is the currently error raised when we gracefull
// close the listener. If this is the case there is no point to log
func errClosing(ln net.Listener) error {
	return fmt.Errorf("accept tcp %v: use of closed network connection", ln.Addr().String())
}

// waitUntilStop waits for close channel
func (s *defaultServer) waitUntilStop(ln net.Listener) {
	defer s.wg.Done()
outer:
	for {
		select {
		case <-s.close:
			// Waits for call to stop
			if err := ln.Close(); err != nil {
				logger.Error("Close()",
					zap.String("server", s.cfg.Name),
					zap.String("id", s.cfg.ID),
					zap.String("format", s.cfg.Format),
					zap.String("port", ln.Addr().String()),
					zap.String("err", err.Error()))
			}
			break outer
		default:
			logger.Info("pulse",
				zap.String("server", s.cfg.Name),
				zap.String("id", s.cfg.ID),
				zap.String("format", s.cfg.Format),
				zap.String("port", ln.Addr().String()))
			time.Sleep(time.Second * 1)
			continue
		}
	}
}

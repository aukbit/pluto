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
	cfg    *Config
	close  chan bool
	wg     *sync.WaitGroup
	logger zap.Logger
}

// NewServer will instantiate a new defaultServer with the given config
func newServer(cfgs ...ConfigFunc) *defaultServer {
	c := newConfig(cfgs...)
	ds := &defaultServer{
		cfg:    c,
		close:  make(chan bool),
		wg:     &sync.WaitGroup{},
		logger: zap.New(zap.NewJSONEncoder())}
	ds.initLog()
	return ds
}

func (ds *defaultServer) initLog() {
	ds.logger = ds.logger.With(
		zap.String("object", "server"),
		zap.String("id", ds.cfg.ID),
		zap.String("name", ds.cfg.Name),
		zap.String("format", ds.cfg.Format),
		zap.String("port", ds.cfg.Addr))
}

// Run Server
func (ds *defaultServer) Run() error {
	if err := ds.start(); err != nil {
		return err
	}
	// wait for go routines to finish
	ds.wg.Wait()
	ds.logger.Info("exit")
	return nil
}

// Stop stops server by sending a message to close the listener via channel
func (ds *defaultServer) Stop() {
	ds.logger.Info("stop")
	ds.close <- true
}

func (ds *defaultServer) Config() *Config {
	cfg := ds.cfg
	return cfg
}

func (ds *defaultServer) start() (err error) {
	ds.logger.Info("start")
	var ln net.Listener

	switch ds.cfg.Format {
	case "https":
		ln, err = ds.listenTLS()
		if err != nil {
			return err
		}
	default:
		ln, err = ds.listen()
		if err != nil {
			return err
		}
	}

	switch ds.cfg.Format {
	case "grpc":
		if err := ds.serveGRPC(ln); err != nil {
			return err
		}
	default:
		if err := ds.serve(ln); err != nil {
			return err
		}
	}

	// add go routine to WaitGroup
	ds.wg.Add(1)
	go ds.waitUntilStop(ln)
	return nil
}

// listen based on http.ListenAndServe
// listens on the TCP network address srv.Addr
// If srv.Addr is blank, ":http" is used.
// returns nil or new listener
func (ds *defaultServer) listen() (net.Listener, error) {

	addr := ds.cfg.Addr
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
func (ds *defaultServer) listenTLS() (net.Listener, error) {

	addr := ds.cfg.Addr
	if addr == "" {
		addr = ":https"
	}

	ln, err := tls.Listen("tcp", addr, ds.cfg.TLSConfig)
	if err != nil {
		return nil, err
	}
	return ln, nil
}

// serve based on http.ListenAndServe
// calls Serve to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
// serve always returns a non-nil error.
func (ds *defaultServer) serve(ln net.Listener) error {

	srv := &http.Server{
		// handler to invoke, http.DefaultServeMux if nil
		Handler: ds.cfg.Mux,

		// ReadTimeout is used by the http server to set a maximum duration before
		// timing out read of the request. The default timeout is 10 seconds.
		ReadTimeout: 10 * time.Second,

		// WriteTimeout is used by the http server to set a maximum duration before
		// timing out write of the response. The default timeout is 10 seconds.
		WriteTimeout: 10 * time.Second,

		TLSConfig: ds.cfg.TLSConfig,
	}
	// add go routine to WaitGroup
	ds.wg.Add(1)
	go func() {
		defer ds.wg.Done()
		if err := srv.Serve(ln); err != nil {
			if err.Error() == errClosing(ln).Error() {
				return
			}
			ds.logger.Error("Serve(ln)", zap.String("err", err.Error()))
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
func (ds *defaultServer) waitUntilStop(ln net.Listener) {
	defer ds.wg.Done()
outer:
	for {
		select {
		case <-ds.close:
			// Waits for call to stop
			if err := ln.Close(); err != nil {
				ds.logger.Error("Close()", zap.String("err", err.Error()))
			}
			break outer
		default:
			ds.logger.Info("pulse")
			time.Sleep(time.Second * 1)
			continue
		}
	}
}

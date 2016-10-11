package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"bitbucket.org/aukbit/pluto/discovery"

	"google.golang.org/grpc"

	"github.com/uber-go/zap"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type defaultServer struct {
	cfg          *Config
	close        chan bool
	wg           *sync.WaitGroup
	logger       zap.Logger
	httpServer   *http.Server
	grpcServer   *grpc.Server
	isDiscovered bool
}

// newServer will instantiate a new defaultServer with the given config
func newServer(cfgs ...ConfigFunc) *defaultServer {
	c := newConfig(cfgs...)
	return &defaultServer{
		cfg:    c,
		close:  make(chan bool),
		wg:     &sync.WaitGroup{},
		logger: zap.New(zap.NewJSONEncoder())}
}

// Run Server
func (ds *defaultServer) Run(cfgs ...ConfigFunc) error {
	// set last configs
	for _, c := range cfgs {
		c(ds.cfg)
	}
	// set logger
	ds.setLogger()

	// register at service discovery
	if err := ds.register(); err != nil {
		return err
	}
	// start server
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
	ds.wg.Add(1)
	go ds.unregister()
	ds.close <- true
}

func (ds *defaultServer) Config() *Config {
	cfg := ds.cfg
	return cfg
}

func (ds *defaultServer) setLogger() {
	ds.logger = ds.logger.With(
		zap.Nest("server",
			zap.String("id", ds.cfg.ID),
			zap.String("name", ds.cfg.Name),
			zap.String("format", ds.cfg.Format),
			zap.String("port", ds.cfg.Addr),
			zap.String("parent", ds.cfg.ParentID)))
}

func (ds *defaultServer) setHttpServer() {
	// append logger
	ds.cfg.Middlewares = append(ds.cfg.Middlewares, loggerMiddleware(ds))
	// wrap Middlewares
	ds.cfg.Mux.WrapperMiddleware(ds.cfg.Middlewares...)
	// initialize http server
	ds.httpServer = &http.Server{
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
}

func (ds *defaultServer) start() (err error) {
	ds.logger.Info("start")
	var ln net.Listener

	switch ds.cfg.Format {
	case "https":
		// append strict security header
		ds.cfg.Middlewares = append(ds.cfg.Middlewares, strictSecurityHeaderMiddleware())
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
		ds.setGRPCServer()
		if err := ds.serveGRPC(ln); err != nil {
			return err
		}
	default:
		ds.setHttpServer()
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

	// add go routine to WaitGroup
	ds.wg.Add(1)
	go func(srv *http.Server) {
		defer ds.wg.Done()
		if err := srv.Serve(ln); err != nil {
			if err.Error() == errClosing(ln).Error() {
				return
			}
			ds.logger.Error("Serve(ln)", zap.String("err", err.Error()))
			return
		}
	}(ds.httpServer)
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
			switch ds.cfg.Format {
			case "grpc":
				ds.grpcServer.GracefulStop()
			default:
				if err := ln.Close(); err != nil {
					ds.logger.Error("Close()", zap.String("err", err.Error()))
				}
			}
			break outer
		default:
			ds.logger.Debug("pulse")
			time.Sleep(time.Second * 1)
			continue
		}
	}
}

// register Server within the service discovery system
func (ds *defaultServer) register() error {
	_, err := discovery.IsAvailable()
	if err != nil {
		ds.logger.Warn("service discovery not available")
		return nil
	}
	s := &discovery.Service{
		ID:   ds.cfg.ID,
		Name: ds.cfg.Name,
		Port: ds.cfg.Port(),
		Tags: []string{ds.cfg.ID, ds.cfg.Version},
	}
	err = discovery.RegisterService(s)
	if err != nil {
		return err
	}
	c := &discovery.Check{
		ID:    fmt.Sprintf("%s_check", ds.cfg.ID),
		Name:  fmt.Sprintf("Service '%s' check", ds.cfg.Name),
		Notes: fmt.Sprintf("Ensure the server is listening on port %s", ds.cfg.Addr),
		DeregisterCriticalServiceAfter: "10m",
		TCP:       ds.cfg.Addr,
		Interval:  "10s",
		Timeout:   "1s",
		ServiceID: ds.cfg.ID,
	}
	err = discovery.RegisterCheck(c)
	if err != nil {
		return err
	}
	ds.isDiscovered = true
	return nil
}

// unregister Server from the service discovery system
func (ds *defaultServer) unregister() {
	defer ds.wg.Done()
	if ds.isDiscovered {
		err := discovery.DeregisterService(ds.cfg.ID)
		if err != nil {
			ds.logger.Error(err.Error())
		}
	}
}

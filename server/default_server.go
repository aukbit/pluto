package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/aukbit/pluto/server/router"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc"
)

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type defaultServer struct {
	cfg   *Config
	close chan bool
	wg    *sync.WaitGroup
	// logger     *zap.Logger
	httpServer *http.Server
	grpcServer *grpc.Server
	health     *health.Server
}

// newServer will instantiate a new defaultServer with the given config
func newServer(cfgs ...ConfigFunc) *defaultServer {
	c := newConfig(cfgs...)
	d := &defaultServer{
		cfg:    c,
		close:  make(chan bool),
		wg:     &sync.WaitGroup{},
		health: health.NewServer(),
	}
	// d.logger, _ = zap.NewProduction()
	return d
}

// Run Server
func (ds *defaultServer) Run(cfgs ...ConfigFunc) error {
	// set last configs
	for _, c := range cfgs {
		c(ds.cfg)
	}
	// register at service discovery
	if err := ds.register(); err != nil {
		return err
	}
	// set logger
	ds.setLogger()
	// start server
	if err := ds.start(); err != nil {
		return err
	}
	// set health
	ds.health.SetServingStatus(ds.cfg.ID, 1)
	// wait for go routines to finish
	ds.wg.Wait()
	// ds.logger.Info("exit")
	return nil
}

// Stop stops server by sending a message to close the listener via channel
func (ds *defaultServer) Stop() {
	// ds.logger.Info("stop")
	// set health as not serving
	ds.health.SetServingStatus(ds.cfg.ID, 2)
	// close listener
	ds.close <- true
}

func (ds *defaultServer) Config() *Config {
	return ds.cfg
}

func (ds *defaultServer) Health() *healthpb.HealthCheckResponse {
	switch ds.cfg.Format {
	case "grpc":
		ds.healthGRPC()
	default:
		ds.healthHTTP()
	}
	hcr, err := ds.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: ds.cfg.ID})
	if err != nil {
		// ds.logger.Error("Health", zap.String("err", err.Error()))
		return &healthpb.HealthCheckResponse{Status: 2}
	}
	return hcr
}

func (ds *defaultServer) setLogger() {
	// ds.logger = ds.logger.With(
	// 	zap.String("type", "server"),
	// 	zap.String("id", ds.cfg.ID),
	// 	zap.String("name", ds.cfg.Name),
	// 	zap.String("format", ds.cfg.Format),
	// 	zap.String("port", ds.cfg.Addr),
	// 	zap.String("parent", ds.cfg.ParentID))
}

func (ds *defaultServer) setHTTPServer() {
	if ds.cfg.Mux == nil {
		ds.cfg.Mux = router.NewMux()
	}
	// set health check handler
	ds.cfg.Mux.GET("/_health", router.Wrap(healthHandler, serverMiddleware(ds)))
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
	// ds.logger.Info("start")
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
		ds.setHTTPServer()
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
			// ds.logger.Error("Serve(ln)", zap.String("err", err.Error()))
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
			ds.unregister()
			switch ds.cfg.Format {
			case "grpc":
				ds.grpcServer.GracefulStop()
			default:
				if err := ln.Close(); err != nil {
					// ds.logger.Error("Close()", zap.String("err", err.Error()))
				}
			}
			break outer
		default:
			// ds.logger.Debug("pulse")
			time.Sleep(time.Second * 1)
			continue
		}
	}
}

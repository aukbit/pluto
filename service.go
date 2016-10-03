package pluto

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/uber-go/zap"

	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/datastore"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
)

// Service
type service struct {
	cfg   *Config
	close chan bool
	wg    *sync.WaitGroup
}

func newService(cfgs ...ConfigFunc) *service {
	c := newConfig(cfgs...)
	s := &service{cfg: c, close: make(chan bool), wg: &sync.WaitGroup{}}
	for _, srv := range c.Servers {
		// Wrap this service to all handlers
		// make it available in handler context
		if srv.Config().Format == "http" {
			srv.Config().Mux.AddMiddleware(middlewareService(s))
		}
	}
	return s
}

// Init TODO should be removed.. redundant makes initialization confusing
func (s *service) Init(cfgs ...ConfigFunc) error {
	for _, c := range cfgs {
		c(s.cfg)
	}
	for _, srv := range s.cfg.Servers {
		// Wrap this service to all handlers
		// make it available in handler context
		if srv.Config().Format == "http" {
			srv.Config().Mux.AddMiddleware(middlewareService(s))
		}
	}
	return nil
}

// Run starts service
func (s *service) Run() error {
	if err := s.start(); err != nil {
		return err
	}
	// wait for all go routines to finish
	s.wg.Wait()
	logger.Info("exit",
		zap.String("service", s.cfg.Name),
		zap.String("id", s.cfg.ID),
	)
	return nil
}

// Stop stops service
func (s *service) Stop() {
	logger.Info("stop",
		zap.String("service", s.cfg.Name),
		zap.String("id", s.cfg.ID),
	)
	s.close <- true
}

// Config service configration options
func (s *service) Config() *Config {
	cfg := s.cfg
	return cfg
}

// Server returns a server instance by name if initialized in service
func (s *service) Server(name string) (srv server.Server) {
	var ok bool
	if srv, ok = s.cfg.Servers[name]; !ok {
		return nil
	}
	return srv
}

// Client returns a client instance by name if initialized in service
func (s *service) Client(name string) (clt client.Client) {
	var ok bool
	if clt, ok = s.cfg.Clients[name]; !ok {
		return nil
	}
	return clt
}

// Datastore TODO there is no need to be public
func (s *service) Datastore() datastore.Datastore {
	return s.cfg.Datastore
}

func (s *service) start() error {
	logger.Info("start",
		zap.String("service", s.cfg.Name),
		zap.String("id", s.cfg.ID),
		zap.Nest("content",
			zap.Int("servers", len(s.cfg.Servers)),
			zap.Int("clients", len(s.cfg.Clients))),
	)

	// connect datastore
	if s.cfg.Datastore != nil {
		s.cfg.Datastore.Connect()
		if err := s.cfg.Datastore.RefreshSession(); err != nil {
			logger.Error("RefreshSession()",
				zap.String("service", s.cfg.Name),
				zap.String("id", s.cfg.ID),
				zap.String("err", err.Error()),
			)
		}
	}

	// TODO: manage errors
	// run servers
	s.startServers()
	// dial clients
	s.startClients()
	// add go routine to WaitGroup
	s.wg.Add(1)
	go s.waitUntilStopOrSig()
	return nil
}

func (s *service) startServers() {
	for _, srv := range s.cfg.Servers {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(ss server.Server) {
			defer s.wg.Done()
			if err := ss.Run(); err != nil {
				logger.Error("Run()",
					zap.String("service", s.cfg.Name),
					zap.String("id", s.cfg.ID),
					zap.String("err", err.Error()),
				)
			}
		}(srv)
	}
}

func (s *service) startClients() {
	for _, clt := range s.cfg.Clients {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(cc client.Client) {
			defer s.wg.Done()
			if err := cc.Dial(); err != nil {
				logger.Error("Dial()",
					zap.String("service", s.cfg.Name),
					zap.String("id", s.cfg.ID),
					zap.String("err", err.Error()),
				)
			}
		}(clt)
	}
}

// waitUntilStopOrSig waits for close channel or syscall Signal
func (s *service) waitUntilStopOrSig() {
	defer s.wg.Done()
	//  Stop also in case of any host signal
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT)

outer:
	for {
		select {
		case <-s.close:
			// Waits for call to stop
			s.stopServers()
			break outer
		case sig := <-sigch:
			// Waits for signal to stop
			logger.Info("signal received",
				zap.String("service", s.cfg.Name),
				zap.String("id", s.cfg.ID),
				zap.String("signal", sig.String()))
			s.stopServers()
			break outer
		default:
			logger.Info("pulse",
				zap.String("service", s.cfg.Name),
				zap.String("id", s.cfg.ID))
			time.Sleep(time.Second * 1)
			continue
		}
	}
}

func (s *service) stopServers() {
	for _, srv := range s.cfg.Servers {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(ss server.Server) {
			defer s.wg.Done()
			ss.Stop()
		}(srv)
	}
}

func middlewareService(s *service) router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, s.cfg.Name, s)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

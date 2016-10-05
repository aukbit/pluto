package pluto

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/uber-go/zap"

	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/datastore"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/server/router"
)

// Service
type service struct {
	cfg    *Config
	close  chan bool
	wg     *sync.WaitGroup
	logger zap.Logger
}

func newService(cfgs ...ConfigFunc) *service {
	c := newConfig(cfgs...)
	s := &service{
		cfg:    c,
		close:  make(chan bool),
		wg:     &sync.WaitGroup{},
		logger: zap.New(zap.NewJSONEncoder())}
	for _, srv := range c.Servers {
		// Wrap this service to all handlers
		// make it available in handler context
		if strings.Contains(srv.Config().Format, "http") {
			srv.Config().Mux.AddMiddleware(middlewareService(s))
		}
	}
	s.initLog()
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
		if strings.Contains(srv.Config().Format, "http") {
			srv.Config().Mux.AddMiddleware(middlewareService(s))
		}
	}
	s.initLog()
	return nil
}

// Run starts service
func (s *service) Run() error {
	if err := s.start(); err != nil {
		return err
	}
	// wait for all go routines to finish
	s.wg.Wait()
	s.logger.Info("exit")
	return nil
}

// Stop stops service
func (s *service) Stop() {
	s.logger.Info("stop")
	s.close <- true
}

// Config service configration options
func (s *service) Config() *Config {
	cfg := s.cfg
	return cfg
}

// Server returns a server instance by name if initialized in service
func (s *service) Server(name string) (srv server.Server, ok bool) {
	if srv, ok = s.cfg.Servers[name]; !ok {
		return
	}
	return srv, true
}

// Client returns a client instance by name if initialized in service
func (s *service) Client(name string) (clt client.Client, ok bool) {
	if clt, ok = s.cfg.Clients[name]; !ok {
		return
	}
	return clt, true
}

// Datastore TODO there is no need to be public
func (s *service) Datastore() datastore.Datastore {
	return s.cfg.Datastore
}

func (s *service) initLog() {
	s.logger = s.logger.With(
		zap.String("object", "service"),
		zap.String("id", s.cfg.ID),
		zap.String("name", s.cfg.Name))
}

func (s *service) start() error {
	s.logger.Info("start",
		zap.Nest("content",
			zap.Int("servers", len(s.cfg.Servers)),
			zap.Int("clients", len(s.cfg.Clients))))

	// connect datastore
	if s.cfg.Datastore != nil {
		s.cfg.Datastore.Connect()
		if err := s.cfg.Datastore.RefreshSession(); err != nil {
			s.logger.Error("RefreshSession()", zap.String("err", err.Error()))
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
		go func(srv server.Server) {
			defer s.wg.Done()
			if err := srv.Run(); err != nil {
				s.logger.Error("Run()", zap.String("err", err.Error()))
			}
		}(srv)
	}
}

func (s *service) startClients() {
	for _, clt := range s.cfg.Clients {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(clt client.Client) {
			defer s.wg.Done()
			if err := clt.Dial(); err != nil {
				s.logger.Error("Dial()", zap.String("err", err.Error()))
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
			s.logger.Info("signal received",
				zap.String("signal", sig.String()))
			s.stopServers()
			break outer
		default:
			s.logger.Info("pulse")
			time.Sleep(time.Second * 1)
			continue
		}
	}
}

func (s *service) stopServers() {
	for _, srv := range s.cfg.Servers {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(srv server.Server) {
			defer s.wg.Done()
			srv.Stop()
		}(srv)
	}
}

func middlewareService(s *service) router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			// set unique event id for every request
			e := uuid.New().String()
			// create new log instance with eventID to be added later to context
			l := s.logger.With(zap.String("event", e))
			l.Info("request",
				zap.String("method", r.Method),
				zap.String("url", r.URL.String()))
			// context
			ctx := r.Context()
			// ctx = context.WithValue(ctx, s.cfg.Name, s)
			// Note: service instance always available in context
			// under the general name > pluto
			ctx = context.WithValue(ctx, "pluto", s)
			ctx = context.WithValue(ctx, "logger", l)
			ctx = context.WithValue(ctx, "event", e)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

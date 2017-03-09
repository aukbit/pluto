package pluto

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	defaultName       = "pluto"
	defaultVersion    = "v1.0.0"
	defaultHealthAddr = ":9090"
)

// Service
type Service struct {
	cfg    Config
	close  chan bool
	wg     *sync.WaitGroup
	health *health.Server
	logger *zap.Logger
}

// New returns a new pluto service with Options passed in
func New(opts ...Option) *Service {
	return newService(opts...)
}

func newService(opts ...Option) *Service {
	s := &Service{
		cfg:    newConfig(),
		close:  make(chan bool),
		wg:     &sync.WaitGroup{},
		health: health.NewServer(),
	}
	s.logger, _ = zap.NewProduction()
	if len(opts) > 0 {
		s = s.WithOptions(opts...)
	}
	return s
}

// WithOptions clones the current Service, applies the supplied Options, and
// returns the resulting Service. It's safe to use concurrently.
func (s *Service) WithOptions(opts ...Option) *Service {
	c := s.clone()
	for _, opt := range opts {
		opt.apply(c)
	}
	return c
}

// clone creates a shallow copy service
func (s *Service) clone() *Service {
	copy := *s
	return &copy
}

// Run starts service
func (s *Service) Run() error {
	// set logger
	s.setLogger()
	// set health server
	s.setHealthServer()
	// register at service discovery
	if err := s.register(); err != nil {
		return err
	}
	// start service
	if err := s.start(); err != nil {
		return err
	}
	// hook run after start
	s.hookAfterStart()
	// wait for all go routines to finish
	s.wg.Wait()
	s.logger.Info("exit")
	return nil
}

// Stop stops service
func (s *Service) Stop() {
	s.logger.Info("stop")
	s.close <- true
}

// Config service configration options
func (s *Service) Config() Config {
	return s.cfg
}

// Server returns a server instance by name if initialized in service
func (s *Service) Server(name string) (srv server.Server, ok bool) {
	name = common.SafeName(name, server.DefaultName)
	if srv, ok = s.cfg.Servers[name]; !ok {
		return
	}
	return srv, true
}

// Client returns a client instance by name if initialized in service
func (s *Service) Client(name string) (clt client.Client, ok bool) {
	name = common.SafeName(name, client.DefaultName)
	if clt, ok = s.cfg.Clients[name]; !ok {
		return
	}
	return clt, true
}

func (s *Service) Health() *healthpb.HealthCheckResponse {
	hcr, err := s.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: s.cfg.ID})
	if err != nil {
		s.logger.Error("Health", zap.String("err", err.Error()))
	}
	return hcr
}

func (s *Service) setHealthServer() {
	s.health.SetServingStatus(s.cfg.ID, 1)
	// Define Router
	mux := router.NewMux()
	mux.GET("/_health/:module/:name", healthHandler)
	// Define server
	srv := server.NewServer(
		server.Name(s.cfg.Name+"_health"),
		server.Addr(s.cfg.HealthAddr),
		server.Mux(mux),
	)
	s.cfg.Servers[srv.Config().Name] = srv
}

func (s *Service) setLogger() {
	s.logger = s.logger.With(
		zap.String("type", "service"),
		zap.String("id", s.cfg.ID),
		zap.String("name", s.cfg.Name))
}

func (s *Service) start() error {
	s.logger.Info("start",
		zap.String("ip4", common.IPaddress()),
		zap.Int("servers", len(s.cfg.Servers)),
		zap.Int("clients", len(s.cfg.Clients)))

	// connect to db
	s.connectDB()
	// run servers
	s.startServers()
	// dial clients
	s.startClients()
	// add go routine to WaitGroup
	s.wg.Add(1)
	go s.waitUntilStopOrSig()
	return nil
}

func (s *Service) hookAfterStart() {
	hooks, ok := s.cfg.Hooks["after_start"]
	if !ok {
		return
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "pluto", s)
	ctx = context.WithValue(ctx, "logger", s.logger)
	for _, h := range hooks {
		h(ctx)
	}
}

func (s *Service) connectDB() {
	// connect datastore
	if _, ok := s.cfg.Datastore.(datastore.Datastore); ok {
		s.cfg.Datastore.Connect(datastore.Discovery(s.Config().Discovery))
		if err := s.cfg.Datastore.RefreshSession(); err != nil {
			s.logger.Error("RefreshSession()", zap.String("err", err.Error()))
		}
	}
}

func (s *Service) startServers() {
	for _, srv := range s.cfg.Servers {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(srv server.Server) {
			defer s.wg.Done()
			for {
				err := srv.Run(
					server.ParentID(s.cfg.ID),
					server.Middlewares(serviceContextMiddleware(s)),
					server.UnaryServerInterceptors(serviceContextUnaryServerInterceptor(s)),
					server.Discovery(s.Config().Discovery))
				if err == nil {
					return
				}
				s.logger.Error(fmt.Sprintf("Run failed on server: %v. Error: %v. On hold by 10s...", srv.Config().Name, err.Error()))
				time.Sleep(time.Second * 10)
				// delete(s.cfg.Servers, srv.Config().Name)
			}
		}(srv)
	}
}

func (s *Service) startClients() {
	for _, clt := range s.cfg.Clients {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(clt client.Client) {
			defer s.wg.Done()
			for {
				err := clt.Dial(
					client.ParentID(s.cfg.ID),
					client.Discovery(s.Config().Discovery))
				if err == nil {
					return
				}
				s.logger.Error(fmt.Sprintf("Dial failed on client: %v. Error: %v. On hold by 10s...", clt.Config().Name, err.Error()))
				time.Sleep(time.Second * 10)
			}
		}(clt)
	}
}

// waitUntilStopOrSig waits for close channel or syscall Signal
func (s *Service) waitUntilStopOrSig() {
	defer s.wg.Done()
	//  Stop also in case of any host signal
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGTERM, syscall.SIGINT)

outer:
	for {
		select {
		case <-s.close:
			// Waits for call to stop
			s.health.SetServingStatus(s.cfg.ID, 2)
			s.unregister()
			s.closeClients()
			s.stopServers()
			break outer
		case sig := <-sigch:
			// Waits for signal to stop
			s.logger.Info("signal received",
				zap.String("signal", sig.String()))
			s.health.SetServingStatus(s.cfg.ID, 2)
			s.unregister()
			s.closeClients()
			s.stopServers()
			break outer
		default:
			s.logger.Debug("pulse")
			time.Sleep(time.Second * 1)
			continue
		}
	}
}

func (s *Service) closeClients() {
	for _, clt := range s.cfg.Clients {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(clt client.Client) {
			defer s.wg.Done()
			clt.Close()
		}(clt)
	}
}

func (s *Service) stopServers() {
	for _, srv := range s.cfg.Servers {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(srv server.Server) {
			defer s.wg.Done()
			srv.Stop()
		}(srv)
	}
}

package pluto

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/aukbit/pluto/client"
	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/datastore"
	"github.com/aukbit/pluto/server"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Service
type service struct {
	cfg   *Config
	close chan bool
	wg    *sync.WaitGroup
	// logger *zap.Logger
	health *health.Server
}

func newService(cfgs ...ConfigFunc) *service {
	c := newConfig(cfgs...)
	s := &service{
		cfg:    c,
		close:  make(chan bool),
		wg:     &sync.WaitGroup{},
		health: health.NewServer(),
	}
	// s.logger, _ = zap.NewProduction()
	return s
}

// Run starts service
func (s *service) Run() error {
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
	// s.logger.Info("exit")
	return nil
}

// Stop stops service
func (s *service) Stop() {
	// s.logger.Info("stop")
	s.close <- true
}

// Config service configration options
func (s *service) Config() *Config {
	return s.cfg
}

// Server returns a server instance by name if initialized in service
func (s *service) Server(name string) (srv server.Server, ok bool) {
	name = common.SafeName(name, server.DefaultName)
	if srv, ok = s.cfg.Servers[name]; !ok {
		return
	}
	return srv, true
}

// Client returns a client instance by name if initialized in service
func (s *service) Client(name string) (clt client.Client, ok bool) {
	name = common.SafeName(name, client.DefaultName)
	if clt, ok = s.cfg.Clients[name]; !ok {
		return
	}
	return clt, true
}

func (s *service) Health() *healthpb.HealthCheckResponse {
	hcr, err := s.health.Check(
		context.Background(), &healthpb.HealthCheckRequest{Service: s.cfg.ID})
	if err != nil {
		// s.logger.Error("Health", zap.String("err", err.Error()))
	}
	return hcr
}

func (s *service) setLogger() {
	// s.logger = s.logger.With(
	// 	zap.String("type", "service"),
	// 	zap.String("id", s.cfg.ID),
	// 	zap.String("name", s.cfg.Name))
}

func (s *service) start() error {
	// s.logger.Info("start",
	// 	zap.String("ip4", common.IPaddress()),
	// 	zap.Int("servers", len(s.cfg.Servers)),
	// 	zap.Int("clients", len(s.cfg.Clients)))

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

func (s *service) hookAfterStart() {
	hooks, ok := s.cfg.Hooks["after_start"]
	if !ok {
		return
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "pluto", s)
	// ctx = context.WithValue(ctx, "logger", s.logger)
	for _, h := range hooks {
		h(ctx)
	}
}

func (s *service) connectDB() {
	// connect datastore
	if _, ok := s.cfg.Datastore.(datastore.Datastore); ok {
		s.cfg.Datastore.Connect(datastore.Discovery(s.Config().Discovery))
		if err := s.cfg.Datastore.RefreshSession(); err != nil {
			// s.logger.Error("RefreshSession()", zap.String("err", err.Error()))
		}
	}
}

func (s *service) startServers() {
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
				// s.logger.Error(fmt.Sprintf("Run failed on server: %v. Error: %v. On hold by 10s...", srv.Config().Name, err.Error()))
				time.Sleep(time.Second * 10)
				// delete(s.cfg.Servers, srv.Config().Name)
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
			for {
				err := clt.Dial(
					client.ParentID(s.cfg.ID),
					client.Discovery(s.Config().Discovery))
				if err == nil {
					return
				}
				// s.logger.Error(fmt.Sprintf("Dial failed on client: %v. Error: %v. On hold by 10s...", clt.Config().Name, err.Error()))
				time.Sleep(time.Second * 10)
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
			s.health.SetServingStatus(s.cfg.ID, 2)
			s.unregister()
			s.closeClients()
			s.stopServers()
			break outer
		case sig := <-sigch:
			// Waits for signal to stop
			// s.logger.Info("signal received",
			// 	zap.String("signal", sig.String()))
			log.Printf("signal %v", sig.String())
			s.health.SetServingStatus(s.cfg.ID, 2)
			s.unregister()
			s.closeClients()
			s.stopServers()
			break outer
		default:
			// s.logger.Debug("pulse")
			time.Sleep(time.Second * 1)
			continue
		}
	}
}

func (s *service) closeClients() {
	for _, clt := range s.cfg.Clients {
		// add go routine to WaitGroup
		s.wg.Add(1)
		go func(clt client.Client) {
			defer s.wg.Done()
			clt.Close()
		}(clt)
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

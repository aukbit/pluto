package pluto

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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
}

func newService(cfgs ...ConfigFunc) *service {
	c := newConfig(cfgs...)
	return &service{cfg: c}

}

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

func middlewareService(s *service) router.Middleware {
	return func(h router.Handler) router.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, s.cfg.Name, s)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func (s *service) Server(name string) (srv server.Server) {
	var ok bool
	if srv, ok = s.cfg.Servers[name]; !ok {
		return nil
	}
	return srv
}

func (s *service) Client(name string) (clt client.Client) {
	var ok bool
	if clt, ok = s.cfg.Clients[name]; !ok {
		return nil
	}
	return clt
}

func (s *service) Datastore() datastore.Datastore {
	return s.cfg.Datastore
}

func (s *service) Run() error {
	if err := s.start(); err != nil {
		return err
	}
	// parse address for host, port
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	sig := <-ch
	logger.Info("signal received",
		zap.String("service", s.cfg.Name),
		zap.String("id", s.cfg.ID),
		zap.String("signal", sig.String()))
	return s.Stop()
}

func (s *service) start() error {
	logger.Info("START",
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

	// run servers
	for _, srv := range s.cfg.Servers {
		go func(ss server.Server) {
			if err := ss.Run(); err != nil {
				logger.Error("Run()",
					zap.String("service", s.cfg.Name),
					zap.String("id", s.cfg.ID),
					zap.String("err", err.Error()),
				)
			}
		}(srv)
	}
	// dial clients
	for _, clt := range s.cfg.Clients {
		go func(cc client.Client) {
			if err := cc.Dial(); err != nil {
				logger.Error("Dial()",
					zap.String("service", s.cfg.Name),
					zap.String("id", s.cfg.ID),
					zap.String("err", err.Error()),
				)
			}
		}(clt)
	}

	return nil
}

func (s *service) Stop() error {
	s.close <- true
	logger.Info("STOP",
		zap.String("service", s.cfg.Name),
		zap.String("id", s.cfg.ID),
	)
	return nil
}

func (s *service) Config() *Config {
	cfg := s.cfg
	return cfg
}

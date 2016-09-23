package pluto

import (
	"log"
	"syscall"
	"os/signal"
	"os"
	"net/http"
	"context"
	"bitbucket.org/aukbit/pluto/server"
	"bitbucket.org/aukbit/pluto/client"
	"bitbucket.org/aukbit/pluto/datastore"
	"bitbucket.org/aukbit/pluto/server/router"
)


// Service
type service struct {
	cfg 			*Config
	close 			chan bool
}

func newService (cfgs ...ConfigFunc) Service {
	c := newConfig(cfgs...)
	return &service{cfg: c}

}

func (s *service) Init(cfgs ...ConfigFunc) error {
	for _, c := range cfgs {
		c(s.cfg)
	}

	for _, srv := range s.Servers(){
		srv.Init()
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

func (s *service) Servers() map[string]server.Server {
	return s.cfg.Servers
}

func (s *service) Datastore() datastore.Datastore {
	return s.cfg.Datastore
}

func (s *service) Clients() map[string]client.Client {
	return s.cfg.Clients
}

func (s *service) Run() error {
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

func (s *service) start() error {
	log.Printf("START %s \t%s", s.cfg.Name, s.cfg.Id)

	// connect datastore
	if s.cfg.Datastore != nil {
		s.cfg.Datastore.Connect()
		if err := s.cfg.Datastore.RefreshSession(); err != nil {
			log.Fatalf("ERROR s.cfg.Datastore.RefreshSession() %v", err.Error())
		}
	}

	// run servers
	for _, srv := range s.Servers(){
		go func(ss server.Server) {
			if err := ss.Run(); err != nil {
				log.Fatalf("ERROR srv.Run() %v", err)
			}
		}(srv)
	}
	// dial clients
	for _, clt := range s.Clients(){
		go func(cc client.Client) {
			_, err := cc.Dial()
			if err != nil {
				log.Fatalf("ERROR cc.Dial() %v", err)
			}
		}(clt)
	}

	return nil
}

func (s *service) Stop() error {
	s.close <-true
	return nil
}

func (s *service) Config() *Config {
	cfg := s.cfg
	return cfg
}
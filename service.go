package pluto

import (
	"pluto/server"
	"log"
	"syscall"
	"os/signal"
	"os"
	"pluto/client"
	"context"
	"pluto/server/router"
	"net/http"
)


// Service
type service struct {
	cfg 			*Config
	ctx			context.Context
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

	// register service in http server handlers context
	for _, srv := range s.Servers(){
		if srv.Config().Format == "http" {
			ctx := context.Background()
			s.ctx = context.WithValue(ctx, "token", "abc123")
			//handlers := srv.Config().Mux.Handlers()
			//// wrap service in context
			//for _, h := range handlers {
			//	h = WrapService(s, h)
			//}
			//srv.Config().Mux
			//router := srv.Config().Router
			srv.Config().Mux.WrapHandlers(s)
		}

	}




	return nil
}

func WrapService(s Service, next router.Handler) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "service", s)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *service) Servers() map[string]server.Server {
	return s.cfg.Servers
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
	return nil
}

func (s *service) start() error {
	log.Printf("START %s \t%s", s.cfg.Name, s.cfg.Id)
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
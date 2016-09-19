package pluto

import (
	"pluto/server"
	"log"
	"syscall"
	"os/signal"
	"os"
	"pluto/client"
	"pluto/datastore"
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

	// Wrap this service to all handlers
	// make it available in handler context
	for _, srv := range s.Servers(){
		if srv.Config().Format == "http" {
			srv.Config().Mux.AddContextWith(s.cfg.Name, s)
		}
	}
	return nil
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
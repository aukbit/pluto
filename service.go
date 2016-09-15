package pluto

import (
	"pluto/server"
	"log"
	"syscall"
	"os/signal"
	"os"
	"fmt"
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
	return nil
}

func (s *service) Server() server.Server {
	return s.cfg.Server
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
	log.Printf("START %s %s", s.cfg.Name, s.cfg.Id)
	//
	go func() {
		if err := s.cfg.Server.Run(); err != nil {
			log.Fatal(fmt.Sprintf("ERROR s.cfg.Server.Run() %v", err))
		}
	}()

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
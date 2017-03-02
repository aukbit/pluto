package pluto

import (
	"fmt"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
)

// register Pluto Service within the service discovery system
func (s *service) register() error {
	if _, ok := s.cfg.Discovery.(discovery.Discovery); ok {
		// define service
		dse := discovery.Service{
			Name: s.cfg.Name,
			Tags: []string{s.cfg.ID, s.cfg.Version},
		}
		// define check
		dck := discovery.Check{
			Name:  fmt.Sprintf("Service '%s' check", s.cfg.Name),
			Notes: fmt.Sprintf("Ensure the Pluto service %s is running", s.cfg.ID),
			DeregisterCriticalServiceAfter: "10m",
			HTTP:      fmt.Sprintf("http://%s:9090/_health/pluto/%s", common.IPaddress(), s.cfg.Name),
			Interval:  "10s",
			Timeout:   "1s",
			ServiceID: s.cfg.Name,
		}
		if err := s.cfg.Discovery.Register(discovery.ServicesCfg(dse), discovery.ChecksCfg(dck)); err != nil {
			return err
		}
	}
	return nil
}

// unregister Server from the service discovery system
func (s *service) unregister() error {
	if _, ok := s.cfg.Discovery.(discovery.Discovery); ok {
		if err := s.cfg.Discovery.Unregister(); err != nil {
			return err
		}
	}
	return nil
}

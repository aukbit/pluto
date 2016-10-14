package pluto

import (
	"fmt"

	"bitbucket.org/aukbit/pluto/common"
	"bitbucket.org/aukbit/pluto/discovery"
)

// register Pluto Service within the service discovery system
func (s *service) register() error {
	if s.cfg.Discovery != nil {
		// define service
		dse := &discovery.Service{
			Name: s.cfg.Name,
			Tags: []string{s.cfg.ID, s.cfg.Version},
		}
		// define check
		dck := &discovery.Check{
			Name:  fmt.Sprintf("Service '%s' check", s.cfg.Name),
			Notes: fmt.Sprintf("Ensure the Pluto service %s is running", s.cfg.ID),
			DeregisterCriticalServiceAfter: "10m",
			HTTP:      fmt.Sprintf("http://%s:9090/_health/pluto/%s", common.IPaddress(), s.cfg.Name),
			Interval:  "10s",
			Timeout:   "1s",
			ServiceID: s.cfg.Name,
		}
		if err := s.cfg.Discovery.Register(discovery.Services(dse), discovery.Checks(dck)); err != nil {
			return err
		}
	}
	return nil

	// _, err := discovery.IsAvailable()
	// if err != nil {
	// 	s.logger.Error("service discovery not available")
	// 	return nil
	// }
	// ds := &discovery.Service{
	// 	Name: s.cfg.Name,
	// 	Tags: []string{s.cfg.Version, s.cfg.ID},
	// }
	// err = discovery.RegisterService(ds)
	// if err != nil {
	// 	return err
	// }
	// c := &discovery.Check{
	// 	Name:  fmt.Sprintf("Service '%s' check", s.cfg.Name),
	// 	Notes: fmt.Sprintf("Ensure the Pluto service %s is running", s.cfg.ID),
	// 	DeregisterCriticalServiceAfter: "10m",
	// 	HTTP:      fmt.Sprintf("http://%s:9090/_health/pluto/%s", common.IPaddress(), s.cfg.Name),
	// 	Interval:  "10s",
	// 	Timeout:   "1s",
	// 	ServiceID: s.cfg.Name,
	// }
	// err = discovery.RegisterCheck(c)
	// if err != nil {
	// 	return err
	// }
	// s.isDiscovered = true
	// return nil
}

// unregister Server from the service discovery system
func (s *service) unregister() error {
	if _, ok := s.cfg.Discovery.(discovery.Discovery); ok {
		if err := s.cfg.Discovery.Unregister(); err != nil {
			return err
		}
	}
	return nil
	// if s.isDiscovered {
	// 	err := discovery.DeregisterService(s.cfg.Name)
	// 	if err != nil {
	// 		s.logger.Error(err.Error())
	// 	}
	// }
}

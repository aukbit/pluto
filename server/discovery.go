package server

import (
	"fmt"

	"bitbucket.org/aukbit/pluto/common"
	"bitbucket.org/aukbit/pluto/discovery"
)

// register Server within the service discovery system
func (ds *defaultServer) register() error {
	if ds.cfg.Discovery != nil {
		// define service
		dse := &discovery.Service{
			Name: ds.cfg.Name,
			Port: ds.cfg.Port(),
			Tags: []string{ds.cfg.ID, ds.cfg.Version},
		}
		// define check
		dck := &discovery.Check{
			Name:  fmt.Sprintf("Service '%s' check", ds.cfg.Name),
			Notes: fmt.Sprintf("Ensure the server is listening on port %s", ds.cfg.Addr),
			DeregisterCriticalServiceAfter: "10m",
			HTTP:      fmt.Sprintf("http://%s:9090/_health/server/%s", common.IPaddress(), ds.cfg.Name),
			Interval:  "10s",
			Timeout:   "1s",
			ServiceID: ds.cfg.Name,
		}
		if err := ds.cfg.Discovery.Register(discovery.Services(dse), discovery.Checks(dck)); err != nil {
			return err
		}
	}
	return nil
	// _, err := discovery.IsAvailable()
	// if err != nil {
	// 	ds.logger.Error("service discovery not available")
	// 	return nil
	// }
	// s := &discovery.Service{
	// 	Name: ds.cfg.Name,
	// 	Port: ds.cfg.Port(),
	// 	Tags: []string{ds.cfg.Version, ds.cfg.ID},
	// }
	// err = discovery.RegisterService(s)
	// if err != nil {
	// 	return err
	// }
	// c := &discovery.Check{
	// 	Name:  fmt.Sprintf("Service '%s' check", ds.cfg.Name),
	// 	Notes: fmt.Sprintf("Ensure the server is listening on port %s", ds.cfg.Addr),
	// 	DeregisterCriticalServiceAfter: "10m",
	// 	HTTP:      fmt.Sprintf("http://%s:9090/_health/server/%s", common.IPaddress(), ds.cfg.Name),
	// 	Interval:  "10s",
	// 	Timeout:   "1s",
	// 	ServiceID: ds.cfg.Name,
	// }
	// err = discovery.RegisterCheck(c)
	// if err != nil {
	// 	return err
	// }
	// ds.isDiscovered = true
	// return nil
}

// unregister Server from the service discovery system
func (ds *defaultServer) unregister() error {
	if _, ok := ds.cfg.Discovery.(discovery.Discovery); ok {
		if err := ds.cfg.Discovery.Unregister(); err != nil {
			return err
		}
	}
	return nil
	// if ds.isDiscovered {
	// 	err := discovery.DeregisterService(ds.cfg.Name)
	// 	if err != nil {
	// 		ds.logger.Error(err.Error())
	// 	}
	// }
}

package server

import (
	"fmt"

	"bitbucket.org/aukbit/pluto/common"
	"bitbucket.org/aukbit/pluto/discovery"
)

// register Server within the service discovery system
func (ds *defaultServer) register() error {
	if _, ok := ds.cfg.Discovery.(discovery.Discovery); ok {
		// define service
		dse := discovery.Service{
			Name:    ds.cfg.Name,
			Address: common.IPaddress(),
			Port:    ds.cfg.Port(),
			Tags:    []string{ds.cfg.ID, ds.cfg.Version},
		}
		// define check
		dck := discovery.Check{
			Name:  fmt.Sprintf("Service '%s' check", ds.cfg.Name),
			Notes: fmt.Sprintf("Ensure the server is listening on port %s", ds.cfg.Addr),
			DeregisterCriticalServiceAfter: "10m",
			HTTP:      fmt.Sprintf("http://%s:9090/_health/server/%s", common.IPaddress(), ds.cfg.Name),
			Interval:  "30s",
			Timeout:   "1s",
			ServiceID: ds.cfg.Name,
		}
		if err := ds.cfg.Discovery.Register(discovery.ServicesCfg(dse), discovery.ChecksCfg(dck)); err != nil {
			return err
		}
	}
	return nil
}

// unregister Server from the service discovery system
func (ds *defaultServer) unregister() error {
	if _, ok := ds.cfg.Discovery.(discovery.Discovery); ok {
		if err := ds.cfg.Discovery.Unregister(); err != nil {
			return err
		}
	}
	return nil
}

package client

import (
	"fmt"

	"bitbucket.org/aukbit/pluto/common"
	"bitbucket.org/aukbit/pluto/discovery"
)

// target get target IP:Port from service discovery system
func (dc *defaultClient) targets() error {
	if _, ok := dc.cfg.Discovery.(discovery.Discovery); ok {
		t, err := dc.cfg.Discovery.Service(dc.cfg.TargetName)
		if err != nil {
			return err
		}
		dc.cfg.Targets = t
	}
	return nil
}

// register Client within the service discovery system
func (dc *defaultClient) register() error {
	if _, ok := dc.cfg.Discovery.(discovery.Discovery); ok {
		// define service
		dse := discovery.Service{
			Name: dc.cfg.Name,
			Tags: []string{dc.cfg.ID, dc.cfg.Version},
		}
		// define check
		dck := discovery.Check{
			Name:  fmt.Sprintf("Service '%s' check", dc.cfg.Name),
			Notes: fmt.Sprintf("Ensure the client is able to connect to service %s", dc.cfg.TargetName),
			DeregisterCriticalServiceAfter: "10m",
			HTTP:      fmt.Sprintf("http://%s:9090/_health/client/%s", common.IPaddress(), dc.cfg.Name),
			Interval:  "10s",
			Timeout:   "1s",
			ServiceID: dc.cfg.Name,
		}
		if err := dc.cfg.Discovery.Register(discovery.ServicesCfg(dse), discovery.ChecksCfg(dck)); err != nil {
			return err
		}
	}
	return nil
}

// unregister Server from the service discovery system
func (dc *defaultClient) unregister() error {
	if _, ok := dc.cfg.Discovery.(discovery.Discovery); ok {
		if err := dc.cfg.Discovery.Unregister(); err != nil {
			return err
		}
	}
	return nil
}

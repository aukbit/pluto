package client

import (
	"fmt"

	"bitbucket.org/aukbit/pluto/common"
	"bitbucket.org/aukbit/pluto/discovery"
)

// target get target IP:Port from service discovery system
func (dc *defaultClient) target() error {
	if _, ok := dc.cfg.Discovery.(discovery.Discovery); ok {
		addr, err := dc.cfg.Discovery.Service(dc.cfg.TargetName)
		if err != nil {
			return err
		}
		dc.cfg.Target = addr
	}
	return nil
}

// register Client within the service discovery system
func (dc *defaultClient) register() error {
	if _, ok := dc.cfg.Discovery.(discovery.Discovery); ok {
		// define service
		dse := &discovery.Service{
			Name: dc.cfg.Name,
			Tags: []string{dc.cfg.ID, dc.cfg.Version},
		}
		// define check
		dck := &discovery.Check{
			Name:  fmt.Sprintf("Service '%s' check", dc.cfg.Name),
			Notes: fmt.Sprintf("Ensure the client is able to connect to %s - %s", dc.cfg.Target, dc.cfg.TargetName),
			DeregisterCriticalServiceAfter: "10m",
			HTTP:      fmt.Sprintf("http://%s:9090/_health/client/%s", common.IPaddress(), dc.cfg.Name),
			Interval:  "10s",
			Timeout:   "1s",
			ServiceID: dc.cfg.Name,
		}
		if err := dc.cfg.Discovery.Register(discovery.Services(dse), discovery.Checks(dck)); err != nil {
			return err
		}
	}
	return nil

	// _, err := discovery.IsAvailable()
	// if err != nil {
	// 	dc.logger.Error("service discovery not available")
	// 	return nil
	// }
	// s := &discovery.Service{
	// 	Name: dc.cfg.Name,
	// 	Tags: []string{dc.cfg.Version, dc.cfg.ID},
	// }
	// err = discovery.RegisterService(s)
	// if err != nil {
	// 	return err
	// }
	// c := &discovery.Check{
	// 	Name:  fmt.Sprintf("Service '%s' check", dc.cfg.Name),
	// 	Notes: fmt.Sprintf("Ensure the client is able to connect to %s - %s", dc.cfg.Target, dc.cfg.TargetDiscovery),
	// 	DeregisterCriticalServiceAfter: "10m",
	// 	HTTP:      fmt.Sprintf("http://%s:9090/_health/client/%s", common.IPaddress(), dc.cfg.Name),
	// 	Interval:  "10s",
	// 	Timeout:   "1s",
	// 	ServiceID: dc.cfg.Name,
	// }
	// err = discovery.RegisterCheck(c)
	// if err != nil {
	// 	return err
	// }
	// dc.isDiscovered = true
	// return nil
}

// unregister Server from the service discovery system
func (dc *defaultClient) unregister() error {
	if _, ok := dc.cfg.Discovery.(discovery.Discovery); ok {
		if err := dc.cfg.Discovery.Unregister(); err != nil {
			return err
		}
	}
	return nil
	// if dc.isDiscovered {
	// 	err := discovery.DeregisterService(dc.cfg.Name)
	// 	if err != nil {
	// 		dc.logger.Error(err.Error())
	// 	}
	// }
}

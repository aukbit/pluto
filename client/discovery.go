package client

import (
	"fmt"

	"bitbucket.org/aukbit/pluto/common"
	"bitbucket.org/aukbit/pluto/discovery"
)

// register Client within the service discovery system
func (dc *defaultClient) register() error {
	_, err := discovery.IsAvailable()
	if err != nil {
		dc.logger.Error("service discovery not available")
		return nil
	}
	s := &discovery.Service{
		Name: dc.cfg.Name,
		Tags: []string{dc.cfg.Version, dc.cfg.ID},
	}
	err = discovery.RegisterService(s)
	if err != nil {
		return err
	}
	c := &discovery.Check{
		Name:  fmt.Sprintf("Service '%s' check", dc.cfg.Name),
		Notes: fmt.Sprintf("Ensure the client is able to connect to %s - %s", dc.cfg.Target, dc.cfg.TargetDiscovery),
		DeregisterCriticalServiceAfter: "10m",
		HTTP:      fmt.Sprintf("http://%s:9090/_health/client/%s", common.IPaddress(), dc.cfg.Name),
		Interval:  "10s",
		Timeout:   "1s",
		ServiceID: dc.cfg.Name,
	}
	err = discovery.RegisterCheck(c)
	if err != nil {
		return err
	}
	dc.isDiscovered = true
	return nil
}

// unregister Server from the service discovery system
func (dc *defaultClient) unregister() {
	if dc.isDiscovered {
		err := discovery.DeregisterService(dc.cfg.Name)
		if err != nil {
			dc.logger.Error(err.Error())
		}
	}
}

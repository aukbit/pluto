package datastore

import (
	"fmt"

	"github.com/aukbit/pluto/common"
	"github.com/aukbit/pluto/discovery"
)

// target get target IP:Port from service discovery system
func (ds *Datastore) target() error {
	if _, ok := ds.cfg.Discovery.(discovery.Discovery); ok {
		targets, err := ds.cfg.Discovery.Service(ds.cfg.TargetName)
		if err != nil {
			return err
		}
		ds.cfg.Target = targets[0]
	}
	return nil
}

// register datastore client within the service discovery system
func (ds *Datastore) register() error {
	if _, ok := ds.cfg.Discovery.(discovery.Discovery); ok {
		// define service
		dse := discovery.Service{
			Name: ds.cfg.Name,
			Tags: []string{ds.cfg.ID, ds.cfg.Version},
		}
		// define check
		dck := discovery.Check{
			Name:  fmt.Sprintf("Service '%s' check", ds.cfg.Name),
			Notes: fmt.Sprintf("Ensure the bd client is able to connect to %s - %s", ds.cfg.Target, ds.cfg.TargetName),
			DeregisterCriticalServiceAfter: "10m",
			HTTP:      fmt.Sprintf("http://%s:9090/_health/db/%s", common.IPaddress(), ds.cfg.Name),
			Interval:  "10s",
			Timeout:   "1s",
			ServiceID: ds.cfg.Name,
		}
		if err := ds.cfg.Discovery.Register(discovery.ServicesCfg(dse), discovery.ChecksCfg(dck)); err != nil {
			return err
		}
	}
	return nil
}

// unregister datastore client from the service discovery system
func (ds *Datastore) unregister() error {
	if _, ok := ds.cfg.Discovery.(discovery.Discovery); ok {
		if err := ds.cfg.Discovery.Unregister(); err != nil {
			return err
		}
	}
	return nil
}

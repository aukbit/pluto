package client

import "bitbucket.org/aukbit/pluto/discovery"

// target get target IP:Port from service discovery system
func (ds *datastore) target() error {
	if _, ok := ds.cfg.Discovery.(discovery.Discovery); ok {
		addr, err := ds.cfg.Discovery.Service(ds.cfg.TargetName)
		if err != nil {
			return err
		}
		ds.cfg.Target = addr
	}
	return nil
}

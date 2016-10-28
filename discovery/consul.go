package discovery

import "github.com/uber-go/zap"

type consulDefault struct {
	cfg          *Config
	logger       zap.Logger
	isDiscovered bool
}

func newConsulDefault(cfgs ...ConfigFunc) *consulDefault {
	c := newConfig(cfgs...)
	return &consulDefault{cfg: c, logger: zap.New(zap.NewJSONEncoder())}
}

// IsAvailable
func (cd *consulDefault) IsAvailable() (bool, error) {
	return isAvailable(cd.cfg.Addr)
}

// Service
func (cd *consulDefault) Service(serviceID string) ([]string, error) {
	return GetServiceTargets(cd.cfg.Addr, serviceID)
}

func (cd *consulDefault) Register(cfgs ...ConfigFunc) error {
	// if _, err := isAvailable(cd.cfg.URL()); err != nil {
	// 	cd.logger.Error("service discovery not available")
	// 	return err
	// }
	cd.isDiscovered = true
	// set last configs
	for _, c := range cfgs {
		c(cd.cfg)
	}
	// register services
	for _, s := range cd.cfg.Services {
		err := DoServiceRegister(&DefaultServiceRegister{}, cd.cfg.Addr, &s)
		if err != nil {
			return err
		}
	}
	// register checks
	for _, c := range cd.cfg.Checks {
		err := DoCheckRegister(&DefaultCheckRegister{}, cd.cfg.Addr, &c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cd *consulDefault) Unregister() error {
	if cd.isDiscovered {
		// unregister services
		for _, s := range cd.cfg.Services {
			if err := DoServiceUnregister(&DefaultServiceRegister{}, cd.cfg.Addr, s.ID); err != nil {
				return err
			}
		}
		// register checks
		for _, c := range cd.cfg.Checks {
			if err := DoCheckUnregister(&DefaultCheckRegister{}, cd.cfg.Addr, c.ID); err != nil {
				return err
			}
		}
		cd.isDiscovered = false
	}
	return nil
}

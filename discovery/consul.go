package discovery

import "go.uber.org/zap"

type consulDefault struct {
	cfg    *Config
	logger *zap.Logger
}

func newConsulDefault(cfgs ...ConfigFunc) *consulDefault {
	c := newConfig(cfgs...)
	d := &consulDefault{cfg: c}
	d.logger, _ = zap.NewProduction()
	return d
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
	// unregister services
	for _, s := range cd.cfg.Services {
		if err := DoServiceUnregister(&DefaultServiceRegister{}, cd.cfg.Addr, s.ID); err != nil {
			return err
		}
	}
	// unregister checks
	for _, c := range cd.cfg.Checks {
		if err := DoCheckUnregister(&DefaultCheckRegister{}, cd.cfg.Addr, c.ID); err != nil {
			return err
		}
	}
	return nil
}

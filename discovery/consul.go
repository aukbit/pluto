package discovery

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/uber-go/zap"
)

const (
	SELF = "/v1/agent/self" // Returns the local node configuration
)

type consulDefault struct {
	cfg          *Config
	logger       zap.Logger
	isDiscovered bool
}

// newServer will instantiate a new defaultServer with the given config
func newConsulDefault(cfgs ...ConfigFunc) *consulDefault {
	c := newConfig(cfgs...)
	return &consulDefault{cfg: c, logger: zap.New(zap.NewJSONEncoder())}
}

// IsAvailable
func (cd *consulDefault) IsAvailable() (bool, error) {
	return isAvailable(cd.cfg.URL())
}

func (cd *consulDefault) Register(cfgs ...ConfigFunc) error {
	_, err := isAvailable(cd.cfg.URL())
	if err != nil {
		cd.logger.Error("service discovery not available")
		return nil
	}
	cd.isDiscovered = true
	// set last configs
	for _, c := range cfgs {
		c(cd.cfg)
	}
	// register services
	for _, s := range cd.cfg.Services {
		err = registerService(cd.cfg.URL(), s)
		if err != nil {
			return err
		}
	}
	// register checks
	for _, c := range cd.cfg.Checks {
		err = registerCheck(cd.cfg.URL(), c)
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
			if err := deregisterService(cd.cfg.URL(), s.ID); err != nil {
				return err
			}
		}
		// register checks
		for _, c := range cd.cfg.Checks {
			if err := deregisterCheck(cd.cfg.URL(), c.ID); err != nil {
				return err
			}
		}
		cd.isDiscovered = false
	}
	return nil
}

func isAvailable(url string) (bool, error) {
	resp, err := http.Get(url + SELF)
	if err != nil {
		return false, err
	}
	_, err = io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		return false, err
	}
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Consul not available at this node")
	}
	return true, nil
}

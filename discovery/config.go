package discovery

import "fmt"

// Config pluto service config
type Config struct {
	Addr string
	// Services []*Service
	Services Services
	// Checks   []*Check
	Checks Checks
}

// ConfigFunc registers the Config
type ConfigFunc func(*Config)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{Addr: "localhost:8500", Services: Services{}, Checks: Checks{}}

	for _, c := range cfgs {
		c(cfg)
	}

	return cfg
}

// Addr service discovery addr
func Addr(a string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Addr = a
	}
}

// ServicesCfg ..
func ServicesCfg(ss ...Service) ConfigFunc {
	return func(cfg *Config) {
		// cfg.Services = append(cfg.Services, s...)
		for _, s := range ss {
			cfg.Services[s.ID] = s
		}
	}
}

// ChecksCfg ..
func ChecksCfg(cc ...Check) ConfigFunc {
	return func(cfg *Config) {
		// cfg.Checks = append(cfg.Checks, c...)
		for _, c := range cc {
			cfg.Checks[c.ID] = c
		}
	}
}

// URL service discovery url
func (c *Config) URL() string {
	return fmt.Sprintf("http://%s", c.Addr)
}

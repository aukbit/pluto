package discovery

import "fmt"

// Config pluto service config
type Config struct {
	Addr     string
	Services []*Service
	Checks   []*Check
}

// ConfigFunc registers the Config
type ConfigFunc func(*Config)

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := &Config{Addr: "localhost:8500"}

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

func Services(s ...*Service) ConfigFunc {
	return func(cfg *Config) {
		cfg.Services = append(cfg.Services, s...)
	}
}

func Checks(c ...*Check) ConfigFunc {
	return func(cfg *Config) {
		cfg.Checks = append(cfg.Checks, c...)
	}
}

func (c *Config) URL() string {
	return fmt.Sprintf("http://%s", c.Addr)
}

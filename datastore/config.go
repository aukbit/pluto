package datastore


type Config struct {
	Keyspace		string
	Addr			string
}

type ConfigFunc func(*Config)

var DefaultConfig = Config{
	Keyspace:		"default",
	Addr:			"127.0.0.1",
}

func newConfig(cfgs ...ConfigFunc) *Config {

	cfg := DefaultConfig

	for _, c := range cfgs {
		c(&cfg)
	}

	return &cfg
}

// Keyspace db keyspace
func Keyspace(ks string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Keyspace = ks
	}
}

// Addr db address
func Addr(a string) ConfigFunc {
	return func(cfg *Config) {
		cfg.Addr = a
	}
}


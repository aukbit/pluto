package discovery

type Discovery interface {
	IsAvailable() (bool, error)
	Register(...ConfigFunc) error
	Unregister() error
}

// NewDiscovery
func NewDiscovery(cfgs ...ConfigFunc) Discovery {
	return newConsulDefault(cfgs...)
}

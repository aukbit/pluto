package discovery

// Discovery ...
type Discovery interface {
	IsAvailable() (bool, error)
	Service(string) (string, error)
	Register(...ConfigFunc) error
	Unregister() error
}

// NewDiscovery ...
func NewDiscovery(cfgs ...ConfigFunc) Discovery {
	return newConsulDefault(cfgs...)
}

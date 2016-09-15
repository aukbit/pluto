package client

// A Client defines parameters for making calls to an HTTP server.
// The zero value for Client is a valid configuration.
type defaultClient struct {
	cfg 			*Config
	close 			chan bool
}

// newDefaultClient will instantiate a new Client with the given config
func newDefaultClient(cfgs ...ConfigFunc) Client {
	c := newConfig(cfgs...)
	return &defaultClient{cfg: c, close: make(chan bool)}
}

func (s *defaultClient) Init(cfgs ...ConfigFunc) error {
	for _, c := range cfgs {
		c(s.cfg)
	}
	return nil
}

func (s *defaultClient) Config() *Config {
	cfg := s.cfg
	return cfg
}

func (s *defaultClient) Run() error {
	// TODO
	return nil
}

func (s *defaultClient) Stop() error {
	s.close <-true
	return nil
}
package client

// Client is an interface to make calls to services
type Client interface {
	Dial() error
	Call() interface{}
	Close() error
	Config() *Config
}

var (
	DefaultName    = "client"
	DefaultVersion = "1.0.0"
)

// NewClient returns a new client with cfg passed in
func NewClient(cfgs ...ConfigFunc) Client {
	return newClient(cfgs...)
}

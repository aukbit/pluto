package server

// Server is the basic interface that defines what to expect from any server.
type Server interface {
	Run(...ConfigFunc) error
	Stop()
	Config() *Config
}

var (
	defaultName    = "server"
	defaultVersion = "1.0.0"
)

// NewServer returns a new http server with cfg passed in
func NewServer(cfgs ...ConfigFunc) Server {
	return newServer(cfgs...)
}

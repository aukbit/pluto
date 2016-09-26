package server

// Server is the basic interface that defines what to expect from any server.
type Server interface {
	Init(...ConfigFunc)		error
	Run() 				error
	Stop() 				error
	Config() 			*Config
}

var (
	DefaultName			= "server"
	DefaultVersion      		= "1.0.0"
	DefaultServer  			= newDefaultServer()
)

// NewServer returns a new http server with cfg passed in
func NewServer(cfgs ...ConfigFunc) Server {
	return newDefaultServer(cfgs...)
}

// NewGRPCServer returns a new grpc server with cfg passed in
func NewGRPCServer(cfgs ...ConfigFunc) Server {
	return newGRPCServer(cfgs...)
}
package client


// Client is an interface to make calls to services
type Client interface {
	Init(...ConfigFunc)		error
	Dial() 				(interface{}, error)
	Close() 			error
	Config() 			*Config
}

var (
	DefaultName			= "client"
	DefaultVersion      		= "1.0.0"
	DefaultClient  			= newGRPCClient()
)

// NewClient returns a new client with cfg passed in
func NewClient(cfgs ...ConfigFunc) Client {
	return newGRPCClient(cfgs...)
}
package client


// Client is an interface to make calls to services
type Client interface {
	Init(...ConfigFunc)						error
	Run() 									error
	Stop() 									error
	Config() 								*Config
}

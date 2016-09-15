package client

// A Server defines parameters for running an HTTP server.
// The zero value for Server is a valid configuration.
type defaultClient struct {
	cfg 			*Config
	close 			chan bool
}

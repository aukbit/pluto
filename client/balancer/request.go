package balancer

// Request requests a connector to make a client call
type Request struct {
	connsCh chan *Connector
}

package balancer

// Request requests a connector to make a client call
type Request struct {
	connsCh chan *Connector
}

func NewRequest(ch chan *Connector) Request {
	return Request{ch}
}

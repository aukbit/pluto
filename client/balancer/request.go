package balancer

// Request requests a connector to make a client call
type Request struct {
	connsCh ConnsCh
}

func NewRequest(ch ConnsCh) Request {
	return Request{ch}
}

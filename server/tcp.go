package server

import (
	"net"
	"time"
)

// TcpKeepAliveListener is a copy of tcpKeepAliveListener
// source: net/http/server.go
//
// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type TcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln TcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// getNewAddr lets the system discover a new port available
func getNewAddr() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}


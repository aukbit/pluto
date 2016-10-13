package common

import (
	"log"
	"net"
)

// IPaddress returns first IP address
// TODO filter by network name...
func IPaddress() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("IPaddress %v", err.Error())
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
				// os.Stdout.WriteString(ipnet.IP.String() + "\n")
			}
		}
	}
	return "localhost"
}

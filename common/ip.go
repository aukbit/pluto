package common

import (
	"fmt"
	"log"
	"net"
	"os"
)

// IPaddress ...
func IPaddress() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}

	for _, a := range addrs {
		log.Printf("TESTE %v\n", a)
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				os.Stdout.WriteString(ipnet.IP.String() + "\n")
			}
		}
	}
}

func IP2() {
	list, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iface := range list {
		fmt.Printf("name=%s %v\n", iface.Name, iface)
		addrs, err := iface.Addrs()
		if err != nil {
			panic(err)
		}
		for j, addr := range addrs {
			fmt.Printf(" %d %v\n", j, addr)
		}
	}
}

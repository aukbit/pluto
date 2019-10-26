package main

import (
	"log"

	"github.com/aukbit/pluto/v6/examples/user/frontend/service"
)

func main() {
	// run frontend service
	if err := frontend.Run(); err != nil {
		log.Fatal(err)
	}
}

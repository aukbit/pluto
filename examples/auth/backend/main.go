package main

import (
	"log"

	"github.com/aukbit/pluto/examples/user/backend/service"
)

func main() {

	// run backend service
	if err := backend.Run(); err != nil {
		log.Fatal(err)
	}
}

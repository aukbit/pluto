package main

import (
	"flag"
	"log"

	"bitbucket.org/aukbit/pluto/examples/user/backend/service"
)

func main() {
	flag.Parse()

	// run backend service
	if err := backend.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"github.com/aukbit/pluto/examples/https/web"
)

func main() {
	// run frontend service
	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}

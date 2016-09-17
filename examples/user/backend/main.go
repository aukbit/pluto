package main

import (
	"pluto/examples/user/backend/service"
)
import "log"

func main(){

	// run backend service
	if err := backend.Run(); err != nil {
		log.Fatal(err)
	}
}



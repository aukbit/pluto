package main

import (
	"log"
	"pluto/examples/user/frontend/service"
)

func main(){
	// run frontend service
	if err := frontend.Run(); err != nil {
		log.Fatal(err)
	}
}
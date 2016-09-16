package main

import (
	"pluto"
	"pluto/server/router"
	"pluto/server"
	"log"
	"pluto/examples/user/frontend/views"
)

func main(){

	// 1. Config service
	s := pluto.NewService(
		pluto.Name("user-frontend-api"),
		pluto.Description("user-frontend-api is responsible to parse all json data to regarding users to internal services"),
	)

	// 2. Set server handlers
	mux := router.NewRouter()
	mux.GET("/user", frontend.GetHandler)
	mux.POST("/user", frontend.PostHandler)
	mux.GET("/user/:id", frontend.GetHandlerDetail)
	mux.PUT("/user/:id", frontend.PutHandler)
	mux.DELETE("/user/:id", frontend.DeleteHandler)
	// 3. Define server Router
	s.Server().Init(server.Router(mux))

	// 4. Run service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}

}



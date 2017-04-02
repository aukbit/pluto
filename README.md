**Note**: This project is still a work in progress, no stable version has yet been released.

# Pluto [![Circle CI](https://circleci.com/gh/aukbit/pluto.svg?style=svg)](https://circleci.com/gh/aukbit/pluto)
An implementation of microservices with golang to tackle some of the challenges in distributed systems.
### Features
- Currently supports a multiplexer HTTP router with dynamic paths and still compatible with the standard net/http library.
- Client/Server implementation with [gRPC](http://www.grpc.io/) for communication between services.
- Service health check.
- Datastore package currently supports [Cassandra](http://cassandra.apache.org/) via [gocql](https://github.com/gocql/gocql) and [MongoDB](https://www.mongodb.com/) via [mgo](https://labix.org/mgo).
- Structured Logs by using [zap](https://github.com/uber-go/zap).


### Inspiration
Projects that had influence in Pluto design and helped to solve technical barriers.
- [go-kit](https://github.com/go-kit/kit)
- [go-micro](https://github.com/myodc/go-micro)
- [gorilla](https://github.com/gorilla/mux)

Books
- [Building Microservices](http://shop.oreilly.com/product/0636920033158.do)
- [Microservice Architecture](http://shop.oreilly.com/product/0636920050308.do)

Articles
- [nginx - Introduction to Microservices](https://www.nginx.com/blog/introduction-to-microservices/?utm_source=event-driven-data-management-microservices&utm_medium=blog&utm_campaign=Microservices)
- [Fault Tolerance in a High Volume, Distributed System](http://techblog.netflix.com/2012/02/fault-tolerance-in-high-volume.html)

### Examples
- [User](https://github.com/aukbit/pluto/tree/master/examples/user)
- [Authentication](https://github.com/aukbit/pluto/tree/master/examples/auth)
- [Distributed deployment](https://github.com/aukbit/pluto/tree/master/examples/dist)
- [HTTPS/TLS](https://github.com/aukbit/pluto/tree/master/examples/https)

### Pluto - Hello World

```go
package main

import (
	"log"
	"net/http"

	"github.com/aukbit/pluto"
	"github.com/aukbit/pluto/reply"
	"github.com/aukbit/pluto/server"
	"github.com/aukbit/pluto/server/router"
)

func main() {
	// Define router
	mux := router.New()
	mux.GET("/", func(w http.ResponseWriter, r *http.Request) {
		reply.Json(w, r, http.StatusOK, "Hello World")
	})

	// Define http server
	srv := server.New(
		server.Mux(mux),
	)

	// Define Pluto service
	s := pluto.New(
		pluto.Servers(srv),
	)

	// Run Pluto service
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}

```

```sh
go run ./examples/hello/main.go
{"level":"info","ts":1489193908.1746852,"caller":"/github.com/aukbit/pluto/service.go:155","msg":"start","id":"plt_5QPRA9","name":"pluto","ip4":"192.168.0.4","servers":2,"clients":0}
{"level":"info","ts":1489193908.1748319,"caller":"/github.com/aukbit/pluto/server/server.go:165","msg":"start","id":"plt_5QPRA9","name":"pluto","id":"srv_R3E4TJ","name":"server","format":"http","port":":8080"}
{"level":"info","ts":1489193908.1748629,"caller":"/github.com/aukbit/pluto/server/server.go:165","msg":"start","id":"plt_5QPRA9","name":"pluto","id":"srv_0XJJCD","name":"pluto_health_server","format":"http","port":":9090"}
{"level":"info","ts":1489193919.5190237,"caller":"/github.com/aukbit/pluto/server/http_middleware.go:25","msg":"request","id":"plt_5QPRA9","name":"pluto","id":"srv_R3E4TJ","name":"server","format":"http","port":":8080","event":"evt_TLPT9N9D69MF","method":"GET","url":"/"}
```

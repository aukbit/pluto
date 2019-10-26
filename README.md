# Pluto [![Circle CI](https://circleci.com/gh/aukbit/pluto.svg?style=svg)](https://circleci.com/gh/aukbit/pluto)
An implementation of microservices with golang to tackle some of the challenges in distributed systems.
### Features
- Currently supports a multiplexer HTTP router with dynamic paths and still compatible with the standard net/http library.
- Client/Server implementation with [gRPC](http://www.grpc.io/) for communication between services.
- Service health check.
- Structured Logs by using [zerolog](https://github.com/rs/zerolog).


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

	"github.com/aukbit/pluto/v6"
	"github.com/aukbit/pluto/v6/reply"
	"github.com/aukbit/pluto/v6/server"
	"github.com/aukbit/pluto/v6/server/router"
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
{"timestamp":"2017-08-15T11:40:44.117833902+01:00","severity":"info","id":"plt_CDVNVF","name":"pluto","ip4":"192.168.15.60","servers":2,"clients":0,"message":"starting pluto, servers: 2 clients: 0"}
{"timestamp":"2017-08-15T11:40:44.118181195+01:00","severity":"info","id":"plt_CDVNVF","name":"pluto","server":{"id":"srv_I3JQ3L","name":"server","format":"http","port":":8080"},"message":"starting http server, listening on :8080"}
{"timestamp":"2017-08-15T11:40:44.118130789+01:00","severity":"info","id":"plt_CDVNVF","name":"pluto","server":{"id":"srv_FP9BC7","name":"pluto_health_server","format":"http","port":":9090"},"message":"starting http pluto_health_server, listening on :9090"}
{"timestamp":"2017-08-15T11:40:55.106279683+01:00","severity":"info","id":"plt_CDVNVF","name":"pluto","server":{"id":"srv_I3JQ3L","name":"server","format":"http","port":":8080"},"eid":"N5G58UTAHSTHEPZQ","method":"GET","url":"/","proto":"HTTP/1.1","remote_addr":"[::1]:50853","header":{"Connection":["keep-alive"],"Upgrade-Insecure-Requests":["1"],"User-Agent":["Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36"],"Accept":["text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"],"Accept-Encoding":["gzip, deflate, br"],"Accept-Language":["en-US,en;q=0.8,pt;q=0.6"]},"message":"GET / HTTP/1.1"}
{"timestamp":"2017-08-15T11:41:02.795156001+01:00","severity":"info","id":"plt_CDVNVF","name":"pluto","message":"shutting down, got signal: interrupt"}
{"timestamp":"2017-08-15T11:41:02.79527628+01:00","severity":"warn","id":"plt_CDVNVF","name":"pluto","server":{"id":"srv_FP9BC7","name":"pluto_health_server","format":"http","port":":9090"},"message":"pluto_health_server has just exited"}
{"timestamp":"2017-08-15T11:41:02.795296844+01:00","severity":"warn","id":"plt_CDVNVF","name":"pluto","server":{"id":"srv_I3JQ3L","name":"server","format":"http","port":":8080"},"message":"server has just exited"}
{"timestamp":"2017-08-15T11:41:02.79531183+01:00","severity":"warn","id":"plt_CDVNVF","name":"pluto","message":"pluto has just exited"}
```

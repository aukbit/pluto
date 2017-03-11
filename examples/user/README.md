# Service

This is an example of creating a simple json backend api using pluto services

## Prerequisites
You should already have a running Cassandra instance. You can follow instructions here [GoCql](https://academy.datastax.com/resources/getting-started-apache-cassandra-and-go)
Also the keyspace **pluto_backend** and schema **users** still set up, as mention her [pluto/README-cassandra.md](../../README-cassandra.md)
README-cassandra.md

You should already have a gRPC installed. You can follow instructions here [gRPC](http://www.grpc.io/docs/quickstart/go.html#prerequisites)

### Compile proto file from proto directory
```
protoc examples/user/proto/user.proto --go_out=plugins=grpc:.
```

### Run Tests
```
$ go test -v ./examples/user
2017-03-11T01:51:58.771Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_VK3ED3", "name": "backend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 0}
2017-03-11T01:51:58.771Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:82	connect	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:51:58.771Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:92	session	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:00.707Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:00.707Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:00.707Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:00.707Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:00.707Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:00.776Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 1}
2017-03-11T01:52:00.777Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:00.777Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:00.777Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:00.777Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:78	dial	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065"}
2017-03-11T01:52:00.777Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:101	watch	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065"}
2017-03-11T01:52:00.777Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_21P3D901ZXVG", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T01:52:00.777Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_21P3D901ZXVG", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T01:52:00.778Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:00.778Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:01.707Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:01.707Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:01.708Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
=== RUN   TestExampleUser
2017-03-11T01:52:01.777Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:01.778Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:01.778Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:01.779Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_8T7PRRKHY9Z7", "method": "POST", "url": "/user"}
2017-03-11T01:52:01.779Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_8T7PRRKHY9Z7", "method": "/user.UserService/CreateUser"}
2017-03-11T01:52:01.779Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_8T7PRRKHY9Z7", "method": "/user.UserService/CreateUser"}
2017-03-11T01:52:01.779Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:92	session	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:02.711Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:02.711Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:02.711Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:02.777Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:02.778Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:02.778Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:03.715Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:03.715Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:03.715Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:03.720Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:108	close	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:03.721Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_LLQV5I9BE1ZY", "method": "GET", "url": "/user/dd2e5f11-eca9-4472-93af-e83e4c4b69bf"}
2017-03-11T01:52:03.721Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_LLQV5I9BE1ZY", "method": "/user.UserService/ReadUser"}
2017-03-11T01:52:03.721Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_LLQV5I9BE1ZY", "method": "/user.UserService/ReadUser"}
2017-03-11T01:52:03.721Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:92	session	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:03.780Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:03.780Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:03.780Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:04.718Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:04.718Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:04.718Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:04.783Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:04.783Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:04.783Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:05.644Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:108	close	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:05.644Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_KMKFA3WEONFX", "method": "GET", "url": "/user/abc"}
2017-03-11T01:52:05.644Z	INFO	/github.com/aukbit/pluto/examples/user/frontend/views/views.go:62	Id abc not found	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_KMKFA3WEONFX"}
2017-03-11T01:52:05.645Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_IHA8EF28N1I7", "method": "PUT", "url": "/user/dd2e5f11-eca9-4472-93af-e83e4c4b69bf"}
2017-03-11T01:52:05.645Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_IHA8EF28N1I7", "method": "/user.UserService/UpdateUser"}
2017-03-11T01:52:05.645Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_IHA8EF28N1I7", "method": "/user.UserService/UpdateUser"}
2017-03-11T01:52:05.645Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:92	session	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:05.723Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:05.723Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:05.723Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:05.787Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:05.787Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:05.787Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:06.725Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:06.725Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:06.725Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:06.791Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:06.791Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:06.791Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:07.558Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:108	close	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:07.558Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_RG447JG7SZF9", "method": "PUT", "url": "/user/abc"}
2017-03-11T01:52:07.558Z	INFO	/github.com/aukbit/pluto/examples/user/frontend/views/views.go:98	Id abc not found	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_RG447JG7SZF9"}
2017-03-11T01:52:07.559Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_4BUPNHAZOWP1", "method": "DELETE", "url": "/user/dd2e5f11-eca9-4472-93af-e83e4c4b69bf"}
2017-03-11T01:52:07.559Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_4BUPNHAZOWP1", "method": "/user.UserService/DeleteUser"}
2017-03-11T01:52:07.559Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_4BUPNHAZOWP1", "method": "/user.UserService/DeleteUser"}
2017-03-11T01:52:07.559Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:92	session	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:07.729Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:07.729Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:07.729Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:07.792Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:07.792Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:07.792Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:08.730Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:08.730Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:08.730Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:08.794Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:08.794Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:08.794Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:09.502Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:108	close	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:09.503Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_DZUVR0BLMZ2Q", "method": "DELETE", "url": "/user/abc"}
2017-03-11T01:52:09.503Z	INFO	/github.com/aukbit/pluto/examples/user/frontend/views/views.go:139	Id abc not found	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_DZUVR0BLMZ2Q"}
2017-03-11T01:52:09.503Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_8TPJU0167EFV", "method": "GET", "url": "/user?name=Gopher"}
2017-03-11T01:52:09.503Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_8TPJU0167EFV", "method": "/user.UserService/FilterUsers"}
2017-03-11T01:52:09.503Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_8TPJU0167EFV", "method": "/user.UserService/FilterUsers"}
2017-03-11T01:52:09.503Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:92	session	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
2017-03-11T01:52:09.735Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:09.735Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:09.735Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:09.794Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:09.794Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:09.794Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:10.737Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:10.737Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:10.737Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:10.797Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:10.797Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
2017-03-11T01:52:10.797Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:11.385Z	INFO	/github.com/aukbit/pluto/datastore/datastore.go:108	close	{"id": "plt_VK3ED3", "name": "backend_pluto", "type": "db", "id": "db_2MUK86", "name": "client_db", "target": "127.0.0.1", "keyspace": "examples_user_backend"}
--- PASS: TestExampleUser (9.61s)
PASS
2017-03-11T01:52:11.738Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:11.738Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:11.738Z	INFO	/github.com/aukbit/pluto/service.go:259	signal received	{"id": "plt_VK3ED3", "name": "backend_pluto", "signal": "interrupt"}
2017-03-11T01:52:11.738Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:11.738Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:11.800Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:11.800Z	INFO	/github.com/aukbit/pluto/service.go:259	signal received	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "signal": "interrupt"}
2017-03-11T01:52:11.800Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:11.800Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:11.800Z	INFO	/github.com/aukbit/pluto/client/client.go:167	close	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc"}
2017-03-11T01:52:11.800Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:11.800Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:137	closed	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "clt_HYPSN4", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065"}
2017-03-11T01:52:12.743Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_Y1ENHL", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T01:52:12.744Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_VK3ED3", "name": "backend_pluto", "id": "srv_4X12EL", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T01:52:12.744Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_VK3ED3", "name": "backend_pluto"}
2017-03-11T01:52:12.802Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TJ9USO", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T01:52:12.802Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_EQ5HLL", "name": "frontend_pluto", "id": "srv_TIO7KY", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T01:52:12.802Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_EQ5HLL", "name": "frontend_pluto"}
ok  	github.com/aukbit/pluto/examples/user	14.043s
```

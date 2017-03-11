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
$ go test -v ./examples/user -run ^TestExampleUser$
2017-03-11T00:24:27.905Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_NHMAAW", "name": "backend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 0}
{"level":"info","ts":1489191867.9051466,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:50","msg":"connect","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
{"level":"info","ts":1489191867.9051673,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:60","msg":"session","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
2017-03-11T00:24:29.840Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:29.840Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:29.840Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto"}
2017-03-11T00:24:29.840Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:29.841Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:29.906Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_IOQLWE", "name": "frontend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 1}
2017-03-11T00:24:29.906Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
2017-03-11T00:24:29.906Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:29.907Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:29.907Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:78	dial	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "clt_Q2CYFH", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065"}
2017-03-11T00:24:29.907Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:101	watch	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "clt_Q2CYFH", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065"}
2017-03-11T00:24:29.907Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "clt_Q2CYFH", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_3W3AL15AY46L", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T00:24:29.907Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:29.907Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_3W3AL15AY46L", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T00:24:29.908Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:30.845Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:30.845Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto"}
2017-03-11T00:24:30.845Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:30.909Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
2017-03-11T00:24:30.909Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:30.909Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
=== RUN   TestExampleUser
2017-03-11T00:24:30.912Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_3C9PWC3UN2I9", "method": "POST", "url": "/user"}
2017-03-11T00:24:30.912Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "clt_Q2CYFH", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_3C9PWC3UN2I9", "method": "/user.UserService/CreateUser"}
2017-03-11T00:24:30.912Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_3C9PWC3UN2I9", "method": "/user.UserService/CreateUser"}
{"level":"info","ts":1489191870.9128742,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:60","msg":"session","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
2017-03-11T00:24:31.846Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:31.846Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:31.846Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto"}
2017-03-11T00:24:31.911Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:31.911Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:31.911Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
2017-03-11T00:24:32.850Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:32.850Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:32.850Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto"}
{"level":"info","ts":1489191872.8528101,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:76","msg":"close","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
2017-03-11T00:24:32.853Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_YR6CVYYCDRB7", "method": "GET", "url": "/user/e8b37754-c4e0-4e37-b3a7-be5df4dc3609"}
2017-03-11T00:24:32.853Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "clt_Q2CYFH", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_YR6CVYYCDRB7", "method": "/user.UserService/ReadUser"}
2017-03-11T00:24:32.854Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_YR6CVYYCDRB7", "method": "/user.UserService/ReadUser"}
{"level":"info","ts":1489191872.8541164,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:60","msg":"session","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
2017-03-11T00:24:32.917Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
2017-03-11T00:24:32.917Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:32.917Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:33.854Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:33.854Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:33.854Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto"}
2017-03-11T00:24:33.921Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
2017-03-11T00:24:33.921Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:33.921Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
{"level":"info","ts":1489191874.7501817,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:76","msg":"close","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
2017-03-11T00:24:34.750Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_YNQBKI4TEFT7", "method": "GET", "url": "/user/abc"}
2017-03-11T00:24:34.750Z	INFO	/github.com/aukbit/pluto/examples/user/frontend/views/views.go:62	Id abc not found	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_YNQBKI4TEFT7"}
2017-03-11T00:24:34.751Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_NG253V1H7PAY", "method": "PUT", "url": "/user/e8b37754-c4e0-4e37-b3a7-be5df4dc3609"}
2017-03-11T00:24:34.751Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "clt_Q2CYFH", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_NG253V1H7PAY", "method": "/user.UserService/UpdateUser"}
2017-03-11T00:24:34.751Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_NG253V1H7PAY", "method": "/user.UserService/UpdateUser"}
{"level":"info","ts":1489191874.7516658,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:60","msg":"session","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
2017-03-11T00:24:34.859Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto"}
2017-03-11T00:24:34.859Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:34.859Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:34.926Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:34.926Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:34.926Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
2017-03-11T00:24:35.860Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto"}
2017-03-11T00:24:35.860Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:35.860Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:35.928Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:35.928Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:35.928Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
{"level":"info","ts":1489191876.6698723,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:76","msg":"close","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
2017-03-11T00:24:36.670Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_340KTHG7L299", "method": "PUT", "url": "/user/abc"}
2017-03-11T00:24:36.670Z	INFO	/github.com/aukbit/pluto/examples/user/frontend/views/views.go:98	Id abc not found	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_340KTHG7L299"}
2017-03-11T00:24:36.671Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_KO57E5FVWIBF", "method": "DELETE", "url": "/user/e8b37754-c4e0-4e37-b3a7-be5df4dc3609"}
2017-03-11T00:24:36.671Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "clt_Q2CYFH", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065", "event": "evt_KO57E5FVWIBF", "method": "/user.UserService/DeleteUser"}
2017-03-11T00:24:36.671Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065", "event": "evt_KO57E5FVWIBF", "method": "/user.UserService/DeleteUser"}
{"level":"info","ts":1489191876.6714473,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:60","msg":"session","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
2017-03-11T00:24:36.863Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:36.863Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:36.863Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto"}
2017-03-11T00:24:36.928Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:36.928Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:36.928Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
2017-03-11T00:24:37.863Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:37.863Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto"}
2017-03-11T00:24:37.863Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:37.929Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:37.929Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:37.929Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
{"level":"info","ts":1489191878.6080215,"caller":"/github.com/aukbit/pluto/datastore/gocql_db.go:76","msg":"close","type":"db","id":"db_KTIUSU","name":"client_db","target":"127.0.0.1","keyspace":"examples_user_backend"}
2017-03-11T00:24:38.608Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_MA6ZLP5OTREE", "method": "DELETE", "url": "/user/abc"}
2017-03-11T00:24:38.608Z	INFO	/github.com/aukbit/pluto/examples/user/frontend/views/views.go:139	Id abc not found	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087", "event": "evt_MA6ZLP5OTREE"}
--- PASS: TestExampleUser (7.70s)
PASS
2017-03-11T00:24:38.865Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:38.865Z	INFO	/github.com/aukbit/pluto/service.go:256	signal received	{"id": "plt_NHMAAW", "name": "backend_pluto", "signal": "interrupt"}
2017-03-11T00:24:38.865Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:38.865Z	INFO	/github.com/aukbit/pluto/server/server.go:110	stop	{"id": "plt_NHMAAW", "name": "backend_pluto"}
{"level":"info","ts":1489191878.8652296,"caller":"/github.com/aukbit/pluto/server/server.go:110","msg":"stop"}
2017-03-11T00:24:38.929Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:38.929Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:38.929Z	INFO	/github.com/aukbit/pluto/service.go:256	signal received	{"id": "plt_IOQLWE", "name": "frontend_pluto", "signal": "interrupt"}
2017-03-11T00:24:38.929Z	INFO	/github.com/aukbit/pluto/server/server.go:110	stop	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
{"level":"info","ts":1489191878.929837,"caller":"/github.com/aukbit/pluto/client/client.go:175","msg":"close"}
{"level":"info","ts":1489191878.9298427,"caller":"/github.com/aukbit/pluto/server/server.go:110","msg":"stop"}
2017-03-11T00:24:38.929Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:137	closed	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "clt_Q2CYFH", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65065"}
2017-03-11T00:24:39.868Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_EQLN0X", "name": "backend_pluto_health_server", "format": "http", "port": ":9096"}
2017-03-11T00:24:39.868Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_NHMAAW", "name": "backend_pluto", "id": "srv_D0IOIE", "name": "server", "format": "grpc", "port": ":65065"}
2017-03-11T00:24:39.868Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_NHMAAW", "name": "backend_pluto"}
2017-03-11T00:24:39.935Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_WJLUA8", "name": "frontend_pluto_health_server", "format": "http", "port": ":9097"}
2017-03-11T00:24:39.935Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_IOQLWE", "name": "frontend_pluto", "id": "srv_3QDHFV", "name": "api_server", "format": "http", "port": ":8087"}
2017-03-11T00:24:39.935Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_IOQLWE", "name": "frontend_pluto"}
ok  	github.com/aukbit/pluto/examples/user	12.045s
```

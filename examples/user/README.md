# Service

This is an example of creating a json backend api using pluto services

## Prerequisites
You should already have a running Cassandra instance. You can follow instructions here [GoCql](https://academy.datastax.com/resources/getting-started-apache-cassandra-and-go) 
Also the keyspace **pluto_backend** and schema **users** still set up, as mention her [pluto/README-cassandra.md](../../README-cassandra.md)
README-cassandra.md

You should already have a gRPC installed. You can follow instructions here [gRPC](http://www.grpc.io/docs/quickstart/go.html#prerequisites)

### 1. Compile proto file from proto directory
```
protoc ./user.proto --go_out=plugins=grpc:.
```
### 2. Run Backend Service
```
$ go run examples/user/backend/main.go
2016/09/19 18:14:35 START pluto_backend 	5014b5df-6d5d-42c0-891f-46c8808ce0aa
2016/09/19 18:14:35 ----- gcql pluto_backend connected on 127.0.0.1
2016/09/19 18:14:35 START grpc server_default 	589f7f38-b345-4675-9ca3-e138a5495cbb
2016/09/19 18:14:35 ----- grpc server_default listening on [::]:65060
```
### 3. Run Tests
```
$ go test -v pluto/examples/user -run ^TestAll$
=== RUN   TestAll
2016/09/19 18:19:33 START pluto_frontend        012cd778-7a15-4250-a4dd-cbaef28568ef
2016/09/19 18:19:33 DIAL  grpc client_user      d080697c-0692-4aa8-b8f6-9410e53f0ffc
2016/09/19 18:19:33 START http server_api       c270e874-3412-45bf-82d4-7b9974688d55
2016/09/19 18:19:33 ----- http server_api listening on [::]:8080
2016/09/19 18:19:33 ----- POST /user
2016/09/19 18:19:33 ----- GET /user/6e50acac-51e4-4f6d-b49c-54a137b1e8b4
2016/09/19 18:19:33 ----- PUT /user/6e50acac-51e4-4f6d-b49c-54a137b1e8b4
2016/09/19 18:19:33 ----- DELETE /user/6e50acac-51e4-4f6d-b49c-54a137b1e8b4
2016/09/19 18:19:33 ----- GET /user?name=Gopher
--- PASS: TestAll (0.02s)
PASS
ok      pluto/examples/user     0.037s
```

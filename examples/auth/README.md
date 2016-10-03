# Service

This is an example of how to use pluto for authentication services

## Prerequisites
You should already have a gRPC installed. You can follow instructions here [gRPC](http://www.grpc.io/docs/quickstart/go.html#prerequisites)

### Compile proto file from proto directory
```
protoc ./auth.proto --go_out=plugins=grpc:.
```
### Generate RSA private and public keys
```
# Key considerations for algorithm "RSA" ≥ 2048-bit
$ openssl genrsa -out auth.rsa 2048
$ openssl rsa -in auth.rsa -pubout > auth.rsa.pub
```

### Run Tests
```
$ go test -v ./examples/auth -run ^TestAll$

```

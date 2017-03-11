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
$ go test -v ./examples/auth -run ^TestExampleAuth$
2017-03-11T00:32:38.486Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 0}
2017-03-11T00:32:38.488Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto"}
2017-03-11T00:32:38.488Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_V1D2CV", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T00:32:38.488Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_WEYRMH", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T00:32:38.488Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_V1D2CV", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T00:32:38.489Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_WEYRMH", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T00:32:38.539Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 1}
2017-03-11T00:32:38.539Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto"}
2017-03-11T00:32:38.539Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_Q9H664", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T00:32:38.539Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_QR02IA", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T00:32:38.540Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:78	dial	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "clt_YD10JR", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T00:32:38.540Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:101	watch	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "clt_YD10JR", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T00:32:38.540Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "clt_YD10JR", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081", "event": "evt_97ROFDTNGIOV", "method": "/grpc.health.v1.Health/Check"}
2017/03/11 00:32:38 grpc: addrConn.resetTransport failed to create client transport: connection error: desc = "transport: dial tcp 127.0.0.1:65081: getsockopt: connection refused"; Reconnecting to {127.0.0.1:65081 <nil>}
2017-03-11T00:32:38.541Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_Q9H664", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T00:32:38.541Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_QR02IA", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T00:32:38.590Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 1}
2017-03-11T00:32:38.590Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_0TLAE7", "name": "auth_backend_pluto"}
2017-03-11T00:32:38.590Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_IV1HA4", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T00:32:38.590Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_C5316A", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T00:32:38.590Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:78	dial	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "clt_L4VHOO", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080"}
2017-03-11T00:32:38.590Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:101	watch	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "clt_L4VHOO", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080"}
2017-03-11T00:32:38.590Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_C5316A", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T00:32:38.590Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "clt_L4VHOO", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080", "event": "evt_X1QNECI43LCF", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T00:32:38.591Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_V1D2CV", "name": "server", "format": "grpc", "port": ":65080", "event": "evt_X1QNECI43LCF", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T00:32:38.591Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_IV1HA4", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T00:32:38.641Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 1}
2017-03-11T00:32:38.641Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto"}
2017-03-11T00:32:38.641Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_2RXZWI", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T00:32:38.641Z	INFO	/github.com/aukbit/pluto/server/server.go:165	start	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_T9BCKU", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T00:32:38.641Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:78	dial	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "clt_P2065U", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T00:32:38.641Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:101	watch	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "clt_P2065U", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T00:32:38.641Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "clt_P2065U", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081", "event": "evt_STT62RESDYLP", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T00:32:38.642Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_C5316A", "name": "server", "format": "grpc", "port": ":65081", "event": "evt_STT62RESDYLP", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T00:32:38.642Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_2RXZWI", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T00:32:38.642Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_T9BCKU", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T00:32:39.490Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto"}
2017-03-11T00:32:39.490Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_WEYRMH", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T00:32:39.490Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_V1D2CV", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T00:32:39.543Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_Q9H664", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T00:32:39.543Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_QR02IA", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T00:32:39.543Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto"}
2017-03-11T00:32:39.590Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_0TLAE7", "name": "auth_backend_pluto"}
2017-03-11T00:32:39.590Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_C5316A", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T00:32:39.591Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_IV1HA4", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
=== RUN   TestExampleAuth
2017-03-11T00:32:39.641Z	DEBUG	/github.com/aukbit/pluto/service.go:263	pulse	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto"}
2017-03-11T00:32:39.642Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_2RXZWI", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T00:32:39.643Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_T9BCKU", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T00:32:39.652Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_T9BCKU", "name": "api_server", "format": "http", "port": ":8089", "event": "evt_IOP0B6EXNFLQ", "method": "POST", "url": "/authenticate"}
2017-03-11T00:32:39.652Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "clt_P2065U", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081", "event": "evt_IOP0B6EXNFLQ", "method": "/auth.AuthService/Authenticate"}
2017-03-11T00:32:39.652Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_C5316A", "name": "server", "format": "grpc", "port": ":65081", "event": "evt_IOP0B6EXNFLQ", "method": "/auth.AuthService/Authenticate"}
2017-03-11T00:32:39.652Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "clt_L4VHOO", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080", "event": "evt_IOP0B6EXNFLQ", "method": "/user.UserService/VerifyUser"}
2017-03-11T00:32:39.652Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_V1D2CV", "name": "server", "format": "grpc", "port": ":65080", "event": "evt_IOP0B6EXNFLQ", "method": "/user.UserService/VerifyUser"}
2017-03-11T00:32:39.657Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_Q9H664", "name": "user_api_server", "format": "http", "port": ":8088", "event": "evt_H7KX3VUYGX1X", "method": "POST", "url": "/user"}
2017-03-11T00:32:39.657Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "clt_YD10JR", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081", "event": "evt_H7KX3VUYGX1X", "method": "/auth.AuthService/Verify"}
2017-03-11T00:32:39.658Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_C5316A", "name": "server", "format": "grpc", "port": ":65081", "event": "evt_H7KX3VUYGX1X", "method": "/auth.AuthService/Verify"}
--- PASS: TestExampleAuth (0.02s)
PASS
2017-03-11T00:32:40.494Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_V1D2CV", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T00:32:40.494Z	INFO	/github.com/aukbit/pluto/service.go:256	signal received	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "signal": "interrupt"}
2017-03-11T00:32:40.494Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_WEYRMH", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T00:32:40.494Z	INFO	/github.com/aukbit/pluto/server/server.go:110	stop	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto"}
{"level":"info","ts":1489192360.49431,"caller":"/github.com/aukbit/pluto/server/server.go:110","msg":"stop"}
2017-03-11T00:32:40.548Z	INFO	/github.com/aukbit/pluto/service.go:256	signal received	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "signal": "interrupt"}
2017-03-11T00:32:40.548Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_Q9H664", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T00:32:40.548Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_QR02IA", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
{"level":"info","ts":1489192360.5489993,"caller":"/github.com/aukbit/pluto/server/server.go:110","msg":"stop"}
{"level":"info","ts":1489192360.5490174,"caller":"/github.com/aukbit/pluto/client/client.go:175","msg":"close"}
2017-03-11T00:32:40.549Z	INFO	/github.com/aukbit/pluto/server/server.go:110	stop	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto"}
2017-03-11T00:32:40.549Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:137	closed	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "clt_YD10JR", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T00:32:40.596Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_C5316A", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T00:32:40.596Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_IV1HA4", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T00:32:40.596Z	INFO	/github.com/aukbit/pluto/service.go:256	signal received	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "signal": "interrupt"}
{"level":"info","ts":1489192360.5964422,"caller":"/github.com/aukbit/pluto/server/server.go:110","msg":"stop"}
{"level":"info","ts":1489192360.596436,"caller":"/github.com/aukbit/pluto/client/client.go:175","msg":"close"}
2017-03-11T00:32:40.596Z	INFO	/github.com/aukbit/pluto/server/server.go:110	stop	{"id": "plt_0TLAE7", "name": "auth_backend_pluto"}
2017-03-11T00:32:40.596Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:137	closed	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "clt_L4VHOO", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080"}
2017-03-11T00:32:40.643Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_T9BCKU", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T00:32:40.643Z	DEBUG	/github.com/aukbit/pluto/server/server.go:286	pulse	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_2RXZWI", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T00:32:40.643Z	INFO	/github.com/aukbit/pluto/service.go:256	signal received	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "signal": "interrupt"}
{"level":"info","ts":1489192360.6438687,"caller":"/github.com/aukbit/pluto/server/server.go:110","msg":"stop"}
2017-03-11T00:32:40.643Z	INFO	/github.com/aukbit/pluto/server/server.go:110	stop	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto"}
{"level":"info","ts":1489192360.6438704,"caller":"/github.com/aukbit/pluto/client/client.go:175","msg":"close"}
2017-03-11T00:32:40.643Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:137	closed	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "clt_P2065U", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T00:32:41.494Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_WEYRMH", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T00:32:41.494Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto", "id": "srv_V1D2CV", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T00:32:41.494Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_6FC4GO", "name": "mockuserbackend_pluto"}
2017-03-11T00:32:41.554Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_QR02IA", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T00:32:41.554Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto", "id": "srv_Q9H664", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T00:32:41.554Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_L4SC3N", "name": "mockuserfrontend_pluto"}
2017-03-11T00:32:41.600Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_IV1HA4", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T00:32:41.601Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_0TLAE7", "name": "auth_backend_pluto", "id": "srv_C5316A", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T00:32:41.601Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_0TLAE7", "name": "auth_backend_pluto"}
2017-03-11T00:32:41.647Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_T9BCKU", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T00:32:41.647Z	INFO	/github.com/aukbit/pluto/server/server.go:104	exit	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto", "id": "srv_2RXZWI", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T00:32:41.647Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_QYMKPP", "name": "auth_frontend_pluto"}
ok  	github.com/aukbit/pluto/examples/auth	3.178s

```

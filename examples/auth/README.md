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
$  go test -v ./examples/auth -run ^TestExampleAuth$
2017-03-11T01:56:41.637Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 0}
2017-03-11T01:56:41.638Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_E0DXVB", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T01:56:41.638Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto"}
2017-03-11T01:56:41.638Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_958RAS", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T01:56:41.638Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_958RAS", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T01:56:41.638Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_E0DXVB", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T01:56:42.140Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 1}
2017-03-11T01:56:42.140Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto"}
2017-03-11T01:56:42.140Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_ZJ4A95", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T01:56:42.140Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_D5O995", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T01:56:42.141Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:78	dial	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "clt_2UYGZD", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T01:56:42.141Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:101	watch	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "clt_2UYGZD", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T01:56:42.141Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "clt_2UYGZD", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081", "event": "evt_FOD7KYZ323BU", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T01:56:42.141Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_ZJ4A95", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T01:56:42.141Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_D5O995", "name": "user_api_server", "format": "http", "port": ":8088"}
2017/03/11 01:56:42 grpc: addrConn.resetTransport failed to create client transport: connection error: desc = "transport: dial tcp 127.0.0.1:65081: getsockopt: connection refused"; Reconnecting to {127.0.0.1:65081 <nil>}
2017-03-11T01:56:42.642Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_E0DXVB", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T01:56:42.642Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto"}
2017-03-11T01:56:42.642Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_958RAS", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T01:56:42.643Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 1}
2017-03-11T01:56:42.643Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_MQIVIF", "name": "auth_backend_pluto"}
2017-03-11T01:56:42.643Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_U4RRCI", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T01:56:42.643Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_89XE21", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T01:56:42.643Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:78	dial	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "clt_24HY2R", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080"}
2017-03-11T01:56:42.643Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_89XE21", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T01:56:42.643Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:101	watch	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "clt_24HY2R", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080"}
2017-03-11T01:56:42.643Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "clt_24HY2R", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080", "event": "evt_0XKNLHICBSK7", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T01:56:42.644Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_958RAS", "name": "server", "format": "grpc", "port": ":65080", "event": "evt_0XKNLHICBSK7", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T01:56:42.644Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_U4RRCI", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T01:56:43.143Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_ZJ4A95", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T01:56:43.143Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_D5O995", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T01:56:43.143Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto"}
2017-03-11T01:56:43.144Z	INFO	/github.com/aukbit/pluto/service.go:155	start	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "ip4": "192.168.0.4", "servers": 2, "clients": 1}
2017-03-11T01:56:43.144Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_SF58X3", "name": "auth_frontend_pluto"}
2017-03-11T01:56:43.144Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_PXDMNS", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T01:56:43.144Z	INFO	/github.com/aukbit/pluto/server/server.go:167	start	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_BO354T", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T01:56:43.144Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:78	dial	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "clt_27BLE6", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T01:56:43.145Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:101	watch	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "clt_27BLE6", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T01:56:43.145Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "clt_27BLE6", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081", "event": "evt_27S9VA4TB338", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T01:56:43.145Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_89XE21", "name": "server", "format": "grpc", "port": ":65081", "event": "evt_27S9VA4TB338", "method": "/grpc.health.v1.Health/Check"}
2017-03-11T01:56:43.145Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_BO354T", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T01:56:43.146Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_PXDMNS", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T01:56:43.645Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_U4RRCI", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T01:56:43.645Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_958RAS", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T01:56:43.645Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_MQIVIF", "name": "auth_backend_pluto"}
2017-03-11T01:56:43.645Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto"}
2017-03-11T01:56:43.645Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_89XE21", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T01:56:43.645Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_E0DXVB", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T01:56:44.145Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_ZJ4A95", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T01:56:44.145Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_SF58X3", "name": "auth_frontend_pluto"}
2017-03-11T01:56:44.145Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_D5O995", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T01:56:44.145Z	DEBUG	/github.com/aukbit/pluto/service.go:266	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto"}
2017-03-11T01:56:44.146Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_BO354T", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
=== RUN   TestExampleAuth
2017-03-11T01:56:44.146Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_PXDMNS", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T01:56:44.149Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_PXDMNS", "name": "api_server", "format": "http", "port": ":8089", "event": "evt_QXL1L4XAAPCT", "method": "POST", "url": "/authenticate"}
2017-03-11T01:56:44.149Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "clt_27BLE6", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081", "event": "evt_QXL1L4XAAPCT", "method": "/auth.AuthService/Authenticate"}
2017-03-11T01:56:44.149Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_89XE21", "name": "server", "format": "grpc", "port": ":65081", "event": "evt_QXL1L4XAAPCT", "method": "/auth.AuthService/Authenticate"}
2017-03-11T01:56:44.149Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "clt_24HY2R", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080", "event": "evt_QXL1L4XAAPCT", "method": "/user.UserService/VerifyUser"}
2017-03-11T01:56:44.150Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_958RAS", "name": "server", "format": "grpc", "port": ":65080", "event": "evt_QXL1L4XAAPCT", "method": "/user.UserService/VerifyUser"}
2017-03-11T01:56:44.154Z	INFO	/github.com/aukbit/pluto/server/http_middleware.go:25	request	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_D5O995", "name": "user_api_server", "format": "http", "port": ":8088", "event": "evt_Y6WF51VZXG0F", "method": "POST", "url": "/user"}
2017-03-11T01:56:44.154Z	INFO	/github.com/aukbit/pluto/client/balancer/interceptor.go:20	call	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "clt_2UYGZD", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081", "event": "evt_Y6WF51VZXG0F", "method": "/auth.AuthService/Verify"}
2017-03-11T01:56:44.154Z	INFO	/github.com/aukbit/pluto/server/grpc_interceptor.go:40	request	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_89XE21", "name": "server", "format": "grpc", "port": ":65081", "event": "evt_Y6WF51VZXG0F", "method": "/auth.AuthService/Verify"}
--- PASS: TestExampleAuth (0.01s)
PASS
2017-03-11T01:56:44.646Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_E0DXVB", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T01:56:44.646Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_U4RRCI", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T01:56:44.646Z	INFO	/github.com/aukbit/pluto/service.go:259	signal received	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "signal": "interrupt"}
2017-03-11T01:56:44.646Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_89XE21", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T01:56:44.646Z	INFO	/Users/paulo/Developmen2017-03-11T01:56:44.646Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_958RAS", "name": "server", "format": "grpc", "port": ":65080"}
t/go/src/github.com/aukbit/pluto/service.go:259	signal received	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "signal": "interrupt"}
2017-03-11T01:56:44.646Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_958RAS", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T01:56:44.646Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_E0DXVB", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T01:56:44.646Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_89XE21", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T01:56:44.646Z	INFO	/github.com/aukbit/pluto/client/client.go:167	close	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "clt_24HY2R", "name": "user_client", "format": "grpc"}
2017-03-11T01:56:44.646Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_U4RRCI", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T01:56:44.646Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:137	closed	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "clt_24HY2R", "name": "user_client", "format": "grpc", "target": "127.0.0.1:65080"}
2017-03-11T01:56:45.146Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_D5O995", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T01:56:45.146Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_BO354T", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/service.go:259	signal received	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "signal": "interrupt"}
2017-03-11T01:56:45.146Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_PXDMNS", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/service.go:259	signal received	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "signal": "interrupt"}
2017-03-11T01:56:45.146Z	DEBUG	/github.com/aukbit/pluto/server/server.go:288	pulse	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_ZJ4A95", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_ZJ4A95", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_BO354T", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/client/client.go:167	close	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "clt_2UYGZD", "name": "auth_client", "format": "grpc"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_PXDMNS", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/server/server.go:112	stop	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_D5O995", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/client/client.go:167	close	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "clt_27BLE6", "name": "auth_client", "format": "grpc"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:137	closed	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "clt_2UYGZD", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T01:56:45.146Z	INFO	/github.com/aukbit/pluto/client/balancer/connector.go:137	closed	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "clt_27BLE6", "name": "auth_client", "format": "grpc", "target": "127.0.0.1:65081"}
2017-03-11T01:56:45.647Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_E0DXVB", "name": "mockuserbackend_pluto_health_server", "format": "http", "port": ":9094"}
2017-03-11T01:56:45.647Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_89XE21", "name": "server", "format": "grpc", "port": ":65081"}
2017-03-11T01:56:45.647Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_MQIVIF", "name": "auth_backend_pluto", "id": "srv_U4RRCI", "name": "auth_backend_pluto_health_server", "format": "http", "port": ":9092"}
2017-03-11T01:56:45.647Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_MQIVIF", "name": "auth_backend_pluto"}
2017-03-11T01:56:45.647Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto", "id": "srv_958RAS", "name": "server", "format": "grpc", "port": ":65080"}
2017-03-11T01:56:45.647Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_LRS4JY", "name": "mockuserbackend_pluto"}
2017-03-11T01:56:46.150Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_ZJ4A95", "name": "mockuserfrontend_pluto_health_server", "format": "http", "port": ":9095"}
2017-03-11T01:56:46.150Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto", "id": "srv_D5O995", "name": "user_api_server", "format": "http", "port": ":8088"}
2017-03-11T01:56:46.150Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_8JARGR", "name": "mockuserfrontend_pluto"}
2017-03-11T01:56:46.150Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_BO354T", "name": "auth_frontend_pluto_health_server", "format": "http", "port": ":9093"}
2017-03-11T01:56:46.150Z	INFO	/github.com/aukbit/pluto/server/server.go:106	exit	{"id": "plt_SF58X3", "name": "auth_frontend_pluto", "id": "srv_PXDMNS", "name": "api_server", "format": "http", "port": ":8089"}
2017-03-11T01:56:46.150Z	INFO	/github.com/aukbit/pluto/service.go:94	exit	{"id": "plt_SF58X3", "name": "auth_frontend_pluto"}
ok  	github.com/aukbit/pluto/examples/auth	4.525s
```

# Service

This is an example of creating a distributed deployment of pluto services using Docker (4 nodes), Consul as service discovery

### Create a docker image for service user_backend on node1
```
$ docker run golang go get -v github.com/aukbit/pluto/examples/dist/...
$ docker commit $(docker ps -lq) pluto
$ docker run --name user_backend --network=bridge -p 65060:65060 -p 9090:9090 -d pluto user_backend -grpc_port=:65060 -db=cassandra -keyspace=pluto_user_backend -name=user_backend -consul_addr=192.168.99.101:8500
```
### Create a docker image for service user_bff on node2
```
$ docker run golang go get -v github.com/aukbit/pluto/examples/dist/...
$ docker commit $(docker ps -lq) pluto
$ docker run --name user_bff --network=bridge -p 8082:8082 -p 9090:9090 -d pluto user_bff -http_port=:8082 -target_name=user_backend -name=user_bff -consul_addr=192.168.99.102:8500
```

# Service

This is an example of creating a distributed deployment of pluto services using Docker (4 nodes), Consul as service discovery

### Create a docker image for service user_backend on node1
```
$ cd ./examples/dist/user_backend
$ docker build -f ./Dockerfile -t user_backend .
$ docker run --name user_backend -p 65060:65060 -d user_backend
```
### Create a docker image for service user_bff on node2
```
$ cd ./examples/dist/user_bff
$ docker build -f ./Dockerfile -t user_bff .
$ docker run --name user_bff -p 8080:8080 -d user_bff
```

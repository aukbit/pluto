# Service

This is an example of creating a distributed deployment of pluto services using Docker (4 nodes), Consul as service discovery

### Create a docker image for service user_backend on default node
```
$ cd ./examples/dist/user_backend
$ docker build -f ./Dockerfile -t user_backend .
$ docker run -d -name user_backend user_backend
```

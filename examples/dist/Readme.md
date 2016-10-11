# Service

This is an example of creating a distributed deployment of pluto services using Docker (4 nodes), Consul as service discovery

### Create a docker image for service user_backend on default node
```
$ docker build -f ./examples/dist/user_backend/Dockerfile .
```

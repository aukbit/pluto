# Local OS X installation of Docker, Docker Machine & Docker Compose development environment

## Prerequisites
On OS X we'll launch Docker inside a boot2docker VM via Docker Machine.
* Make sure you have the latest VirtualBox (v5) correctly installed on your system.
* You will need to have Docker installed (1.12.1 or newer), Docker Machine would be installed with Docker engine (v0.8.1 or newer).

## Create a Docker VM machine
You can follow the steps on how to create a [VM machine here](https://docs.docker.com/machine/get-started/)
### Create
```
$ docker-machine create --driver virtualbox default
```
### List available machines again to see your newly minted machine.
```
$ docker-machine ls
```
### Get the environment commands for your new VM.
```
$ docker-machine env default
```
### Connect your shell to the new machine.
```
$ eval "$(docker-machine env default)"
```

## Install Docker Cassandra
[Docker Cassandra](https://hub.docker.com/_/cassandra/)

### Create docker container
Make it broadcast the ip address and exporting ports 7000 (to be visible to other cassandra nodes), also export the port 9042 for services to connect

```
$ docker run --name cassandra -d -e CASSANDRA_BROADCAST_ADDRESS=192.168.99.100 -p 7000:7000 -p 9042:9042 cassandra:3.7
```
### Note: before the service pluto_backend to initiate run the following

Connect to cassandra from cqlsh
```
$ docker exec -ti cassandra cqlsh 192.168.99.100
```

```
cqlsh>
    DESCRIBE keyspaces;
    CREATE KEYSPACE examples_user_backend WITH replication = {'class':'SimpleStrategy', 'replication_factor':1};
    USE examples_user_backend;
    CREATE TABLE users (id uuid, name text, email text, password text, PRIMARY KEY (id));
    DESCRIBE tables;
```

## Run Docker compose
```
$ docker-compose up
Attaching to plutosample_backend_1, plutosample_frontend_1
frontend_1  | 2016/09/22 12:09:43 START pluto_frontend 	aa10b956-1345-4b8f-aa31-db9d246b0db7
frontend_1  | 2016/09/22 12:09:43 START http server_api 	229f13a3-783a-4ee7-8758-524db435522e
backend_1   | 2016/09/22 12:09:43 START pluto_backend 	dcf42621-32b0-4b6d-b309-dda66dcfe8b5
backend_1   | 2016/09/22 12:09:43 ----- gcql pluto_backend connected on 192.168.99.100
frontend_1  | 2016/09/22 12:09:43 ----- http server_api listening on [::]:8080
frontend_1  | 2016/09/22 12:09:43 DIAL  grpc client_user 	187adf81-e810-4cb4-a11c-7ff38f3adff5
backend_1   | 2016/09/22 12:09:43 START grpc server_default 	4789d93b-1b55-4f60-80ad-576883068b9b
backend_1   | 2016/09/22 12:09:43 ----- grpc server_default listening on [::]:65060
```

## Make some requests

```
curl -H "Content-Type: application/json" http://192.168.99.100:8080/user
curl -H "Content-Type: application/json" -X POST -d '{"name":"Gopher", "email": "gopher@email.com", "password":"123456"}' http://192.168.99.100:8080/user
curl -H "Content-Type: application/json" -X PUT -d '{"name":"Super Gopher", "email": "gopher@email.com"}' http://192.168.99.100:8080/user/e6bcd635-38f5-4386-ab24-a310d092f65d
curl -H "Content-Type: application/json" -X GET http://192.168.99.100:8080/user/e6bcd635-38f5-4386-ab24-a310d092f65d
curl -H "Content-Type: application/json" -X DELETE http://192.168.99.100:8080/user/e6bcd635-38f5-4386-ab24-a310d092f65d
```

# How to set up a Consul cluster for real testing

### Create at least 3 nodes with docker-machine
<!-- 2 consul servers 2 consul clients -->
```
$ docker-machine create --driver virtualbox default
$ docker-machine create --driver virtualbox node1
$ docker-machine create --driver virtualbox node2
$ docker-machine create --driver virtualbox node3
```
### List nodes
```
$ docker-machine ls
NAME      ACTIVE   DRIVER       STATE     URL                         SWARM   DOCKER    ERRORS
default   -        virtualbox   Running   tcp://192.168.99.100:2376           v1.12.1   
node1     -        virtualbox   Running   tcp://192.168.99.101:2376           v1.12.1   
node2     -        virtualbox   Running   tcp://192.168.99.102:2376           v1.12.1   
node3     -        virtualbox   Running   tcp://192.168.99.103:2376           v1.12.1
```
#### Open 4 shell terminals, in each one run the below to find the respectiv environment variables for each node
e.g. for node 3
```
$ docker-machine env node3
export DOCKER_TLS_VERIFY="1"
export DOCKER_HOST="tcp://192.168.99.103:2376"
export DOCKER_CERT_PATH="/Users/paulo/.docker/machine/machines/node3"
export DOCKER_MACHINE_NAME="node3"
# Run this command to configure your shell:
# eval $(docker-machine env node3)
```
#### Run the mentioned command to configure the shell
```
$ eval $(docker-machine env node3)
```

### Running Consul Agent in Server Mode in default node
```
$ docker run -d --net=host --name=consul -e 'CONSUL_LOCAL_CONFIG={"skip_leave_on_interrupt": true}' consul agent -ui -server -bind=192.168.99.100 -retry-join=192.168.99.101 -bootstrap-expect=2
```
#### For web ui to be available on the host machine create a configuration file with the HTTP address listening on the public IP
```
{
    "addresses" : {
    "http": "192.168.99.100"
  }
}
```
#### Copy file into container and restart
```
$ docker cp ./basic.json consul:consul/config/
$ docker restart consul
$ docker exec -t consul netstat -lpt
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name    
tcp        0      0 default:8400            0.0.0.0:*               LISTEN      -
tcp        0      0 192.168.99.100:8500     0.0.0.0:*               LISTEN      -
tcp        0      0 0.0.0.0:ssh             0.0.0.0:*               LISTEN      -
tcp        0      0 default:8600            0.0.0.0:*               LISTEN      -
tcp        0      0 192.168.99.100:8300     0.0.0.0:*               LISTEN      -
tcp        0      0 192.168.99.100:8301     0.0.0.0:*               LISTEN      -
tcp        0      0 192.168.99.100:8302     0.0.0.0:*               LISTEN      -
tcp        0      0 :::ssh                  :::*                    LISTEN      -
tcp        0      0 :::2376                 :::*                    LISTEN      -
```

### Running Consul Agent in Server Mode in node1
```
$ docker run -d --net=host --name=consul -e 'CONSUL_LOCAL_CONFIG={"skip_leave_on_interrupt": true}' consul agent -server -bind=192.168.99.101 -retry-join=192.168.99.100 -bootstrap-expect=2
```

### Running Consul Agent in Client Mode in node2
```
$ docker run -d --net=host --name=consul -e 'CONSUL_LOCAL_CONFIG={"leave_on_terminate": true}' consul agent -ui -bind=192.168.99.102 -retry-join=192.168.99.101
```

### Running Consul Agent in Client Mode in node3
```
$ docker run -d --net=host --name=consul -e 'CONSUL_LOCAL_CONFIG={"leave_on_terminate": true}' consul agent -bind=192.168.99.103 -retry-join=192.168.99.101
```

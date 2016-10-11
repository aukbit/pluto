# How to set up a Cassandra testing

## Install Docker Cassandra
[Docker Cassandra](https://hub.docker.com/_/cassandra/)

### Create a docker image
```
$ docker build --name cassandra -d cassandra:3.7
```
### Create docker container
Make it broadcast the ip address and exporting ports 7000 (to be visible to other cassandra nodes), also export the port 9042 for services to connect

```
$ docker run --name cassandra -d -e CASSANDRA_BROADCAST_ADDRESS=192.168.99.100 -p 7000:7000 -p 9042:9042 cassandra:3.7
```

### Set cassandra instance as a service in Consul, creating a config file like
```
{
  "service": {
    "name": "cassandra",
    "tags": ["master", "v3.7"],
    "port": 9042,
    "enableTagOverride": false,
    "checks": [
      {
        "notes": "Ensure cassandra is listening on port 9042",
        "tcp": ":9042",
        "interval": "10s",
        "timeout": "1s"
    },
    {
      "notes": "Ensure cassandra is listening on port 7000 (other cassandra nodes to connect)",
      "tcp": ":7000",
      "interval": "10s",
      "timeout": "1s"
    }
    ]
  }
}
```
### Copy file into consul container and restart
```
$ docker cp ./cassandra.json consul:consul/config/
$ docker restart consul
```

### Note: before running service pluto_user_backend set the following

Connect to cassandra from cqlsh
```
$ docker exec -ti cassandra cqlsh 192.168.99.100
```

```
cqlsh>
    DESCRIBE keyspaces;
    CREATE KEYSPACE pluto_user_backend WITH replication = {'class':'SimpleStrategy','replication_factor':1};
    USE pluto_user_backend;
    CREATE TABLE users (id uuid, name text, email text, password text, PRIMARY KEY (id));
    DESCRIBE tables;
```

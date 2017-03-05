# Run Tests against a cassandra docker instance

### Prerequisites
On OS X we'll launch Docker inside a boot2docker VM via Docker Machine.
* Make sure you have the latest VirtualBox (v5) correctly installed on your system.
* You will need to have Docker installed (1.12.1 or newer), Docker Machine would be installed with Docker engine (v0.8.1 or newer).

## Start a cassandra server instance
To get more information on how to run a cassandra instance with docker go to [Docker Cassandra](https://hub.docker.com/_/cassandra/)

Start a cassandra instance, export port 7000 to allow other cassandra nodes to connect and export port 9042 to allow services to connect:
```
$ docker run --name cassandra -d -p 7000:7000 -p 9042:9042 cassandra:3.7
```
The following command starts another Cassandra container instance and runs cqlsh against your original Cassandra container to execute CQL statements on /config/commands.cql against your database instance:
```
$ docker run --link cassandra:cassandra -v $PWD/test/config/cassandra:/config --rm cassandra:3.7 /config/wait-for-cassandra.sh cassandra cqlsh -f /config/commands.cql cassandra
```

# Run Cassandra locally
### 1. Check Java version
 (minimum 8u40)
  ```
  java -version
  ```

### 2. Install the [Java SDK](http://www.oracle.com/technetwork/java/javase/downloads/jdk8-downloads-2133151.html).

### 3. Download the [Cassandra SDK](http://www.apache.org/dyn/closer.lua/cassandra/3.7/apache-cassandra-3.7-bin.tar.gz); [Cassandra tutorial](http://www.planetcassandra.org/try-cassandra/);

### 4. Run Cassandra on terminal
  ```
  cd Cassandra
  ./apache-cassandra-3.7/bin/cassandra -f
  ```

### 5. Or run Cassandra in background
  ```
  sudo ./apache-cassandra-3.7/bin/cassandra
  ```

### 6. Start cqlsh using the command cqlsh as shown below. It gives the Cassandra cqlsh prompt as output.
  ```
  ./apache-cassandra-3.7/bin/cqlsh
  ```

### 7. Confirm your cluster is running
  ```
  cqlsh> SELECT cluster_name, listen_address FROM system.local;
  ```

### 8. List Keyspaces
  ```
  cqlsh> DESCRIBE keyspaces;
  ```

### 9. Create a development KEYSPACE eg. dev
  ```
  cqlsh> CREATE KEYSPACE dev WITH replication = {'class':'SimpleStrategy', 'replication_factor':1};
  ```

### 10. To use the keyspace we’ve created, type:
  ```
  cqlsh> USE dev;
  ```

### 11. Create Users table:
  ```
  cqlsh> CREATE TABLE users (id timeuuid, name text, email text, password text, PRIMARY KEY (id));
  ```

### 12. Information about data objects in the cluster
 ```
 cqlsh> DESCRIBE tables;
 ```




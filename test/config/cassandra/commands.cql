DESCRIBE keyspaces;
CREATE KEYSPACE IF NOT EXISTS default WITH replication = {'class':'SimpleStrategy','replication_factor':1};
CREATE KEYSPACE examples_user_backend WITH replication = {'class':'SimpleStrategy', 'replication_factor':1};
USE examples_user_backend;
CREATE TABLE users (id uuid, name text, email text, password text, PRIMARY KEY (id));
DESCRIBE tables;
DESCRIBE keyspaces;

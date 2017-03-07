DESCRIBE keyspaces;
CREATE KEYSPACE IF NOT EXISTS default WITH replication = {'class':'SimpleStrategy','replication_factor':1};
USE default;

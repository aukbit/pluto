DESCRIBE keyspaces;
CREATE KEYSPACE IF NOT EXISTS gs_event_source_pluto WITH replication = {'class':'SimpleStrategy','replication_factor':1};
USE gs_event_source_pluto;

DROP TABLE IF EXISTS events;
DROP TYPE IF EXISTS payload;

CREATE TYPE IF NOT EXISTS payload (
    id bigint,
    object text,
    version int,
    data blob
);

CREATE TABLE IF NOT EXISTS events (
    id bigint,
    topic int,
    payload frozen <payload>,
    priority int,
    source text,
    time timestamp,
    status int,
    error_description text,
    PRIMARY KEY (id)
);

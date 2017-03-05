#!/bin/bash

set -e

host="$1"
shift
cmd="$@"

until cqlsh "$host"; do
  >&2 echo "Cassandra is unavailable - sleeping"
  sleep 15
done

>&2 echo "Cassandra is up - executing command"
exec $cmd

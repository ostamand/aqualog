#!/bin/sh
/app/scripts/wait-for.sh postgres:5432
migrate -path db/migration -database "$DB_SOURCE" -verbose up
exec "$1"
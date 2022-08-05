#!/bin/sh
if [ ! -v $APPLY_MIGRATION ]; then
   /app/scripts/wait-for.sh postgres:5432
   migrate -path db/migration -database $DB_SOURCE -verbose up
fi
exec $1
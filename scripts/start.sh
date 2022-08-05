#!/bin/sh
if [ ! -z ${APPLY_MIGRATION+x} ]; then
   /app/scripts/wait-for.sh postgres:5432
   migrate -path db/migration -database $DB_SOURCE -verbose up
fi
exec $1
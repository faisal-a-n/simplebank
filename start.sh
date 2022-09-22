#!/bin/sh

set -e

echo "Run DB migrations"
source /app/config.env
cat /app/config.env

/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Start the app"
exec "$@"
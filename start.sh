#!/bin/sh

set -e

echo "Run DB migrations"
source config.env
/app/migrate -path /app/migrations -database "$DB_SOURCE" -verbose up

echo "Start the app"
exec "$@"
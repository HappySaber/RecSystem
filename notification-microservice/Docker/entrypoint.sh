#!/bin/sh
set -e

echo "running migrations..."
/app/migrator --migrations-path=/app/migrations

echo "migrations completed"

echo "starting application..."
exec /app/app --config=/app/config/prod.yaml
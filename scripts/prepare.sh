#!/bin/bash

set -e

echo "Preparing database and dependencies..."

# Подстраховка на случай запуска без переменных (не в GH Actions)
: "${POSTGRES_HOST:=localhost}"
: "${POSTGRES_PORT:=5432}"
: "${POSTGRES_DB:=project-sem-1}"
: "${POSTGRES_USER:=validator}"
: "${POSTGRES_PASSWORD:=val1dat0r}"

echo "Waiting for PostgreSQL at $POSTGRES_HOST:$POSTGRES_PORT..."
until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q' 2>/dev/null; do
    echo "PostgreSQL is unavailable yet - retry in 2s"
    sleep 2
done
echo "Successfully connetcted to PostgreSQL"

echo "Creating table prices if not exists..."
PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    create_date DATE NOT NULL,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL
);
"

echo "Installing Go dependencies..."
go mod tidy
go mod download

echo "Environment is ready"
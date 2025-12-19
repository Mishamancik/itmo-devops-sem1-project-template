#!/bin/bash

set -e

echo "Preparing database and dependencies..."

# Экспорт для локального запуска (без переменных GH Actions)
export POSTGRES_HOST="${POSTGRES_HOST:=localhost}"
export POSTGRES_PORT="${POSTGRES_PORT:=5432}"
export POSTGRES_DB="${POSTGRES_DB:=project-sem-1}"
export POSTGRES_USER="${POSTGRES_USER:=validator}"
export POSTGRES_PASSWORD="${POSTGRES_PASSWORD:=val1dat0r}"

echo "Waiting for PostgreSQL at $POSTGRES_HOST:$POSTGRES_PORT..."

until PGPASSWORD=$POSTGRES_PASSWORD psql -h "$POSTGRES_HOST" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q' 2>/dev/null; do
    echo "PostgreSQL is unavailable yet - retry in 2s"
    sleep 2
done

echo "Successfully connetcted to PostgreSQL"

echo "Creating table 'prices' if not exists..."
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

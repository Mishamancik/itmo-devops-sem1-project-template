#!/bin/bash

set -e

echo "Preparing project-sem-1-db container with PostgreSQL..."

if ! command -v docker &> /dev/null; then
    echo "Docker not found on host"
    exit 1
fi

# Пока секреты явно в коде
if ! docker ps -a --format '{{.Names}}' | grep -q '^project-sem-1-db$'; then
    docker run --name project-sem-1-db \
      -e POSTGRES_USER=validator \
      -e POSTGRES_PASSWORD=val1dat0r \
      -e POSTGRES_DB=project-sem-1 \
      -p 5432:5432 \
      -d postgres:15
else
    echo "Container project-sem-1-db already exists"
    docker start project-sem-1-db
fi

sleep 5

# Пока пароль тут тоже явно в коде
PGPASSWORD=val1dat0r psql -h localhost -U validator -d project-sem-1 -c "
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    create_date DATE NOT NULL,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL
);
"

echo "Database is ready to work"

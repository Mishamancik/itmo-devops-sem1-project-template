#!/bin/bash

set -e

echo "Starting application in background..."
go run main.go &
APP_PID=$!

echo "Waiting for API to be ready..."

for i in {1..20}; do
  if curl -s http://localhost:8080/health >/dev/null 2>&1; then
    echo "API is ready"
    exit 0
  fi
  echo "API not ready yet, retrying..."
  sleep 1
done

echo "API did not start in time"
exit 1
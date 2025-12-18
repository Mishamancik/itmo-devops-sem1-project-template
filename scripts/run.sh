#!/bin/bash

set -e

echo "Starting application in background..."

go run main.go &

APP_PID=$!

echo "Application started with PID $APP_PID"

sleep 3
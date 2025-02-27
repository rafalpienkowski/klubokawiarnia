#!/bin/zsh

# Navigate to the directory containing podman-compose.yml
cd "$(dirname "$0")"

echo "🚀 Starting infrastructure"
podman-compose up -d

# Wait a bit to ensure Kafka starts properly
sleep 5

echo "✅ Infrastructure should be running!"
podman ps --format "table {{.ID}}\t{{.Names}}\t{{.Status}}"


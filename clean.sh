#!/bin/zsh

echo "🛑 Stopping..."
podman-compose down

echo "🧹 Removing resources..."
podman system prune -f
podman volume prune -f

echo "✅ Cleanup complete!"


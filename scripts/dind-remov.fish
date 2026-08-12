#!/usr/bin/env fish
# dind-teardown.fish
# Fully removes the isolated dind test container, its storage volume,
# and unsets DOCKER_HOST so this shell goes back to talking to your
# real host Docker daemon.
#
# Usage: source dind-teardown.fish
# (Must be `source`d, not run directly, so DOCKER_HOST is unset in
#  your current shell session.)

set CONTAINER_NAME dind
set VOLUME_NAME dind-storage

# IMPORTANT: unset DOCKER_HOST FIRST. If it's still pointed at the dind
# daemon itself, "docker rm dind" gets sent *into* the isolated daemon
# instead of to your real host daemon, and silently does nothing.
if set -q DOCKER_HOST
    echo "→ Unsetting DOCKER_HOST (was: $DOCKER_HOST)"
    set -e DOCKER_HOST
else
    echo "→ DOCKER_HOST was not set."
end

if docker ps -a --format '{{.Names}}' | grep -qx $CONTAINER_NAME
    echo "→ Stopping '$CONTAINER_NAME'..."
    docker stop $CONTAINER_NAME
    echo "→ Removing '$CONTAINER_NAME'..."
    docker rm $CONTAINER_NAME
else
    echo "→ No container named '$CONTAINER_NAME' found on host daemon, skipping."
end

if docker volume ls --format '{{.Name}}' | grep -qx $VOLUME_NAME
    echo "→ Removing volume '$VOLUME_NAME'..."
    docker volume rm $VOLUME_NAME
else
    echo "→ No volume named '$VOLUME_NAME' found on host daemon, skipping."
end

echo ""
echo "✓ Cleaned up. You're back on your host Docker daemon:"
docker ps -a

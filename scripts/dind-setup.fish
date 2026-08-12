#!/usr/bin/env fish
# dind-setup.fish
# Spins up an isolated Docker-in-Docker daemon for testing pulls/images
# with a completely empty, separate image cache from your host Docker.
#
# Usage: source dind-setup.fish
# (Must be `source`d, not run directly, so DOCKER_HOST persists in your
#  current shell session and your Go program picks it up automatically.)

set CONTAINER_NAME dind
set VOLUME_NAME dind-storage
set PORT 2375

echo "→ Checking for existing '$CONTAINER_NAME' container..."
if docker ps -a --format '{{.Names}}' | grep -qx $CONTAINER_NAME
    echo "  Found an existing '$CONTAINER_NAME' container. Removing it first..."
    docker rm -f $CONTAINER_NAME > /dev/null
end

echo "→ Starting isolated dind container..."
docker run -d --privileged --cgroupns=host --name $CONTAINER_NAME \
    -p $PORT:2375 \
    -e DOCKER_TLS_CERTDIR="" \
    -v $VOLUME_NAME:/var/lib/docker \
    docker:dind > /dev/null

if test $status -ne 0
    echo "✗ Failed to start dind container. Aborting."
    exit 1
end

# The inner dockerd takes ~15-20s to fully initialize (containerd,
# buildkit, and an intentional ~1s TLS-warning delay), so we poll
# with a generous timeout instead of assuming it's ready immediately.
echo "→ Waiting for inner Docker daemon to be ready (this can take ~15-20s)..."
set -l max_tries 60
set -l tries 0
while true
    set tries (math $tries + 1)
    if curl -s -o /dev/null http://localhost:$PORT/_ping
        echo "  Daemon is ready. (took ~"(math $tries \* 0.5)"s)"
        break
    end
    if test $tries -ge $max_tries
        echo "✗ Timed out waiting for dind daemon to respond after "(math $max_tries \* 0.5)"s."
        echo "  Check logs with: docker logs $CONTAINER_NAME"
        exit 1
    end
    sleep 0.5
end

set -x DOCKER_HOST tcp://localhost:$PORT

echo "→ Confirming isolated daemon is empty..."
docker images

echo ""
echo "✓ Isolated Docker environment is ready."
echo "  DOCKER_HOST is now set to $DOCKER_HOST for this shell session."
echo "  Run your program normally, e.g.:"
echo "    go run main.go init"
echo ""
echo "  When done, run: source dind-teardown.fish"

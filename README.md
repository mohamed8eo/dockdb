# DockDB

DockDB provisions and manages local PostgreSQL and MySQL Docker containers from
an interactive terminal UI or a non-interactive CLI.

## Requirements

- Docker Engine running and available to the current user
- `lazydocker` only when using `dockdb list --ui`

## Install

Download the archive matching your operating system and CPU from the GitHub
releases page, extract it, and put `dockdb` on your `PATH`.

## Usage

```sh
# Interactive setup
dockdb init

# Non-interactive PostgreSQL setup
dockdb init --db postgres --name pg-dev --port 5432 --password "$POSTGRES_PASSWORD"

# Manage containers by name or ID
dockdb list --all
dockdb up pg-dev
dockdb down pg-dev
dockdb delete pg-dev
```

`dockdb init` prints a connection URL containing the supplied password. Avoid
running it where terminal output is retained if that credential is sensitive.

## Releasing

The release workflow publishes on a pushed semantic-version tag such as
`v1.0.0`. It builds Linux, macOS, and Windows binaries for amd64 and arm64,
uploads archives plus `checksums.txt`, and creates GitHub artifact attestations.

```sh
git tag v1.0.0
git push origin v1.0.0
```

Run a local packaging check before tagging when GoReleaser is installed:

```sh
goreleaser check
goreleaser release --snapshot --clean
```

## Testing

Run the fast unit suite with `go test ./...`. Docker lifecycle integration tests
are opt-in because they create, start, stop, and remove a temporary container:

```sh
DOCKDB_RUN_INTEGRATION=1 go test -tags=integration ./internal/docker -run TestContainerLifecycleIntegration
```

# DockDB

DockDB provisions and manages local PostgreSQL and MySQL Docker containers from
an interactive terminal UI or a non-interactive CLI.

## Requirements

- Docker Engine running and available to the current user
- `lazydocker` only when using `dockdb list --ui`

## Install

Build from source with Go, or download a release archive when one is available.

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

GitHub Actions is intentionally not configured. To package a release locally,
install GoReleaser and run:

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

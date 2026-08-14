package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/mohamed8eo/dockdb/internal/logger"
	"github.com/mohamed8eo/dockdb/internal/render"
)

func ListContainer(ctx context.Context, cli *client.Client, all bool) error {
	containers, err := cli.ContainerList(
		ctx,
		client.ContainerListOptions{
			// INFO: true  => return all containers
			// INFO: false => return only the running containers
			All: all,
		},
	)
	if err != nil {
		return err
	}

	rows := make([]render.ContainerRow, 0, len(containers.Items))
	for _, c := range containers.Items {
		name := ""
		if len(c.Names) > 0 {
			name = strings.TrimPrefix(c.Names[0], "/")
		}

		rows = append(rows, render.ContainerRow{
			ID:      c.ID,
			Image:   c.Image,
			Command: c.Command,
			Created: c.Created,
			Status:  c.Status,
			Ports:   formatPorts(c.Ports), // see note below
			Name:    name,
			Running: c.State == "running",
		})
	}

	render.Containers(rows)
	return nil
}

func UpContainer(ctx context.Context, cli *client.Client, id string) error {
	id = strings.TrimSpace(id)

	if id == "" {
		return fmt.Errorf("container ID cannot be empty")
	}

	result, err := cli.ContainerInspect(
		ctx,
		id,
		client.ContainerInspectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to inspect container %q: %w", id, err)
	}

	if result.Container.State == nil {
		return fmt.Errorf("container %q has no state information", id)
	}

	name := strings.TrimPrefix(result.Container.Name, "/")
	state := result.Container.State

	if state.Running {
		logger.Info("Container is already running")
		return nil
	}

	if state.Dead {
		return fmt.Errorf("container %q is dead", name)
	}

	if err := startContainer(ctx, cli, name, id); err != nil {
		return fmt.Errorf("failed to start container %q: %w", name, err)
	}

	logger.Info("Container started", "name", name)

	return nil
}

func DownContainer(ctx context.Context, cli *client.Client, id string) error {
	id = strings.TrimSpace(id)

	if id == "" {
		return fmt.Errorf("container ID cannot be empty")
	}

	result, err := cli.ContainerInspect(
		ctx,
		id,
		client.ContainerInspectOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to inspect container %q: %w", id, err)
	}

	if result.Container.State == nil {
		return fmt.Errorf("container %q has no state information", id)
	}

	name := strings.TrimPrefix(result.Container.Name, "/")
	state := result.Container.State

	if !state.Running {
		logger.Info("Container is already Stoped")
		return nil
	}

	if state.Dead {
		return fmt.Errorf("container %q is dead", name)
	}

	if err := stopContainer(ctx, cli, name, id); err != nil {
		return fmt.Errorf("failed to stop container %q: %w", name, err)
	}

	logger.Info("Container stoped", "name", name)

	return nil
}

// formatPorts converts the SDK's port slice into docker-ps-style text
// like "0.0.0.0:5432->5432/tcp". Adjust field names (PublicPort,
// PrivatePort, IP, Type) to match whatever your client package's
// port struct actually calls them — check with `go doc` if unsure.
func formatPorts(ports []container.PortSummary) string {
	if len(ports) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(ports))
	for _, p := range ports {
		if p.PublicPort != 0 {
			parts = append(parts, fmt.Sprintf("%s:%d->%d/%s", p.IP, p.PublicPort, p.PrivatePort, p.Type))
		} else {
			parts = append(parts, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		}
	}
	return strings.Join(parts, ", ")
}

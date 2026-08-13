package docker

import (
	"context"
	"fmt"
	"net/netip"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
	"github.com/mohamed8eo/dockdb/internal/logger"
	"github.com/mohamed8eo/dockdb/internal/system"
	"github.com/pterm/pterm"
)

type ContainerSpec struct {
	Name        string
	Image       string
	Env         []string
	Labels      map[string]string
	ExposedPort string // container-side port, e.g. "5432/tcp"
	HostPort    string // host-side port, e.g. "5432"
	Restart     bool
}

func CreateAndStart(ctx context.Context, cli *client.Client, spec ContainerSpec) (containerID string, err error) {
	if err = ensureImage(ctx, cli, spec.Image); err != nil {
		return "", err
	}

	// Separate the completed image pull from the container lifecycle output.
	fmt.Println()

	exposedPort, err := network.ParsePort(spec.ExposedPort)
	if err != nil {
		logger.Error("invalid exposed port", "port", spec.ExposedPort, "error", err)
		return "", fmt.Errorf("invalid exposed port %q: %w", spec.ExposedPort, err)
	}

	resp, err := createContainer(ctx, cli, spec, exposedPort)
	if err != nil {
		return "", err
	}

	if err = startContainer(ctx, cli, spec.Name, resp.ID); err != nil {
		return "", err
	}

	printSummary(spec, resp.ID)
	return resp.ID, nil
}

func createContainer(
	ctx context.Context,
	cli *client.Client,
	spec ContainerSpec,
	exposedPort network.Port,
) (client.ContainerCreateResult, error) {
	spinner, _ := pterm.DefaultSpinner.
		WithText(fmt.Sprintf("Creating container %q...", spec.Name)).
		Start()

	containerConfi := &container.Config{
		Image:  spec.Image,
		Env:    spec.Env,
		Labels: spec.Labels,
		ExposedPorts: network.PortSet{
			exposedPort: struct{}{},
		},
	}

	var hostConfig *container.HostConfig
	if spec.Restart {
		enabled, err := system.IsDockerEnabledOnBoot()
		if err != nil {
			logger.Warn("could not check docker boot status", "error", err)
		}
		if !enabled {
			logger.Warn("docker is not enabled to start on boot; container restart policy won't survive a reboot ")
		}

		fmt.Println("Enable Docker in Boot")
		system.EnableDockerOnBoot()

		hostConfig = &container.HostConfig{
			PortBindings: network.PortMap{
				exposedPort: []network.PortBinding{
					{
						HostIP:   netip.MustParseAddr("0.0.0.0"),
						HostPort: spec.HostPort,
					},
				},
			},
			// auto restart the container after reboot
			RestartPolicy: container.RestartPolicy{
				Name: "always",
			},
		}
	} else {
		hostConfig = &container.HostConfig{
			PortBindings: network.PortMap{
				exposedPort: []network.PortBinding{
					{
						HostIP:   netip.MustParseAddr("0.0.0.0"),
						HostPort: spec.HostPort,
					},
				},
			},
			RestartPolicy: container.RestartPolicy{
				Name:              "on-failure",
				MaximumRetryCount: 2,
			},
		}
	}

	resp, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config:           containerConfi,
		HostConfig:       hostConfig,
		NetworkingConfig: &network.NetworkingConfig{},
		Name:             spec.Name,
	})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to create container %q: %v", spec.Name, err))
		return resp, fmt.Errorf("creating container %q: %w", spec.Name, err)
	}

	spinner.Success(fmt.Sprintf("Created container %q (%s)", spec.Name, shortID(resp.ID)))
	return resp, nil
}

func startContainer(
	ctx context.Context,
	cli *client.Client,
	name,
	id string,
) error {
	spinner, _ := pterm.DefaultSpinner.
		WithText(fmt.Sprintf("Starting container %q...", name)).
		Start()

	if _, err := cli.ContainerStart(ctx, id, client.ContainerStartOptions{}); err != nil {
		spinner.Fail(fmt.Sprintf("Failed to start container %q: %v", name, err))
		return fmt.Errorf("starting container %q: %w", name, err)
	}

	spinner.Success(fmt.Sprintf("Started container %q", name))
	return nil
}

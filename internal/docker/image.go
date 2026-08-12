package docker

import (
	"context"
	"fmt"

	"github.com/containerd/errdefs"
	"github.com/moby/moby/client"
	"github.com/pterm/pterm"
)

func ensureImage(ctx context.Context, cli *client.Client, image string) error {
	exists, err := imageExists(ctx, cli, image)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return pullImageWithProgress(ctx, cli, image)
}

func imageExists(ctx context.Context, cli *client.Client, image string) (bool, error) {
	_, err := cli.ImageInspect(ctx, image)
	if err == nil {
		return true, nil
	}
	if errdefs.IsNotFound(err) {
		return false, nil
	}
	return false, err
}

func pullImageWithProgress(ctx context.Context, cli *client.Client, image string) error {
	multi := pterm.DefaultMultiPrinter

	spinner, err := pterm.DefaultSpinner.
		WithText(fmt.Sprintf("Pulling image %q...", image)).
		Start()
	if err != nil {
		return fmt.Errorf("starting spinner: %w", err)
	}

	if _, err = multi.Start(); err != nil {
		spinner.Fail(fmt.Sprintf("Failed to start terminal UI: %v", err))
		return err
	}

	resp, err := cli.ImagePull(ctx, image, client.ImagePullOptions{})
	if err != nil {
		spinner.Fail(fmt.Sprintf("Failed to pull image %q: %v", image, err))
		return fmt.Errorf("pulling image %q: %w", image, err)
	}
	defer resp.Close()

	if err := parseProgress(resp, &multi); err != nil {
		spinner.Fail(fmt.Sprintf("Failed to parse progress for %q: %v", image, err))
		return err
	}

	spinner.Success(fmt.Sprintf("Pulled image %q", image))
	return nil
}

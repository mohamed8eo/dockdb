package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/errdefs"
	"github.com/moby/moby/client"
	"github.com/pterm/pterm"
)

func ensureImage(
	ctx context.Context,
	cli *client.Client,
	image string,
) error {
	multi := pterm.DefaultMultiPrinter.WithUpdateDelay(
		100 * time.Millisecond,
	)

	if _, err := multi.Start(); err != nil {
		return fmt.Errorf("starting terminal UI: %w", err)
	}
	defer multi.Stop()

	spinner, err := pterm.DefaultSpinner.
		WithWriter(multi.NewWriter()).
		WithText(fmt.Sprintf(
			"Checking for image %q...",
			image,
		)).Start()
	if err != nil {
		return fmt.Errorf("starting spinner: %w", err)
	}

	exists, err := imageExists(ctx, cli, image)
	if err != nil {
		spinner.Fail(
			fmt.Sprintf(
				"Failed to check image %q: %v",
				image,
				err,
			),
		)

		return err
	}

	if exists {
		spinner.Success(
			fmt.Sprintf(
				"Image %q already exists locally",
				image,
			),
		)

		return nil
	}

	pterm.Warning.Printfln(
		"Image %q not found locally",
		image,
	)

	pterm.Info.Printfln("Pulling %q image", image)

	return pullImage(
		ctx,
		cli,
		image,
		multi,
	)
}

func pullImage(
	ctx context.Context,
	cli *client.Client,
	image string,
	multi *pterm.MultiPrinter,
) error {
	reader, err := cli.ImagePull(
		ctx,
		image,
		client.ImagePullOptions{},
	)
	if err != nil {
		return fmt.Errorf(
			"pulling image %q: %w",
			image,
			err,
		)
	}
	defer reader.Close()

	return renderPullProgress(
		reader,
		image,
		multi,
	)
}

func imageExists(
	ctx context.Context,
	cli *client.Client,
	image string,
) (bool, error) {
	_, err := cli.ImageInspect(ctx, image)

	if err == nil {
		return true, nil
	}

	if errdefs.IsNotFound(err) {
		return false, nil
	}

	return false, err
}

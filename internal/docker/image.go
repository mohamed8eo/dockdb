package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/containerd/errdefs"
	"github.com/moby/moby/client"
	"github.com/pterm/pterm"
)

// minCheckSpinnerDuration keeps the "Checking for image" spinner
// visible for at least this long. The underlying check is a local
// Docker API call and normally resolves in a few milliseconds,
// which reads as if no check happened at all.
const minCheckSpinnerDuration = 350 * time.Millisecond

// waitForMinDuration sleeps just enough so that at least min has
// elapsed since start, doing nothing if it already has.
func waitForMinDuration(start time.Time, min time.Duration) {
	if remaining := min - time.Since(start); remaining > 0 {
		time.Sleep(remaining)
	}
}

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

	checkStarted := time.Now()

	exists, err := imageExists(ctx, cli, image)
	if err != nil {
		waitForMinDuration(checkStarted, minCheckSpinnerDuration)

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
		waitForMinDuration(checkStarted, minCheckSpinnerDuration)

		spinner.Success(
			fmt.Sprintf(
				"Image %q already exists locally",
				image,
			),
		)

		return nil
	}

	waitForMinDuration(checkStarted, minCheckSpinnerDuration)

	// Resolves the spinner into a warning line and removes the
	// spinner animation, the pull progress bars take over from
	// here in the same live area.
	spinner.Warning(
		fmt.Sprintf(
			"Image %q not found locally",
			image,
		),
	)

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
	spinner, err := pterm.DefaultSpinner.
		WithWriter(multi.NewWriter()).
		WithText(fmt.Sprintf(
			"Pulling %q image...",
			image,
		)).Start()
	if err != nil {
		return fmt.Errorf("starting spinner: %w", err)
	}

	reader, err := cli.ImagePull(
		ctx,
		image,
		client.ImagePullOptions{},
	)
	if err != nil {
		spinner.Fail(
			fmt.Sprintf(
				"Failed to pull image %q: %v",
				image,
				err,
			),
		)

		return fmt.Errorf(
			"pulling image %q: %w",
			image,
			err,
		)
	}
	defer reader.Close()

	spinner.Success(
		fmt.Sprintf(
			"Pulling %q image",
			image,
		),
	)

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

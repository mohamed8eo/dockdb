package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/pterm/pterm"
)

type resettingWriter struct {
	buf *bytes.Buffer
}

func (rw resettingWriter) Write(p []byte) (n int, err error) {
	rw.buf.Reset()
	return rw.buf.Write(p)
}

func makeResettingWriter(multi *pterm.MultiPrinter) io.Writer {
	rawWriter := multi.NewWriter()
	if buf, ok := rawWriter.(*bytes.Buffer); ok {
		return resettingWriter{buf: buf}
	}
	return rawWriter
}

// pullMessage represents a single JSON message returned by
// Docker's ImagePull API.
type pullMessage struct {
	ID     string `json:"id"`
	Status string `json:"status"`

	ProgressDetail struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`

	Progress string `json:"progress"`

	Error       string `json:"error"`
	ErrorDetail struct {
		Message string `json:"message"`
	} `json:"errorDetail"`
}

// renderPullProgress reads Docker's pull stream and converts
// each layer into a PTerm progress bar.
func renderPullProgress(
	reader io.Reader,
	image string,
	multi *pterm.MultiPrinter,
) error {
	progressbars := make(map[string]*pterm.ProgressbarPrinter)

	// Docker reports absolute byte values.
	// PTerm's Add() expects a delta, so we keep the previous
	// value for every layer.
	previous := make(map[string]int64)

	// Docker's "Pull complete" / "Already exists" messages
	// arrive with an empty progressDetail (no Total), so we
	// can't rely on the per-message total to finish the bar.
	// Remember each layer's Total when its bar is created.
	totals := make(map[string]int64)

	scanner := bufio.NewScanner(reader)

	// Docker normally sends small JSON lines, but increase the
	// scanner buffer so a large response cannot fail because
	// of Scanner's default 64 KB limit.
	scanner.Buffer(
		make([]byte, 64*1024),
		1024*1024,
	)

	for scanner.Scan() {
		var msg pullMessage

		if err := json.Unmarshal(
			scanner.Bytes(),
			&msg,
		); err != nil {
			return fmt.Errorf(
				"decoding pull response: %w",
				err,
			)
		}

		// Docker can send an error inside the pull stream.
		if msg.Error != "" {
			return fmt.Errorf(
				"docker pull failed: %s",
				msg.Error,
			)
		}

		if msg.ErrorDetail.Message != "" {
			return fmt.Errorf(
				"docker pull failed: %s",
				msg.ErrorDetail.Message,
			)
		}

		// Messages such as:
		//
		// Pulling from library/postgres
		// Digest: ...
		// Status: ...
		//
		// don't have a layer ID.
		if msg.ID == "" {
			continue
		}

		total := msg.ProgressDetail.Total

		// Create a progress bar when Docker tells us the total size
		// of the layer. Completed layers remain active in the live
		// area so every layer stays visible through the end of the pull.
		if _, exists := progressbars[msg.ID]; !exists &&
			total > 0 {
			pb, err := startLayerProgress(multi, msg.ID, total)
			if err != nil {
				return err
			}

			progressbars[msg.ID] = pb
			previous[msg.ID] = 0
			totals[msg.ID] = total
		}

		// Cached layers do not include a byte total, but they still deserve
		// a persistent completed row alongside the downloaded layers.
		if _, exists := progressbars[msg.ID]; !exists &&
			msg.Status == "Already exists" {
			pb, err := startLayerProgress(multi, msg.ID, 1)
			if err != nil {
				return err
			}

			progressbars[msg.ID] = pb
			totals[msg.ID] = 1
		}

		pb, exists := progressbars[msg.ID]
		if !exists {
			continue
		}

		current := msg.ProgressDetail.Current
		last := previous[msg.ID]

		// Docker sends the absolute current byte count.
		// PTerm.Add() needs the difference.
		if current > last {
			current = min(current, totals[msg.ID])
			delta := current - last

			if current < totals[msg.ID] && delta > int64(maxInt()) {
				return fmt.Errorf(
					"progress delta for layer %q is too large",
					msg.ID,
				)
			}

			if current < totals[msg.ID] {
				pb.Add(int(delta))
			} else {
				// PTerm stops and clears a bar when Add reaches its total.
				// Set Current directly instead so the 100% row stays in place.
				pb.Current = int(current)
			}

			previous[msg.ID] = current

			// Replace the raw byte counter with a human-readable
			// value such as:
			//
			// 5.00 MB / 5.68 MB
			pb.UpdateTitle(
				fmt.Sprintf(
					"%s [%s / %s]",
					shortLayerID(msg.ID),
					formatBytes(current),
					formatBytes(total),
				),
			)
		}

		switch msg.Status {
		case "Pull complete":
			completeProgressbar(
				pb,
				msg.ID,
				totals[msg.ID],
			)

		case "Already exists":
			completeProgressbar(
				pb,
				msg.ID,
				totals[msg.ID],
			)

		case "Download complete":
			// The layer has finished downloading, but Docker
			// may still be extracting it.
			//
			// Don't remove the progress bar yet.

		case "Extracting":
			// Docker is extracting the layer.
			//
			// Keep the progress bar visible.
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf(
			"reading pull output for %q: %w",
			image,
			err,
		)
	}

	// Some registries finish the stream without emitting a final
	// "Pull complete" event for every layer. A completed image is the
	// authoritative signal, so settle any remaining layer rows at 100%.
	completeRemainingProgressbars(progressbars, totals)

	return nil
}

func completeRemainingProgressbars(
	progressbars map[string]*pterm.ProgressbarPrinter,
	totals map[string]int64,
) {
	for id, pb := range progressbars {
		completeProgressbar(pb, id, totals[id])
	}
}

func startLayerProgress(
	multi *pterm.MultiPrinter,
	id string,
	total int64,
) (*pterm.ProgressbarPrinter, error) {
	if total > int64(maxInt()) {
		return nil, fmt.Errorf(
			"layer %q is too large for pterm progressbar",
			id,
		)
	}

	pb, err := pterm.DefaultProgressbar.
		WithTotal(int(total)).
		WithWriter(makeResettingWriter(multi)).
		WithShowElapsedTime(false).
		WithShowPercentage(true).
		WithShowCount(false).
		WithRemoveWhenDone(false).
		Start(layerProgressTitle(id, 0, total))
	if err != nil {
		return nil, fmt.Errorf(
			"creating progressbar for layer %q: %w",
			id,
			err,
		)
	}

	return pb, nil
}

// completeProgressbar makes sure the progress bar reaches 100%.
//
// Docker can occasionally send "Pull complete" without the
// immediately preceding Current value being the exact Total.
func completeProgressbar(
	pb *pterm.ProgressbarPrinter,
	id string,
	total int64,
) {
	if total <= 0 {
		return
	}

	// Do not use Add here: PTerm stops a progress bar when it reaches
	// its total, which removes the completed layer from MultiPrinter.
	pb.Current = int(total)
	pb.UpdateTitle(layerProgressTitle(id, total, total))
}

func layerProgressTitle(id string, current, total int64) string {
	return fmt.Sprintf(
		"%s [%s / %s]",
		shortLayerID(id),
		formatBytes(current),
		formatBytes(total),
	)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// shortLayerID keeps Docker's long layer ID from taking too
// much terminal space.
func shortLayerID(id string) string {
	const length = 12

	if len(id) <= length {
		return id
	}

	return id[:length]
}

// formatBytes converts Docker's byte counts into human-readable
// values.
//
// Examples:
//
//	1024       -> 1.00 KB
//	5242880    -> 5.00 MB
//	1073741824 -> 1.00 GB
func formatBytes(bytes int64) string {
	const (
		KB = int64(1024)
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf(
			"%.2f GB",
			float64(bytes)/float64(GB),
		)

	case bytes >= MB:
		return fmt.Sprintf(
			"%.2f MB",
			float64(bytes)/float64(MB),
		)

	case bytes >= KB:
		return fmt.Sprintf(
			"%.2f KB",
			float64(bytes)/float64(KB),
		)

	default:
		return fmt.Sprintf(
			"%d B",
			bytes,
		)
	}
}

// maxInt returns the maximum int value for the current
// architecture.
func maxInt() int {
	return int(^uint(0) >> 1)
}

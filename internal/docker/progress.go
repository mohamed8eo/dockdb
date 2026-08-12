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

		// Create a progress bar only when Docker tells us
		// the total size of the layer.
		if _, exists := progressbars[msg.ID]; !exists &&
			total > 0 {

			if total > int64(maxInt()) {
				return fmt.Errorf(
					"layer %q is too large for pterm progressbar",
					msg.ID,
				)
			}

			pb, err := pterm.DefaultProgressbar.
				WithTotal(int(total)).
				WithWriter(makeResettingWriter(multi)).
				WithShowElapsedTime(false).
				WithShowPercentage(true).
				WithShowCount(false).
				WithRemoveWhenDone(false).
				Start(
					fmt.Sprintf(
						"%s [%s / %s]",
						shortLayerID(msg.ID),
						formatBytes(0),
						formatBytes(total),
					),
				)
			if err != nil {
				return fmt.Errorf(
					"creating progressbar for layer %q: %w",
					msg.ID,
					err,
				)
			}

			progressbars[msg.ID] = pb
			previous[msg.ID] = 0
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
			delta := current - last

			if delta > int64(maxInt()) {
				return fmt.Errorf(
					"progress delta for layer %q is too large",
					msg.ID,
				)
			}

			pb.Add(int(delta))

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
				current,
				total,
			)

			delete(previous, msg.ID)

		case "Already exists":
			completeProgressbar(
				pb,
				current,
				total,
			)

			delete(previous, msg.ID)

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

	return nil
}

// completeProgressbar makes sure the progress bar reaches 100%.
//
// Docker can occasionally send "Pull complete" without the
// immediately preceding Current value being the exact Total.
func completeProgressbar(
	pb *pterm.ProgressbarPrinter,
	current int64,
	total int64,
) {
	if total <= 0 {
		return
	}

	if current < total {
		remaining := total - current

		if remaining <= int64(maxInt()) {
			pb.Add(int(remaining))
		}
	}

	pb.UpdateTitle(
		fmt.Sprintf(
			"%s [%s / %s]",
			"Layer",
			formatBytes(total),
			formatBytes(total),
		),
	)
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

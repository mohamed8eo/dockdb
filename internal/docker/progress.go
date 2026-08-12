package docker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pterm/pterm"
)

type jsonMessage struct {
	ID             string       `json:"id"`
	Status         string       `json:"status"`
	Progress       string       `json:"progress"`
	ProgressDetail jsonProgress `json:"progressDetail"`
	Error          string       `json:"error"`
	ErrorDetail    struct {
		Message string `json:"message"`
	} `json:"errorDetail"`
}

type jsonProgress struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

func parseProgress(r io.Reader, multi *pterm.MultiPrinter) error {
	scanner := bufio.NewScanner(r)
	progressbars := make(map[string]*pterm.ProgressbarPrinter)

	for scanner.Scan() {
		var msg jsonMessage
		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			return fmt.Errorf("decoding pull response: %w", err)
		}

		if msg.Error != "" {
			return fmt.Errorf("docker pull failed: %s", msg.Error)
		}
		if msg.ErrorDetail.Message != "" {
			return fmt.Errorf("docker pull failed: %s", msg.ErrorDetail.Message)
		}

		if msg.ID == "" {
			continue
		}

		pb, exists := progressbars[msg.ID]
		if !exists {
			var err error
			pb, err = pterm.DefaultProgressbar.
				WithTotal(int(msg.ProgressDetail.Total)).
				WithTitle(msg.ID).
				WithWriter(multi.NewWriter()).
				Start()
			if err != nil {
				return fmt.Errorf("starting progressbar: %w", err)
			}
			progressbars[msg.ID] = pb
		}

		if msg.ProgressDetail.Total > 0 {
			pb.Total = int(msg.ProgressDetail.Total)
			current := int(msg.ProgressDetail.Current)
			if current > pb.Total {
				current = pb.Total
			}
			pb.Current = current
			pb.UpdateTitle(fmt.Sprintf("%s: %s", msg.ID, msg.Status))
		} else {
			pb.UpdateTitle(fmt.Sprintf("%s: %s", msg.ID, msg.Status))
		}

		if strings.Contains(strings.ToLower(msg.Status), "download complete") ||
			strings.Contains(strings.ToLower(msg.Status), "extract complete") ||
			strings.Contains(strings.ToLower(msg.Status), "already exist") {
			if pb.Total > 0 {
				pb.Current = pb.Total
			}
			pb.Stop()
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading pull response: %w", err)
	}

	return nil
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	if bytes >= GB {
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	}
	if bytes >= MB {
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	}
	if bytes >= KB {
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	}
	return fmt.Sprintf("%d B", bytes)
}

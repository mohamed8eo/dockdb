package docker

import (
	"fmt"

	"github.com/pterm/pterm"
)

// printSummary prints a boxed summary of the created/started container.
func printSummary(spec ContainerSpec, id string) {
	pterm.DefaultBox.WithTitle("Container Ready").WithTitleTopCenter().Println(
		fmt.Sprintf(
			"Name:  %s\nImage: %s\nID:    %s\nPort:  %s -> %s",
			spec.Name, spec.Image, shortID(id), spec.HostPort, spec.ExposedPort,
		),
	)
}

// shortID trims a full container ID down to Docker's conventional
// 12-character short form for cleaner output.
func shortID(id string) string {
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

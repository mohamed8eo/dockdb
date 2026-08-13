package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// ContainerRow is a plain, client-agnostic view of a container
// used purely for rendering.
type ContainerRow struct {
	ID      string
	Image   string
	Command string
	Created int64  // unix timestamp
	Status  string // docker's raw status text, e.g. "Up 2 hours"
	Ports   string // already formatted, e.g. "0.0.0.0:5432->5432/tcp"
	Name    string
	Running bool
}

func Containers(rows []ContainerRow) {
	// newHeader := pterm.HeaderPrinter{
	// 	TextStyle:       pterm.NewStyle(pterm.FgLightBlue),
	// 	BackgroundStyle: pterm.NewStyle(),
	// 	Margin:          20,
	// }
	// newHeader.Println("🐳 Containers")
	fmt.Println()

	if len(rows) == 0 {
		pterm.Info.Println("No containers found.")
		return
	}

	// Column widths, minimums matching docker ps roughly
	idW, imgW, cmdW, createdW, statusW, portsW, nameW, stateW := 12, 5, 20, 14, 16, 5, 4, 9

	type prepped struct {
		id, image, command, created, status, ports, name string
		running                                          bool
	}
	prows := make([]prepped, 0, len(rows))

	for _, r := range rows {
		id := r.ID
		if len(id) > 12 {
			id = id[:12]
		}
		command := truncate(r.Command, 20)
		created := humanizeSince(r.Created)
		ports := r.Ports
		if ports == "" {
			ports = "-"
		}

		if len(id) > idW {
			idW = len(id)
		}
		if len(r.Image) > imgW {
			imgW = len(r.Image)
		}
		if len(command) > cmdW {
			cmdW = len(command)
		}
		if len(created) > createdW {
			createdW = len(created)
		}
		if len(r.Status) > statusW {
			statusW = len(r.Status)
		}
		if len(ports) > portsW {
			portsW = len(ports)
		}
		if len(r.Name) > nameW {
			nameW = len(r.Name)
		}

		prows = append(prows, prepped{
			id:      id,
			image:   r.Image,
			command: command,
			created: created,
			status:  r.Status,
			ports:   ports,
			name:    r.Name,
			running: r.Running,
		})
	}

	headerStyle := pterm.NewStyle(pterm.FgCyan, pterm.Bold)
	fmt.Println(headerStyle.Sprintf(
		"%-*s  %-*s  %-*s  %-*s  %-*s  %-*s  %-*s  %-*s",
		idW, "CONTAINER ID",
		imgW, "IMAGE",
		cmdW, "COMMAND",
		createdW, "CREATED",
		statusW, "STATUS",
		portsW, "PORTS",
		nameW, "NAMES",
		stateW, "STATE",
	))

	for _, r := range prows {
		var state string
		if r.running {
			state = pterm.Green("● RUNNING")
		} else {
			state = pterm.Red("○ STOPPED")
		}

		fmt.Printf(
			"%s  %-*s  %-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
			pterm.Cyan(fmt.Sprintf("%-*s", idW, r.id)),
			imgW, r.image,
			cmdW, r.command,
			createdW, r.created,
			statusW, r.status,
			portsW, r.ports,
			nameW, r.name,
			state,
		)
	}

	fmt.Println()
	pterm.Info.Printf("Have: %d container", len(rows))
}

// truncate shortens s to max chars, appending "…" if cut.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	if max <= 1 {
		return s[:max]
	}
	return s[:max-1] + "…"
}

// humanizeSince turns a unix timestamp into a docker-ps-style
// relative duration, e.g. "2 hours ago".
func humanizeSince(created int64) string {
	if created == 0 {
		return "-"
	}
	d := time.Since(time.Unix(created, 0))

	switch {
	case d < time.Minute:
		return "seconds ago"
	case d < time.Hour:
		m := int(d.Minutes())
		return pluralize(m, "minute") + " ago"
	case d < 24*time.Hour:
		h := int(d.Hours())
		return pluralize(h, "hour") + " ago"
	default:
		days := int(d.Hours() / 24)
		return pluralize(days, "day") + " ago"
	}
}

func pluralize(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", unit)
	}
	return fmt.Sprintf("%d %ss", n, unit)
}

// formatPorts is a helper you can call from the docker package if your
// port type isn't already a plain string. Adjust the type param to
// whatever your client package returns.
func FormatPorts(ports []string) string {
	if len(ports) == 0 {
		return "-"
	}
	return strings.Join(ports, ", ")
}

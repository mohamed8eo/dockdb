package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func IsDockerEnabledOnBoot() (bool, error) {
	out, err := exec.Command("systemctl", "is-enabled", "docker").Output()
	if err != nil {
		if strings.TrimSpace(string(out)) == "disabled" {
			return false, nil
		}
		return false, fmt.Errorf("checking docker enablement: %w", err)
	}

	status := strings.TrimSpace(string(out))
	return status == "enabled", nil
}

// EnableDockerOnBoot enables docker.service to start on boot. Requires root/sudo.
func EnableDockerOnBoot() error {
	cmd := exec.Command("sudo", "systemctl", "enable", "docker")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("enabling docker on boot: %w", err)
	}
	return nil
}

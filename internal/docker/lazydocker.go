package docker

import (
	"fmt"
	"os"
	"os/exec"
)

// RunLazydocker launches lazydocker as an interactive TUI attached to the
// current terminal. Users must install lazydocker through their trusted package
// manager before invoking this command.
func RunLazydocker() error {
	if !lazydockerInstalled() {
		return fmt.Errorf("lazydocker is not installed; install it from https://github.com/jesseduffield/lazydocker#installation")
	}
	return execLazydocker()
}

// lazydockerInstalled reports whether the lazydocker binary is on PATH.
func lazydockerInstalled() bool {
	_, err := exec.LookPath("lazydocker")
	return err == nil
}

// execLazydocker runs lazydocker with the current process's stdin/stdout/stderr
// attached, so the TUI renders directly in the user's terminal.
func execLazydocker() error {
	cmd := exec.Command("lazydocker")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

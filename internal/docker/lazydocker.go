package docker

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// RunLazydocker checks if lazydocker is installed, installs it if missing,
// then launches it as an interactive TUI attached to the current terminal.
func RunLazydocker() error {
	if !lazydockerInstalled() {
		fmt.Println("lazydocker not found, installing...")
		if err := installLazydocker(); err != nil {
			return fmt.Errorf("failed to install lazydocker: %w", err)
		}
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

// installLazydocker installs lazydocker using the best available method
// for the current OS.
func installLazydocker() error {
	switch runtime.GOOS {
	case "linux", "darwin":
		return installUnix()
	case "windows":
		return installWindows()
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// installUnix runs the official install script (works for both Linux and macOS).
func installUnix() error {
	script := "curl https://raw.githubusercontent.com/jesseduffield/lazydocker/master/scripts/install_update_linux.sh | bash"
	cmd := exec.Command("bash", "-c", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// installWindows tries scoop first, falling back to choco if scoop isn't available.
func installWindows() error {
	if _, err := exec.LookPath("scoop"); err == nil {
		cmd := exec.Command("scoop", "install", "lazydocker")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	if _, err := exec.LookPath("choco"); err == nil {
		cmd := exec.Command("choco", "install", "lazydocker", "-y")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	return fmt.Errorf("neither scoop nor choco found; install one first, or install lazydocker manually: https://github.com/jesseduffield/lazydocker#installation")
}

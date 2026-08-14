package cmd

import (
	"os"

	"github.com/mohamed8eo/dockdb/internal/logger"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dockdb",
	Short: "DockDB is a CLI tool to easily provision and manage local database containers",
	Long: `DockDB is a modern, lightweight command-line utility for spinning up,
managing, and inspecting local database containers (PostgreSQL, MySQL) via Docker
with both interactive TUI and headless CLI workflows.`,
	Run: func(cmd *cobra.Command, args []string) { _ = cmd.Help() },
}

// Version is set at build time by GoReleaser.
var Version = "dev"

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Version
}

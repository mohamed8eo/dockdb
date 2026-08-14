package cmd

import (
	"errors"

	"github.com/mohamed8eo/dockdb/internal/app"
	"github.com/mohamed8eo/dockdb/internal/tui"
	"github.com/spf13/cobra"
)

var (
	name     string
	dbType   string
	port     string
	password string
	restart  bool
)

var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"create", "new"},
	Short:   "Initialize and start a new database container",
	Long: `Initialize a database by creating and starting a Docker container.
Supports PostgreSQL and MySQL.

If run without flags, launches an interactive configuration wizard (TUI).
When flags are provided, runs in headless/CI mode.`,
	Example: `  # Launch interactive TUI wizard
  dockdb init

  # Quick start PostgreSQL
  dockdb init --db postgres --name pg-db --port 5432 --password secret

  # Quick start MySQL with auto-restart enabled
  dockdb init --db mysql --name mysql-db --port 3306 --password secret --restart`,

	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := app.ResolveConfig(cmd, restart, dbType, name, password, port)
		if err != nil {
			if errors.Is(err, tui.ErrPromptCancelled) {
				return nil
			}
			return err
		}

		return app.Start(cmd.Context(), cfg)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&name, "name", "n", "postgresql", "Container name")
	initCmd.Flags().StringVarP(&port, "port", "p", "5432", "Database port")
	initCmd.Flags().StringVarP(&password, "password", "e", "postgresql", "Database root password")
	initCmd.Flags().StringVarP(&dbType, "db", "d", "postgres", "Database type (postgres or mysql)")
	initCmd.Flags().BoolVarP(&restart, "restart", "r", false, "Enable automatic container restart on reboot")
}

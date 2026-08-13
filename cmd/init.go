/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/mohamed8eo/dockdb/internal/app"
	"github.com/mohamed8eo/dockdb/internal/logger"
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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create and start a database container",
	Long: `Initialize a database by creating and starting a Docker container.

Supported databases:
  postgres, mysql`,

	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := app.ResolveConfig(cmd, restart, dbType, name, password, port)
		if err != nil {
			if errors.Is(err, tui.ErrPromptCancelled) {
				return
			}
			logger.Error("failed to configure database", "error", err)
			return
		}

		app.Start(cfg)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags
	initCmd.Flags().StringVarP(&name, "name", "n", "postgresql", "db name")
	initCmd.Flags().StringVarP(&port, "port", "p", "5432", "db port")
	initCmd.Flags().StringVarP(&password, "password", "e", "postgresql", "db password")
	initCmd.Flags().StringVarP(&dbType, "db", "t", "postgres", "database type (postgres or mysql)")
	initCmd.Flags().BoolVarP(&restart, "restart", "r", false, "Rerun db after reboot")
}

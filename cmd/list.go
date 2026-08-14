package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamed8eo/dockdb/internal/docker"
	"github.com/spf13/cobra"
)

var (
	all bool
	ui  bool
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all database containers",
	Long:    `List all provisioned database containers (running and stopped), or launch the interactive TUI dashboard.`,
	Example: `  dockdb list
  dockdb ls --all
  dockdb list --ui`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !ui {
			ctx, cancel := context.WithTimeout(cmd.Context(), 10*time.Second)
			defer cancel()
			cli, err := docker.NewClient()
			if err != nil {
				return fmt.Errorf("create Docker client: %w", err)
			}
			defer cli.Close()
			return docker.ListContainer(ctx, cli, all)
		}
		return docker.RunLazydocker()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers (including stopped ones)")
	listCmd.Flags().BoolVarP(&ui, "ui", "i", false, "Launch interactive lazydocker TUI dashboard")
}

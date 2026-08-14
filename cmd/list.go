package cmd

import (
	"context"
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
	Run: func(cmd *cobra.Command, args []string) {
		if !ui {
			ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
			defer cancel()
			cli, _ := docker.NewClient()
			defer cli.Close()
			docker.ListContainer(ctx, cli, all)
			return
		}
		docker.RunLazydocker()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&all, "all", "a", false, "Show all containers (including stopped ones)")
	listCmd.Flags().BoolVarP(&ui, "ui", "i", false, "Launch interactive lazydocker TUI dashboard")
}

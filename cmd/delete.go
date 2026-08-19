package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamed8eo/dockdb/internal/docker"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <container>...",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete database containers",
	Long:    `Permanently remove one or more Docker database containers and their associated resources by ID or name.`,
	Example: `  dockdb delete postgresql
  dockdb rm postgresql mysql`,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(
			cmd.Context(),
			10*time.Second,
		)
		defer cancel()

		cli, err := docker.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Docker client: %w", err)
		}
		defer cli.Close()
		if err := docker.DeleteContainer(ctx, cli, args[0:]); err != nil {
			return fmt.Errorf("failed to delete container: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

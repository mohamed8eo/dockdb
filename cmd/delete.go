package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamed8eo/dockdb/internal/docker"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <container>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a database container",
	Long:    `Permanently remove a Docker database container and its associated resources by ID or name.`,
	Example: `  dockdb delete postgresql
  dockdb rm postgresql`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(
			context.TODO(),
			1*time.Second,
		)
		defer cancel()

		cli, err := docker.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Docker client: %w", err)
		}
		defer cli.Close()
		if err := docker.DeleteContainer(ctx, cli, args[0]); err != nil {
			return fmt.Errorf("failed to delete container: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

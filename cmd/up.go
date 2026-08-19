package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamed8eo/dockdb/internal/docker"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:     "up <container>...",
	Aliases: []string{"start"},
	Short:   "Start a stopped database container",
	Long:    `Start an existing, stopped Docker database container by its ID or name.`,
	Example: `  dockdb up postgresql`,
	Args:    cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		cli, err := docker.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create Docker client: %w", err)
		}
		defer cli.Close()

		ctx, cancel := context.WithTimeout(
			cmd.Context(),
			10*time.Second,
		)
		defer cancel()

		if err := docker.UpContainer(ctx, cli, args[0:]); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}

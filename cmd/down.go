/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamed8eo/dockdb/internal/docker"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down <container ID>",
	Short: "Stop container",
	Args:  cobra.ExactArgs(1),
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

		if err := docker.DownContainer(ctx, cli, args[0]); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}

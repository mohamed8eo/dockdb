/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List directory contents",
	Run: func(cmd *cobra.Command, args []string) {
		if !ui {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			cli, _ := docker.NewClient()
			docker.ListContainer(ctx, cli, all)
			defer cancel()
			return
		}
		docker.RunLazydocker()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&all, "all", "a", false, "show hidden entries")
	listCmd.Flags().BoolVarP(&ui, "ui", "i", false, "show TUI Lazy Docker")
}

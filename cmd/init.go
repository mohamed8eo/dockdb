/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamed8eo/dockdb/internal/database"
	"github.com/mohamed8eo/dockdb/internal/docker"
	"github.com/mohamed8eo/dockdb/internal/logger"
	"github.com/spf13/cobra"
)

var (
	name     string
	dbType   string
	port     string
	password string
	detached bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create and start a database container",
	Long: `Initialize a database by creating and starting a Docker container.

Supported databases:
  postgres, pgsql, pg
  mysql

The database runs in detached mode by default only when the --detached flag
is provided.

Examples:
  # Create a PostgreSQL container
  dockdb init --db postgres

  # Create a PostgreSQL container using an alias
  dockdb init --db pg

  # Create a MySQL container
  dockdb init --db mysql

  # Customize the container
  dockdb init --db postgres --name my-db --port 5432 --password secret

  # Run the container in the background
  dockdb init --db postgres --detached`,
	Run: func(cmd *cobra.Command, args []string) {
		parsedDBType, err := database.ParseType(dbType)
		if err != nil {
			logger.Error("failed to parse database type", "error", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
		defer cancel()

		cli, err := docker.NewClient()
		if err != nil {
			logger.Fatal("failed to connect to docker", "error", err)
		}
		defer cli.Close()

		cfg := database.Config{
			Name:     name,
			Type:     parsedDBType,
			Port:     port,
			Password: password,
			Detached: detached,
		}

		id, err := docker.CreateAndStart(ctx, cli, cfg.ToContainerSpec())
		if err != nil {
			logger.Fatal("failed to start database", "error", err)
		}

		fmt.Printf("Started %s container: %s\n", parsedDBType, id[:12])
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Flags
	initCmd.Flags().StringVarP(&name, "name", "n", "postgresql", "db name")
	initCmd.Flags().StringVarP(&port, "port", "p", "5432", "db port")
	initCmd.Flags().StringVarP(&password, "password", "e", "postgresql", "db password")
	initCmd.Flags().StringVarP(&dbType, "db", "t", "postgres", "database type (postgres or mysql)")
	initCmd.Flags().BoolVarP(&detached, "detached", "d", false, "db detached")
}

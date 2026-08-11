/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mohamed8eo/dockdb/internal/database"
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
  mycli init --db postgres

  # Create a PostgreSQL container using an alias
  mycli init --db pg

  # Create a MySQL container
  mycli init --db mysql

  # Customize the container
  mycli init --db postgres --name my-db --port 5432 --password secret

  # Run the container in the background
  mycli init --db postgres --detached`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := database.ParseType(dbType)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("init called")
		fmt.Println(name)
		fmt.Println(port)
		fmt.Println(password)
		fmt.Println(detached)
		fmt.Println(db)

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	//Flags
	initCmd.Flags().StringVarP(&name, "name", "n", "postgresql", "db name")
	initCmd.Flags().StringVarP(&port, "port", "p", "5432", "db port")
	initCmd.Flags().StringVarP(&password, "password", "e", "postgresql", "db password")
	initCmd.Flags().StringVarP(&dbType, "db", "t", "postgres", "database type (postgres or mysql)")
	//run on backgroun
	initCmd.Flags().BoolVarP(&detached, "detached", "d", false, "db detached")
}

package app

import (
	"fmt"
	"strconv"

	"github.com/mohamed8eo/dockdb/internal/database"
	"github.com/mohamed8eo/dockdb/internal/render"
	"github.com/mohamed8eo/dockdb/internal/tui"
	"github.com/spf13/cobra"
)

// ResolveConfig determines the database configuration either from the
// interactive TUI (when no flags were passed) or from CLI flags.
func ResolveConfig(cmd *cobra.Command, restart bool, dbType, name, password, port string) (database.Config, error) {
	// TUI
	if cmd.Flags().NFlag() == 0 {
		render.PrintBanner("DOCKERDB")

		answers, err := tui.Create()
		if err != nil {
			return database.Config{}, err
		}

		// Leave a visual boundary between the completed form and Docker's output.
		fmt.Println()

		parsedDBType, err := database.ParseType(answers.DBType)
		if err != nil {
			return database.Config{}, err
		}

		return database.Config{
			Name:     answers.Name,
			Type:     parsedDBType,
			Port:     strconv.Itoa(answers.Port),
			Password: answers.Password,
			Restart:  restart,
		}, nil
	}

	// CLI FLAGS
	parsedDBType, err := database.ParseType(dbType)
	if err != nil {
		return database.Config{}, err
	}

	return database.Config{
		Name:     name,
		Type:     parsedDBType,
		Port:     port,
		Password: password,
		Restart:  restart,
	}, nil
}

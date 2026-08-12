package ui

import (
	"fmt"
	"strconv"

	"github.com/pterm/pterm"
)

type CreateConfig struct {
	Name     string
	Port     int
	Password string
	DBType   string
	Detach   bool
}

func Create() (CreateConfig, error) {
	name, err := pterm.DefaultInteractiveTextInput.
		WithDefaultValue("my-database").
		Show("Database name")
	if err != nil {
		return CreateConfig{}, err
	}

	portString, err := pterm.DefaultInteractiveTextInput.
		WithDefaultValue("5432").
		Show("Port")
	if err != nil {
		return CreateConfig{}, err
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return CreateConfig{}, fmt.Errorf("invalid port: %w", err)
	}

	passwordInput := pterm.DefaultInteractiveTextInput.WithMask("*")
	passwordInput.Show("Password")
	if err != nil {
		return CreateConfig{}, err
	}

	dbType, err := pterm.DefaultInteractiveSelect.
		WithOptions(
			[]string{
				"PostgreSQL",
				"MySQL",
			},
		).
		Show("Database type")
	if err != nil {
		return CreateConfig{}, err
	}

	detach, err := pterm.DefaultInteractiveConfirm.
		WithDefaultValue(true).
		Show("Run container in detached mode?")
	if err != nil {
		return CreateConfig{}, err
	}

	return CreateConfig{
		Name:     name,
		Port:     port,
		Password: passwordInput.DefaultText,
		DBType:   dbType,
		Detach:   detach,
	}, nil
}

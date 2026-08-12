package ui

import (
	"errors"
	"fmt"
	"strconv"
)

type CreateConfig struct {
	Name     string
	Port     int
	Password string
	DBType   string
	Detach   bool
}

func Create() (CreateConfig, error) {
	result, err := newCreateForm().run()
	if err != nil {
		return CreateConfig{}, err
	}
	if result.cancelled {
		return CreateConfig{}, ErrPromptCancelled
	}

	port, err := strconv.Atoi(result.port())
	if err != nil {
		return CreateConfig{}, fmt.Errorf("invalid port: %w", err)
	}

	return CreateConfig{
		Name:     result.name(),
		Port:     port,
		Password: result.inputs[2].Value(),
		DBType:   result.dbType(),
		Detach:   result.detached(),
	}, nil
}

var ErrPromptCancelled = errors.New("setup cancelled")

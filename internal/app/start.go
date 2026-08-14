package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamed8eo/dockdb/internal/database"
	"github.com/mohamed8eo/dockdb/internal/docker"
)

// Start creates and starts the database container described by cfg,
// then prints the resulting connection URL.
func Start(ctx context.Context, cfg database.Config) error {
	ctx, cancel := context.WithTimeout(ctx, 240*time.Second)
	defer cancel()

	cli, err := docker.NewClient()
	if err != nil {
		return fmt.Errorf("connect to Docker: %w", err)
	}
	defer cli.Close()

	_, err = docker.CreateAndStart(ctx, cli, cfg.ToContainerSpec())
	if err != nil {
		return fmt.Errorf("start database: %w", err)
	}

	// Keep the final command result visually distinct from the summary card.
	fmt.Println()
	fmt.Printf("DBURL= \"%s\"\n", cfg.DBURL())
	return nil
}

package app

import (
	"context"
	"fmt"
	"time"

	"github.com/mohamed8eo/dockdb/internal/database"
	"github.com/mohamed8eo/dockdb/internal/docker"
	"github.com/mohamed8eo/dockdb/internal/logger"
)

// Start creates and starts the database container described by cfg,
// then prints the resulting connection URL.
func Start(cfg database.Config) {
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()

	cli, err := docker.NewClient()
	if err != nil {
		logger.Fatal("failed to connect to docker", "error", err)
	}
	defer cli.Close()

	_, err = docker.CreateAndStart(ctx, cli, cfg.ToContainerSpec())
	if err != nil {
		logger.Fatal("failed to start database", "error", err)
	}

	// Keep the final command result visually distinct from the summary card.
	fmt.Println()
	fmt.Printf("DBURL= \"%s\"\n", cfg.DBURL())
}

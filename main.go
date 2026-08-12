/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"github.com/mohamed8eo/dockdb/cmd"
	"github.com/mohamed8eo/dockdb/internal/logger"
)

func main() {
	logger.Init()
	cmd.Execute()
}

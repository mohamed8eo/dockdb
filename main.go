/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/mohamed8eo/dockdb/cmd"
	"github.com/mohamed8eo/dockdb/internal/logger"
	"github.com/mohamed8eo/dockdb/internal/render"
)

func main() {
	render.ApplyTheme()
	logger.Init()
	cmd.Execute()
}

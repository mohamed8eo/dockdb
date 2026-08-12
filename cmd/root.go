/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/mohamed8eo/dockdb/internal/logger"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dockdb",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		printDynamicBanner("DOCKERDB")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logger.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dockdb.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func printDynamicBanner(text string) {
	bigText, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromString(text),
	).Srender()

	lines := strings.Split(bigText, "\n")

	// Exact Laravel gradient colors (top to bottom)
	colors := []string{
		"#06B6D4", // Cyan
		"#38BDF8", // Light Blue
		"#3B82F6", // Blue
		"#2563EB", // Dark Blue
		"#7C3AED", // Purple
	}

	fmt.Println()
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		colorHex := colors[i%len(colors)]
		rgb, err := pterm.NewRGBFromHEX(colorHex)
		if err == nil {
			fmt.Println(rgb.Sprint(line))
		} else {
			fmt.Println(line)
		}
	}
	fmt.Println()
}



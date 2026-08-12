/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/mohamed8eo/dockdb/internal/logger"
	"github.com/mohamed8eo/dockdb/internal/ui"
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

		fmt.Println("\n\n")

		config, err := ui.Create()
		if err != nil {
			logger.Error("Error: %q", err)
		}

		fmt.Println()
		fmt.Println("Configuration:")
		fmt.Printf("Name: %s\n", config.Name)
		fmt.Printf("Port: %d\n", config.Port)
		fmt.Printf("Password: %s\n", config.Password)
		fmt.Printf("Database: %s\n", config.DBType)
		fmt.Printf("Detached: %t\n", config.Detach)
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

// bannerThemes are Laravel-style top-to-bottom gradients. One is
// picked at random each run so the banner isn't always the same
// color.
var bannerThemes = [][]string{
	{ // Cyan -> Purple
		"#06B6D4",
		"#38BDF8",
		"#3B82F6",
		"#2563EB",
		"#7C3AED",
	},
	{ // Red -> Maroon
		"#FCA5A5",
		"#F87171",
		"#EF4444",
		"#B91C1C",
		"#7F1D1D",
	},
	{ // Mint -> Teal
		"#6EE7B7",
		"#34D399",
		"#10B981",
		"#059669",
		"#065F46",
	},
	{ // Orange -> Brown (Sunset)
		"#FDBA74",
		"#FB923C",
		"#F97316",
		"#C2410C",
		"#7C2D12",
	},
	{ // Pink -> Magenta -> Deep Purple
		"#F9A8D4",
		"#F472B6",
		"#EC4899",
		"#BE185D",
		"#831843",
	},
	{ // Gold -> Amber -> Brown
		"#FDE68A",
		"#FBBF24",
		"#F59E0B",
		"#B45309",
		"#78350F",
	},
}

func printDynamicBanner(text string) {
	bigText, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromString(text),
	).Srender()

	lines := strings.Split(bigText, "\n")

	// Pick a random gradient theme for this run.
	colors := bannerThemes[rand.Intn(len(bannerThemes))]

	fmt.Println()
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// BigText's own default style is already baked into the
		// line as ANSI codes, which would override our gradient
		// color if we just wrapped the line as-is. Strip it first
		// so only our color applies.
		plainLine := pterm.RemoveColorFromString(line)

		colorHex := colors[i%len(colors)]
		rgb, err := pterm.NewRGBFromHEX(colorHex)
		if err == nil {
			fmt.Println(rgb.Sprint(plainLine))
		} else {
			fmt.Println(plainLine)
		}
	}
	fmt.Println()
}

package render

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

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

// PrintBanner renders the DockDB banner using a randomly picked gradient theme.
func PrintBanner(text string) {
	bigText, _ := pterm.DefaultBigText.WithLetters(
		putils.LettersFromString(text),
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
		rgb, err := putils.RGBFromHEX(colorHex)
		if err == nil {
			fmt.Println(rgb.Sprint(plainLine))
		} else {
			fmt.Println(plainLine)
		}
	}
	fmt.Println()
}

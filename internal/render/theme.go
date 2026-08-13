package render

import "github.com/pterm/pterm"

// ApplyTheme installs the palette before DockDB renders terminal output.
func ApplyTheme() {
	dockdbTheme := pterm.ThemeDefault.
		WithPrimaryStyle(*pterm.NewStyle(pterm.FgLightCyan)).
		WithSecondaryStyle(*pterm.NewStyle(pterm.FgCyan)).
		WithHighlightStyle(*pterm.NewStyle(pterm.FgYellow, pterm.Bold)).
		WithInfoMessageStyle(*pterm.NewStyle(pterm.FgLightCyan)).
		WithInfoPrefixStyle(*pterm.NewStyle(pterm.FgBlack, pterm.BgCyan)).
		WithSuccessMessageStyle(*pterm.NewStyle(pterm.FgLightGreen)).
		WithSuccessPrefixStyle(*pterm.NewStyle(pterm.FgBlack, pterm.BgLightGreen)).
		WithWarningMessageStyle(*pterm.NewStyle(pterm.FgLightYellow)).
		WithWarningPrefixStyle(*pterm.NewStyle(pterm.FgBlack, pterm.BgLightYellow)).
		WithErrorMessageStyle(*pterm.NewStyle(pterm.FgLightRed)).
		WithErrorPrefixStyle(*pterm.NewStyle(pterm.FgWhite, pterm.BgRed, pterm.Bold)).
		WithFatalMessageStyle(*pterm.NewStyle(pterm.FgLightRed)).
		WithFatalPrefixStyle(*pterm.NewStyle(pterm.FgWhite, pterm.BgRed, pterm.Bold))

	dockdbTheme.ProgressbarBarStyle = *pterm.NewStyle(pterm.FgLightCyan)
	dockdbTheme.ProgressbarTitleStyle = *pterm.NewStyle(pterm.FgLightCyan)
	pterm.ThemeDefault = dockdbTheme
}

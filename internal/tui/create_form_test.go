package tui

import (
	"testing"

	"charm.land/lipgloss/v2"
)

func TestActionButtonsKeepTheSameSizeWhenSelected(t *testing.T) {
	normalWidth, normalHeight := lipgloss.Size(buttonStyle.Render("Submit"))
	selectedWidth, selectedHeight := lipgloss.Size(selectedButtonStyle.Render("Submit"))
	cancelWidth, cancelHeight := lipgloss.Size(cancelButtonStyle.Render("Cancel"))

	if normalWidth != selectedWidth || normalHeight != selectedHeight {
		t.Fatalf("submit size changed from %dx%d to %dx%d", normalWidth, normalHeight, selectedWidth, selectedHeight)
	}
	if normalWidth != cancelWidth || normalHeight != cancelHeight {
		t.Fatalf("cancel size changed from %dx%d to %dx%d", normalWidth, normalHeight, cancelWidth, cancelHeight)
	}
}

func TestActionButtonsRenderSideBySide(t *testing.T) {
	form := createForm{active: actionField}
	if got, want := lipgloss.Height(form.actionsView()), 4; got != want {
		t.Fatalf("action buttons rendered across %d lines; want %d", got, want)
	}
}

package docker

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pterm/pterm"
)

func TestCompleteProgressbarKeepsCompletedLayerVisible(t *testing.T) {
	var output bytes.Buffer
	pb, err := pterm.DefaultProgressbar.
		WithTotal(100).
		WithWriter(&output).
		WithShowElapsedTime(false).
		Start("downloading")
	if err != nil {
		t.Fatalf("starting progress bar: %v", err)
	}
	defer pb.Stop()

	completeProgressbar(pb, "layer-123", 100)

	if !pb.IsActive {
		t.Fatal("completed progress bar was stopped and would disappear")
	}
	if pb.Current != 100 {
		t.Fatalf("progress bar current = %d; want 100", pb.Current)
	}
	if !strings.Contains(output.String(), "100 B / 100 B") {
		t.Fatalf("completed progress was not rendered at 100%%: %q", output.String())
	}
}

func TestCompleteRemainingProgressbarsFinishesEveryLayer(t *testing.T) {
	first := newTestProgressbar(t, 100)
	second := newTestProgressbar(t, 200)

	first.Current = 88
	second.Current = 150

	completeRemainingProgressbars(
		map[string]*pterm.ProgressbarPrinter{
			"first-layer":  first,
			"second-layer": second,
		},
		map[string]int64{
			"first-layer":  100,
			"second-layer": 200,
		},
	)

	if first.Current != 100 || second.Current != 200 {
		t.Fatalf("unfinished layers were not completed: %d, %d", first.Current, second.Current)
	}
}

func newTestProgressbar(t *testing.T, total int) *pterm.ProgressbarPrinter {
	t.Helper()

	pb, err := pterm.DefaultProgressbar.
		WithTotal(total).
		WithWriter(&bytes.Buffer{}).
		WithShowElapsedTime(false).
		Start("downloading")
	if err != nil {
		t.Fatalf("starting progress bar: %v", err)
	}
	t.Cleanup(func() { _, _ = pb.Stop() })

	return pb
}

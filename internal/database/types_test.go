package database

import (
	"strings"
	"testing"
)

func TestParseTypeExplainsUnsupportedDatabase(t *testing.T) {
	_, err := ParseType("sqlite")
	if err == nil {
		t.Fatal("ParseType() returned nil error for an unsupported database")
	}

	if got := err.Error(); !strings.Contains(got, `unsupported database "sqlite"`) ||
		!strings.Contains(got, "supported values: postgres, mysql") {
		t.Fatalf("unexpected error: %q", got)
	}
}

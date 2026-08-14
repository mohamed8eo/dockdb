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

func TestParseTypeAcceptsSupportedAliases(t *testing.T) {
	tests := map[string]DBType{
		"postgres":   Postgres,
		"PostgreSQL": Postgres,
		"pg":         Postgres,
		"mysql":      MySQL,
		"MariaDB":    MySQL,
	}

	for input, want := range tests {
		t.Run(input, func(t *testing.T) {
			got, err := ParseType(input)
			if err != nil {
				t.Fatalf("ParseType(%q) returned an error: %v", input, err)
			}
			if got != want {
				t.Fatalf("ParseType(%q) = %q, want %q", input, got, want)
			}
		})
	}
}

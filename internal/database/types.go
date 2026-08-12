package database

import (
	"fmt"
	"strings"
)

type DBType string

const (
	Postgres DBType = "postgres"
	MySQL    DBType = "mysql"
)

func ParseType(s string) (DBType, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "postgres", "postgresql", "pg", "pgsql":
		return Postgres, nil
	case "mysql", "mariadb":
		return MySQL, nil
	default:
		return "", fmt.Errorf("unsupported database: %s", s)
	}
}

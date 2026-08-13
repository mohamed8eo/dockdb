package database

import (
	"fmt"

	"github.com/mohamed8eo/dockdb/internal/docker"
)

type Config struct {
	Name     string
	Type     DBType
	Port     string
	Password string
	Restart  bool
}

func (c Config) ToContainerSpec() docker.ContainerSpec {
	switch c.Type {
	case Postgres:
		return c.postgresSpec()
	case MySQL:
		return c.mySQLSpec()
	default:
		return docker.ContainerSpec{}
	}
}

func (c Config) DBURL() string {
	var dbURL string
	switch c.Type {
	case Postgres:
		dbURL = fmt.Sprintf("%s://postgres:%s@localhost:%s/%s", PostgresImage, c.Password, c.Port, PostgresImage)
		return dbURL
	case MySQL:
		dbURL = fmt.Sprintf("%s://root:%s@localhost:%s/%s", MySqlImage, c.Password, c.Port, MySqlImage)
		return dbURL

	default:
		return ""
	}
}

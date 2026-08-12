package database

import "github.com/mohamed8eo/dockdb/internal/docker"

type Config struct {
	Name     string
	Type     DBType
	Port     string
	Password string
	Detached bool
}

func (c Config) ToContainerSpec() docker.ContainerSpec {
	switch c.Type {
	case Postgres:
		return c.postgresSpec()
	case MySQL:
		return c.mySqlSpac()
	default:
		return docker.ContainerSpec{}
	}
}

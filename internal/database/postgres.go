package database

import "github.com/mohamed8eo/dockdb/internal/docker"

const (
	PostgresImage       = "postgres"
	postgresDefaultPort = "5432/tcp"
)

func (c Config) postgresSpec() docker.ContainerSpec {
	return docker.ContainerSpec{
		Name:        c.Name,
		Image:       PostgresImage,
		ExposedPort: postgresDefaultPort,
		HostPort:    c.Port,

		Env: []string{
			"POSTGRES_PASSWORD=" + c.Password,
		},
		Labels: map[string]string{
			"managed-by": "dockdb",
			"db-type":    "postgres",
		},
		Restart: c.Restart,
	}
}

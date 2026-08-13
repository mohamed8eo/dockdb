package database

import "github.com/mohamed8eo/dockdb/internal/docker"

const (
	MySqlImage       = "mysql"
	MySqlDefaultPort = "3306/tcp"
)

func (c Config) mySQLSpec() docker.ContainerSpec {
	return docker.ContainerSpec{
		Name:        c.Name,
		Image:       MySqlImage,
		ExposedPort: MySqlDefaultPort,
		HostPort:    c.Port,

		Env: []string{
			"MYSQL_ROOT_PASSWORD=" + c.Password,
		},
		Labels: map[string]string{
			"managed-by": "dockdb",
			"db-type":    "mysql",
		},
		Restart: c.Restart,
	}
}

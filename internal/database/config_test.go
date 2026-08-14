package database

import (
	"reflect"
	"testing"

	"github.com/mohamed8eo/dockdb/internal/docker"
)

func TestConfigToContainerSpec(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want docker.ContainerSpec
	}{
		{
			name: "postgres",
			cfg:  Config{Name: "pg-dev", Type: Postgres, Port: "5433", Password: "secret", Restart: true},
			want: docker.ContainerSpec{
				Name: "pg-dev", Image: PostgresImage, ExposedPort: "5432/tcp", HostPort: "5433", Restart: true,
				Env:    []string{"POSTGRES_PASSWORD=secret"},
				Labels: map[string]string{"managed-by": "dockdb", "db-type": "postgres"},
			},
		},
		{
			name: "mysql",
			cfg:  Config{Name: "mysql-dev", Type: MySQL, Port: "3307", Password: "secret"},
			want: docker.ContainerSpec{
				Name: "mysql-dev", Image: MySqlImage, ExposedPort: "3306/tcp", HostPort: "3307",
				Env:    []string{"MYSQL_ROOT_PASSWORD=secret"},
				Labels: map[string]string{"managed-by": "dockdb", "db-type": "mysql"},
			},
		},
		{name: "unsupported database", cfg: Config{Type: "unknown"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.ToContainerSpec(); !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("ToContainerSpec() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestConfigDBURL(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want string
	}{
		{"postgres", Config{Type: Postgres, Port: "5433", Password: "secret"}, "postgres://postgres:secret@localhost:5433/postgres"},
		{"mysql", Config{Type: MySQL, Port: "3307", Password: "secret"}, "mysql://root:secret@localhost:3307/mysql"},
		{"unsupported database", Config{Type: "unknown"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.DBURL(); got != tt.want {
				t.Fatalf("DBURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

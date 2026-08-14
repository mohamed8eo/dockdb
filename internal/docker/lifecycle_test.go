package docker

import (
	"net/netip"
	"testing"

	"github.com/moby/moby/api/types/container"
)

func TestFormatPorts(t *testing.T) {
	tests := []struct {
		name  string
		ports []container.PortSummary
		want  string
	}{
		{name: "no ports", want: "-"},
		{
			name:  "published port",
			ports: []container.PortSummary{{IP: netip.MustParseAddr("0.0.0.0"), PublicPort: 5432, PrivatePort: 5432, Type: "tcp"}},
			want:  "0.0.0.0:5432->5432/tcp",
		},
		{
			name:  "unpublished port",
			ports: []container.PortSummary{{PrivatePort: 3306, Type: "tcp"}},
			want:  "3306/tcp",
		},
		{
			name: "multiple ports",
			ports: []container.PortSummary{
				{IP: netip.MustParseAddr("127.0.0.1"), PublicPort: 5432, PrivatePort: 5432, Type: "tcp"},
				{PrivatePort: 5433, Type: "tcp"},
			},
			want: "127.0.0.1:5432->5432/tcp, 5433/tcp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatPorts(tt.ports); got != tt.want {
				t.Fatalf("formatPorts() = %q, want %q", got, tt.want)
			}
		})
	}
}

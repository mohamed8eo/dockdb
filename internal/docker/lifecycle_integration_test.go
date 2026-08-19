//go:build integration

package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/mohamed8eo/dockdb/internal/logger"
)

func TestContainerLifecycleIntegration(t *testing.T) {
	if os.Getenv("DOCKDB_RUN_INTEGRATION") != "1" {
		t.Skip("set DOCKDB_RUN_INTEGRATION=1 to run Docker integration tests")
	}

	logger.Init()
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	cli, err := NewClient()
	if err != nil {
		t.Fatalf("create Docker client: %v", err)
	}
	defer cli.Close()

	if _, err := cli.Ping(ctx, client.PingOptions{}); err != nil {
		t.Fatalf("Docker daemon is unavailable: %v", err)
	}

	image := os.Getenv("DOCKDB_TEST_IMAGE")
	if image == "" {
		image = "alpine:3.21"
	}
	ensureTestImage(t, ctx, cli, image)

	name := fmt.Sprintf("dockdb-integration-%d", time.Now().UnixNano())
	created, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: image,
			Cmd:   []string{"sh", "-c", "sleep 300"},
		},
		Name: name,
	})
	if err != nil {
		t.Fatalf("create test container: %v", err)
	}
	t.Cleanup(func() {
		_, _ = cli.ContainerRemove(context.Background(), created.ID, client.ContainerRemoveOptions{Force: true})
	})

	if err := UpContainer(ctx, cli, []string{created.ID}); err != nil {
		t.Fatalf("start test container: %v", err)
	}
	assertContainerRunning(t, ctx, cli, created.ID, true)

	if err := DownContainer(ctx, cli, []string{created.ID}); err != nil {
		t.Fatalf("stop test container: %v", err)
	}
	assertContainerRunning(t, ctx, cli, created.ID, false)

	if err := DeleteContainer(ctx, cli, []string{created.ID}); err != nil {
		t.Fatalf("delete test container: %v", err)
	}
	if _, err := cli.ContainerInspect(ctx, created.ID, client.ContainerInspectOptions{}); err == nil {
		t.Fatal("test container still exists after deletion")
	}
}

func ensureTestImage(t *testing.T, ctx context.Context, cli *client.Client, image string) {
	t.Helper()
	if exists, err := imageExists(ctx, cli, image); err != nil {
		t.Fatalf("inspect test image: %v", err)
	} else if exists {
		return
	}

	reader, err := cli.ImagePull(ctx, image, client.ImagePullOptions{})
	if err != nil {
		t.Fatalf("pull test image %q: %v", image, err)
	}
	defer reader.Close()
	if _, err := io.Copy(io.Discard, reader); err != nil {
		t.Fatalf("read image-pull output: %v", err)
	}
}

func assertContainerRunning(t *testing.T, ctx context.Context, cli *client.Client, id string, want bool) {
	t.Helper()
	inspected, err := cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
	if err != nil {
		t.Fatalf("inspect test container: %v", err)
	}
	if inspected.Container.State == nil || inspected.Container.State.Running != want {
		t.Fatalf("container running = %v, want %v", inspected.Container.State != nil && inspected.Container.State.Running, want)
	}
}

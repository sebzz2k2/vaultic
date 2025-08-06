//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Helper to run a command in the container and get output
func execCommand(ctx context.Context, container testcontainers.Container, cmd []string) (string, error) {
	exitCode, stdout, stderr, err := container.Exec(ctx, cmd)
	if err != nil {
		return "", err
	}
	if exitCode != 0 {
		return stderr.String(), nil
	}
	return stdout.String(), nil
}

func TestVaulticIntegration(t *testing.T) {
	ctx := context.Background()

	containerReq := testcontainers.ContainerRequest{
		Image:        "golang:1.21",             // Use official Go image, mount local vaultic binary
		Cmd:          []string{"sleep", "3600"}, // Keep container alive
		WaitingFor:   wait.ForListeningPort("80").WithStartupTimeout(2 * time.Second),
		AutoRemove:   true,
		ExposedPorts: []string{"80"},
		Mounts: testcontainers.Mounts(
			testcontainers.BindMount("/vaultic/vaultic", "/vaultic/vaultic"), // Mount built binary
		),
		Entrypoint: []string{"/bin/sh", "-c"},
		Env:        map[string]string{},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer func() {
		_ = container.Terminate(ctx)
	}()

	// Helper to run vaultic commands
	run := func(args ...string) string {
		cmd := append([]string{"/vaultic/vaultic"}, args...)
		out, err := execCommand(ctx, container, cmd)
		if err != nil {
			t.Fatalf("Failed to exec %v: %v", args, err)
		}
		return strings.TrimSpace(out)
	}

	// 1. No vaultic file
	if _, err := execCommand(ctx, container, []string{"rm", "-f", "/vaultic/vaultic"}); err != nil {
		t.Fatalf("Failed to remove vaultic file: %v", err)
	}
	t.Run("No vaultic file", func(t *testing.T) {
		if got := run("get", "a"); got != "(nil)" {
			t.Errorf("expected (nil), got %q", got)
		}
		if got := run("set", "a", "b"); !strings.Contains(got, "OK") {
			t.Errorf("expected OK, got %q", got)
		}
		if got := run("get", "a"); got != "b" {
			t.Errorf("expected b, got %q", got)
		}
	})

	// 2. Restart (simulate by stopping/starting container)
	_ = container.StopLogProducer()
	_ = container.Stop(ctx, nil)
	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to restart container: %v", err)
	}

	t.Run("After restart, vaultic file present", func(t *testing.T) {
		if got := run("get", "a"); got != "b" {
			t.Errorf("expected b, got %q", got)
		}
		if got := run("set", "a", "x"); !strings.Contains(got, "OK") {
			t.Errorf("expected OK, got %q", got)
		}
		if got := run("set", "b", "e"); !strings.Contains(got, "OK") {
			t.Errorf("expected OK, got %q", got)
		}
		if got := run("get", "a"); got != "x" {
			t.Errorf("expected x, got %q", got)
		}
		if got := run("get", "b"); got != "e" {
			t.Errorf("expected e, got %q", got)
		}
	})

	// 3. Restart again
	_ = container.StopLogProducer()
	_ = container.Stop(ctx, nil)
	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to restart container: %v", err)
	}

	t.Run("After restart, del b", func(t *testing.T) {
		if got := run("get", "a"); got != "x" {
			t.Errorf("expected x, got %q", got)
		}
		if got := run("get", "b"); got != "e" {
			t.Errorf("expected e, got %q", got)
		}
		if got := run("del", "b"); !strings.Contains(got, "OK") {
			t.Errorf("expected OK, got %q", got)
		}
		if got := run("get", "b"); got != "(nil)" {
			t.Errorf("expected (nil), got %q", got)
		}
	})

	// 4. Final restart
	_ = container.StopLogProducer()
	_ = container.Stop(ctx, nil)
	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to restart container: %v", err)
	}

	t.Run("After restart, keys and exists", func(t *testing.T) {
		if got := run("get", "b"); got != "(nil)" {
			t.Errorf("expected (nil), got %q", got)
		}
		if got := run("keys"); got != "a" {
			t.Errorf("expected a, got %q", got)
		}
		if got := run("exists", "a"); got != "true" {
			t.Errorf("expected true, got %q", got)
		}
		if got := run("exists", "b"); got != "false" {
			t.Errorf("expected false, got %q", got)
		}
	})
}

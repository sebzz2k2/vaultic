package integration

// function to create build test container with no-cache flag and run it. It should be buillt from docker-compose file
import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestContainer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	composeFile := "../../test-container.compose.yml"
	stack, err := compose.NewDockerComposeWith(
		compose.StackIdentifier("test-vaultic"),
		compose.WithStackFiles(composeFile),
	)
	if err != nil {
		panic(err)
	}

	err = stack.WaitForService("vaultic",
		wait.ForAll(
			wait.ForExposedPort(),
			wait.ForLog("Starting Vaultic server"),
			wait.ForLog("Building index"),
			wait.ForLog("Finished building index"),
		)).Up(ctx, compose.Wait(true))
	if err != nil {
		t.Fatalf("Failed to start test container: %v", err)
	}
	defer func() {
		err = stack.Down(
			context.Background(),
			compose.RemoveOrphans(true),
			compose.RemoveVolumes(true),
			compose.RemoveImagesLocal,
		)
		if err != nil {
			t.Fatalf("Failed to tear down test container: %v", err)
		}
	}()
	serviceNames := stack.Services()

	if len(serviceNames) == 0 {
		t.Fatal("No services found in the test container stack")
	}

	fmt.Println("Test container started successfully with services:")

	// Test with the vaultic service
	service := "vaultic"
	ctr, err := stack.ServiceContainer(ctx, service)
	if err != nil {
		t.Fatalf("Failed to get container for service %s: %v", service, err)
	}

	// Helper function to execute commands and get response
	execCommand := func(cmd string) string {
		_, reader, err := ctr.Exec(ctx, []string{"/bin/sh", "-c", cmd})
		require.NoError(t, err)

		buf := new(strings.Builder)
		_, err = io.Copy(buf, reader)
		require.NoError(t, err)

		return strings.TrimSpace(buf.String())
	}

	// Helper function to restart the container
	restartContainer := func() {
		err := ctr.Stop(ctx, nil)
		require.NoError(t, err)

		err = ctr.Start(ctx)
		require.NoError(t, err)

		// Wait for service to be ready again
		err = stack.WaitForService("vaultic",
			wait.ForAll(
				wait.ForExposedPort(),
				wait.ForLog("Starting Vaultic server"),
				wait.ForLog("Building index"),
				wait.ForLog("Finished building index"),
			)).Up(ctx, compose.Wait(true))
		require.NoError(t, err)
	}

	// Helper functions for common operations
	get := func(key string) string {
		return execCommand(fmt.Sprintf(`echo "GET %s" | nc -w 1 localhost 5381`, key))
	}

	set := func(key, value string) string {
		return execCommand(fmt.Sprintf(`echo "SET %s %s" | nc -w 1 localhost 5381`, key, value))
	}

	del := func(key string) string {
		return execCommand(fmt.Sprintf(`echo "DEL %s" | nc -w 1 localhost 5381`, key))
	}

	keys := func() string {
		return execCommand(`echo "KEYS" | nc -w 1 localhost 5381`)
	}

	exists := func(key string) string {
		return execCommand(fmt.Sprintf(`echo "EXISTS %s" | nc -w 1 localhost 5381`, key))
	}

	// Test initial state
	require.Contains(t, get("a"), "(nil)")
	require.Contains(t, set("a", "b"), "OK")
	require.Contains(t, get("a"), "b")

	// Test persistence after restart
	restartContainer()
	require.Contains(t, get("a"), "b")

	// Update values
	require.Contains(t, set("a", "x"), "OK")
	require.Contains(t, set("b", "e"), "OK")
	require.Contains(t, get("a"), "x")
	require.Contains(t, get("b"), "e")

	// Test persistence after second restart
	restartContainer()
	require.Contains(t, get("a"), "x")
	require.Contains(t, get("b"), "e")

	// Test deletion
	require.Contains(t, del("b"), "OK")
	require.Contains(t, get("b"), "(nil)")

	// Test deletion persistence after restart
	restartContainer()
	require.Contains(t, get("b"), "(nil)")
	require.Contains(t, keys(), "a")
	require.NotContains(t, keys(), "b")
	require.Contains(t, exists("a"), "true")
	require.Contains(t, exists("b"), "false")
}

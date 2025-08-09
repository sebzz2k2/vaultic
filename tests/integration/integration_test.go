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
		t.Log("Restarting container...")
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

	t.Log("=== Phase 1: Initial state without vaultic file ===")

	// Get value of variable a => should return (nil)
	response := execCommand(`echo "GET a" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "(nil)", "Expected GET a to return (nil)")

	// Set value of variable a as b => should return OK
	response = execCommand(`echo "SET a b" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "OK", "Expected SET a b to return OK")

	// Get value of variable a => should return b
	response = execCommand(`echo "GET a" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "b", "Expected GET a to return b")

	t.Log("=== Phase 2: First restart - vaultic file should be present ===")
	restartContainer()

	// Verify vaultic file is present
	response = execCommand(`ls -la | grep vaultic`)
	require.NotEmpty(t, response, "Expected vaultic file to be present after restart")

	// Get value of variable a => should return b
	response = execCommand(`echo "GET a" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "b", "Expected GET a to return b after restart")

	// Set value of variable a as x => should return OK
	response = execCommand(`echo "SET a x" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "OK", "Expected SET a x to return OK")

	// Set value of variable b as e => should return OK
	response = execCommand(`echo "SET b e" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "OK", "Expected SET b e to return OK")

	// Get value of variable a => should return x
	response = execCommand(`echo "GET a" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "x", "Expected GET a to return x")

	// Get value of variable b => should return e
	response = execCommand(`echo "GET b" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "e", "Expected GET b to return e")

	t.Log("=== Phase 3: Second restart - verify persistence ===")
	restartContainer()

	// Get value of variable a => should return x
	response = execCommand(`echo "GET a" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "x", "Expected GET a to return x after second restart")

	// Get value of variable b => should return e
	response = execCommand(`echo "GET b" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "e", "Expected GET b to return e after second restart")

	// Delete b => should return OK
	response = execCommand(`echo "DEL b" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "OK", "Expected DEL b to return OK")

	// Get value of variable b => should return (nil)
	response = execCommand(`echo "GET b" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "(nil)", "Expected GET b to return (nil) after deletion")

	t.Log("=== Phase 4: Final restart - verify deletion persistence ===")
	restartContainer()

	// Get value of variable b => should return (nil)
	response = execCommand(`echo "GET b" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "(nil)", "Expected GET b to return (nil) after final restart")

	// KEYS should return a
	response = execCommand(`echo "KEYS" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "a", "Expected KEYS to return a")
	require.NotContains(t, response, "b", "Expected KEYS to not contain b")

	// EXISTS a should return true
	response = execCommand(`echo "EXISTS a" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "true", "Expected EXISTS a to return true")

	// EXISTS b should return false
	response = execCommand(`echo "EXISTS b" | nc -w 1 localhost 5381`)
	require.Contains(t, response, "false", "Expected EXISTS b to return false")

	t.Log("=== All tests completed successfully ===")
}

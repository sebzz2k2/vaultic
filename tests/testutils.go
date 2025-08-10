package tests

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestStack wraps a Docker Compose stack for testing
type TestStack struct {
	Stack     *compose.DockerCompose
	Container testcontainers.Container
	ctx       context.Context
}

// SetupTestContainer creates and starts a test container stack
func SetupTestContainer(t *testing.T) (*TestStack, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	composeFile := "../../test-container.compose.yml"
	stack, err := compose.NewDockerComposeWith(
		compose.StackIdentifier("test-vaultic"),
		compose.WithStackFiles(composeFile),
	)
	if err != nil {
		cancel()
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
		cancel()
		t.Fatalf("Failed to start test container: %v", err)
	}

	serviceNames := stack.Services()
	if len(serviceNames) == 0 {
		cancel()
		t.Fatal("No services found in the test container stack")
	}

	fmt.Println("Test container started successfully with services:")

	// Get the vaultic service container
	ctr, err := stack.ServiceContainer(ctx, "vaultic")
	if err != nil {
		cancel()
		t.Fatalf("Failed to get container for service vaultic: %v", err)
	}

	testStack := &TestStack{
		Stack:     stack,
		Container: ctr,
		ctx:       ctx,
	}

	// Cleanup function
	cleanup := func() {
		cancel()
		err := stack.Down(
			context.Background(),
			compose.RemoveOrphans(true),
			compose.RemoveVolumes(true),
			compose.RemoveImagesLocal,
		)
		if err != nil {
			t.Fatalf("Failed to tear down test container: %v", err)
		}
	}

	return testStack, cleanup
}

// ExecCommand executes a command in the test container and returns the output
func (ts *TestStack) ExecCommand(cmd string) (string, error) {
	_, reader, err := ts.Container.Exec(ts.ctx, []string{"/bin/sh", "-c", cmd})
	if err != nil {
		return "", err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), nil
}

// ExecCommandWithRequire executes a command and fails the test if there's an error
func (ts *TestStack) ExecCommandWithRequire(t *testing.T, cmd string) string {
	result, err := ts.ExecCommand(cmd)
	require.NoError(t, err)
	return result
}

// RestartContainer restarts the container and waits for it to be ready
func (ts *TestStack) RestartContainer(t *testing.T) {
	err := ts.Container.Stop(ts.ctx, nil)
	require.NoError(t, err)

	err = ts.Container.Start(ts.ctx)
	require.NoError(t, err)

	// Wait for service to be ready again
	err = ts.Stack.WaitForService("vaultic",
		wait.ForAll(
			wait.ForLog("Starting Vaultic server"),
			wait.ForLog("Building index"),
			wait.ForLog("Finished building index"),
		)).Up(ts.ctx, compose.Wait(true))
	require.NoError(t, err)
}

// VaulticCommands provides helper methods for common Vaultic operations
type VaulticCommands struct {
	testStack *TestStack
}

// NewVaulticCommands creates a new VaulticCommands instance
func NewVaulticCommands(testStack *TestStack) *VaulticCommands {
	return &VaulticCommands{testStack: testStack}
}

// Get executes a GET command
func (vc *VaulticCommands) Get(key string) string {
	result, _ := vc.testStack.ExecCommand(fmt.Sprintf(`echo "GET %s" | nc -w 1 localhost 5381`, key))
	return result
}

// Set executes a SET command
func (vc *VaulticCommands) Set(key, value string) string {
	result, _ := vc.testStack.ExecCommand(fmt.Sprintf(`echo "SET %s %s" | nc -w 1 localhost 5381`, key, value))
	return result
}

// Del executes a DEL command
func (vc *VaulticCommands) Del(key string) string {
	result, _ := vc.testStack.ExecCommand(fmt.Sprintf(`echo "DEL %s" | nc -w 1 localhost 5381`, key))
	return result
}

// Keys executes a KEYS command
func (vc *VaulticCommands) Keys() string {
	result, _ := vc.testStack.ExecCommand(`echo "KEYS" | nc -w 1 localhost 5381`)
	return result
}

// Exists executes an EXISTS command
func (vc *VaulticCommands) Exists(key string) string {
	result, _ := vc.testStack.ExecCommand(fmt.Sprintf(`echo "EXISTS %s" | nc -w 1 localhost 5381`, key))
	return result
}

// BenchmarkHelper provides utilities for benchmark testing
type BenchmarkHelper struct {
	testStack *TestStack
}

// NewBenchmarkHelper creates a new BenchmarkHelper instance
func NewBenchmarkHelper(testStack *TestStack) *BenchmarkHelper {
	return &BenchmarkHelper{testStack: testStack}
}

// ExecuteAsync executes a command asynchronously with concurrency control
func (bh *BenchmarkHelper) ExecuteAsync(command string, wg *sync.WaitGroup, semaphore chan struct{}) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Acquire semaphore (limit concurrency)
		semaphore <- struct{}{}
		defer func() { <-semaphore }() // Release semaphore

		bh.testStack.ExecCommand(command)
	}()
}

package integration

import (
	"testing"

	"github.com/sebzz2k2/vaultic/tests"
	"github.com/stretchr/testify/require"
)

func TestContainer(t *testing.T) {
	testStack, cleanup := tests.SetupTestContainer(t)
	defer cleanup()

	// Create command helpers
	vaultic := tests.NewVaulticCommands(testStack)

	// Test initial state
	require.Contains(t, vaultic.Get("a"), "(nil)")
	require.Contains(t, vaultic.Set("a", "b"), "OK")
	require.Contains(t, vaultic.Get("a"), "b")

	// Test persistence after restart
	testStack.RestartContainer(t)
	require.Contains(t, vaultic.Get("a"), "b")

	// Update values
	require.Contains(t, vaultic.Set("a", "x"), "OK")
	require.Contains(t, vaultic.Set("b", "e"), "OK")
	require.Contains(t, vaultic.Get("a"), "x")
	require.Contains(t, vaultic.Get("b"), "e")

	// Test persistence after second restart
	testStack.RestartContainer(t)
	require.Contains(t, vaultic.Get("a"), "x")
	require.Contains(t, vaultic.Get("b"), "e")

	// Test deletion
	require.Contains(t, vaultic.Del("b"), "OK")
	require.Contains(t, vaultic.Get("b"), "(nil)")

	// Test deletion persistence after restart
	testStack.RestartContainer(t)
	require.Contains(t, vaultic.Get("b"), "(nil)")
	require.Contains(t, vaultic.Keys(), "a")
	require.NotContains(t, vaultic.Keys(), "b")
	require.Contains(t, vaultic.Exists("a"), "true")
	require.Contains(t, vaultic.Exists("b"), "false")
}

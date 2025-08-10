package benchmark

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sebzz2k2/vaultic/tests"
)

func TestContainer(t *testing.T) {
	testStack, cleanup := tests.SetupTestContainer(t)
	defer cleanup()

	// Create benchmark helper
	benchHelper := tests.NewBenchmarkHelper(testStack)

	// Helper function to run benchmark for any operation
	runBenchmark := func(t *testing.T, operationName string, keys int, maxConcurrency int, commandGenerator func(int) string) {
		var wg sync.WaitGroup
		semaphore := make(chan struct{}, maxConcurrency)

		start := time.Now()
		for i := 0; i < keys; i++ {
			if i%10000 == 0 {
				t.Logf("%s key %d", operationName, i)
			}
			command := commandGenerator(i)
			benchHelper.ExecuteAsync(command, &wg, semaphore)
		}
		sendDuration := time.Since(start)

		t.Logf("All %s commands sent in %v", operationName, sendDuration)
		t.Logf("Waiting for all %s operations to complete (max %d concurrent)...", operationName, maxConcurrency)

		// Wait for all operations to complete
		waitStart := time.Now()
		wg.Wait()
		waitDuration := time.Since(waitStart)
		totalDuration := time.Since(start)

		// Calculate ops/s
		totalOpsPerSecond := float64(keys) / totalDuration.Seconds()
		sendOpsPerSecond := float64(keys) / sendDuration.Seconds()

		t.Logf("%s commands sent in: %v", operationName, sendDuration)
		t.Logf("Wait time for completion: %v", waitDuration)
		t.Logf("Total time: %v", totalDuration)
		t.Logf("Send ops/s: %.2f", sendOpsPerSecond)
		t.Logf("Total %s ops/s: %.2f", strings.ToLower(operationName), totalOpsPerSecond)
	}

	//  need to know how much time it takes to set 100 keys
	t.Run("Benchmark Set 100 Keys", func(t *testing.T) {
		KEYS := 100
		maxConcurrency := 100
		runBenchmark(t, "Setting", KEYS, maxConcurrency, func(i int) string {
			return fmt.Sprintf(`echo "SET key%d value" | socat - TCP:localhost:5381`, i)
		})
	})

	// Benchmark getting multiple keys with parallelism
	t.Run("Benchmark Get 100 Keys", func(t *testing.T) {
		KEYS := 100
		maxConcurrency := 100
		runBenchmark(t, "Getting", KEYS, maxConcurrency, func(i int) string {
			return fmt.Sprintf(`echo "GET key%d" | socat - TCP:localhost:5381`, i)
		})
	})

	// Benchmark deleting multiple keys with parallelism
	t.Run("Benchmark Delete 100 Keys", func(t *testing.T) {
		KEYS := 100
		maxConcurrency := 100
		runBenchmark(t, "Deleting", KEYS, maxConcurrency, func(i int) string {
			return fmt.Sprintf(`echo "DEL key%d" | socat - TCP:localhost:5381`, i)
		})
	})
}

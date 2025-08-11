package kvstore

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSkipListConcurrent(t *testing.T) {
	const goroutines = 50
	const operations = 5000

	sl := NewSkipList(10)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func(id int) {
			defer wg.Done()
			for i := 0; i < operations; i++ {
				key := fmt.Sprintf("key-%d-%d", id, i)

				switch i % 3 {
				case 0:

					sl.Insert(uint64(time.Now().UnixNano()), true, key, string(fmt.Sprintf("value-%d-%d", id, i)))
				case 1:
					sl.Get(key)
				case 2:
					sl.Delete(key, uint64(time.Now().UnixNano()))
				}
			}
		}(g)
	}

	wg.Wait()
}

func TestSkipListFuzzConcurrent(t *testing.T) {
	sl := NewSkipList(10)
	const goroutines = 100
	const opsPerGoroutine = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				op := rand.Intn(3)
				key := fmt.Sprintf("key-%d", rand.Intn(500))
				switch op {
				case 0:
					sl.Insert(uint64(time.Now().UnixNano()), true, key, string(fmt.Sprintf("value-%d-%d", g, i)))
				case 1:
					sl.Get(key)
				case 2:
					sl.Delete(key, uint64(time.Now().UnixNano()))
				}
			}
		}()
	}
	wg.Wait()
}

func TestSkipListBasic(t *testing.T) {
	sl := NewSkipList(10)
	const goroutines = 10
	const opsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < opsPerGoroutine; i++ {
				op := rand.Intn(3)
				key := fmt.Sprintf("%d-%d", g, rand.Intn(500))
				switch op {
				case 0:
					sl.Insert(uint64(time.Now().UnixNano()), true, key, string(fmt.Sprintf("value-%d-%d", g, i)))
				case 1:
					sl.Delete(key, uint64(time.Now().UnixNano()))
				}
			}
		}()
	}
	wg.Wait()

	for node := range sl.Iterator() {
		if node != nil {
			t.Logf("Key: %s, Value: %s\n", node.Key, node.Value)
		}
	}
}

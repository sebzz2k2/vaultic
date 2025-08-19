package index

import (
	"fmt"
)

type IndexValue struct {
	Start uint32
	End   uint32
}

func (idx *Index) Set(key string, start, end uint32) {
	idx.index.Store(key, IndexValue{Start: start, End: end})
}

func (idx *Index) Get(key string) (uint32, uint32, bool) {
	val, exists := idx.index.Load(key)
	if exists {
		return val.(IndexValue).Start, val.(IndexValue).End, true
	}
	return 0, 0, false
}

func (idx *Index) Del(key string) {
	idx.index.Delete(key)
}

func (idx *Index) Exists(key string) bool {
	_, exists := idx.index.Load(key)
	return exists
}

func (idx *Index) Keys() []string {
	keys := []string{}
	idx.index.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}

func (idx *Index) print() {
	idx.index.Range(func(key, value interface{}) bool {
		fmt.Printf("%s: %d %d\n", key.(string), value.(IndexValue).Start, value.(IndexValue).End)
		return true
	})
}

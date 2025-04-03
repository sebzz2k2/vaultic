package utils

import (
	"fmt"
	"sync"
)

var Data sync.Map

type IndexValue struct {
	Start uint32
	End   uint32
}

func SetIndexKey(key string, start, end uint32) {
	Data.Store(key, IndexValue{Start: start, End: end})
}

func GetIndexVal(key string) (int, bool) {
	val, exists := Data.Load(key)
	if exists {
		return val.(int), true
	}
	return 0, false
}

func PrintIndexMap() {
	Data.Range(func(key, value interface{}) bool {
		fmt.Printf("%s: %d %d\n", key.(string), value.(IndexValue).Start, value.(IndexValue).End)
		return true
	})
}

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

func GetIndexVal(key string) (uint32, uint32, bool) {
	val, exists := Data.Load(key)
	if exists {
		return val.(IndexValue).Start, val.(IndexValue).End, true
	}
	return 0, 0, false
}

func DeleteIndexKey(key string) {
	Data.Delete(key)
}

func IsPresent(key string) bool {
	_, exists := Data.Load(key)
	return exists
}

func PrintIndexMap() {
	Data.Range(func(key, value interface{}) bool {
		fmt.Printf("%s: %d %d\n", key.(string), value.(IndexValue).Start, value.(IndexValue).End)
		return true
	})
}

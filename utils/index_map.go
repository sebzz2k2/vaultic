package utils

import "sync"

var Data sync.Map

func SetIndexKey(key string, value int) {
	Data.Store(key, value)
}

func GetIndexVal(key string) (int, bool) {
	val, exists := Data.Load(key)
	if exists {
		return val.(int), true
	}
	return 0, false
}

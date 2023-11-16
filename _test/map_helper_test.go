package test

import (
	"github.com/mangenotwork/gathertool"
	"testing"
)

// go test -v
func TestOderMap(t *testing.T) {
	intMap := gathertool.OrderMap[int, int]()
	intMap.Add(1, 1)
	t.Log(intMap.Get(1))
	intMap.Del(1)
	t.Log(intMap.Get(1))
	intMap.Add(1, 1).Add(2, 2).Add(3, 3)
	intMap.RangeAt(func(id, k, v int) {
		t.Log(id, k, v)
	})
	stringMap := gathertool.OrderMap[string, int]()
	stringMap.Add("a", 1)
	t.Log(stringMap.Get("a"))
	stringMap.Del("a")
	t.Log(stringMap.Get("a"))
	stringMap.Add("a", 1)
	stringMap.AddMap(map[string]int{"b": 2, "c": 3})
	stringMap.RangeAt(func(id int, k string, v int) {
		t.Log(id, k, v)
	})
	stringMap.Reverse()
	stringMap.RangeAt(func(id int, k string, v int) {
		t.Log(id, k, v)
	})
	t.Log(stringMap.KeyList())
}

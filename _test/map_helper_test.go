package test

import (
	"github.com/mangenotwork/gathertool"
	"testing"
)

// go test -v -run=TestMapHelperOderMap
func TestMapHelperOderMap(t *testing.T) {
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

	t.Log(stringMap.Json())

	t.Log(stringMap.Move("a", 2))
	stringMap.DebugPrint()
	t.Log(stringMap.Move("a", 3))
	t.Log(stringMap.Insert("d", 4, 1))
	stringMap.DebugPrint()

	t.Log("洗牌  ------> ")
	stringMap.Shuffle()
	stringMap.DebugPrint()

	t.Log(stringMap.GetAtPosition(2))
	t.Log(stringMap.Pop())
	stringMap.DebugPrint()
	t.Log(stringMap.BackPop())
	stringMap.DebugPrint()

}

// go test -v -run=TestMapHelperSet
func TestMapHelperSet(t *testing.T) {
	s := gathertool.NewSet()
	s.Add(1)
	s.Add(1)
	s.DebugPrint()
	s.Add(2)
	s.DebugPrint()
	t.Log(s.Has(3))
	t.Log(s.Has(1))
	s.Delete(1)
	s.DebugPrint()
}

// go test -v -run=TestMapHelperStack
func TestMapHelperStack(t *testing.T) {
	s := gathertool.NewStack[string]()
	s.Push("a")
	s.Push("b")
	s.DebugPrint()
	t.Log(s.Pop())
	s.DebugPrint()
}

// go test -v -run=TestMapHelper
func TestMapHelper(t *testing.T) {
	map1 := map[int]string{1: "a", 2: "b"}
	map2 := gathertool.MapCopy(map1)
	t.Log(map2)
	t.Log(gathertool.MapMergeCopy(map1, map[int]string{3: "c", 2: "b"}))

	t.Log(gathertool.Map2Slice(map1))

	t.Log(gathertool.Slice2Map([]any{"a", 1, "b", 2}))
}

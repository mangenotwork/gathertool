/*
	Description : 切片工具
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"math/rand"
	"sync"
	"time"
)

type SliceType interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64 | string | bool |
		*int | *uint | *int8 | *uint8 | *int16 | *uint16 | *int32 | *uint32 | *int64 | *uint64 | *float32 | *float64 |
		*string | *bool | chan int | chan uint | chan int8 | chan uint8 | chan int16 | chan uint16 | chan int32 |
		chan uint32 | chan int64 | chan uint64 | chan float32 | chan float64 | chan string | chan bool | chan *int |
		chan *uint | chan *int8 | chan *uint8 | chan *int16 | chan *uint16 | chan *int32 | chan *uint32 | chan *int64 |
		chan *uint64 | chan *float32 | chan *float64 | chan *string | chan *bool
}

type sliceTool struct{}

var st *sliceTool
var stOnce sync.Once

// SliceTool use : SliceTool().CopyInt64(a)
func SliceTool() *sliceTool {
	stOnce.Do(func() {
		st = &sliceTool{}
	})
	return st
}

// CopyInt64 copy int64
func (sliceTool) CopyInt64(a []int64) []int64 {
	return append(a[:0:0], a...)
}

// CopyStr copy string
func (sliceTool) CopyStr(a []string) []string {
	return append(a[:0:0], a...)
}

// CopyInt copy int
func (sliceTool) CopyInt(a []int) []int {
	return append(a[:0:0], a...)
}

// ContainsByte contains byte
func (sliceTool) ContainsByte(a []byte, x byte) bool {
	l := len(a)
	if l == 0 {
		return false
	}
	for i := 0; i < l; i++ {
		if a[i] == x {
			return true
		}
	}
	return false
}

// ContainsInt contains int
func (sliceTool) ContainsInt(a []int, x int) bool {
	l := len(a)
	if l == 0 {
		return false
	}
	for i := 0; i < l; i++ {
		if a[i] == x {
			return true
		}
	}
	return false
}

// ContainsInt64  contains int64
func (sliceTool) ContainsInt64(a []int64, x int64) bool {
	l := len(a)
	if l == 0 {
		return false
	}
	for i := 0; i < l; i++ {
		if a[i] == x {
			return true
		}
	}
	return false
}

// ContainsStr contains str
func (sliceTool) ContainsStr(a []string, x string) bool {
	l := len(a)
	if l == 0 {
		return false
	}
	for i := 0; i < l; i++ {
		if a[i] == x {
			return true
		}
	}
	return false
}

func (sliceTool) DeduplicateInt(a []int) []int {
	l := len(a)
	if l < 2 {
		return a
	}
	seen := make(map[int]struct{})
	j := 0
	for i := 0; i < l; i++ {
		if _, ok := seen[a[i]]; ok {
			continue
		}
		seen[a[i]] = struct{}{}
		a[j] = a[i]
		j++
	}
	return a[:j]
}

// DeduplicateInt64 deduplicate int64
func (sliceTool) DeduplicateInt64(a []int64) []int64 {
	l := len(a)
	if l < 2 {
		return a
	}
	seen := make(map[int64]struct{})
	j := 0
	for i := 0; i < l; i++ {
		if _, ok := seen[a[i]]; ok {
			continue
		}
		seen[a[i]] = struct{}{}
		a[j] = a[i]
		j++
	}
	return a[:j]
}

// DeduplicateStr  deduplicate string
func (sliceTool) DeduplicateStr(a []string) []string {
	l := len(a)
	if l < 2 {
		return a
	}
	seen := make(map[string]struct{})
	j := 0
	for i := 0; i < l; i++ {
		if _, ok := seen[a[i]]; ok {
			continue
		}
		seen[a[i]] = struct{}{}
		a[j] = a[i]
		j++
	}
	return a[:j]
}

// DelInt del int
func (sliceTool) DelInt(a []int, i int) []int {
	l := len(a)
	if l == 0 {
		return nil
	}
	if i < 0 || i > l-1 {
		return nil
	}
	return append(a[:i], a[i+1:]...)
}

// DelInt64 del int64
func (sliceTool) DelInt64(a []int64, i int) []int64 {
	l := len(a)
	if l == 0 {
		return nil
	}
	if i < 0 || i > l-1 {
		return nil
	}
	return append(a[:i], a[i+1:]...)
}

// DelStr delete str
func (sliceTool) DelStr(a []string, i int) []string {
	l := len(a)
	if l == 0 {
		return nil
	}
	if i < 0 || i > l-1 {
		return nil
	}
	return append(a[:i], a[i+1:]...)
}

func (sliceTool) MaxInt(a []int) int {
	l := len(a)
	if l == 0 {
		return 0
	}
	max := a[0]
	for k := 1; k < l; k++ {
		if a[k] > max {
			max = a[k]
		}
	}
	return max
}

func (sliceTool) MaxInt64(a []int64) int64 {
	l := len(a)
	if l == 0 {
		return 0
	}
	max := a[0]
	for k := 1; k < l; k++ {
		if a[k] > max {
			max = a[k]
		}
	}
	return max
}

func (sliceTool) MinInt(a []int) int {
	l := len(a)
	if l == 0 {
		return 0
	}
	min := a[0]
	for k := 1; k < l; k++ {
		if a[k] < min {
			min = a[k]
		}
	}
	return min
}

func (sliceTool) MinInt64(a []int64) int64 {
	l := len(a)
	if l == 0 {
		return 0
	}
	min := a[0]
	for k := 1; k < l; k++ {
		if a[k] < min {
			min = a[k]
		}
	}
	return min
}

func (sliceTool) PopInt(a []int) (int, []int) {
	if len(a) == 0 {
		return 0, nil
	}
	return a[len(a)-1], a[:len(a)-1]
}

func (sliceTool) PopInt64(a []int64) (int64, []int64) {
	if len(a) == 0 {
		return 0, nil
	}
	return a[len(a)-1], a[:len(a)-1]
}

func (sliceTool) PopStr(a []string) (string, []string) {
	if len(a) == 0 {
		return "", nil
	}
	return a[len(a)-1], a[:len(a)-1]
}

// ReverseInt  反转
func (sliceTool) ReverseInt(a []int) []int {
	l := len(a)
	if l == 0 {
		return a
	}
	for s, e := 0, len(a)-1; s < e; {
		a[s], a[e] = a[e], a[s]
		s++
		e--
	}
	return a
}

// ReverseInt64 reverse int64
func (sliceTool) ReverseInt64(a []int64) []int64 {
	l := len(a)
	if l == 0 {
		return a
	}
	for s, e := 0, len(a)-1; s < e; {
		a[s], a[e] = a[e], a[s]
		s++
		e--
	}
	return a
}

// ReverseStr  reverseStr
func (sliceTool) ReverseStr(a []string) []string {
	l := len(a)
	if l == 0 {
		return a
	}
	for s, e := 0, len(a)-1; s < e; {
		a[s], a[e] = a[e], a[s]
		s++
		e--
	}
	return a
}

// ShuffleInt 洗牌
func (sliceTool) ShuffleInt(a []int) []int {
	l := len(a)
	if l <= 1 {
		return a
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(l, func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}

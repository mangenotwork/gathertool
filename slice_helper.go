/*
*	Description : 切片相关
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"math/rand"
	"time"
)

var randObj = rand.New(rand.NewSource(time.Now().UnixNano()))

func SliceContains[V comparable](a []V, v V) bool {
	l := len(a)
	if l == 0 {
		return false
	}
	for i := 0; i < l; i++ {
		if a[i] == v {
			return true
		}
	}
	return false
}

func SliceDeduplicate[V comparable](a []V) []V {
	l := len(a)
	if l < 2 {
		return a
	}
	seen := make(map[V]struct{})
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

func SliceDel[V comparable](a []V, i int) []V {
	l := len(a)
	if l == 0 {
		return a
	}
	if i < 0 || i > l-1 {
		return a
	}
	return append(a[:i], a[i+1:]...)
}

type number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func SliceMax[V number](a []V) V {
	l := len(a)
	if l == 0 {
		var none V
		return none
	}
	maxV := a[0]
	for k := 1; k < l; k++ {
		if a[k] > maxV {
			maxV = a[k]
		}
	}
	return maxV
}

func SliceMin[V number](a []V) V {
	l := len(a)
	if l == 0 {
		return 0
	}
	minV := a[0]
	for k := 1; k < l; k++ {
		if a[k] < minV {
			minV = a[k]
		}
	}
	return minV
}

func SlicePop[V comparable](a []V) (V, []V) {
	if len(a) == 0 {
		var none V
		return none, a
	}
	return a[len(a)-1], a[:len(a)-1]
}

func SliceReverse[V comparable](a []V) []V {
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

func SliceShuffle[V comparable](a []V) []V {
	l := len(a)
	if l <= 1 {
		return a
	}
	randObj.Shuffle(l, func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})
	return a
}

func SliceCopy[V comparable](a []V) []V {
	return append(a[:0:0], a...)
}

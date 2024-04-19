package test

import (
	"github.com/mangenotwork/gathertool"
	"testing"
)

func TestSliceHelper(t *testing.T) {
	t1 := []string{"1", "2", "3", "4", "4", "3"}
	t.Log(gathertool.SliceContains[string](t1, "1"))

	t.Log(gathertool.SliceDeduplicate[string](t1))

	t2 := []int{1, 2, 3, 4}
	t.Log(gathertool.SliceDel[int](t2, 1))

	t.Log(gathertool.SliceMax[int](t2))

	t.Log(gathertool.SliceMin[int](t2))

	t.Log(gathertool.SlicePop[int](t2))

	t.Log(t1)
	t.Log(gathertool.SliceReverse[string](t1))

	t.Log(gathertool.SliceShuffle[string](t1))

	t.Log(gathertool.SliceCopy[string](t1))
}

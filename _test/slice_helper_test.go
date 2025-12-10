package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
)

func TestSliceHelper(t *testing.T) {
	t1 := []string{"1", "2", "3", "4", "4", "3"}
	t.Log(gt.SliceContains[string](t1, "1"))

	t.Log(gt.SliceDeduplicate[string](t1))

	t2 := []int{1, 2, 3, 4}
	t.Log(gt.SliceDel[int](t2, 1))

	t.Log(gt.SliceMax[int](t2))

	t.Log(gt.SliceMin[int](t2))

	t.Log(gt.SlicePop[int](t2))

	t.Log(t1)
	t.Log(gt.SliceReverse[string](t1))

	t.Log(gt.SliceShuffle[string](t1))

	t.Log(gt.SliceCopy[string](t1))
}

package test

import (
	"github.com/mangenotwork/gathertool"
	"testing"
)

func TestStringHelperTypeConversionTest(t *testing.T) {
	t.Log(gathertool.StringValue(1))
	t.Log(gathertool.StringValue(2.22))
	t.Log(gathertool.StringValue(true))
	t.Log(gathertool.Any2String([]int{1, 2}))

	t.Log(gathertool.Json2Map("{\"a\":2,\"3\":4}"))
	t.Log(gathertool.Any2Map("{\"a\":2,\"3\":4}"))

	t.Log(gathertool.Any2Int("123"))
	t.Log(gathertool.Any2Int("123a"))
	t.Log(gathertool.Any2Int64("123"))
	t.Log(gathertool.Any2Arr("123,123"))
	t.Log(gathertool.Any2Float64("123.123"))

	t.Log(gathertool.Int2Hex(123))
	t.Log(gathertool.Int642Hex(123))
	t.Log(gathertool.Hex2Int("7b"))
	t.Log(gathertool.Hex2Int64("7b"))
	t.Log(gathertool.Str2Int("123"))
	t.Log(gathertool.Str2Int32("123"))
	t.Log(gathertool.Str2Float64("123.123"))
	t.Log(gathertool.Str2Float32("123.123"))
	t.Log(gathertool.Uint82Str([]uint8{1, 2, 3}))
	t.Log(gathertool.Str2Byte("abc"))
	t.Log(gathertool.Byte2Str([]byte("abc")))
	t.Log(gathertool.Bool2Byte(true))
	t.Log(gathertool.Byte2Bool([]byte("0")))
	t.Log(gathertool.Int2Byte(123))
	t.Log(gathertool.Byte2Int([]byte("123")))
	t.Log(gathertool.Int642Byte(123))
	t.Log(gathertool.Byte2Int64([]byte("123")))
	t.Log(gathertool.Float322Byte(123.123))
	t.Log(gathertool.Float322Uint32(123.123))
	t.Log(gathertool.Float322Uint32(123.123))
	t.Log(gathertool.Float642Byte(123.123))
	t.Log(gathertool.Float642Uint64(123.123))
	t.Log(gathertool.Byte2Float64([]byte("123")))

	t.Log(gathertool.CleaningStr(" sada\n\t\r "))

}

func TestStringHelperTest(t *testing.T) {
	t.Log(gathertool.MD5("asdsad"))
}

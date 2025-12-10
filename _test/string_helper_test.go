package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
)

func TestStringHelperTypeConversion(t *testing.T) {
	t.Log(gt.Any2String(1))
	t.Log(gt.Any2String(2.22))
	t.Log(gt.Any2String(true))
	t.Log(gt.Any2String([]int{1, 2}))

	t.Log(gt.Json2Map("{\"a\":2,\"3\":4}"))
	t.Log(gt.Any2Map("{\"a\":2,\"3\":4}"))

	t.Log(gt.Any2Int("123"))
	t.Log(gt.Any2Int("123a"))
	t.Log(gt.Any2Int64("123"))
	t.Log(gt.Any2Arr("123,123"))
	t.Log(gt.Any2Float64("123.123"))

	t.Log(gt.Int2Hex(123))
	t.Log(gt.Int642Hex(123))
	t.Log(gt.Hex2Int("7b"))
	t.Log(gt.Hex2Int64("7b"))
	t.Log(gt.Str2Int("123"))
	t.Log(gt.Str2Int32("123"))
	t.Log(gt.Str2Float64("123.123"))
	t.Log(gt.Str2Float32("123.123"))
	t.Log(gt.Uint82Str([]uint8{1, 2, 3}))
	t.Log(gt.Str2Byte("abc"))
	t.Log(gt.Byte2Str([]byte("abc")))
	t.Log(gt.Bool2Byte(true))
	t.Log(gt.Byte2Bool([]byte("0")))
	t.Log(gt.Int2Byte(123))
	t.Log(gt.Byte2Int([]byte("1111")))
	t.Log(gt.Int642Byte(123))
	t.Log(gt.Byte2Int64(gt.Int642Byte(123)))
	t.Log(gt.Float322Byte(123.123))
	t.Log(gt.Float322Uint32(123.123))
	t.Log(gt.Byte2Float32(gt.Float322Byte(123.123)))
	t.Log(gt.Float642Byte(123.123))
	t.Log(gt.Float642Uint64(123.123))
	t.Log(gt.Byte2Float64(gt.Float642Byte(123.123)))
	t.Log(gt.Float642Uint64(123.123))
	t.Log(gt.Byte2Float64(gt.Float642Byte(123.123)))
	type A struct {
		C string
	}
	a1 := &A{C: "1"}
	t.Log(gt.Struct2Map(a1, true))

	t.Log(gt.CleaningStr(" sada\n\t\r "))
	t.Log(gt.StrDeleteSpace(" sada "))
	t1 := gt.EncodeByte("test1")
	t.Log(t1)
	t.Log(gt.DecodeByte(t1))
	t.Log(gt.Byte2Bit([]byte("aaa")))
	t.Log(gt.Bit2Byte(gt.Byte2Bit([]byte("aaa"))))

	a2 := &A{}
	t.Log(gt.DeepCopy[*A](a2, a1))
	t.Log(a2)

	t.Log(gt.PanicToError(func() {
		panic(1)
	}))

	t.Log(gt.ByteToBinaryString(byte('a')))
	t.Log(gt.UnicodeDec("aaa"))
	t.Log(gt.UnicodeDecByte([]byte("aaa")))
	t.Log(gt.UnescapeUnicode(gt.UnicodeDecByte([]byte("aaa"))))
	t.Log(gt.Base64Encode("一二三"))
	t.Log(gt.Base64Decode(gt.Base64Encode("一二三")))
	t.Log(gt.Base64UrlEncode("一二三"))
	t.Log(gt.Base64UrlDecode(gt.Base64UrlEncode("一二三")))

	t2, _ := gt.ToUTF8("HZGB2312", "一二三")
	t.Log(t2)
	t.Log(gt.UTF8To("HZGB2312", t2))

	t3, _ := gt.ToUTF16("HZGB2312", "一二三")
	t.Log(t3)
	t.Log(gt.UTF16To("HZGB2312", t3))

	t4, _ := gt.ToBIG5("HZGB2312", "一二三")
	t.Log(t4)
	t.Log(gt.BIG5To("HZGB2312", t4))

	t5, _ := gt.ToGDK("HZGB2312", "一二三")
	t.Log(t5)
	t.Log(gt.GDKTo("HZGB2312", t5))

	t6, _ := gt.ToGB18030("HZGB2312", "一二三")
	t.Log(t6)
	t.Log(gt.GB18030To("HZGB2312", t6))

	t7, _ := gt.ToGB2312("HZGB2312", "一二三")
	t.Log(t7)
	t.Log(gt.GB2312To("HZGB2312", t7))

	t8, _ := gt.ToHZGB2312("HZGB2312", "一二三")
	t.Log(t8)
	t.Log(gt.HZGB2312To("HZGB2312", t8))

	t.Log(gt.IF[string](true, "a", "b"))

	t.Log(gt.ReplaceAllToOne("aaacvccc", []string{"a", "c"}, ""))

	t.Log(gt.HumanFriendlyTraffic(12354646))
	t.Log(gt.StrToSize("123415kb"))

	t.Log(gt.IP2Binary("127.0.0.1"))
	t.Log(gt.UInt32ToIP(12700001))
	t.Log(gt.IP2Int64("127.0.0.1"))

	t.Log(gt.GzipCompress([]byte("abc")))
	t.Log(gt.GzipDecompress(gt.GzipCompress([]byte("abc"))))
	t.Log(gt.ConvertStr2GBK("一二三"))
	t.Log(gt.ConvertGBK2Str("一二三"))
	t.Log(gt.ByteToGBK([]byte("一二三")))
	t.Log(gt.ByteToUTF8([]byte("一二三")))
	t.Log(gt.ID())
	t.Log(gt.IDMd5())

}

func TestStringHelper(t *testing.T) {
	t.Log(gt.MD5("asdsad"))
}

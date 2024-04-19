package test

import (
	"github.com/mangenotwork/gathertool"
	"testing"
)

func TestStringHelperTypeConversion(t *testing.T) {
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
	t.Log(gathertool.Byte2Int([]byte("1111")))
	t.Log(gathertool.Int642Byte(123))
	t.Log(gathertool.Byte2Int64(gathertool.Int642Byte(123)))
	t.Log(gathertool.Float322Byte(123.123))
	t.Log(gathertool.Float322Uint32(123.123))
	t.Log(gathertool.Byte2Float32(gathertool.Float322Byte(123.123)))
	t.Log(gathertool.Float642Byte(123.123))
	t.Log(gathertool.Float642Uint64(123.123))
	t.Log(gathertool.Byte2Float64(gathertool.Float642Byte(123.123)))
	t.Log(gathertool.Float642Uint64(123.123))
	t.Log(gathertool.Byte2Float64(gathertool.Float642Byte(123.123)))
	type A struct {
		C string
	}
	a1 := &A{C: "1"}
	t.Log(gathertool.Struct2Map(a1, true))

	t.Log(gathertool.CleaningStr(" sada\n\t\r "))
	t.Log(gathertool.StrDeleteSpace(" sada "))
	t1 := gathertool.EncodeByte("test1")
	t.Log(t1)
	t.Log(gathertool.DecodeByte(t1))
	t.Log(gathertool.Byte2Bit([]byte("aaa")))
	t.Log(gathertool.Bit2Byte(gathertool.Byte2Bit([]byte("aaa"))))

	a2 := &A{}
	t.Log(gathertool.DeepCopy[*A](a2, a1))
	t.Log(a2)

	t.Log(gathertool.PanicToError(func() {
		panic(1)
	}))

	t.Log(gathertool.ByteToBinaryString(byte('a')))
	t.Log(gathertool.UnicodeDec("aaa"))
	t.Log(gathertool.UnicodeDecByte([]byte("aaa")))
	t.Log(gathertool.UnescapeUnicode(gathertool.UnicodeDecByte([]byte("aaa"))))
	t.Log(gathertool.Base64Encode("一二三"))
	t.Log(gathertool.Base64Decode(gathertool.Base64Encode("一二三")))
	t.Log(gathertool.Base64UrlEncode("一二三"))
	t.Log(gathertool.Base64UrlDecode(gathertool.Base64UrlEncode("一二三")))

	t2, _ := gathertool.ToUTF8("HZGB2312", "一二三")
	t.Log(t2)
	t.Log(gathertool.UTF8To("HZGB2312", t2))

	t3, _ := gathertool.ToUTF16("HZGB2312", "一二三")
	t.Log(t3)
	t.Log(gathertool.UTF16To("HZGB2312", t3))

	t4, _ := gathertool.ToBIG5("HZGB2312", "一二三")
	t.Log(t4)
	t.Log(gathertool.BIG5To("HZGB2312", t4))

	t5, _ := gathertool.ToGDK("HZGB2312", "一二三")
	t.Log(t5)
	t.Log(gathertool.GDKTo("HZGB2312", t5))

	t6, _ := gathertool.ToGB18030("HZGB2312", "一二三")
	t.Log(t6)
	t.Log(gathertool.GB18030To("HZGB2312", t6))

	t7, _ := gathertool.ToGB2312("HZGB2312", "一二三")
	t.Log(t7)
	t.Log(gathertool.GB2312To("HZGB2312", t7))

	t8, _ := gathertool.ToHZGB2312("HZGB2312", "一二三")
	t.Log(t8)
	t.Log(gathertool.HZGB2312To("HZGB2312", t8))

	t.Log(gathertool.IF[string](true, "a", "b"))

	t.Log(gathertool.ReplaceAllToOne("aaacvccc", []string{"a", "c"}, ""))

	t.Log(gathertool.HumanFriendlyTraffic(12354646))
	t.Log(gathertool.StrToSize("123415kb"))

	t.Log(gathertool.IP2Binary("127.0.0.1"))
	t.Log(gathertool.UInt32ToIP(12700001))
	t.Log(gathertool.IP2Int64("127.0.0.1"))

	t.Log(gathertool.GzipCompress([]byte("abc")))
	t.Log(gathertool.GzipDecompress(gathertool.GzipCompress([]byte("abc"))))
	t.Log(gathertool.ConvertStr2GBK("一二三"))
	t.Log(gathertool.ConvertGBK2Str("一二三"))
	t.Log(gathertool.ByteToGBK([]byte("一二三")))
	t.Log(gathertool.ByteToUTF8([]byte("一二三")))
	t.Log(gathertool.ID())
	t.Log(gathertool.IDMd5())

}

func TestStringHelper(t *testing.T) {
	t.Log(gathertool.MD5("asdsad"))
}

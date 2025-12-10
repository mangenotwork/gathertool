/*
*	Description : 字符串，字符编码等等相关的操作方法
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
)

// CleaningStr 清理字符串前后空白，回车，换行符号
func CleaningStr(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\\n", "", -1)
	str = strings.TrimSpace(str)
	str = StrDeleteSpace(str)
	return str
}

// StrDeleteSpace 删除字符串前后的空格
func StrDeleteSpace(str string) string {
	strList := []byte(str)
	spaceCount, count := 0, len(strList)
	for i := 0; i <= len(strList)-1; i++ {
		if strList[i] == 32 {
			spaceCount++
		} else {
			break
		}
	}

	strList = strList[spaceCount:]
	spaceCount, count = 0, len(strList)
	for i := count - 1; i >= 0; i-- {
		if strList[i] == 32 {
			spaceCount++
		} else {
			break
		}
	}

	return string(strList[:count-spaceCount])
}

// EncodeByte 任意类型编码为[]byte
func EncodeByte(v any) []byte {
	switch value := v.(type) {
	case int, int8, int16, int32:
		return Int2Byte(value.(int))
	case int64:
		return Int642Byte(value)
	case string:
		return Str2Byte(value)
	case bool:
		return Bool2Byte(value)
	case float32:
		return Float322Byte(value)
	case float64:
		return Float642Byte(value)
	}
	return []byte("")
}

// DecodeByte  decode byte
func DecodeByte(b []byte) ([]byte, error) {
	rse := make([]byte, 0)
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, rse)
	return rse, err
}

func deepCopy[T any](dst, src T) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// DeepCopy 深copy
func DeepCopy[T any](dst, src T) error {
	return deepCopy(dst, src)
}

// PanicToError panic -> error
func PanicToError(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic error: %v", r)
		}
	}()
	fn()
	return
}

// P2E panic -> error
func P2E() {
	defer func() {
		if r := recover(); r != nil {
			Error("Panic error: %v", r)
		}
	}()
}

// Charset 字符集类型
type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GBK     = Charset("GBK")
	GB2312  = Charset("GB2312")
)

// ConvertByte2String 编码转换
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)

	case GBK:
		var decodeBytes, _ = simplifiedchinese.GBK.NewDecoder().Bytes(byte)
		str = string(decodeBytes)

	case GB2312:
		var decodeBytes, _ = simplifiedchinese.HZGB2312.NewDecoder().Bytes(byte)
		str = string(decodeBytes)

	case UTF8:
		fallthrough

	default:
		str = string(byte)
	}

	return str
}

// UnicodeDec Unicode编码
func UnicodeDec(raw string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(raw), `\\u`, `\u`, -1))
	if err != nil {
		return ""
	}
	return str
}

// UnicodeDecByte Unicode编码输出[]byte
func UnicodeDecByte(raw []byte) []byte {
	rawStr := string(raw)
	return []byte(UnicodeDec(rawStr))
}

// UnescapeUnicode Unicode 转码
func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

// Base64Encode base64 编码
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64 解码
func Base64Decode(str string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(str)
	return string(b), err
}

// Base64UrlEncode base64 url 编码
func Base64UrlEncode(str string) string {
	return base64.URLEncoding.EncodeToString([]byte(str))
}

// Base64UrlDecode base64 url 解码
func Base64UrlDecode(str string) (string, error) {
	b, err := base64.URLEncoding.DecodeString(str)
	return string(b), err
}

func convert(dstCharset string, srcCharset string, src string) (dst string, err error) {
	if dstCharset == srcCharset {
		return src, nil
	}
	dst = src

	if srcCharset != "UTF-8" {
		if e := getEncoding(srcCharset); e != nil {
			tmp, err := io.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewDecoder()),
			)
			if err != nil {
				return "", err
			}
			src = string(tmp)
		} else {
			return dst, err
		}
	}

	if dstCharset != "UTF-8" {
		if e := getEncoding(dstCharset); e != nil {
			tmp, err := io.ReadAll(
				transform.NewReader(bytes.NewReader([]byte(src)), e.NewEncoder()),
			)
			if err != nil {
				return "", err
			}
			dst = string(tmp)
		} else {
			return dst, err
		}
	} else {
		dst = src
	}
	return dst, nil
}

var (
	// Alias for charsets.
	charsetAlias = map[string]string{
		"HZGB2312": "HZ-GB-2312",
		"hzgb2312": "HZ-GB-2312",
		"GB2312":   "HZ-GB-2312",
		"gb2312":   "HZ-GB-2312",
	}
)

func getEncoding(charset string) encoding.Encoding {
	if c, ok := charsetAlias[charset]; ok {
		charset = c
	}
	enc, err := ianaindex.MIB.Encoding(charset)
	if err != nil {
		log.Println(err)
	}
	return enc
}

// ToUTF8  to utf8
func ToUTF8(srcCharset string, src string) (dst string, err error) {
	return convert("UTF-8", srcCharset, src)
}

// UTF8To utf8 to
func UTF8To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "UTF-8", src)
}

// ToUTF16 to utf16
func ToUTF16(srcCharset string, src string) (dst string, err error) {
	return convert("UTF-16", srcCharset, src)
}

// UTF16To utf16 to
func UTF16To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "UTF-16", src)
}

// ToBIG5  to big5
func ToBIG5(srcCharset string, src string) (dst string, err error) {
	return convert("big5", srcCharset, src)
}

// BIG5To  big to
func BIG5To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "big5", src)
}

// ToGDK to gdk
func ToGDK(srcCharset string, src string) (dst string, err error) {
	return convert("gbk", srcCharset, src)
}

// GDKTo gdk to
func GDKTo(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "gbk", src)
}

// ToGB18030  to gb18030
func ToGB18030(srcCharset string, src string) (dst string, err error) {
	return convert("gb18030", srcCharset, src)
}

// GB18030To gb18030 to
func GB18030To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "gb18030", src)
}

// ToGB2312 to gb2312
func ToGB2312(srcCharset string, src string) (dst string, err error) {
	return convert("GB2312", srcCharset, src)
}

// GB2312To gb2312 to
func GB2312To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "GB2312", src)
}

// ToHZGB2312 to hzgb2312
func ToHZGB2312(srcCharset string, src string) (dst string, err error) {
	return convert("HZGB2312", srcCharset, src)
}

// HZGB2312To hzgb2312 to
func HZGB2312To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "HZGB2312", src)
}

// IF 三元表达式
func IF[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// ReplaceAllToOne 批量统一替换字符串
func ReplaceAllToOne(str string, from []string, to string) string {
	arr := make([]string, len(from)*2)
	for i, s := range from {
		arr[i*2] = s
		arr[i*2+1] = to
	}
	r := strings.NewReplacer(arr...)
	return r.Replace(str)
}

// 字节换算
const (
	KiB = 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
	TiB = GiB * 1024
)

// HumanFriendlyTraffic 字节换算
func HumanFriendlyTraffic(bytes uint64) string {
	if bytes <= KiB {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes <= MiB {
		return fmt.Sprintf("%.2f KiB", float32(bytes)/KiB)
	}
	if bytes <= GiB {
		return fmt.Sprintf("%.2f MiB", float32(bytes)/MiB)
	}
	if bytes <= TiB {
		return fmt.Sprintf("%.2f GiB", float32(bytes)/GiB)
	}
	return fmt.Sprintf("%.2f TiB", float32(bytes)/TiB)
}

// StrToSize 字节换算
func StrToSize(sizeStr string) int64 {
	i := 0
	for ; i < len(sizeStr); i++ {
		if sizeStr[i] == '.' || (sizeStr[i] >= '0' && sizeStr[i] <= '9') {
			continue
		} else {
			break
		}
	}
	var (
		unit      = sizeStr[i:]
		number, _ = strconv.ParseFloat(sizeStr[:i], 64)
	)
	if unit == "" {
		return int64(number)
	}
	switch strings.ToLower(unit) {
	case "b", "bytes":
		return int64(number)
	case "k", "kb", "ki", "kib", "kilobyte":
		return int64(number * 1024)
	case "m", "mb", "mi", "mib", "mebibyte":
		return int64(number * 1024 * 1024)
	case "g", "gb", "gi", "gib", "gigabyte":
		return int64(number * 1024 * 1024 * 1024)
	case "t", "tb", "ti", "tib", "terabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024)
	case "p", "pb", "pi", "pib", "petabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024)
	case "e", "eb", "ei", "eib", "exabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024 * 1024)
	case "z", "zb", "zi", "zib", "zettabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024)
	case "y", "yb", "yi", "yib", "yottabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024)
	case "bb", "brontobyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024)
	}
	return -1
}

// GzipCompress Gzip压缩
func GzipCompress(src []byte) []byte {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	_, _ = w.Write(src)
	_ = w.Close()
	return in.Bytes()
}

// GzipDecompress Gzip解压
func GzipDecompress(src []byte) []byte {
	dst := make([]byte, 0)
	br := bytes.NewReader(src)
	gr, err := gzip.NewReader(br)
	if err != nil {
		return dst
	}
	defer func() {
		_ = gr.Close()
	}()
	tmp, err := io.ReadAll(gr)
	if err != nil {
		return dst
	}
	dst = tmp
	return dst
}

// ConvertStr2GBK 将utf-8编码的字符串转换为GBK编码
func ConvertStr2GBK(str string) string {
	ret, err := simplifiedchinese.GBK.NewEncoder().String(str)
	if err != nil {
		ret = str
	}
	return ret
}

// ConvertGBK2Str 将GBK编码的字符串转换为utf-8编码
func ConvertGBK2Str(gbkStr string) string {
	ret, err := simplifiedchinese.GBK.NewDecoder().String(gbkStr)
	if err != nil {
		ret = gbkStr
	}
	return ret
}

// ByteToGBK   byte -> gbk byte
func ByteToGBK(strBuf []byte) []byte {
	if IsUtf8(strBuf) {
		if GBKBuf, err := simplifiedchinese.GBK.NewEncoder().Bytes(strBuf); err == nil {
			if IsUtf8(GBKBuf) {
				return GBKBuf
			}
		}
		if GB18030Buf, err := simplifiedchinese.GB18030.NewEncoder().Bytes(strBuf); err == nil {
			if IsUtf8(GB18030Buf) {
				return GB18030Buf
			}
		}
		if HZGB2312Buf, err := simplifiedchinese.HZGB2312.NewEncoder().Bytes(strBuf); err == nil {
			if IsUtf8(HZGB2312Buf) {
				return HZGB2312Buf
			}
		}
		return strBuf
	} else {
		return strBuf
	}
}

// ByteToUTF8   byte -> utf8 byte
func ByteToUTF8(strBuf []byte) []byte {
	if IsUtf8(strBuf) {
		return strBuf
	} else {
		if GBKBuf, err := simplifiedchinese.GBK.NewDecoder().Bytes(strBuf); err == nil {
			if IsUtf8(GBKBuf) {
				return GBKBuf
			}
		}
		if GB18030Buf, err := simplifiedchinese.GB18030.NewDecoder().Bytes(strBuf); err == nil {
			if IsUtf8(GB18030Buf) {
				return GB18030Buf
			}
		}
		if HZGB2312Buf, err := simplifiedchinese.HZGB2312.NewDecoder().Bytes(strBuf); err == nil {
			fmt.Println("3")
			if IsUtf8(HZGB2312Buf) {
				return HZGB2312Buf
			}
		}
		return strBuf
	}
}

// IsUtf8 是否是utf8编码
func IsUtf8(buf []byte) bool {
	return utf8.Valid(buf)
}

// Get16MD5Encode 返回一个16位md5加密后的字符串
func Get16MD5Encode(data string) string {
	return GetMD5Encode(data)[8:24]
}

// GetMD5Encode 获取Md5编码
func GetMD5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

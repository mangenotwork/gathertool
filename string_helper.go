/*
	Description : 数据类型相关的操作
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
	"unsafe"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// StringValue 任何类型返回值字符串形式
func StringValue(i interface{}) string {
	if i == nil {
		return ""
	}
	if reflect.ValueOf(i).Kind() == reflect.String {
		return i.(string)
	}
	var buf bytes.Buffer
	stringValue(reflect.ValueOf(i), 0, &buf)
	return buf.String()
}

// StringValueMysql 用于mysql字符拼接使用
func StringValueMysql(i interface{}) string {
	if i == nil {
		return ""
	}
	if reflect.ValueOf(i).Kind() == reflect.String {
		str := i.(string)
		str = strings.Replace(str, `"`, `\"`, -1)
		if len(str) > 1 && string(str[len(str)-1]) == `\` {
			str += `\`
		}
		return `"` + str + `"`
	}
	var buf bytes.Buffer
	stringValue(reflect.ValueOf(i), 0, &buf)
	return buf.String()
}

// stringValue 任何类型返回值字符串形式的实现方法，私有
func stringValue(v reflect.Value, indent int, buf *bytes.Buffer) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Struct:
		buf.WriteString("{\n")
		for i := 0; i < v.Type().NumField(); i++ {
			ft := v.Type().Field(i)
			fv := v.Field(i)
			if ft.Name[0:1] == strings.ToLower(ft.Name[0:1]) {
				continue
			}
			if (fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Slice) && fv.IsNil() {
				continue
			}
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(ft.Name + ": ")
			if tag := ft.Tag.Get("sensitive"); tag == "true" {
				buf.WriteString("<sensitive>")
			} else {
				stringValue(fv, indent+2, buf)
			}
			buf.WriteString(",\n")
		}
		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")

	case reflect.Slice:
		nl, id, id2 := "", "", ""
		if v.Len() > 3 {
			nl, id, id2 = "\n", strings.Repeat(" ", indent), strings.Repeat(" ", indent+2)
		}
		buf.WriteString("[" + nl)
		for i := 0; i < v.Len(); i++ {
			buf.WriteString(id2)
			stringValue(v.Index(i), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString("," + nl)
			}
		}
		buf.WriteString(nl + id + "]")

	case reflect.Map:
		buf.WriteString("{\n")
		for i, k := range v.MapKeys() {
			buf.WriteString(strings.Repeat(" ", indent+2))
			buf.WriteString(k.String() + ": ")
			stringValue(v.MapIndex(k), indent+2, buf)

			if i < v.Len()-1 {
				buf.WriteString(",\n")
			}
		}
		buf.WriteString("\n" + strings.Repeat(" ", indent) + "}")

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatUint(v.Uint(), 10))

	case reflect.Float32, reflect.Float64:
		result := fmt.Sprintf("%f", v.Float())
		// 去除result末尾的0
		for strings.HasSuffix(result, "0") {
			result = strings.TrimSuffix(result, "0")
		}
		if strings.HasSuffix(result, ".") {
			result = strings.TrimSuffix(result, ".")
		}
		buf.WriteString(result)

	default:
		format := "%v"
		switch v.Interface().(type) {
		case string:
			format = "%q"
		}
		_, _ = fmt.Fprintf(buf, format, v.Interface())
	}
}

// OSLine  系统对应换行符
func OSLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// MD5 MD5
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Json2Map json -> map
func Json2Map(str string) (map[string]interface{}, error) {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		return nil, err
	}
	return tempMap, nil
}

// Map2Json map -> json
func Map2Json(m interface{}) (string, error) {
	jsonStr, err := json.Marshal(m)
	return string(jsonStr), err
}

// Any2Map interface{} -> map[string]interface{}
func Any2Map(data interface{}) map[string]interface{} {
	if v, ok := data.(map[string]interface{}); ok {
		return v
	}
	if reflect.ValueOf(data).Kind() == reflect.String {
		dataMap, err := Json2Map(data.(string))
		if err == nil {
			return dataMap
		}
	}
	return nil
}

// Any2String interface{} -> string
func Any2String(data interface{}) string {
	return StringValue(data)
}

// Any2Int interface{} -> int
func Any2Int(data interface{}) int {
	var t2 int
	switch data.(type) {
	case uint:
		t2 = int(data.(uint))
		break
	case int8:
		t2 = int(data.(int8))
		break
	case uint8:
		t2 = int(data.(uint8))
		break
	case int16:
		t2 = int(data.(int16))
		break
	case uint16:
		t2 = int(data.(uint16))
		break
	case int32:
		t2 = int(data.(int32))
		break
	case uint32:
		t2 = int(data.(uint32))
		break
	case int64:
		t2 = int(data.(int64))
		break
	case uint64:
		t2 = int(data.(uint64))
		break
	case float32:
		t2 = int(data.(float32))
		break
	case float64:
		t2 = int(data.(float64))
		break
	case string:
		t2, _ = strconv.Atoi(data.(string))
		break
	default:
		t2 = data.(int)
		break
	}
	return t2
}

// Any2Int64 interface{} -> int64
func Any2Int64(data interface{}) int64 {
	return int64(Any2Int(data))
}

// Any2Arr interface{} -> []interface{}
func Any2Arr(data interface{}) []interface{} {
	if v, ok := data.([]interface{}); ok {
		return v
	}
	return nil
}

// Any2Float64 interface{} -> float64
func Any2Float64(data interface{}) float64 {
	if v, ok := data.(float64); ok {
		return v
	}
	if v, ok := data.(float32); ok {
		return float64(v)
	}
	return 0
}

// Any2Strings interface{} -> []string
func Any2Strings(data interface{}) []string {
	listValue, ok := data.([]interface{})
	if !ok {
		return nil
	}
	keyStringValues := make([]string, len(listValue))
	for i, arg := range listValue {
		keyStringValues[i] = arg.(string)
	}
	return keyStringValues
}

// Any2Json interface{} -> json string
func Any2Json(data interface{}) (string, error) {
	jsonStr, err := json.Marshal(data)
	return string(jsonStr), err
}

// Int2Hex int -> hex
func Int2Hex(i int64) string {
	return fmt.Sprintf("%x", i)
}

// Int642Hex int64 -> hex
func Int642Hex(i int64) string {
	return fmt.Sprintf("%x", i)
}

// Hex2Int hex -> int
func Hex2Int(s string) int {
	n, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		panic("Parse Error")
	}
	n2 := uint8(n)
	return int(*(*int8)(unsafe.Pointer(&n2)))
}

// Hex2Int64 hex -> int
func Hex2Int64(s string) int64 {
	n, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		panic("Parse Error")
	}
	n2 := uint8(n)
	return int64(*(*int8)(unsafe.Pointer(&n2)))
}

// CleaningStr 清理字符串前后空白 和回车 换行符号
func CleaningStr(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	str = strings.Replace(str, "\\n", "", -1)
	//str = strings.Replace(str, "\"", "", -1)
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

// Str2Int64 string -> int64
func Str2Int64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// Str2Int string -> int
func Str2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// Str2Int32 string -> int32
func Str2Int32(str string) int32 {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int32(i)
}

// Str2Float64 string -> float64
func Str2Float64(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return i
}

// Str2Float32 string -> float32
func Str2Float32(str string) float32 {
	i, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0
	}
	return float32(i)
}

// Uint82Str []uint8 -> string
func Uint82Str(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

// Str2Byte string -> []byte
func Str2Byte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Byte2Str []byte -> string
func Byte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Bool2Byte bool -> []byte
func Bool2Byte(b bool) []byte {
	if b == true {
		return []byte{1}
	}
	return []byte{0}
}

// Byte2Bool []byte -> bool
func Byte2Bool(b []byte) bool {
	if len(b) == 0 || bytes.Compare(b, make([]byte, len(b))) == 0 {
		return false
	}
	return true
}

// Int2Byte int -> []byte
func Int2Byte(i int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

// Byte2Int []byte -> int
func Byte2Int(b []byte) int {
	return int(binary.LittleEndian.Uint32(b))
}

// Int642Byte int64 -> []byte
func Int642Byte(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

// Byte2Int64 []byte -> int64
func Byte2Int64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}

// Float322Byte float32 -> []byte
func Float322Byte(f float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, Float322Uint32(f))
	return b
}

// Float322Uint32 float32 -> uint32
func Float322Uint32(f float32) uint32 {
	return math.Float32bits(f)
}

// Byte2Float32 []byte -> float32
func Byte2Float32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

// Float642Byte float64 -> []byte
func Float642Byte(f float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, Float642Uint64(f))
	return b
}

// Float642Uint64 float64 -> uint64
func Float642Uint64(f float64) uint64 {
	return math.Float64bits(f)
}

// Byte2Float64 []byte -> float64
func Byte2Float64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(b))
}

// Struct2Map  struct -> map[string]interface{}
func Struct2Map(obj interface{}) map[string]interface{} {
	rt, rv := reflect.TypeOf(obj), reflect.ValueOf(obj)
	if rt != nil && rt.Kind() != reflect.Struct {
		return make(map[string]interface{}, 0)
	}

	out := make(map[string]interface{}, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// Unexported fields, access not allowed
		if field.PkgPath != "" {
			continue
		}

		var fieldName string
		if tagVal, ok := field.Tag.Lookup("json"); ok {
			// Honor the special "-" in json attribute
			if strings.HasPrefix(tagVal, "-") {
				continue
			}
			fieldName = tagVal
		} else {
			fieldName = field.Name
		}

		val := valueToInterface(rv.Field(i))
		if val != nil {
			out[fieldName] = val
		}
	}

	return out
}

func valueToInterface(value reflect.Value) interface{} {
	if !value.IsValid() {
		return nil
	}

	switch value.Type().Kind() {
	case reflect.Struct:
		return Struct2Map(value.Interface())

	case reflect.Ptr:
		if !value.IsNil() {
			return valueToInterface(value.Elem())
		}

	case reflect.Array:
	case reflect.Slice:
		arr := make([]interface{}, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			val := valueToInterface(value.Index(i))
			if val != nil {
				arr = append(arr, val)
			}
		}
		return arr

	case reflect.Map:
		m := make(map[string]interface{}, value.Len())
		for _, k := range value.MapKeys() {
			v := value.MapIndex(k)
			m[k.String()] = valueToInterface(v)
		}
		return m

	default:
		return value.Interface()
	}

	return nil
}

// EncodeByte encode byte
func EncodeByte(v interface{}) []byte {
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
func DecodeByte(b []byte) (interface{}, error) {
	var values interface{}
	buf := bytes.NewBuffer(b)
	err := binary.Read(buf, binary.BigEndian, values)
	return values, err
}

// Byte2Bit []byte -> []uint8 (bit)
func Byte2Bit(b []byte) []uint8 {
	bits := make([]uint8, 0)
	for _, v := range b {
		bits = bits2Uint(bits, uint(v), 8)
	}
	return bits
}

// bits2Uint bits2Uint
func bits2Uint(bits []uint8, ui uint, l int) []uint8 {
	a := make([]uint8, l)
	for i := l - 1; i >= 0; i-- {
		a[i] = uint8(ui & 1)
		ui >>= 1
	}
	if bits != nil {
		return append(bits, a...)
	}
	return a
}

// Bit2Byte []uint8 -> []byte
func Bit2Byte(b []uint8) []byte {
	if len(b)%8 != 0 {
		for i := 0; i < len(b)%8; i++ {
			b = append(b, 0)
		}
	}
	by := make([]byte, 0)
	for i := 0; i < len(b); i += 8 {
		by = append(b, byte(bitsToUint(b[i:i+8])))
	}
	return by
}

// bitsToUint bitsToUint
func bitsToUint(bits []uint8) uint {
	v := uint(0)
	for _, i := range bits {
		v = v<<1 | uint(i)
	}
	return v
}

// FileSizeFormat 字节的单位转换 保留两位小数
func FileSizeFormat(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// deepCopy 深copy
func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func DeepCopy(dst, src interface{}) error {
	return deepCopy(dst, src)
}

// Struct2MapV2 Struct  ->  map
// hasValue=true表示字段值不管是否存在都转换成map
// hasValue=false表示字段为空或者不为0则转换成map
func Struct2MapV2(obj interface{}, hasValue bool) (map[string]interface{}, error) {
	mp := make(map[string]interface{})
	value := reflect.ValueOf(obj).Elem()
	typeOf := reflect.TypeOf(obj).Elem()
	for i := 0; i < value.NumField(); i++ {
		switch value.Field(i).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if hasValue {
				if value.Field(i).Int() != 0 {
					mp[typeOf.Field(i).Name] = value.Field(i).Int()
				}
			} else {
				mp[typeOf.Field(i).Name] = value.Field(i).Int()
			}

		case reflect.String:
			if hasValue {
				if len(value.Field(i).String()) != 0 {
					mp[typeOf.Field(i).Name] = value.Field(i).String()
				}
			} else {
				mp[typeOf.Field(i).Name] = value.Field(i).String()
			}

		case reflect.Float32, reflect.Float64:
			if hasValue {
				if len(value.Field(i).String()) != 0 {
					mp[typeOf.Field(i).Name] = value.Field(i).Float()
				}
			} else {
				mp[typeOf.Field(i).Name] = value.Field(i).Float()
			}

		case reflect.Bool:
			if hasValue {
				if len(value.Field(i).String()) != 0 {
					mp[typeOf.Field(i).Name] = value.Field(i).Bool()
				}
			} else {
				mp[typeOf.Field(i).Name] = value.Field(i).Bool()
			}

		default:
			return mp, fmt.Errorf("数据类型不匹配")
		}
	}

	return mp, nil
}

// Struct2MapV3 struct -> map
func Struct2MapV3(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)
	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

// PanicToError panic -> error
func PanicToError(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic error: %v", r)
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

// ============================================  转码

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

func UnicodeDec(raw string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(raw), `\\u`, `\u`, -1))
	if err != nil {
		return ""
	}
	return str
}

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

// convert
func convert(dstCharset string, srcCharset string, src string) (dst string, err error) {
	if dstCharset == srcCharset {
		return src, nil
	}
	dst = src
	// Converting `src` to UTF-8.
	if srcCharset != "UTF-8" {
		if e := getEncoding(srcCharset); e != nil {
			tmp, err := ioutil.ReadAll(
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
	// Do the converting from UTF-8 to `dstCharset`.
	if dstCharset != "UTF-8" {
		if e := getEncoding(dstCharset); e != nil {
			tmp, err := ioutil.ReadAll(
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

// ================================================ Set 集合

// 可以用于去重
type Set map[string]struct{}

func (s Set) Has(key string) bool {
	_, ok := s[key]
	return ok
}

func (s Set) Add(key string) {
	s[key] = struct{}{}
}

func (s Set) Delete(key string) {
	delete(s, key)
}

// ================================================ Stack 栈

type Stack struct {
	data map[int]interface{}
}

func New() *Stack {
	s := new(Stack)
	s.data = make(map[int]interface{})
	return s
}

func (s *Stack) Push(data interface{}) {
	s.data[len(s.data)] = data
}

func (s *Stack) Pop() {
	delete(s.data, len(s.data)-1)
}

func (s *Stack) String() string {
	info := ""
	for i := 0; i < len(s.data); i++ {
		info = info + "[" + StringValue(s.data[i]) + "]"
	}
	return info
}

// IsContainStr  字符串是否等于items中的某个元素
func IsContainStr(items []string, item string) bool {
	for i := 0; i < len(items); i++ {
		if items[i] == item {
			return true
		}
	}
	return false
}

// FileMd5  file md5   文件md5
func FileMd5(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	md5hash := md5.New()
	if _, err := io.Copy(md5hash, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5hash.Sum(nil)), nil
}

// PathExists 目录不存在则创建
func PathExists(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}
}

// ByteToBinaryString  字节 -> 二进制字符串
func ByteToBinaryString(data byte) (str string) {
	var a byte
	for i := 0; i < 8; i++ {
		a = data
		data <<= 1
		data >>= 1
		switch a {
		case data:
			str += "0"
		default:
			str += "1"
		}
		data <<= 1
	}
	return str
}

// StrDuplicates  数组，切片去重和去空串
func StrDuplicates(a []string) []string {
	m := make(map[string]struct{})
	ret := make([]string, 0, len(a))
	for i := 0; i < len(a); i++ {
		if a[i] == "" {
			continue
		}
		if _, ok := m[a[i]]; !ok {
			m[a[i]] = struct{}{}
			ret = append(ret, a[i])
		}
	}
	return ret
}

// IsElementStr 判断字符串是否与数组里的某个字符串相同
func IsElementStr(listData []string, element string) bool {
	for _, k := range listData {
		if k == element {
			return true
		}
	}
	return false
}

// windowsPath windows平台需要转一下
func windowsPath(path string) string {
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	return path
}

// GetNowPath 获取当前运行路径
func GetNowPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return windowsPath(path)
}

// FileMd5sum 文件 Md5
func FileMd5sum(fileName string) string {
	fin, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		Info(fileName, err)
		return ""
	}
	defer fin.Close()
	Buf, buferr := ioutil.ReadFile(fileName)
	if buferr != nil {
		Info(fileName, buferr)
		return ""
	}
	m := md5.Sum(Buf)
	return hex.EncodeToString(m[:16])
}

// SearchBytesIndex []byte 字节切片 循环查找
func SearchBytesIndex(bSrc []byte, b byte) int {
	for i := 0; i < len(bSrc); i++ {
		if bSrc[i] == b {
			return i
		}
	}
	return -1
}

// IF 三元表达式
// use: IF(a>b, a, b).(int)
func IF(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

// CopySlice Copy slice
func CopySlice(s []interface{}) []interface{} {
	return append(s[:0:0], s...)
}

func CopySliceStr(s []string) []string {
	return append(s[:0:0], s...)
}

func CopySliceInt(s []int) []int {
	return append(s[:0:0], s...)
}

func CopySliceInt64(s []int64) []int64 {
	return append(s[:0:0], s...)
}

func IsInSlice(s []interface{}, v interface{}) bool {
	for i := range s {
		if s[i] == v {
			return true
		}
	}
	return false
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

// MapStr2Any map[string]string -> map[string]interface{}
func MapStr2Any(m map[string]string) map[string]interface{} {
	dest := make(map[string]interface{})
	for k, v := range m {
		dest[k] = interface{}(v)
	}
	return dest
}

func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func Pwd() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

func Chdir(dir string) error {
	return os.Chdir(dir)
}

func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// 字节换算
const (
	KiB = 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
	TiB = GiB * 1024
)

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

func MapCopy(data map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{}, len(data))
	for k, v := range data {
		copy[k] = v
	}
	return
}

func MapMergeCopy(src ...map[string]interface{}) (copy map[string]interface{}) {
	copy = make(map[string]interface{})
	for _, m := range src {
		for k, v := range m {
			copy[k] = v
		}
	}
	return
}

// Map2Slice Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]
func Map2Slice(data interface{}) []interface{} {
	var (
		reflectValue = reflect.ValueOf(data)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Map:
		array := make([]interface{}, 0)
		for _, key := range reflectValue.MapKeys() {
			array = append(array, key.Interface())
			array = append(array, reflectValue.MapIndex(key).Interface())
		}
		return array
	}
	return nil
}

func SliceCopy(data []interface{}) []interface{} {
	newData := make([]interface{}, len(data))
	copy(newData, data)
	return newData
}

// Slice2Map ["K1", "v1", "K2", "v2"] => {"K1": "v1", "K2": "v2"}
// ["K1", "v1", "K2"]       => nil
func Slice2Map(slice interface{}) map[string]interface{} {
	var (
		reflectValue = reflect.ValueOf(slice)
		reflectKind  = reflectValue.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	switch reflectKind {
	case reflect.Slice, reflect.Array:
		length := reflectValue.Len()
		if length%2 != 0 {
			return nil
		}
		data := make(map[string]interface{})
		for i := 0; i < reflectValue.Len(); i += 2 {
			data[Any2String(reflectValue.Index(i).Interface())] = reflectValue.Index(i + 1).Interface()
		}
		return data
	}
	return nil
}

func GzipCompress(src []byte) []byte {
	var in bytes.Buffer
	w := gzip.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

func GzipDecompress(src []byte) []byte {
	dst := make([]byte, 0)
	br := bytes.NewReader(src)
	gr, err := gzip.NewReader(br)
	if err != nil {
		return dst
	}
	defer gr.Close()
	tmp, err := ioutil.ReadAll(gr)
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

// AbPathByCaller 获取当前执行文件绝对路径（go run）
func AbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return path.Join(abPath, "../../../")
}

// GetWD 获取当前工作目录
func GetWD() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
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
		//if HZGB2312Buf, err := simplifiedchinese.HZGB2312.NewEncoder().Bytes(strBuf); err == nil {
		//	if IsUtf8(HZGB2312Buf) {
		//		return HZGB2312Buf
		//	}
		//}
		return strBuf
	} else {
		return strBuf
	}
}

// ByteToUTF8    byte -> utf8 byte
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
		//if HZGB2312Buf, err := simplifiedchinese.HZGB2312.NewDecoder().Bytes(strBuf); err == nil {
		//	fmt.Println("3")
		//	if IsUtf8(HZGB2312Buf) {
		//		return HZGB2312Buf
		//	}
		//}
		return strBuf
	}
}

func IsUtf8(buf []byte) bool {
	return utf8.Valid(buf)
}

// ====================================  json find

// JsonFind 按路径寻找指定json值
// 用法参考  ./_examples/json/main.go
// @find : 寻找路径，与目录的url类似， 下面是一个例子：
// json:  {a:[{b:1},{b:2}]}
// find=/a/[0]  =>   {b:1}
// find=a/[0]/b  =>   1
func JsonFind(jsonStr, find string) (interface{}, error) {
	if !IsJson(jsonStr) {
		return nil, fmt.Errorf("不是标准的Json格式")
	}
	jxList := strings.Split(find, "/")
	jxLen := len(jxList)
	var (
		data  = Any2Map(jsonStr)
		value interface{}
		err   error
	)
	for i := 0; i < jxLen; i++ {
		l := len(jxList[i])
		if l > 2 && string(jxList[i][0]) == "[" && string(jxList[i][l-1]) == "]" {
			numStr := jxList[i][1 : l-1]
			dataList := Any2Arr(value)
			value = dataList[Any2Int(numStr)]
			data, err = interface2Map(value)
			if err != nil {
				continue
			}
		} else {
			if IsHaveKey(data, jxList[i]) {
				value = data[jxList[i]]
				data, err = interface2Map(value)
				if err != nil {
					continue
				}
			} else {
				value = nil
			}
		}
	}
	return value, nil
}

// JsonFind2Json 寻找json,输出 json格式字符串
func JsonFind2Json(jsonStr, find string) (string, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return "", err
	}
	return Map2Json(value)
}

// JsonFind2Map 寻找json,输出 map[string]interface{}
func JsonFind2Map(jsonStr, find string) (map[string]interface{}, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return nil, err
	}
	return Any2Map(value), nil
}

// JsonFind2Arr 寻找json,输出 []interface{}
func JsonFind2Arr(jsonStr, find string) ([]interface{}, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return nil, err
	}
	return Any2Arr(value), nil
}

// JsonFind2Str 寻找json,输出字符串
func JsonFind2Str(jsonStr, find string) (string, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return "", err
	}
	return Any2String(value), nil
}

// JsonFind2Int 寻找json,输出int
func JsonFind2Int(jsonStr, find string) (int, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return 0, err
	}
	return Any2Int(value), nil
}

// JsonFind2Int64 寻找json,输出int64
func JsonFind2Int64(jsonStr, find string) (int64, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return 0, err
	}
	return Any2Int64(value), nil
}

// IsJson 是否是json格式
func IsJson(str string) bool {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		return false
	}
	return true
}

// IsHaveKey map[string]interface{} 是否存在 输入的key
func IsHaveKey(data map[string]interface{}, key string) bool {
	_, ok := data[key]
	return ok
}

// Any2Map interface{} -> map[string]interface{}
func interface2Map(data interface{}) (map[string]interface{}, error) {
	if v, ok := data.(map[string]interface{}); ok {
		return v, nil
	}
	if reflect.ValueOf(data).Kind() == reflect.String {
		return Json2Map(data.(string))
	}
	return nil, fmt.Errorf("not map type")
}

// Int642Str int64 -> string
func Int642Str(i int64) string {
	return strconv.FormatInt(i, 10)
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

// ========================================================  文件相关处理

// GetAllFile 获取目录下的所有文件
func GetAllFile(pathname string) ([]string, error) {
	s := make([]string, 0)
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		Error("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)
	if start < 0 || start > length {
		Error("start is wrong")
		return ""
	}
	if end < start || end > length {
		Error("end is wrong")
		return ""
	}
	return string(rs[start:end])
}

// DeCompressZIP zip解压文件
func DeCompressZIP(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		filename := dest + file.Name
		err = os.MkdirAll(subString(filename, 0, strings.LastIndex(filename, "/")), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		_ = w.Close()
		_ = rc.Close()
	}
	return nil
}

// DeCompressTAR tar 解压文件
func DeCompressTAR(tarFile, dest string) error {
	file, err := os.Open(tarFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		_ = file.Close()
	}()
	tr := tar.NewReader(file)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		filename := dest + hdr.Name
		err = os.MkdirAll(subString(filename, 0, strings.LastIndex(filename, "/")), 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, tr)
		if err != nil {
			return err
		}
		_ = w.Close()
	}
	return nil
}

// DecompressionZipFile zip压缩文件
func DecompressionZipFile(src, dest string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		_ = reader.Close()
	}()
	for _, file := range reader.File {
		filePath := path.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := file.Open()
			if err != nil {
				return err
			}
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
			_ = inFile.Close()
			_ = outFile.Close()
		}
	}
	return nil
}

// CompressFiles 压缩很多文件
// files 文件数组，可以是不同dir下的文件或者文件夹
// dest 压缩文件存放地址
func CompressFiles(files []string, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compressFiles(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compressFiles(filePath string, prefix string, zw *zip.Writer) error {
	file, err := os.Open(filePath)
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f := file.Name() + "/" + fi.Name()
			if err != nil {
				return err
			}
			err = compressFiles(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// CompressDirZip 压缩目录
func CompressDirZip(src, outFile string) error {
	// 预防：旧文件无法覆盖
	_ = os.RemoveAll(outFile)
	// 创建：zip文件
	zipFile, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	// 打开：zip文件
	archive := zip.NewWriter(zipFile)
	defer archive.Close()
	// 遍历路径信息
	filepath.Walk(src, func(path string, info os.FileInfo, _ error) error {
		// 如果是源路径，提前进行下一个遍历
		if path == src {
			return nil
		}
		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src+`/`)
		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}
		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
	return nil
}

// OutJsonFile 将data输出到json文件
func OutJsonFile(data interface{}, fileName string) error {
	var (
		f   *os.File
		err error
	)
	if Exists(fileName) { //如果文件存在
		f, err = os.OpenFile(fileName, os.O_APPEND, 0666) //打开文件
	} else {
		f, err = os.Create(fileName) //创建文件
	}
	if err != nil {
		return err
	}
	str, err := Any2Json(data)
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, str)
	if err != nil {
		return err
	}
	return nil
}

// ===================================================  雪花Id

type IdWorker struct {
	startTime             int64
	workerIdBits          uint
	datacenterIdBits      uint
	maxWorkerId           int64
	maxDatacenterId       int64
	sequenceBits          uint
	workerIdLeftShift     uint
	datacenterIdLeftShift uint
	timestampLeftShift    uint
	sequenceMask          int64
	workerId              int64
	datacenterId          int64
	sequence              int64
	lastTimestamp         int64
	signMask              int64
	idLock                *sync.Mutex
}

func (idw *IdWorker) InitIdWorker(workerId, datacenterId int64) error {
	var baseValue int64 = -1
	idw.startTime = 1463834116272
	idw.workerIdBits = 5
	idw.datacenterIdBits = 5
	idw.maxWorkerId = baseValue ^ (baseValue << idw.workerIdBits)
	idw.maxDatacenterId = baseValue ^ (baseValue << idw.datacenterIdBits)
	idw.sequenceBits = 12
	idw.workerIdLeftShift = idw.sequenceBits
	idw.datacenterIdLeftShift = idw.workerIdBits + idw.workerIdLeftShift
	idw.timestampLeftShift = idw.datacenterIdBits + idw.datacenterIdLeftShift
	idw.sequenceMask = baseValue ^ (baseValue << idw.sequenceBits)
	idw.sequence = 0
	idw.lastTimestamp = -1
	idw.signMask = ^baseValue + 1
	idw.idLock = &sync.Mutex{}
	if idw.workerId < 0 || idw.workerId > idw.maxWorkerId {
		return fmt.Errorf("workerId[%v] is less than 0 or greater than maxWorkerId[%v].",
			workerId, datacenterId)
	}
	if idw.datacenterId < 0 || idw.datacenterId > idw.maxDatacenterId {
		return fmt.Errorf("datacenterId[%d] is less than 0 or greater than maxDatacenterId[%d].",
			workerId, datacenterId)
	}
	idw.workerId = workerId
	idw.datacenterId = datacenterId
	return nil
}

// NextId 返回一个唯一的 INT64 ID
func (idw *IdWorker) NextId() (int64, error) {
	idw.idLock.Lock()
	timestamp := time.Now().UnixNano()
	if timestamp < idw.lastTimestamp {
		return -1, fmt.Errorf(fmt.Sprintf("Clock moved backwards.  Refusing to generate id for %d milliseconds",
			idw.lastTimestamp-timestamp))
	}
	if timestamp == idw.lastTimestamp {
		idw.sequence = (idw.sequence + 1) & idw.sequenceMask
		if idw.sequence == 0 {
			timestamp = idw.tilNextMillis()
			idw.sequence = 0
		}
	} else {
		idw.sequence = 0
	}
	idw.lastTimestamp = timestamp
	idw.idLock.Unlock()
	id := ((timestamp - idw.startTime) << idw.timestampLeftShift) |
		(idw.datacenterId << idw.datacenterIdLeftShift) |
		(idw.workerId << idw.workerIdLeftShift) |
		idw.sequence
	if id < 0 {
		id = -id
	}
	return id, nil
}

// tilNextMillis
func (idw *IdWorker) tilNextMillis() int64 {
	timestamp := time.Now().UnixNano()
	if timestamp <= idw.lastTimestamp {
		timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	}
	return timestamp
}

func ID64() (int64, error) {
	currWorker := &IdWorker{}
	err := currWorker.InitIdWorker(1000, 2)
	if err != nil {
		return 0, err
	}
	return currWorker.NextId()
}

func ID() int64 {
	id, _ := ID64()
	return id
}

func IDStr() string {
	currWorker := &IdWorker{}
	err := currWorker.InitIdWorker(1000, 2)
	if err != nil {
		return ""
	}
	id, err := currWorker.NextId()
	if err != nil {
		return ""
	}
	return Int642Str(id)
}

func IDMd5() string {
	return Get16MD5Encode(IDStr())
}

// IP2Int64 IP str ==> int64
func IP2Int64(ip string) int64 {
	address := net.ParseIP(ip)
	if address == nil {
		Error("ip地址不正确")
		return 0
	}
	bits := strings.Split(ip, ".")
	b0, b1, b2, b3 := 0, 0, 0, 0
	if len(bits) >= 1 {
		b0, _ = strconv.Atoi(bits[0])
	}
	if len(bits) >= 2 {
		b1, _ = strconv.Atoi(bits[1])
	}
	if len(bits) >= 3 {
		b2, _ = strconv.Atoi(bits[2])
	}
	if len(bits) >= 4 {
		b3, _ = strconv.Atoi(bits[3])
	}
	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

// IP2Binary IP str ==> binary int64
func IP2Binary(ip string) string {
	rse := IP2Int64(ip)
	return strconv.FormatInt(rse, 2)
}

// UInt32ToIP  uint32 ==> net.IP
func UInt32ToIP(ip uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ip & 0xFF)
	bytes[1] = byte((ip >> 8) & 0xFF)
	bytes[2] = byte((ip >> 16) & 0xFF)
	bytes[3] = byte((ip >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

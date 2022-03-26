/*
	Description : 数据类型相关的操作
	Author : ManGe
*/

package gathertool

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// 全局json
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// StringValue 任何类型返回值字符串形式
func StringValue(i interface{}) string {
	if i == nil {
		return ""
	}
	if reflect.ValueOf(i).Kind() == reflect.String{
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
	if reflect.ValueOf(i).Kind() == reflect.String{
		return "'"+i.(string)+"'"
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

	default:
		format := "%v"
		switch v.Interface().(type) {
		case string:
			format = "%q"
		}
		_,_= fmt.Fprintf(buf, format, v.Interface())
	}
}

// OSLine  系统对应换行符
func OSLine() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

// MD5
func MD5(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// Json2Map json -> map
func Json2Map(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		Error(err)
		return nil
	}
	return tempMap
}

// Map2Json map -> json
func Map2Json(m map[string]interface{}) (string, error) {
	jsonStr,err :=json.Marshal(m)
	return string(jsonStr), err
}

// Any2Map interface{} -> map[string]interface{}
func Any2Map(data interface{}) map[string]interface{}{
	if v, ok := data.(map[string]interface{}); ok {
		return v
	}
	if reflect.ValueOf(data).Kind() == reflect.String{
		return Json2Map(data.(string))
	}
	return nil
}

// Any2String interface{} -> string
func Any2String(data interface{}) string {
	return StringValue(data)
}

// Any2Int interface{} -> int
func Any2Int(data interface{}) int {
	if v, ok := data.(float64); ok {
		return int(v)
	}
	if v, ok := data.(float32); ok {
		return int(v)
	}
	if v, ok := data.(string); ok {
		return Str2Int(v)
	}
	if v, ok := data.(int64); ok {
		return int(v)
	}
	if v, ok := data.(int); ok {
		return v
	}
	return 0
}

// Any2int64 interface{} -> int64
func Any2int64(data interface{}) int64 {
	if v, ok := data.(float64); ok {
		return int64(v)
	}
	if v, ok := data.(float32); ok {
		return int64(v)
	}
	if v, ok := data.(string); ok {
		return Str2Int64(v)
	}
	if v, ok := data.(int32); ok {
		return int64(v)
	}
	if v, ok := data.(int64); ok {
		return v
	}
	return 0
}

// Any2Arr interface{} -> []interface{}
func Any2Arr(data interface{}) []interface{} {
	if v,ok := data.([]interface{}); ok {
		return v
	}
	return nil
}

// Any2Float64 interface{} -> float64
func Any2Float64(data interface{}) float64 {
	if v,ok := data.(float64); ok {
		return v
	}
	if v,ok := data.(float32); ok {
		return float64(v)
	}
	return 0
}

// Any2Strings interface{} -> []string
func Any2Strings(data interface{}) []string{
	listValue,ok := data.([]interface{})
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
	jsonStr,err :=json.Marshal(data)
	return string(jsonStr), err
}

// Int2Hex int -> hex
func Int2Hex(i int64) string {
	return fmt.Sprintf("%x",i)
}

// Int642Hex int64 -> hex
func Int642Hex(i int64) string {
	return fmt.Sprintf("%x",i)
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
func CleaningStr(str string) string{
	str = strings.Replace(str, "\n","", -1)
	str = strings.Replace(str, "\r","", -1)
	str = strings.Replace(str, "\\n","", -1)
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

// Str2Int string -> int32
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

// Str2Bytes string -> []byte
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

// EncodeByte encode byte
func EncodeByte(v interface{}) []byte {
	switch value := v.(type) {
	case int, int8, int16, int32 :
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

// Struct2Map Struct  ->  map
// hasValue=true表示字段值不管是否存在都转换成map
// hasValue=false表示字段为空或者不为0则转换成map
func Struct2Map(obj interface{}, hasValue bool) (map[string]interface{}, error) {
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
	GBK = Charset("GBK")
	GB2312 = Charset("GB2312")
)

// ConvertByte2String 编码转换
func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes,_=simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str= string(decodeBytes)

	case GBK:
		var decodeBytes,_=simplifiedchinese.GBK.NewDecoder().Bytes(byte)
		str= string(decodeBytes)

	case GB2312:
		var decodeBytes,_=simplifiedchinese.HZGB2312.NewDecoder().Bytes(byte)
		str= string(decodeBytes)

	case UTF8:
		fallthrough

	default:
		str = string(byte)
	}

	return str
}

// UnicodeDec
func UnicodeDec(raw string) string {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(raw), `\\u`, `\u`, -1))
	if err != nil {
		return ""
	}
	return str
}

// UnicodeDecByte
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
func Base64Encode(str string) string{
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64 解码
func Base64Decode(str string) (string,error){
	b, err := base64.StdEncoding.DecodeString(str)
	return string(b), err
}

// Base64UrlEncode base64 url 编码
func Base64UrlEncode(str string) string{
	return base64.URLEncoding.EncodeToString([]byte(str))
}

// Base64UrlDecode base64 url 解码
func Base64UrlDecode(str string) (string,error){
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

// ToUTF8
func ToUTF8(srcCharset string, src string) (dst string, err error) {
	return convert("UTF-8", srcCharset, src)
}

// UTF8To
func UTF8To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "UTF-8", src)
}

// ToUTF16
func ToUTF16(srcCharset string, src string) (dst string, err error) {
	return convert("UTF-16", srcCharset, src)
}

// UTF16To
func UTF16To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "UTF-16", src)
}

// ToBIG5
func ToBIG5(srcCharset string, src string) (dst string, err error) {
	return convert("big5", srcCharset, src)
}

// BIG5To
func BIG5To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "big5", src)
}

// ToGDK
func ToGDK(srcCharset string, src string) (dst string, err error) {
	return convert("gbk", srcCharset, src)
}

// GDKTo
func GDKTo(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "gbk", src)
}

// ToGB18030
func ToGB18030(srcCharset string, src string) (dst string, err error) {
	return convert("gb18030", srcCharset, src)
}

// GB18030To
func GB18030To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "gb18030", src)
}

// ToGB2312
func ToGB2312(srcCharset string, src string) (dst string, err error) {
	return convert("GB2312", srcCharset, src)
}

// GB2312To
func GB2312To(dstCharset string, src string) (dst string, err error) {
	return convert(dstCharset, "GB2312", src)
}

// ToHZGB2312
func ToHZGB2312(srcCharset string, src string) (dst string, err error) {
	return convert("HZGB2312", srcCharset, src)
}

// HZGB2312To
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

func New() *Stack{
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
	for i:=0; i<len(items); i++ {
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
func StrDuplicates(a []string) []string{
	m := make(map[string]struct{})
	ret := make([]string, 0, len(a))
	for i:=0; i < len(a); i++{
		if a[i] == ""{
			continue
		}
		if _,ok := m[a[i]]; !ok {
			m[a[i]] = struct{}{}
			ret = append(ret, a[i])
		}
	}
	return ret
}

// IsElementStr 判断字符串是否与数组里的某个字符串相同
func IsElementStr(listData []string, element string) bool{
	for _,k := range listData{
		if k == element{ return true }
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

// CopySliceStr
func CopySliceStr(s []string) []string{
	return append(s[:0:0], s...)
}

// CopySliceInt
func CopySliceInt(s []int) []int {
	return append(s[:0:0], s...)
}

// CopySliceInt64
func CopySliceInt64(s []int64) []int64 {
	return append(s[:0:0], s...)
}

// IsInSlice
func IsInSlice(s []interface{}, v interface{})  bool {
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
	for k,v := range m {
		dest[k] = interface{}(v)
	}
	return dest
}

// Exists
func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// IsDir
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Pwd
func Pwd() string {
	path, err := os.Getwd()
	if err != nil {
		return ""
	}
	return path
}

// Chdir
func Chdir(dir string) error {
	return os.Chdir(dir)
}

// IsFile
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

// HumanFriendlyTraffic
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

// StrToSize
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


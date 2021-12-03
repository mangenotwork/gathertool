/*
	Description : 数据类型相关的操作
	Author : ManGe
	Version : v0.2
	Date : 2021-10-13
*/

package gathertool

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
				// ignore unexported fields
				continue
			}
			if (fv.Kind() == reflect.Ptr || fv.Kind() == reflect.Slice) && fv.IsNil() {
				// ignore unset fields
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
		loger(err)
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
	if v, ok := data.(int); ok {
		return v
	}
	return 0
}

// Any2int64 interface{} -> int64
func Any2int64(data interface{}) int64 {
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

// Str2Float64 string -> float64
func Str2Float64(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return i
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
func Str2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Bytes2Str []byte -> string
func Bytes2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
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
			loger("Panic error: %v", r)
		}
	}()
}

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

// ================================================ Set
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

// ================================================ Stack
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

// Byte2Str []byte -> string
func Byte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
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

// StrDuplicates数组，切片去重和去空串
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
		log.Println(fileName, err)
		return ""
	}
	defer fin.Close()
	Buf, buferr := ioutil.ReadFile(fileName)
	if buferr != nil {
		log.Println(fileName, buferr)
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


const (
	KiB = 1024
	MiB = KiB * 1024
	GiB = MiB * 1024
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
	return fmt.Sprintf("%.2f GiB", float32(bytes)/GiB)
}



// TODO 二进制字符串 -> 字符串



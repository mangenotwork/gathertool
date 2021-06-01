/*
	Description : 数据类型相关的操作
	Author : ManGe
	Version : v0.1
	Date : 2021-04-26
*/

package gathertool

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"unsafe"
)

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
		fmt.Fprintf(buf, format, v.Interface())
	}
}


// MD5
func MD5(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// json转map函数，通用
func Json2Map(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		loger(err)
		return nil
	}
	return tempMap
}

// interface{} -> map[string]interface{}
func Any2Map(data interface{}) map[string]interface{}{
	return data.(map[string]interface{})
}

// interface{} -> string
func Any2String(data interface{}) string {
	return StringValue(data)
}

// interface{} -> int
func Any2Int(data interface{}) int {
	return data.(int)
}

// interface{} -> int64
func Any2int64(data interface{}) int64 {
	return data.(int64)
}

// interface{} -> []interface{}
func Any2AnyList(data interface{}) []interface{} {
	return data.([]interface{})
}

// interface{} -> float64
func Any2Float64(data interface{}) float64 {
	return data.(float64)
}

// interface{} -> []string
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

// CleaningStr 清理字符串前后空白 和回车 换行符号
func CleaningStr(str string) string{
	str = strings.Replace(str, "\n","", -1)
	str = strings.Replace(str, "\r","", -1)
	str = strings.Replace(str, "\\n","", -1)
	str = strings.Replace(str, "\"", "", -1)
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

// string -> int64
func Str2Int64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// string -> float64
func Str2Float64(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return i
}


// []uint8 -> string
func Uint82Str(bs []uint8) string {
	ba := []byte{}
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

// 字节的单位转换 保留两位小数
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

// 深copy
func deepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// Struct  ->  map
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

// panic -> error
func PanicToError(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Panic error: %v", r)
		}
	}()

	fn()
	return
}


type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
	GBK = Charset("GBK")
	GB2312 = Charset("GB2312")
)

// 编码转换
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


// Unicode 转码
func UnescapeUnicode(raw []byte) ([]byte, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}


// base64 编码
func Base64Encode(str string) string{
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// base64 解码
func Base64Decode(str string) (string,error){
	b, err := base64.StdEncoding.DecodeString(str)
	return string(b), err
}

// base64 url 编码
func Base64UrlEncode(str string) string{
	return base64.URLEncoding.EncodeToString([]byte(str))
}

// base64 url 解码
func Base64UrlDecode(str string) (string,error){
	b, err := base64.URLEncoding.DecodeString(str)
	return string(b), err
}


// ======== Set
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


// ======== Stack
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

// 目录不存在则创建
func PathExists(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(path, 0777)
	}
}

// []byte -> string
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


// 数组，切片去重和去空串
func StrsDuplicates(a []string) []string{
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


//判断字符串是否与数组里的某个字符串相同
func IsElementStr(list []string, element string) bool{
	fmt.Println(list)
	fmt.Println(element)
	for _,k := range(list){
		if k == element{ return true }
	}
	return false
}

//windows平台需要转一下
func windowsPath(path string) string {
	if runtime.GOOS == "windows" {
		path = strings.Replace(path, "\\", "/", -1)
	}
	return path
}


// 获取当前运行路径
func GetNowPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return windowsPath(path)
}


// 文件 Md5
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


// TODO 二进制字符串 -> 字符串
/*
*	Description : 数据类型转换
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

// Json2Map 数据类型转换 json -> map
func Json2Map(str string) (map[string]any, error) {
	var tempMap map[string]any
	err := json.Unmarshal([]byte(str), &tempMap)
	return tempMap, err
}

// Any2Map 数据类型转换 any -> map[string]any
func Any2Map(data any) map[string]any {
	if v, ok := data.(map[string]any); ok {
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

// Any2String 数据类型转换 any -> string
func Any2String(i any) string {
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
		buf.WriteString(strconv.FormatInt(v.Int(), 10))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
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

// Any2Int 数据类型转换 any -> int
func Any2Int(data any) int {
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

// Any2Int64 数据类型转换 any -> int64
func Any2Int64(data any) int64 {
	return int64(Any2Int(data))
}

// Any2Arr 数据类型转换 any -> []any
func Any2Arr(data any) []any {
	if v, ok := data.([]any); ok {
		return v
	}
	return nil
}

// Any2Float64 数据类型转换 any -> float64
func Any2Float64(data any) float64 {
	if v, ok := data.(float64); ok {
		return v
	}
	if v, ok := data.(float32); ok {
		return float64(v)
	}
	return 0
}

// Any2Strings 数据类型转换 any -> []string
func Any2Strings(data any) []string {
	listValue, ok := data.([]any)
	if !ok {
		return nil
	}
	keyStringValues := make([]string, len(listValue))
	for i, arg := range listValue {
		keyStringValues[i] = arg.(string)
	}
	return keyStringValues
}

// Any2Json 数据类型转换 any -> json string
func Any2Json(data any) (string, error) {
	jsonStr, err := json.Marshal(data)
	return string(jsonStr), err
}

// Int2Hex 数据类型转换 int -> hex
func Int2Hex(i int64) string {
	return fmt.Sprintf("%x", i)
}

// Int642Hex 数据类型转换 int64 -> hex
func Int642Hex(i int64) string {
	return fmt.Sprintf("%x", i)
}

// Hex2Int 数据类型转换 hex -> int
func Hex2Int(s string) int {
	n, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		Error(err)
		return 0
	}

	n2 := uint8(n)
	return int(*(*int8)(unsafe.Pointer(&n2)))
}

// Hex2Int64 数据类型转换 hex -> int
func Hex2Int64(s string) int64 {
	n, err := strconv.ParseUint(s, 16, 8)
	if err != nil {
		Error(err)
		return 0
	}

	n2 := uint8(n)
	return int64(*(*int8)(unsafe.Pointer(&n2)))
}

// Str2Int64 数据类型转换 string -> int64
func Str2Int64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// Str2Int 数据类型转换 string -> int
func Str2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// Str2Int32 数据类型转换 string -> int32
func Str2Int32(str string) int32 {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return int32(i)
}

// Str2Float64 数据类型转换 string -> float64
func Str2Float64(str string) float64 {
	i, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}
	return i
}

// Str2Float32 数据类型转换 string -> float32
func Str2Float32(str string) float32 {
	i, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return 0
	}
	return float32(i)
}

// Uint82Str 数据类型转换 []uint8 -> string
func Uint82Str(bs []uint8) string {
	ba := make([]byte, 0)
	for _, b := range bs {
		ba = append(ba, b)
	}
	return string(ba)
}

// Str2Byte 数据类型转换 string -> []byte
func Str2Byte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// Byte2Str 数据类型转换 []byte -> string
func Byte2Str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Bool2Byte 数据类型转换 bool -> []byte
func Bool2Byte(b bool) []byte {
	if b == true {
		return []byte{1}
	}
	return []byte{0}
}

// Byte2Bool 数据类型转换 []byte -> bool
func Byte2Bool(b []byte) bool {
	if len(b) == 0 || bytes.Compare(b, make([]byte, len(b))) == 0 {
		return false
	}
	return true
}

// Int2Byte 数据类型转换 int -> []byte
func Int2Byte(i int) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}

// Byte2Int 数据类型转换 []byte -> int
func Byte2Int(b []byte) int {
	return int(binary.LittleEndian.Uint32(b))
}

// Int642Byte 数据类型转换 int64 -> []byte
func Int642Byte(i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return b
}

// Byte2Int64 数据类型转换 []byte -> int64
func Byte2Int64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}

// Float322Byte 数据类型转换 float32 -> []byte
func Float322Byte(f float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, Float322Uint32(f))
	return b
}

// Float322Uint32 数据类型转换 float32 -> uint32
func Float322Uint32(f float32) uint32 {
	return math.Float32bits(f)
}

// Byte2Float32 数据类型转换 []byte -> float32
func Byte2Float32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(b))
}

// Float642Byte 数据类型转换 float64 -> []byte
func Float642Byte(f float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, Float642Uint64(f))
	return b
}

// Float642Uint64 数据类型转换 float64 -> uint64
func Float642Uint64(f float64) uint64 {
	return math.Float64bits(f)
}

// Byte2Float64 数据类型转换 []byte -> float64
func Byte2Float64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(b))
}

// Struct2Map 数据类型转换 Struct ->  map
// 参数说明:
//   - hasValue=true表示字段值不管是否存在都转换成map
//   - hasValue=false表示字段为空或者不为0则转换成map
func Struct2Map(obj any, hasValue bool) (map[string]any, error) {
	mp := make(map[string]any)
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

// Byte2Bit 数据类型转换 []byte -> []uint8 (bit)
func Byte2Bit(b []byte) []uint8 {
	bits := make([]uint8, 0)
	for _, v := range b {
		bits = bits2Uint(bits, uint(v), 8)
	}
	return bits
}

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

// Bit2Byte 数据类型转换 []uint8 -> []byte
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

func bitsToUint(bits []uint8) uint {
	v := uint(0)
	for _, i := range bits {
		v = v<<1 | uint(i)
	}
	return v
}

// ByteToBinaryString 数据类型转换 字节 -> 二进制字符串
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

// IP2Binary 数据类型转换 IP str -> binary int64
func IP2Binary(ip string) string {
	rse := IP2Int64(ip)
	return strconv.FormatInt(rse, 2)
}

// UInt32ToIP  数据类型转换 uint32 -> net.IP
func UInt32ToIP(ip uint32) net.IP {
	var b [4]byte
	b[0] = byte(ip & 0xFF)
	b[1] = byte((ip >> 8) & 0xFF)
	b[2] = byte((ip >> 16) & 0xFF)
	b[3] = byte((ip >> 24) & 0xFF)
	return net.IPv4(b[3], b[2], b[1], b[0])
}

// IP2Int64 数据类型转换 IP str -> int64
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

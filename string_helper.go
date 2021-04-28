package gathertool

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

//StringValue 任何类型返回值字符串形式
func StringValue(i interface{}) string {
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
func JSON2Map(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
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
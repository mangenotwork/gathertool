/*
*	Description : HTML内容提取； json内容提取  TODO 测试
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/net/html"
)

// GetPointHTML  获取指定位置的HTML， 用标签， 标签属性， 属性值来定位
func GetPointHTML(htmlStr, label, attr, val string) ([]string, error) {
	rse := make([]string, 0)
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return rse, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			//log.Println("attr = ", n.Attr, n.Namespace, n.Data)
			if n.Data == label {
				if attr == "" && val == "" {
					rse = add(rse, n)
				} else {
					for _, a := range n.Attr {
						if a.Key == attr && a.Val == val {
							rse = add(rse, n)
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return rse, nil
}

func add(rse []string, n *html.Node) []string {
	var buf bytes.Buffer
	err := html.Render(&buf, n)
	if err == nil {
		rse = append(rse, buf.String())
	}
	return rse
}

// GetPointIDHTML 获取指定标签id属性的html
func GetPointIDHTML(htmlStr, label, val string) ([]string, error) {
	return GetPointHTML(htmlStr, label, "id", val)
}

// GetPointClassHTML 获取指定标签class属性的html
func GetPointClassHTML(htmlStr, label, val string) ([]string, error) {
	return GetPointHTML(htmlStr, label, "class", val)
}

// ====================================  json find

// JsonFind 按路径寻找指定json值
// 用法参考  ./_examples/json/main.go
// @find : 寻找路径，与目录的url类似， 下面是一个例子：
// json:  {a:[{b:1},{b:2}]}
// find=/a/[0]  =>   {b:1}
// find=a/[0]/b  =>   1
func JsonFind(jsonStr, find string) (any, error) {
	if !IsJson(jsonStr) {
		return nil, fmt.Errorf("不是标准的Json格式")
	}
	jxList := strings.Split(find, "/")
	jxLen := len(jxList)
	var (
		data  = Any2Map(jsonStr)
		value any
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
	return Any2Json(value)
}

// JsonFind2Map 寻找json,输出 map[string]any
func JsonFind2Map(jsonStr, find string) (map[string]any, error) {
	value, err := JsonFind(jsonStr, find)
	if err != nil {
		return nil, err
	}
	return Any2Map(value), nil
}

// JsonFind2Arr 寻找json,输出 []any
func JsonFind2Arr(jsonStr, find string) ([]any, error) {
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
	var tempMap map[string]any
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		return false
	}
	return true
}

// IsHaveKey map[string]any 是否存在 输入的key
func IsHaveKey[T comparable](data map[T]any, key T) bool {
	_, ok := data[key]
	return ok
}

// Any2Map any -> map[string]any
func interface2Map(data any) (map[string]any, error) {
	if v, ok := data.(map[string]any); ok {
		return v, nil
	}
	if reflect.ValueOf(data).Kind() == reflect.String {
		return Json2Map(data.(string))
	}
	return nil, fmt.Errorf("not map type")
}

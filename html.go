/*
	Description : HTML内容提取
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"bytes"
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

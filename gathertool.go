/*

▀██    ██▀                      ▄▄█▀▀▀▄█
 ███  ███   ▄▄▄▄   ▄▄ ▄▄▄      ▄█▀     ▀    ▄▄▄▄
 █▀█▄▄▀██  ▀▀ ▄██   ██  ██     ██    ▄▄▄▄ ▄█▄▄▄██
 █ ▀█▀ ██  ▄█▀ ██   ██  ██     ▀█▄    ██  ██
▄█▄ █ ▄██▄ ▀█▄▄▀█▀ ▄██▄ ██▄     ▀▀█▄▄▄▀█   ▀█▄▄▄▀

	Description : gathertool 轻量级爬虫，接口测试，压力测试框架
	Author : ManGe
	Version : v0.2
	Date : 20211227

*/

package gathertool

import (
	"bytes"
	"net/http"
	"net/url"
	"strings"
)

// Get 请求
func Get(url string, vs ...interface{}) (*Context, error) {
	ctx := NewGet(url, vs...)
	ctx.Do()
	return ctx,nil
}

// NewGet
func NewGet(url string, vs ...interface{}) *Context {
	request, err := http.NewRequest(GET, urlStr(url), nil)
	if err != nil{
		panic(err)
	}
	return	Req(request, vs...)
}

// POST 请求
func Post(url string, data []byte, contentType string, vs ...interface{}) (*Context, error) {
	cxt := NewPost(url, data, contentType, vs...)
	cxt.Do()
	return cxt, nil
}

// NewPost
func NewPost(url string, data []byte, contentType string, vs ...interface{}) *Context {
	request, err := http.NewRequest(POST, urlStr(url), bytes.NewBuffer(data))
	if err != nil{
		panic(err)
	}
	if contentType == "" {
		request.Header.Set("Content-Type", "application/json;")
	} else {
		request.Header.Set("Content-Type", contentType)
	}
	return	Req(request, vs...)
}

// POST json 请求
func PostJson(url string, jsonStr string, vs ...interface{}) (*Context, error) {
	cxt := NewPost(url, []byte(jsonStr), "application/json; charset=UTF-8", vs...)
	cxt.Do()
	return cxt, nil
}

// FormData
type FormData map[string]string

// POST Form
func PostForm(url string, data map[string]string, vs ...interface{}) (*Context, error){
	cxt := NewPostForm(url, data, vs...)
	cxt.Do()
	return cxt, nil
}

// POST NewForm
func NewPostForm(u string, data map[string]string, vs ...interface{}) *Context{
	postData := url.Values{}
	for k,v := range data {
		postData.Add(k,v)
	}
	request, err := http.NewRequest(POST, urlStr(u), strings.NewReader(postData.Encode()))
	if err != nil{
		return nil
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	return Req(request, vs...)
}

// Put
func Put(url string, data []byte, contentType string, vs ...interface{}) (*Context, error){
	cxt :=	NewPut(url, data, contentType, vs...)
	cxt.Do()
	return cxt, nil
}

// NewPut
func NewPut(url string, data []byte, contentType string, vs ...interface{}) *Context{
	request, err := http.NewRequest(PUT, urlStr(url), bytes.NewBuffer(data))
	if err != nil{
		panic(err)
	}
	request.Header.Set("Content-Type", contentType)
	return Req(request, vs...)
}

// Delete
func Delete(url string, vs ...interface{}) (*Context, error) {
	cxt :=	NewDelete(url, vs...)
	cxt.Do()
	return cxt, nil
}

// NewDelete
func NewDelete(url string, vs ...interface{}) *Context {
	request, err := http.NewRequest(DELETE, urlStr(url), nil)
	if err != nil{
		panic(err)
	}
	return	Req(request, vs...)
}

// Options
func Options(url string, vs ...interface{}) (*Context, error) {
	cxt := NewOptions(url, vs...)
	cxt.Do()
	return cxt, nil
}

// NewOptions
func NewOptions(url string, vs ...interface{}) *Context {
	request, err := http.NewRequest(OPTIONS, urlStr(url), nil)
	if err != nil{
		panic(err)
	}
	return	Req(request, vs...)
}


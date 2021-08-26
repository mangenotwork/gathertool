/*
	Description : gathertool 网络请求，爬虫，测试 实用库
	Author : ManGe
	Version : v0.1
	Date : 2021-0828
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
	if !isUrl(url) {
		return nil,UrlBad
	}
	request, err := http.NewRequest(GET, url, nil)
	if err != nil{
		return nil, err
	}
	ctx := Req(request, vs...)
	ctx.Do()
	return ctx,nil
}

func NewGet(url string, vs ...interface{}) *Context {
	if !isUrl(url) {
		loger(UrlBad)
		return nil
	}
	request, err := http.NewRequest(GET, url, nil)
	if err != nil{
		loger("NewGet err->", err)
		return nil
	}
	return	Req(request, vs...)
}


// POST 请求
func Post(url string, data []byte, contentType string, vs ...interface{}) (*Context, error) {
	if !isUrl(url) {
		return nil, UrlBad
	}

	request, err := http.NewRequest(POST, url, bytes.NewBuffer(data))
	if err != nil{
		return nil, err
	}

	if contentType == "" {
		request.Header.Set("Content-Type", "application/json;")
	} else {
		request.Header.Set("Content-Type", contentType)
	}

	cxt := Req(request, vs...)
	cxt.Do()
	return cxt, nil
}

func NewPost(url string, data []byte, contentType string, vs ...interface{}) *Context {
	if !isUrl(url) {
		loger(UrlBad)
		return nil
	}

	request, err := http.NewRequest(POST, url, bytes.NewBuffer(data))
	if err != nil{
		loger("NewPost err->", err)
		return nil
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
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest(POST, url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil{
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	cxt := Req(request, vs...)
	cxt.Do()
	return cxt, nil
}


// POST Form
func PostForm(url string, data url.Values, vs ...interface{}) (*Context, error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest(POST, url, strings.NewReader(data.Encode()))
	if err != nil{
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	cxt := Req(request, vs...)
	cxt.Do()
	return cxt, nil
}


// Put
func Put(url string, data []byte, contentType string, vs ...interface{}) (*Context, error){
	if !isUrl(url) {
		loger(UrlBad)
		return nil, UrlBad
	}
	request, err := http.NewRequest(PUT, url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		loger("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	cxt :=	Req(request, vs...)
	cxt.Do()
	return cxt, nil
}

// NewPut
func NewPut(url string, data []byte, contentType string, vs ...interface{}) *Context{
	if !isUrl(url) {
		loger(UrlBad)
		return nil
	}
	request, err := http.NewRequest(PUT, url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		loger("err->", err)
		return nil
	}
	request.Header.Set("Content-Type", contentType)
	return Req(request, vs...)
}

// NewDelete
func NewDelete(url string, vs ...interface{}) *Context {
	if !isUrl(url) {
		loger(UrlBad)
		return nil
	}
	request, err := http.NewRequest(DELETE, url, nil)
	if err != nil{
		loger("err->", err)
		return nil
	}
	return	Req(request, vs...)
}


// Delete
func Delete(url string, vs ...interface{}) (*Context, error) {
	if !isUrl(url) {
		loger(UrlBad)
		return nil, UrlBad
	}
	request, err := http.NewRequest(DELETE, url, nil)
	if err != nil{
		loger("err->", err)
		return nil, err
	}
	cxt :=	Req(request, vs...)
	cxt.Do()
	return cxt, nil
}

// NewOptions
func NewOptions(url string, vs ...interface{}) *Context {
	if !isUrl(url) {
		loger(UrlBad)
		return nil
	}
	request, err := http.NewRequest(OPTIONS, url, nil)
	if err != nil{
		loger("err->", err)
		return nil
	}
	return	Req(request, vs...)
}

// Options
func Options(url string, vs ...interface{}) (*Context, error) {
	if !isUrl(url) {
		loger(UrlBad)
		return nil, UrlBad
	}
	request, err := http.NewRequest(OPTIONS, url, nil)
	if err != nil{
		loger("err->", err)
		return nil, err
	}
	cxt :=	Req(request, vs...)
	cxt.Do()
	return cxt, nil
}

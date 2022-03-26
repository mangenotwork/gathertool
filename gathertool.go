/*

▀██    ██▀                      ▄▄█▀▀▀▄█
 ███  ███   ▄▄▄▄   ▄▄ ▄▄▄      ▄█▀     ▀    ▄▄▄▄
 █▀█▄▄▀██  ▀▀ ▄██   ██  ██     ██    ▄▄▄▄ ▄█▄▄▄██
 █ ▀█▀ ██  ▄█▀ ██   ██  ██     ▀█▄    ██  ██
▄█▄ █ ▄██▄ ▀█▄▄▀█▀ ▄██▄ ██▄     ▀▀█▄▄▄▀█   ▀█▄▄▄▀

	Description : gathertool 轻量级爬虫，接口测试，压力测试框架, 提高开发对应场景的golang程序的效率。
	Author : ManGe

*/

package gathertool

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Get(url string, vs ...interface{}) (*Context, error) {
	ctx := NewGet(url, vs...)
	ctx.Do()
	return ctx,nil
}

func NewGet(url string, vs ...interface{}) *Context {
	request, err := http.NewRequest(GET, urlStr(url), nil)
	if err != nil{
		panic(err)
	}
	return	Req(request, vs...)
}

func Post(url string, data []byte, contentType string, vs ...interface{}) (*Context, error) {
	cxt := NewPost(url, data, contentType, vs...)
	cxt.Do()
	return cxt, nil
}

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

func PostJson(url string, jsonStr string, vs ...interface{}) (*Context, error) {
	cxt := NewPost(url, []byte(jsonStr), "application/json; charset=UTF-8", vs...)
	cxt.Do()
	return cxt, nil
}

// FormData
type FormData map[string]string

func PostForm(url string, data map[string]string, vs ...interface{}) (*Context, error){
	cxt := NewPostForm(url, data, vs...)
	cxt.Do()
	return cxt, nil
}

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

func PostFile(url, paramName, filePath string, vs ...interface{}) *Context {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	uploadBody := &bytes.Buffer{}
	writer := multipart.NewWriter(uploadBody)

	fWriter, err := writer.CreateFormFile("uploadFile", filePath)
	if err != nil {
		Info("copy file writer %v", err)
	}

	_, err = io.Copy(fWriter, f)
	if err != nil {
		panic(err)
	}

	fieldMap := map[string]string{
		paramName: filePath,
	}

	for k, v := range fieldMap {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close()
	if err != nil {
		panic(err)
	}

	request, err := http.NewRequest(POST, urlStr(url), uploadBody)
	if err != nil{
		return nil
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	return Req(request, vs...)
}

func Put(url string, data []byte, contentType string, vs ...interface{}) (*Context, error){
	cxt :=	NewPut(url, data, contentType, vs...)
	cxt.Do()
	return cxt, nil
}

func NewPut(url string, data []byte, contentType string, vs ...interface{}) *Context{
	request, err := http.NewRequest(PUT, urlStr(url), bytes.NewBuffer(data))
	if err != nil{
		panic(err)
	}
	request.Header.Set("Content-Type", contentType)
	return Req(request, vs...)
}

func Delete(url string, vs ...interface{}) (*Context, error) {
	cxt :=	NewDelete(url, vs...)
	cxt.Do()
	return cxt, nil
}

func NewDelete(url string, vs ...interface{}) *Context {
	request, err := http.NewRequest(DELETE, urlStr(url), nil)
	if err != nil{
		panic(err)
	}
	return	Req(request, vs...)
}

func Options(url string, vs ...interface{}) (*Context, error) {
	cxt := NewOptions(url, vs...)
	cxt.Do()
	return cxt, nil
}

func NewOptions(url string, vs ...interface{}) *Context {
	request, err := http.NewRequest(OPTIONS, urlStr(url), nil)
	if err != nil{
		panic(err)
	}
	return	Req(request, vs...)
}

func Upload(url, savePath string, vs ...interface{}) (*Context, error) {
	c := NewGet(urlStr(url), vs)
	c.Upload(savePath)
	if c.Err != nil {
		return c, c.Err
	}
	return c, nil
}


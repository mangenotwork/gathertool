/*

▀██    ██▀                      ▄▄█▀▀▀▄█
 ███  ███   ▄▄▄▄   ▄▄ ▄▄▄      ▄█▀     ▀    ▄▄▄▄
 █▀█▄▄▀██  ▀▀ ▄██   ██  ██     ██    ▄▄▄▄ ▄█▄▄▄██
 █ ▀█▀ ██  ▄█▀ ██   ██  ██     ▀█▄    ██  ██
▄█▄ █ ▄██▄ ▀█▄▄▀█▀ ▄██▄ ██▄     ▀▀█▄▄▄▀█   ▀█▄▄▄▀

	Description : gathertool是golang脚本化开发库，目的是提高对应场景程序开发的效率；轻量级爬虫库，接口测试&压力测试库，DB操作库等。
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Get 执行get请求，请求内容在上下文 *Context 里
func Get(url string, vs ...interface{}) (*Context, error) {
	ctx := NewGet(url, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewGet 定义一个get请求对象，请求需执行 XX.Do()
func NewGet(url string, vs ...interface{}) *Context {
	request, err := http.NewRequest(GET, urlStr(url), nil)
	if err != nil {
		Error(err)
		return &Context{
			Err:       fmt.Errorf(url, " 请求错误， err: ", err.Error()),
			StateCode: 404,
		}
	}
	return Req(request, vs...)
}

// Post 执行post请求，请求内容在上下文 *Context 里
func Post(url string, data []byte, contentType string, vs ...interface{}) (*Context, error) {
	ctx := NewPost(url, data, contentType, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewPost 定义一个post请求对象，请求需执行 XX.Do()
func NewPost(url string, data []byte, contentType string, vs ...interface{}) *Context {
	request, err := http.NewRequest(POST, urlStr(url), bytes.NewBuffer(data))
	if err != nil {
		Error(err)
		return &Context{
			Err:       fmt.Errorf(url, " 请求错误， err: ", err.Error()),
			StateCode: 404,
		}
	}
	if contentType == "" {
		request.Header.Set("Content-Type", "application/json;")
	} else {
		request.Header.Set("Content-Type", contentType)
	}
	return Req(request, vs...)
}

// PostJson 执行post请求，请求参数为json，请求内容在上下文 *Context 里
func PostJson(url string, jsonStr string, vs ...interface{}) (*Context, error) {
	ctx := NewPost(url, []byte(jsonStr), "application/json; charset=UTF-8", vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// FormData form data
type FormData map[string]string

// PostForm 执行post请求，请求参数为fromdata，请求内容在上下文 *Context 里
func PostForm(url string, data map[string]string, vs ...interface{}) (*Context, error) {
	ctx := NewPostForm(url, data, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewPostForm 定义一个post请求对象，请求参数为fromdata，请求需执行 XX.Do()
func NewPostForm(u string, data map[string]string, vs ...interface{}) *Context {
	postData := url.Values{}
	for k, v := range data {
		postData.Add(k, v)
	}
	request, err := http.NewRequest(POST, urlStr(u), strings.NewReader(postData.Encode()))
	if err != nil {
		return nil
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	return Req(request, vs...)
}

// PostFile 定义一个post上传文件的请求对象，请求需执行 XX.Do()
func PostFile(url, paramName, filePath string, vs ...interface{}) *Context {
	f, err := os.Open(filePath)
	if err != nil {
		Error(err)
		return &Context{
			Err:       fmt.Errorf(url, " 请求错误， err: ", err.Error()),
			StateCode: 404,
		}
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
		Error(err)
	}

	fieldMap := map[string]string{
		paramName: filePath,
	}

	for k, v := range fieldMap {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close()
	if err != nil {
		Error(err)
	}

	request, err := http.NewRequest(POST, urlStr(url), uploadBody)
	if err != nil {
		return nil
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	return Req(request, vs...)
}

// Put 执行put请求，请求内容在上下文 *Context 里
func Put(url string, data []byte, contentType string, vs ...interface{}) (*Context, error) {
	ctx := NewPut(url, data, contentType, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewPut 定义一个put请求对象，请求需执行 XX.Do()
func NewPut(url string, data []byte, contentType string, vs ...interface{}) *Context {
	request, err := http.NewRequest(PUT, urlStr(url), bytes.NewBuffer(data))
	if err != nil {
		Error(err)
		return &Context{
			Err:       fmt.Errorf(url, " 请求错误， err: ", err.Error()),
			StateCode: 404,
		}
	}
	request.Header.Set("Content-Type", contentType)
	return Req(request, vs...)
}

// Delete 执行delete请求，请求内容在上下文 *Context 里
func Delete(url string, vs ...interface{}) (*Context, error) {
	ctx := NewDelete(url, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewDelete 定义一个delete请求对象，请求需执行 XX.Do()
func NewDelete(url string, vs ...interface{}) *Context {
	request, err := http.NewRequest(DELETE, urlStr(url), nil)
	if err != nil {
		Error(err)
		return &Context{
			Err:       fmt.Errorf(url, " 请求错误， err: ", err.Error()),
			StateCode: 404,
		}
	}
	return Req(request, vs...)
}

// Options 执行options请求，请求内容在上下文 *Context 里
func Options(url string, vs ...interface{}) (*Context, error) {
	ctx := NewOptions(url, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewOptions 定义一个options请求对象，请求需执行 XX.Do()
func NewOptions(url string, vs ...interface{}) *Context {
	request, err := http.NewRequest(OPTIONS, urlStr(url), nil)
	if err != nil {
		Error(err)
		return &Context{
			Err:       fmt.Errorf(url, " 请求错误， err: ", err.Error()),
			StateCode: 404,
		}
	}
	return Req(request, vs...)
}

// Upload 下载文件
func Upload(url, savePath string, vs ...interface{}) (*Context, error) {
	c := NewGet(urlStr(url), vs)
	c.Upload(savePath)
	if c.Err != nil {
		return c, c.Err
	}
	return c, nil
}

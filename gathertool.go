/*
*	Description : 根据method提供对应的http请求方法，该库的请求方法会生成请求上下文，可以使用上下文的扩展功能来满足你当前的场景。
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

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

// Get 执行GET请求，请求结果在 Context 对象里，请求错误会直接抛出
//
// caseUrl : 请求链接地址
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	ctx,err := Get("https://xxx.xxx")
func Get(caseUrl string, vs ...any) (*Context, error) {
	ctx := NewGet(caseUrl, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewGet 定义一个GET请求的 Context ，但是不执行，执行需要调用 Context.Do
//
// caseUrl : 请求链接地址
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	NewGet("https://xxx.xxx").Do()
func NewGet(caseUrl string, vs ...any) *Context {
	request, err := http.NewRequest(GET, urlStr(caseUrl), nil)
	if err != nil {
		Error(err)
		return &Context{
			Err:       reqErr(caseUrl, err),
			StateCode: RequestStateError,
		}
	}
	return Req(request, vs...)
}

// Post 执行POST请求，请求结果在 Context 对象里，请求错误会直接抛出
//
// caseUrl : 请求链接地址
//
// data : 请求的body
//
// contentType : 请求的content-type,如application/json; charset=UTF-8
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	ctx,err := Post("https://xxx.xxx", []byte("xx"), "application/json")
func Post(caseUrl string, data []byte, contentType string, vs ...any) (*Context, error) {
	ctx := NewPost(caseUrl, data, contentType, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewPost 定义一个POST请求的 Context ，但是不执行，执行需要调用 Context.Do
//
// caseUrl : 请求链接地址
//
// data : 请求的body
//
// contentType : 请求的content-type,如application/json; charset=UTF-8
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	ctx := NewPost("https://xxx.xxx", []byte("xx"), "application/json")
func NewPost(caseUrl string, data []byte, contentType string, vs ...any) *Context {
	request, err := http.NewRequest(POST, urlStr(caseUrl), bytes.NewBuffer(data))
	if err != nil {
		Error(err)
		return &Context{
			Err:       reqErr(caseUrl, err),
			StateCode: RequestStateError,
		}
	}
	if contentType == "" {
		request.Header.Set("Content-Type", "application/json;")
	} else {
		request.Header.Set("Content-Type", contentType)
	}
	return Req(request, vs...)
}

// PostJson 执行Content-Type=json的POST请求，请求结果在 Context 对象里，请求错误会直接抛出
//
// caseUrl : 请求链接地址
//
// jsonStr : 请求body,应为json字符串
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	ctx,err := PostJson("https://xxx.xxx", `{"xx":"xx"}`)
func PostJson(caseUrl string, jsonStr string, vs ...any) (*Context, error) {
	ctx := NewPost(caseUrl, []byte(jsonStr), "application/json; charset=UTF-8", vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// PostForm 执行Content-Type=from-data的POST请求，请求结果在 Context 对象里，请求错误会直接抛出
//
// caseUrl : 请求链接地址
//
// data : 请求body,应为map[string]string
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	ctx,err := PostForm("https://xxx.xxx", map[string]string{"xx":"xx"})
func PostForm(caseUrl string, data map[string]string, vs ...any) (*Context, error) {
	postData := url.Values{}
	for k, v := range data {
		postData.Add(k, v)
	}
	request, err := http.NewRequest(POST, urlStr(caseUrl), strings.NewReader(postData.Encode()))
	if err != nil {
		Error(err)
		ctx := &Context{
			Err:       reqErr(caseUrl, err),
			StateCode: RequestStateError,
		}
		return ctx, ctx.Err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	ctx := Req(request, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// PostFile 执行一个POST上传文件的请求，请求结果在 Context 对象里，请求错误会直接抛出
//
// caseUrl : 请求链接地址
//
// paramName : from-data的参数名
//
// filePath : 文件路径
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	ctx,err := PostFile("https://xxx.xxx","file","./test.txt")
func PostFile(caseUrl, paramName, filePath string, vs ...any) (*Context, error) {
	f, err := os.Open(filePath)
	if err != nil {
		Error(err)
		ctx := &Context{
			Err:       reqErr(caseUrl, err),
			StateCode: RequestStateError,
		}
		return ctx, ctx.Err
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

	request, err := http.NewRequest(POST, urlStr(caseUrl), uploadBody)
	if err != nil {
		Error(err)
		ctx := &Context{
			Err:       fmt.Errorf(caseUrl, " 请求错误， err: ", err.Error()),
			StateCode: RequestStateError,
		}
		return ctx, ctx.Err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	ctx := Req(request, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// Put 执行一个PUT请求，请求结果在 Context 对象里，请求错误会直接抛出
//
// caseUrl : 请求链接地址
//
// data : 请求的body
//
// contentType : 请求的 content-type
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	ctx,err := Put("https://xxx.xxx",[]byte("xxx"),"")
func Put(caseUrl string, data []byte, contentType string, vs ...any) (*Context, error) {
	ctx := NewPut(caseUrl, data, contentType, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewPut 定义一个PUT请求的 Context ，但是不执行，执行需要调用 Context.Do
//
// caseUrl : 请求链接地址
//
// data : 请求的body
//
// contentType : 请求的content-type,如application/json; charset=UTF-8
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	ctx := NewPut("https://xxx.xxx", []byte("xx"), "application/json")
func NewPut(caseUrl string, data []byte, contentType string, vs ...any) *Context {
	request, err := http.NewRequest(PUT, urlStr(caseUrl), bytes.NewBuffer(data))
	if err != nil {
		Error(err)
		return &Context{
			Err:       reqErr(caseUrl, err),
			StateCode: RequestStateError,
		}
	}
	request.Header.Set("Content-Type", contentType)
	return Req(request, vs...)
}

// Delete 执行一个DELETE请求，请求结果在 Context 对象里，请求错误会直接抛出
func Delete(caseUrl string, vs ...any) (*Context, error) {
	ctx := NewDelete(caseUrl, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewDelete 定义一个DELETE请求的 Context ，但是不执行，执行需要调用 Context.Do
func NewDelete(caseUrl string, vs ...any) *Context {
	request, err := http.NewRequest(DELETE, urlStr(caseUrl), nil)
	if err != nil {
		Error(err)
		return &Context{
			Err:       reqErr(caseUrl, err),
			StateCode: RequestStateError,
		}
	}
	return Req(request, vs...)
}

// Options 执行一个OPTIONS请求，请求结果在 Context 对象里，请求错误会直接抛出
func Options(caseUrl string, vs ...any) (*Context, error) {
	ctx := NewOptions(caseUrl, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// NewOptions 定义一个OPTIONS请求的 Context ，但是不执行，执行需要调用 Context.Do
func NewOptions(caseUrl string, vs ...any) *Context {
	request, err := http.NewRequest(OPTIONS, urlStr(caseUrl), nil)
	if err != nil {
		Error(err)
		return &Context{
			Err:       reqErr(caseUrl, err),
			StateCode: RequestStateError,
		}
	}
	return Req(request, vs...)
}

// Head 执行一个HEAD请求，请求结果在 Context 对象里，请求错误会直接抛出
func Head(caseUrl string, vs ...any) (*Context, error) {
	request, err := http.NewRequest(HEAD, urlStr(caseUrl), nil)
	if err != nil {
		Error(err)
		return &Context{
			Err:       reqErr(caseUrl, err),
			StateCode: RequestStateError,
		}, err
	}
	ctx := Req(request, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// Patch 执行一个PATCH请求，请求结果在 Context 对象里，请求错误会直接抛出
func Patch(caseUrl string, vs ...any) (*Context, error) {
	request, err := http.NewRequest(PATCH, urlStr(caseUrl), nil)
	if err != nil {
		Error(err)
		return &Context{
			Err:       reqErr(caseUrl, err),
			StateCode: RequestStateError,
		}, err
	}
	ctx := Req(request, vs...)
	ctx.Do()
	return ctx, ctx.Err
}

// Upload 执行GET下载请求，将下载的内容保存到指定路径终端会默认显示下载进度，并返回 Context 对象，请求错误会直接抛出
//
// caseUrl : 请求链接地址
//
// savePath : 下载内容保存路径
//
// vs : 更多可选参数见 Req
//
// Example:
//
//	ctx := Upload("https://xxx.xxx", "./test.txt")
func Upload(caseUrl, savePath string, vs ...any) (*Context, error) {
	ctx := NewGet(urlStr(caseUrl), vs)
	ctx.Upload(savePath)
	if ctx.Err != nil {
		return ctx, ctx.Err
	}
	return ctx, nil
}

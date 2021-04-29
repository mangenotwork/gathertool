/*
	Description : 对外提供的方法
	Author : ManGe
	Version : v0.1
	Date : 2021-04-23
*/

package gathertool

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"time"
)

var (
	UrlBad error = errors.New("url is bad.") // 错误的url
)

type ReqTimeOut int
type ReqTimeOutMs int

// Get 请求, 当请求失败或状态码是失败的则会先执行 ff 再回调
func Get(url string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	return	Req(request, vs...)
}

// POST 请求
func Post(url string, data []byte, contentType string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return	Req(request, vs...)
}

// POST json 请求
func PostJson(url string, jsonStr string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return	Req(request, vs...)
}

// Put
func Put(url string, data []byte, contentType string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return	Req(request, vs...)
}

// Delete
func Delete(url string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	return	Req(request, vs...)
}

// Options
func Options(url string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	return	Req(request, vs...)
}


// Request 请求
func Request(url, method string, data []byte, contentType string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return	Req(request, vs...)
}


// isUrl 验证是否是有效的 url
func isUrl(url string) bool {
	if url == ""{
		return false
	}
	return true
}


// Req 初始化请求
// @url  请求链接
// @maxTimes  重试次数
// @sf  请求成功后做的事情, 200等
// @ff  请求失败后做的事情, 403等，502等
// @vs  可变参数
// @vs UserAgentType  设置指定类型 user agent 如 AndroidAgent
func Req(request *http.Request, vs ...interface{}) (*Context,error){
	var (
		client *http.Client
		maxTimes RetryTimes = 10
		task *Task
		start StartFunc
		succeed SucceedFunc
		failed FailedFunc
		retry RetryFunc
		end EndFunc
		reqTimeOut ReqTimeOut
		reqTimeOutMs ReqTimeOutMs
	)

	//添加默认的Header
	request.Header.Set("Connection","close")
	request.Header.Set("User-Agent", GetAgent(PCAgent))

	//解析可变参
	for _, v := range vs {
		switch vv := v.(type) {
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					request.Header.Add(key, value)
				}
			}
		case *http.Client:
			client = vv
		case UserAgentType:
			request.Header.Set("User-Agent", GetAgent(vv))
		case *http.Cookie:
			request.AddCookie(vv)
		case RetryTimes:
			maxTimes = vv
		case *Task:
			task = vv
		case StartFunc:
			start = vv
		case SucceedFunc:
			succeed = vv
		case FailedFunc:
			failed = vv
		case RetryFunc:
			retry = vv
		case EndFunc:
			end = vv
		case ReqTimeOut:
			reqTimeOut = vv
		case ReqTimeOutMs:
			reqTimeOutMs = vv
		}
	}

	// 如果使用方未传入Client，  初始化 Client
	if client == nil{
		//log.Println("使用方未传入Client， 默认 client")
		client = &http.Client{
			Timeout: 60*time.Second,
		}
	}

	if reqTimeOut > 0 {
		client.Timeout =  time.Duration(reqTimeOut) * time.Second
	}

	if reqTimeOutMs > 0 {
		client.Timeout =  time.Duration(reqTimeOutMs) * time.Millisecond
	}

	// 创建对象
	return &Context{
		Client: client,
		Req : request,
		times : 0,
		MaxTimes : maxTimes,
		Task: task,
		StartFunc: start,
		SucceedFunc: succeed,
		FailedFunc: failed,
		RetryFunc: retry,
		EndFunc: end,
	},nil
}




/*
	Description : 请求上下文
	Author : ManGe
	Version : v0.6
	Date : 2022-03-16
*/

package gathertool

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// RetryTimes 重试次数
type RetryTimes int

// StartFunc 请求开始前的方法类型
type StartFunc func(c *Context)

// SucceedFunc 请求成功后的方法类型
type SucceedFunc func(c *Context)

// FailedFunc 请求失败的方法类型
type FailedFunc func(c *Context)

// RetryFunc 指定重试状态码重试前的方法类型
type RetryFunc func(c *Context)

// EndFunc 请求流程结束后的方法类型
type EndFunc func(c *Context)

// IsLog 全局是否开启日志
type IsLog bool

// ProxyUrl 全局代理地址
type ProxyUrl string

// Context 请求上下文
type Context struct {
	// Token
	Token string

	// http client
	Client *http.Client

	// http Request
	Req *http.Request

	// http Response
	Resp *http.Response

	// Error
	Err error

	// Ctx context.Context
	Ctx context.Context

	// 执行的次数 初始化都是0
	times RetryTimes

	// 最大允许重试次数
	MaxTimes RetryTimes

	// 请求成功了需要处理的事件
	SucceedFunc SucceedFunc

	// 请求失败了需要做的事
	FailedFunc FailedFunc

	// 请求状态码设置了重试，在重试前的事件
	RetryFunc RetryFunc

	// 请求开始前的方法
	StartFunc StartFunc

	// 请求完成后的方法
	EndFunc EndFunc

	// 本次请求的任务
	// 用于有步骤的请求和并发执行请求
	Task *Task

	// 请求返回的结果
	RespBody []byte

	// job编号
	// 在执行多并发执行抓取任务，每个并发都有一个编号
	// 这个编号是递增分配的
	JobNumber int

	// 请求的响应时间 单位ms
	Ms time.Duration

	// 是否显示日志, 默认是显示的
	IsLog IsLog

	// 指定失败执行重试事件
	err2retry bool

	// 是否关闭重试
	isRetry bool

	// 休眠时间
	sleep Sleep

	// 输出字符串
	Text string

	// 输出Json
	Json string

	// 输出xml
	Xml string

	// 输出HTML
	Html string

	// 请求上下文参数
	Param map[string]interface{}
}

// SetSucceedFunc 设置成功后的方法
func (c *Context) SetSucceedFunc(successFunc func(c *Context)) *Context {
	c.SucceedFunc = successFunc
	return c
}

// SetFailedFunc 设置错误后的方法
func (c *Context) SetFailedFunc(failedFunc func(c *Context)) *Context {
	c.FailedFunc = failedFunc
	return c
}

// SetRetryFunc 设置重试，在重试前的方法
func (c *Context) SetRetryFunc(retryFunc func(c *Context)) *Context {
	c.RetryFunc = retryFunc
	return c
}

// SetRetryTimes 设置重试次数
func (c *Context) SetRetryTimes(times int) *Context {
	c.MaxTimes = RetryTimes(times)
	return c
}

// Do 执行请求
func (c *Context) Do() func() {
	var bodyBytes []byte

	if c == nil{
		return nil
	}

	//执行 start
	if c.times == 0 && c.StartFunc != nil{
		c.StartFunc(c)
	}

	//执行 end
	if c.times == c.MaxTimes && c.EndFunc != nil {
		c.EndFunc(c)
	}

	//重试验证
	c.times++
	if c.times > c.MaxTimes{
		Error( "【日志】 请求失败操过", c.MaxTimes, " 次了,结束重试操作！")
		// 超过了重试次数，就算失败，则执行失败方法
		if c.FailedFunc != nil{
			c.FailedFunc(c)
		}
		return nil
	}

	// 复用 Req.Body
	if c.Req.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Req.Body)
	}
	c.Req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// 设置的休眠时间
	if c.sleep != 0 {
		time.Sleep(time.Duration(c.sleep))
	}

	// 执行请求
	before := time.Now()
	c.Resp,c.Err = c.Client.Do(c.Req)
	c.Ms = time.Now().Sub(before)

	// 是否超时
	if c.Err != nil && (
		strings.Contains(c.Err.Error(), "(Client.Timeout exceeded while awaiting headers)") ||
		strings.Contains(c.Err.Error(), ("Too Many Requests")) ||
		strings.Contains(c.Err.Error(), ("To Many Requests")) ||
		strings.Contains(c.Err.Error(), ("EOF")) ||
		strings.Contains(c.Err.Error(), ("connection timed out")) ){

		Error("【日志】 请求 超时 = ", c.Err)
		if c.RetryFunc != nil && !c.isRetry {
			InfoTimes(4, "【日志】 执行 retry 事件： 第", c.times, "次， 总： ",  c.MaxTimes)
			c.Req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			c.RetryFunc(c)
			return c.Do()
		}
		return nil
	}

	// 其他错误
	if c.Err != nil {
		Error("【日志】 请求 err = ", c.Err)
		// 指定的失败都执行 retry
		if c.err2retry && c.RetryFunc != nil && !c.isRetry {
			InfoTimes(4, "【日志】 执行 retry 事件： 第", c.times, "次， 总： ",  c.MaxTimes)
			c.Req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			c.RetryFunc(c)
			return c.Do()
		}

		if c.FailedFunc != nil{
			c.FailedFunc(c)
		}
		return nil
	}

	defer func(c *Context){
		if c.Resp != nil {
			_=c.Resp.Body.Close()
		}
	}(c)

	if c.Resp.Header.Get("Content-Encoding") == "gzip" {
		c.Resp.Body, _ = gzip.NewReader(c.Resp.Body)
	}
	InfoTimes( 4,"【日志】 请求状态码：", c.Resp.StatusCode, " | 用时 ： ", c.Ms)

	// 根据状态码配置的事件了类型进行该事件的方法
	v,ok := StatusCodeMap[c.Resp.StatusCode]
	if !ok {
		body, err := ioutil.ReadAll(c.Resp.Body)
		if err != nil{
			Error(err)
			return nil
		}
		c.RespBody = body
	}

	switch v {
	case "success":
		InfoTimes(4, "【日志】 执行 success 事件")
		// 请求后的结果
		body, err := ioutil.ReadAll(c.Resp.Body)
		if err != nil{
			Error(err)
			return nil
		}
		c.RespBody = body
		// 执行成功方法
		if c.SucceedFunc != nil {
			c.SucceedFunc(c)
		}

	case "retry":
		c.Req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		if c.RetryFunc != nil && !c.isRetry {
			InfoTimes(4, "【日志】 执行 retry 事件： 第", c.times, "次， 总： ",  c.MaxTimes)
			c.Req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			c.RetryFunc(c)
			return c.Do()
		}

	case "fail":
		if c.FailedFunc != nil{
			InfoTimes(4, "【日志】 执行 failed 事件")
			c.FailedFunc(c)
		}

	case "start":
		if c.StartFunc != nil {
			InfoTimes(4, "【日志】 执行 请求前的方法")
			c.StartFunc(c)
		}

	case "end":
		if c.EndFunc != nil {
			InfoTimes(4, "【日志】 执行 请求结束后的方法")
			c.EndFunc(c)
		}

	default:
		ok = false
	}

	c.Text = c.RespBodyString()
	c.Json = c.RespBodyString()
	c.Xml = c.RespBodyString()
	c.Html = c.RespBodyHtml()
	return nil
}

// RespBodyString Body -> String
func (c *Context) RespBodyString() string {
	if c.RespBody != nil {
		return string(c.RespBody)
	}
	return ""
}

// RespBodyHtml Body -> html string
func (c *Context) RespBodyHtml() string {
	html := c.RespBodyString()
	return strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&#34;", `"`,
		"&#39;", "'",
	).Replace(html)
}

// RespBodyMap Body -> Map
func (c *Context) RespBodyMap() map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal(c.RespBody, &tempMap)
	if err != nil {
		Error(err)
		return nil
	}
	return tempMap
}

// RespBodyArr Body -> Arr
func (c *Context) RespBodyArr() []interface{} {
	var tempArr []interface{}
	err := json.Unmarshal(c.RespBody, &tempArr)
	if err != nil {
		Error(err)
		return nil
	}
	return tempArr
}

// GetRespHeader
func (c *Context) GetRespHeader() string {
	header, _ := httputil.DumpResponse(c.Resp, false)
	return string(header)
}

// RespContentLength
func (c *Context) RespContentLength() int64 {
	return c.Resp.ContentLength
}

// CheckReqMd5 check req Md5
func (c *Context) CheckReqMd5() string {
	var buffer bytes.Buffer
	urlStr := c.Req.URL.String()
	reqBodyBytes, _ := ioutil.ReadAll(c.Req.Body)
	method := c.Req.Method
	buffer.WriteString(urlStr)
	buffer.Write(reqBodyBytes)
	buffer.WriteString(method)
	h := md5.New()
	h.Write(buffer.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

// CheckMd5 check Md5
func (c *Context) CheckMd5() string {
	var buffer bytes.Buffer
	urlStr := c.Req.URL.String()
	reqBodyBytes, _ := ioutil.ReadAll(c.Req.Body)
	method := c.Req.Method
	buffer.WriteString(urlStr)
	buffer.Write(reqBodyBytes)
	buffer.WriteString(method)
	buffer.WriteString(c.Resp.Status)
	buffer.Write(c.RespBody)
	h := md5.New()
	h.Write(buffer.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

// AddHeader add header
func (c *Context) AddHeader(k,v string) *Context {
	c.Req.Header.Add(k,v)
	return c
}

// AddCookie add Cookie
func (c *Context) AddCookie(k, v string) *Context {
	cookie := &http.Cookie{Name: k, Value: v, HttpOnly: true}
	c.Req.AddCookie(cookie)
	return c
}

// SetProxy set proxy
func (c *Context) SetProxy(proxyUrl string) *Context {
	proxy, _ := url.Parse(proxyUrl)
	c.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	return c
}

// SetProxyFunc set proxy func
func (c *Context) SetProxyFunc(f func() *http.Transport) *Context {
	c.Client.Transport = f()
	return c
}

// Upload 下载
func (c *Context) Upload(filePath string) func(){
	//空验证
	if c == nil{
		return nil
	}

	//重试验证
	c.times++
	if c.times > c.MaxTimes{
		Infof("请求失败操过", c.MaxTimes, "次了")
		return nil
	}

	//执行请求
	c.Resp,c.Err = c.Client.Do(c.Req)

	// 是否超时
	// Go 1.6 以下
	if c.Err != nil && strings.Contains(c.Err.Error(), "(Client.Timeout exceeded while awaiting headers)"){
		if c.RetryFunc != nil {
			c.RetryFunc(c)
			return c.Do()
		}
		return nil
	}

	// 自 Go 1.6开始， 所有的超时导致的网络错误都可以通过net.Error的Timeout()方法检查。
	if err, ok := c.Err.(net.Error); ok && err.Timeout() {
		if c.RetryFunc != nil {
			c.RetryFunc(c)
			return c.Do()
		}
		return nil
	}

	// 其他错误
	if c.Err != nil {
		Error(c.Err)
		if c.FailedFunc != nil{
			c.FailedFunc(c)
		}
		return nil
	}

	defer func(cxt *Context){
		if cxt.Resp != nil {
			_=cxt.Resp.Body.Close()
		}
	}(c)

	f, err := os.Create(filePath)
	if err != nil {
		c.Err = err
		return nil
	}

	defer func(){
		f.Close()
	}()

	contentLength := Str2Float64(c.Resp.Header.Get("Content-Length"))
	var sum int64 = 0
	buf := make([]byte, 1024*100)
	st := time.Now()
	i := 0
	for {
		i++
		n, err := c.Resp.Body.Read(buf)
		sum=sum+int64(n)
		_,_=f.Write(buf[:n])
		if err != nil || n == 0{
			break
		}
		if i%9 == 0{
			log.Println("[下载] ", filePath, " : ", FileSizeFormat(sum),"/", FileSizeFormat(int64(contentLength)),
				" |\t ", math.Floor((float64(sum)/contentLength)*100),"%")
		}
	}
	ct := time.Now().Sub(st)
	log.Println("[下载] ", filePath, " : ", FileSizeFormat(sum),"/", FileSizeFormat(int64(contentLength)),
		" |\t ", math.Floor((float64(sum)/contentLength)*100), "%", "|\t ", ct )
	//loger(" rep header ", c.Resp.ContentLength)
	return nil
}


// CookieNext
func (c *Context) CookieNext() error {
	if c.Resp == nil{
		return fmt.Errorf("response is nil.")
	}
	if c.Req == nil {
		return fmt.Errorf("request is nil.")
	}
	// 上下文cookies
	for _,cookie := range c.Resp.Cookies(){
		c.Req.AddCookie(cookie)
	}
	return nil
}

// CloseLog close log
func (c *Context) CloseLog() {
	c.IsLog = false
}

// OpenErr2Retry 开启请求失败都执行retry
func (c *Context) OpenErr2Retry() {
	c.err2retry = true
}

// CloseRetry 关闭重试
func (c *Context) CloseRetry() {
	c.isRetry = true
}

// GetParam 获取上下文参数
func (c *Context) GetParam(key string) interface{} {
	return c.Param[key]
}

// AddParam 添加上下文参数
func (c *Context) AddParam(key string, val interface{}) {
	c.Param[key] = val
}

// DelParam 删除上下文参数
func (c *Context) DelParam(key string) {
	delete(c.Param, key)
}

// CookiePool  cookie池
type cookiePool struct {
	cookie []*http.Cookie
	mux sync.Mutex
}

var CookiePool *cookiePool
var _cookiePoolOnce sync.Once

func NewCookiePool() *cookiePool {
	_cookiePoolOnce.Do(func() {
		CookiePool = &cookiePool{
			cookie : make([]*http.Cookie, 0),
		}
	})
	return CookiePool
}

func (c *cookiePool) Add(cookie *http.Cookie){
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cookie = append(c.cookie, cookie)
}

func (c *cookiePool) Get() *http.Cookie {
	c.mux.Lock()
	defer c.mux.Unlock()
	n := rand.Int63n(int64(len(c.cookie)))
	return c.cookie[n]
}



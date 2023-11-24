/*
*	Description : 请求上下文  TODO 测试
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
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
	"sync/atomic"
	"time"
)

// RetryTimes 重试次数
type RetryTimes int

// StartFunc 请求开始前执行的方法类型
type StartFunc func(c *Context)

// SucceedFunc 请求成功后执行的方法类型
type SucceedFunc func(c *Context)

// FailedFunc 请求失败执行的方法类型
type FailedFunc func(c *Context)

// RetryFunc 重试请求前执行的方法类型，是否重试来自上次请求的状态码来确定，见StatusCodeMap;
type RetryFunc func(c *Context)

// EndFunc 请求流程结束后执行的方法类型
type EndFunc func(c *Context)

// IsLog 是否开启全局日志
type IsLog bool

// IsRetry 是否关闭重试
type IsRetry bool

// ProxyUrl 全局代理地址，一个代理，多个代理请使用代理池ProxyPool
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

	// 请求失败了需要处理的事件
	FailedFunc FailedFunc

	// 重试请求处理的事件，可以是更换代理，设置等待时间
	RetryFunc RetryFunc

	// 请求开始前处理的事件
	StartFunc StartFunc

	// 请求完成后处理的事件
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
	isRetry IsRetry

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
	Param map[string]any

	StateCode int
}

// SetSucceedFunc 设置成功后执行的方法
func (c *Context) SetSucceedFunc(successFunc func(c *Context)) *Context {
	c.SucceedFunc = successFunc
	return c
}

// SetFailedFunc 设置错误后执行的方法
func (c *Context) SetFailedFunc(failedFunc func(c *Context)) *Context {
	c.FailedFunc = failedFunc
	return c
}

// SetRetryFunc 重试请求前执行的方法
func (c *Context) SetRetryFunc(retryFunc func(c *Context)) *Context {
	c.RetryFunc = retryFunc
	return c
}

// SetRetryTimes 设置重试最大次数，如果超过这个次数还没有请求成功，请求结束，返回一个错误；
func (c *Context) SetRetryTimes(times int) *Context {
	c.MaxTimes = RetryTimes(times)
	return c
}

// GetRetryTimes 获取当前重试最大次数
func (c *Context) GetRetryTimes() int {
	return int(c.times)
}

// SetSleep 设置请求延迟时间，单位秒
func (c *Context) SetSleep(i int) *Context {
	c.sleep = Sleep(time.Duration(i) * time.Second)
	return c
}

// SetSleepRand 设置延迟随机时间，单位秒
func (c *Context) SetSleepRand(min, max int) *Context {
	r := randEr.Intn(max) + min
	c.sleep = Sleep(time.Duration(r) * time.Second)
	return c
}

// Do 执行请求
func (c *Context) Do() func() {
	var bodyBytes []byte
	if c == nil {
		return nil
	}
	//执行 start
	if c.times == 0 && c.StartFunc != nil {
		c.StartFunc(c)
	}
	//执行 end
	if c.times == c.MaxTimes && c.EndFunc != nil {
		c.EndFunc(c)
	}
	//重试验证
	c.times++
	if c.times > c.MaxTimes {
		Error("【日志】 请求失败操过", c.MaxTimes, " 次了,结束重试操作！")
		// 超过了重试次数，就算失败，则执行失败方法
		if c.FailedFunc != nil {
			c.FailedFunc(c)
		}
		return nil
	}
	// 复用 Req.Body
	if c.Req.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Req.Body)
	}
	c.Req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// 设置的休眠时间
	if c.sleep != 0 {
		time.Sleep(time.Duration(c.sleep))
	}
	// 执行请求
	before := time.Now()
	c.Resp, c.Err = c.Client.Do(c.Req)
	c.Ms = time.Now().Sub(before)
	if c.Resp != nil {
		c.StateCode = c.Resp.StatusCode
	}
	// 是否超时
	if c.Err != nil && (strings.Contains(c.Err.Error(), "(Client.Timeout exceeded while awaiting headers)") ||
		strings.Contains(c.Err.Error(), "Too Many Requests") ||
		strings.Contains(c.Err.Error(), "To Many Requests") ||
		strings.Contains(c.Err.Error(), "EOF") ||
		strings.Contains(c.Err.Error(), "connection timed out")) {
		Error("【日志】 请求 超时 = ", c.Err)
		if c.RetryFunc != nil && c.isRetry == true {
			InfoTimes(4, "【日志】 执行 retry 事件： 第", c.times, "次， 总： ", c.MaxTimes)
			c.Req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			c.RetryFunc(c)
			return c.Do()
		}
		return nil
	}

	// 其他错误
	if c.Err != nil {
		Error("【日志】 请求 err = ", c.Err)
		// 指定的失败都执行 retry
		if c.err2retry && c.RetryFunc != nil && c.isRetry == true {
			InfoTimes(4, "【日志】 执行 retry 事件： 第", c.times, "次， 总： ", c.MaxTimes)
			c.Req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			c.RetryFunc(c)
			return c.Do()
		}
		if c.FailedFunc != nil {
			c.FailedFunc(c)
		}
		return nil
	}

	defer func(c *Context) {
		if c.Resp != nil {
			_ = c.Resp.Body.Close()
		}
	}(c)

	if c.Resp.Header.Get("Content-Encoding") == "gzip" {
		c.Resp.Body, _ = gzip.NewReader(c.Resp.Body)
	}
	HTTPTimes(4, "[", c.Req.Method, "] ", c.Req.URL.String(), " \t| s: ", c.Resp.StatusCode, "| t: ", c.Ms)

	// 根据状态码配置的事件了类型进行该事件的方法
	v, ok := StatusCodeMap[c.Resp.StatusCode]
	if !ok {
		body, err := io.ReadAll(c.Resp.Body)
		if err != nil {
			Error(err)
			return nil
		}
		c.RespBody = body
	}

	switch v {
	case "success":
		body, err := io.ReadAll(c.Resp.Body)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			Error(err)
			return nil
		}
		c.RespBody = body
		if c.SucceedFunc != nil {
			c.SucceedFunc(c)
		}

	case "retry":
		c.Req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		if c.RetryFunc != nil && c.isRetry {
			InfoTimes(4, "[日志] 执行 retry 事件： 第", c.times, "次， 总： ", c.MaxTimes)
			c.Req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			c.RetryFunc(c)
			return c.Do()
		}

	case "fail":
		if c.FailedFunc != nil {
			InfoTimes(4, "[日志] 执行 failed 事件")
			c.FailedFunc(c)
		}

	case "start":
		if c.StartFunc != nil {
			InfoTimes(4, "[日志] 执行 请求前的方法")
			c.StartFunc(c)
		}

	case "end":
		if c.EndFunc != nil {
			InfoTimes(4, "[日志] 执行 请求结束后的方法")
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

// RespBodyString 请求到的Body转换为string类型
func (c *Context) RespBodyString() string {
	if c.RespBody != nil {
		return string(c.RespBody)
	}
	return ""
}

// RespBodyHtml 请求到的Body转换为string类型的html格式
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

// RespBodyMap 请求到的Body转换为map[string]any类型
func (c *Context) RespBodyMap() map[string]any {
	var tempMap map[string]any
	err := json.Unmarshal(c.RespBody, &tempMap)
	if err != nil {
		Error(err)
		return nil
	}
	return tempMap
}

// RespBodyArr 请求到的Body转换为[]any类型
func (c *Context) RespBodyArr() []any {
	var tempArr []any
	err := json.Unmarshal(c.RespBody, &tempArr)
	if err != nil {
		Error(err)
		return nil
	}
	return tempArr
}

// GetRespHeader 获取Response的Header
func (c *Context) GetRespHeader() string {
	header, _ := httputil.DumpResponse(c.Resp, false)
	return string(header)
}

// RespContentLength 获取Response ContentLength的值
func (c *Context) RespContentLength() int64 {
	return c.Resp.ContentLength
}

// CheckReqMd5 本次请求的md5值， url+reqBodyBytes+method
func (c *Context) CheckReqMd5() string {
	var buffer bytes.Buffer
	urlStr := c.Req.URL.String()
	reqBodyBytes, _ := io.ReadAll(c.Req.Body)
	method := c.Req.Method
	buffer.WriteString(urlStr)
	buffer.Write(reqBodyBytes)
	buffer.WriteString(method)
	h := md5.New()
	h.Write(buffer.Bytes())
	return hex.EncodeToString(h.Sum(nil))
}

// CheckMd5 本次请求上下文包含输出结果的md5值
func (c *Context) CheckMd5() string {
	var buffer bytes.Buffer
	urlStr := c.Req.URL.String()
	reqBodyBytes, _ := io.ReadAll(c.Req.Body)
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

// AddHeader 给请求添加Header
func (c *Context) AddHeader(k, v string) *Context {
	c.Req.Header.Add(k, v)
	return c
}

// AddCookie 给请求添加Cookie
func (c *Context) AddCookie(k, v string) *Context {
	cookie := &http.Cookie{Name: k, Value: v, HttpOnly: true}
	c.Req.AddCookie(cookie)
	return c
}

// SetProxy 给请求设置代理
func (c *Context) SetProxy(proxyUrl string) *Context {
	proxy, _ := url.Parse(proxyUrl)
	c.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	return c
}

// SetProxyFunc 给请求设置代理函数 f func() *http.Transport
func (c *Context) SetProxyFunc(f func() *http.Transport) *Context {
	c.Client.Transport = f()
	return c
}

// SetProxyPool 给请求设置代理池
func (c *Context) SetProxyPool(pool *ProxyPool) *Context {
	ip, _ := pool.Get()
	InfoFTimes(3, "[日志] 使用代理: %s", ip)
	proxy, _ := url.Parse(ip)
	c.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	return c
}

// ProxyIP 代理IP
type ProxyIP struct {
	IP    string
	Post  int
	User  string
	Pass  string // 密码
	IsTLS bool   // 是否是 https
}

// NewProxyIP 实例化代理IP
func NewProxyIP(ip string, port int, user, pass string, isTls bool) *ProxyIP {
	return &ProxyIP{
		IP:    ip,
		Post:  port,
		User:  user,
		Pass:  pass,
		IsTLS: isTls,
	}
}

// String 代理IP输出
func (p *ProxyIP) String() string {
	h := "http://"
	if p.IsTLS {
		h = "https://"
	}
	return fmt.Sprintf("%s%s:%s@%s:%d", h, p.User, p.Pass, p.IP, p.Post)
}

// ProxyPool 代理池
type ProxyPool struct {
	ip   []string
	num  int32
	len  int32
	lock sync.Mutex
}

// NewProxyPool 实例化代理池
func NewProxyPool() *ProxyPool {
	return &ProxyPool{
		ip:   make([]string, 0),
		num:  0,
		len:  0,
		lock: sync.Mutex{},
	}
}

// Add 代理池添加代理
func (p *ProxyPool) Add(proxyIP *ProxyIP) {
	p.ip = append(p.ip, proxyIP.String())
	p.len++
}

// Del 代理池删除代理
func (p *ProxyPool) Del(n int) {
	if int32(n+1) > p.len {
		return
	}
	p.ip = append(p.ip[:n], p.ip[n+1:]...)
	p.len--
}

// Get 代理池按获取顺序获取一个代理
func (p *ProxyPool) Get() (string, int) {
	p.lock.Lock()
	if p.num >= p.len {
		p.num = 0
	}
	data := p.ip[p.num]
	p.lock.Unlock()
	n := p.num
	atomic.AddInt32(&p.num, int32(1))
	return data, int(n)
}

// Upload 下载
func (c *Context) Upload(filePath string) func() {
	//空验证
	if c == nil {
		return nil
	}

	//重试验证
	c.times++
	if c.times > c.MaxTimes {
		InfoF("请求失败操过", c.MaxTimes, "次了")
		return nil
	}

	//执行请求
	c.Resp, c.Err = c.Client.Do(c.Req)

	// 是否超时
	// Go 1.6 以下
	if c.Err != nil && strings.Contains(c.Err.Error(), "(Client.Timeout exceeded while awaiting headers)") {
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
		if c.FailedFunc != nil {
			c.FailedFunc(c)
		}
		return nil
	}

	defer func(cxt *Context) {
		if cxt.Resp != nil {
			_ = cxt.Resp.Body.Close()
		}
	}(c)

	f, err := os.Create(filePath)
	if err != nil {
		c.Err = err
		return nil
	}

	defer func() {
		_ = f.Close()
	}()

	contentLength := Str2Float64(c.Resp.Header.Get("Content-Length"))
	var sum int64 = 0
	buf := make([]byte, 1024*100)
	sTime := time.Now()
	i := 0
	for {
		i++
		n, err := c.Resp.Body.Read(buf)
		sum = sum + int64(n)
		_, _ = f.Write(buf[:n])
		if err != nil || n == 0 {
			break
		}
		if i%9 == 0 {
			log.Println("[下载] ", filePath, " : ", FileSizeFormat(sum), "/", FileSizeFormat(int64(contentLength)),
				" |\t ", math.Floor((float64(sum)/contentLength)*100), "%")
		}
	}
	ct := time.Now().Sub(sTime)
	log.Println("[下载] ", filePath, " : ", FileSizeFormat(sum), "/", FileSizeFormat(int64(contentLength)),
		" |\t ", math.Floor((float64(sum)/contentLength)*100), "%", "|\t ", ct)
	return nil
}

// CookieNext 使用上下文cookies
func (c *Context) CookieNext() error {
	if c.Resp == nil {
		return fmt.Errorf("response is nil")
	}
	if c.Req == nil {
		return fmt.Errorf("request is nil")
	}
	// 上下文cookies
	for _, cookie := range c.Resp.Cookies() {
		c.Req.AddCookie(cookie)
	}
	return nil
}

// CloseLog 关闭日志打印
func (c *Context) CloseLog() {
	c.IsLog = false
}

// OpenErr2Retry 开启请求失败都执行retry
func (c *Context) OpenErr2Retry() {
	c.err2retry = true
}

// CloseRetry 关闭重试
func (c *Context) CloseRetry() {
	c.isRetry = false
}

// GetParam 获取上下文参数
func (c *Context) GetParam(key string) any {
	return c.Param[key]
}

// AddParam 添加上下文参数
func (c *Context) AddParam(key string, val any) {
	c.Param[key] = val
}

// DelParam 删除上下文参数
func (c *Context) DelParam(key string) {
	delete(c.Param, key)
}

// CookiePool  cookie池
type cookiePool struct {
	cookie []*http.Cookie
	mux    sync.Mutex
}

// CookiePool  Cookie池
var CookiePool *cookiePool
var _cookiePoolOnce sync.Once

// NewCookiePool 实例化cookie池
func NewCookiePool() *cookiePool {
	_cookiePoolOnce.Do(func() {
		CookiePool = &cookiePool{
			cookie: make([]*http.Cookie, 0),
		}
	})
	return CookiePool
}

// Add cookie池添加cookie
func (c *cookiePool) Add(cookie *http.Cookie) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cookie = append(c.cookie, cookie)
}

// Get 在cookie池获取一个cookie
func (c *cookiePool) Get() *http.Cookie {
	c.mux.Lock()
	defer c.mux.Unlock()
	n := rand.Int63n(int64(len(c.cookie)))
	return c.cookie[n]
}

// TODO 高并发下载

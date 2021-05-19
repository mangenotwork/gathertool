/*
	Description : 请求上下文
	Author : ManGe
	Version : v0.3
	Date : 2021-05-19
*/

package gathertool

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// 重试次数
type  RetryTimes int

// 请求开始前的方法类型
type StartFunc func(c *Context)

// 成功后的方法类型
type SucceedFunc func(c *Context)

// 失败后的方法类型
type FailedFunc func(c *Context)

// 重试前的方法类型
type RetryFunc func(c *Context)

// 请求结束后的方法类型
type EndFunc func(c *Context)


// 请求上下文
type Context struct {
	// Token
	Token string

	// client
	Client *http.Client

	// Request
	Req *http.Request

	// Response
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
	IsLog bool

	// 指定失败执行重试事件
	err2retry bool

	// 是否关闭重试
	isRetry bool
}

// SetSucceedFunc 设置成功后的方法
func (c *Context) SetSucceedFunc(successFunc func(c *Context)){
	c.SucceedFunc = successFunc
}

// SetFailed 设置错误后的方法
func (c *Context) SetFailedFunc(failedFunc func(c *Context)) {
	c.FailedFunc = failedFunc
}

// SetRetryFunc 设置重试，在重试前的方法
func (c *Context) SetRetryFunc(retryFunc func(c *Context)) {
	c.RetryFunc = retryFunc
}

// SetRetryTimes 设置重试次数
func (c *Context) SetRetryTimes(times int) {
	c.MaxTimes = RetryTimes(times)
}

// Do 执行请求
func (c *Context) Do() func(){

	var bodyBytes []byte

	//空验证
	if c == nil{
		log.Println("空对象")
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
		logerTimes(2 + int(c.times), "【日志】 请求失败操过", c.MaxTimes, "次了,结束重试操作；")

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

	//执行请求
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
		loger("【日志】 请求 超时 = ", c.Err)

		if c.RetryFunc != nil && !c.isRetry {
			logerTimes(2 + int(c.times), "【日志】 执行 retry 事件： 第", c.times, "次， 总： ",  c.MaxTimes)
			c.Req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			c.RetryFunc(c)
			return c.Do()
		}

		return nil
	}

	// 其他错误
	if c.Err != nil {
		logerTimes(2 + int(c.times), "【日志】 请求 err = ", c.Err)

		// 指定的失败都执行 retry
		if c.err2retry && c.RetryFunc != nil && !c.isRetry {
			logerTimes(2 + int(c.times), "【日志】 执行 retry 事件： 第", c.times, "次， 总： ",  c.MaxTimes)
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
			//loger("【日志】 请求 Resp close")
			c.Resp.Body.Close()
		}
	}(c)

	logerTimes(2 + int(c.times), "【日志】 请求状态码：", c.Resp.StatusCode, " | 用时 ： ", c.Ms)

	// 根据状态码配置的事件了类型进行该事件的方法
	if v,ok := StatusCodeMap[c.Resp.StatusCode]; ok{
		switch v {

		case "success":
			logerTimes(2 + int(c.times), "【日志】 执行 success 事件")
			//请求后的结果
			body, err := ioutil.ReadAll(c.Resp.Body)
			if err != nil{
				log.Println(err)
				return nil
			}

			c.RespBody = body
			//执行成功方法
			if c.SucceedFunc != nil {
				c.SucceedFunc(c)
			}

		case "retry":
			c.Req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			if c.RetryFunc != nil && !c.isRetry {
				logerTimes(2 + int(c.times), "【日志】 执行 retry 事件： 第", c.times, "次， 总： ",  c.MaxTimes)
				c.Req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
				c.RetryFunc(c)
				return c.Do()
			}

		case "fail":
			if c.FailedFunc != nil{
				logerTimes(2 + int(c.times), "【日志】 执行 failed 事件")
				c.FailedFunc(c)
			}

		case "start":
			if c.StartFunc != nil {
				logerTimes(2 + int(c.times), "【日志】 执行 请求前的方法")
				c.StartFunc(c)
			}

			case "end":
				if c.EndFunc != nil {
					logerTimes(2 + int(c.times), "【日志】 执行 请求结束后的方法")
					c.EndFunc(c)
				}
		}
	}

	return nil
}

// add header
func (c *Context) AddHeader(k,v string) {
	c.Req.Header.Add(k,v)
}

// add Cookie
func (c *Context) AddCookie(cookie *http.Cookie){
	c.Req.AddCookie(cookie)
}

// Upload 下载
func (c *Context) Upload(filePath string) func(){
	//空验证
	if c == nil{
		log.Println("空对象")
		return nil
	}

	//重试验证
	c.times++
	if c.times > c.MaxTimes{
		log.Println("请求失败操过", c.MaxTimes, "次了")
		return nil
	}

	//执行请求
	c.Resp,c.Err = c.Client.Do(c.Req)

	// 是否超时
	if c.Err != nil && strings.Contains(c.Err.Error(), "(Client.Timeout exceeded while awaiting headers)"){
		if c.RetryFunc != nil {
			c.RetryFunc(c)
			return c.Do()
		}
		return nil
	}

	// 其他错误
	if c.Err != nil {
		log.Println("err = ", c.Err)
		if c.FailedFunc != nil{
			c.FailedFunc(c)
		}
		return nil
	}
	defer func(cxt *Context){
		if cxt.Resp != nil {
			cxt.Resp.Body.Close()
		}
	}(c)

	f, err := os.Create(filePath)
	if err != nil {
		c.Err = err
		return nil
	}
	defer f.Close()

	contentLength := Str2Float64(c.Resp.Header.Get("Content-Length"))
	var sum int64 = 0
	buf := make([]byte, 1024*100)
	st := time.Now()
	i := 0
	for {
		i++
		n, err := c.Resp.Body.Read(buf)
		sum=sum+int64(n)
		if err != nil || n == 0{
			f.Write(buf[:n])
			break
		}
		f.Write(buf[:n])
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
		return errors.New("Response is nil")
	}
	if c.Req == nil {
		return errors.New("Request is nil")
	}
	// 上下文cookies
	for _,cookie := range c.Resp.Cookies(){
		c.Req.AddCookie(cookie)
	}
	return nil
}

// close log
func (c *Context) CloseLog() {
	c.IsLog = false
}

// 开启请求失败都执行retry
func (c *Context) OpenErr2Retry() {
	c.err2retry = true
}

// 关闭重试
func (c *Context) CloseRetry() {
	c.isRetry = true
}

// CookiePool   cookie池
type cookiePool struct {
	cookie []*http.Cookie
	mux sync.Mutex
}

var CookiePool = &cookiePool{}

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



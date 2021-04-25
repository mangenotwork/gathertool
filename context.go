/*
	Description : 请求上下文
	Author : ManGe
	Version : v0.2
	Date : 2021-04-25
*/

package gathertool

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// 重试次数
type  RetryTimes int

// 请求开始前的方法
type StartFunc func(c *Context)

// 成功后的方法
type SucceedFunc func(c *Context)

// 失败后的方法
type FailedFunc func(c *Context)

// 重试前的方法
type RetryFunc func(c *Context)

// 请求结束后的方法
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
	Task *Task

	RespBody []byte

	// job 编号
	JobNumber int

	// 请求的响应时间 单位ms
	Ms time.Duration
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

// Do 执行请求
func (c *Context) Do() func(){

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
		log.Println("请求失败操过", c.MaxTimes, "次了")
		return nil
	}

	//执行请求
	time_start := time.Now()
	c.Resp,c.Err = c.Client.Do(c.Req)
	c.Ms = time.Since(time_start)
	if c.Err != nil {
		log.Println("err = ", c.Err)
		if c.FailedFunc != nil{
			c.FailedFunc(c)
		}
		return nil
	}

	defer func(){
		if c.Resp != nil {
			c.Resp.Body.Close()
		}
	}()

	//log.Println("状态码：", c.Resp.StatusCode)

	// 根据状态码配置的事件了类型进行该事件的方法
	if v,ok := StatusCodeMap[c.Resp.StatusCode]; ok{
		switch v {
		case "success":
			//log.Println("执行 success 事件")
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

			return nil
		case "retry":
			//log.Println("执行 retry 事件")

			log.Println("第", c.times, "请求失败,状态码： ", c.Resp.StatusCode, ".")
			//执行重试前的方法
			if c.RetryFunc != nil{
				c.RetryFunc(c)
			}
			return c.Do()
		case "file":
			//log.Println("执行 file 事件")
			if c.FailedFunc != nil{
				c.FailedFunc(c)
			}
			return nil
		}
	}

	return nil
}

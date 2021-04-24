/*
	Description : 请求实现
	Author : ManGe
	Version : v0.1
	Date : 2021-04-23
*/

package gathertool

import (
	"io/ioutil"
	"log"
	"net/http"
)

type  RetryTimes int64

type Req struct {

	// client
	Client *http.Client

	// 请求
	Req *http.Request

	// 执行的次数 初始化都是0
	times RetryTimes

	// 最大允许重试次数
	MaxTimes RetryTimes

	// 请求成功了需要处理的事件
	SuccessFunc func([]byte)

	// 请求失败了需要做的事
	FailFunc func()

	// 请求状态码设置了重试，在重试前的事件
	RetryFunc func(req *Req)
}

func (r *Req) Succeed(successFunc func([]byte)){
	r.SuccessFunc = successFunc
}

// Failed 设置错误处理
func (r *Req) Failed(failedFunc func()) {
	r.FailFunc = failedFunc
}

// Retry 请求状态码设置了重试，在重试前的事件
func (r *Req) Retry(retryFunc func(c *Req)) {
	r.RetryFunc = retryFunc
}

// Do 执行请求
func (r *Req) Do() func(){

	//空验证
	if r == nil{
		log.Println("空对象")
		return nil
	}

	//重试验证
	r.times++
	if r.times > r.MaxTimes{
		log.Println("请求失败操过", r.MaxTimes, "次了")
		return nil
	}

	//执行请求
	resp,err := r.Client.Do(r.Req)
	if err != nil {
		log.Println("请求出现错误: ", err)
		return nil
	}
	defer resp.Body.Close()

	log.Println("状态码：", resp.StatusCode)
	if v,ok := StatusCodeMap[resp.StatusCode]; ok{
		switch v {
		case "success":
			//请求后的结果
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil{
				log.Println(err)
				return nil
			}
			//执行成功方法
			r.SuccessFunc(body)
			return nil
		case "retry":
			log.Println("第", r.times, "请求失败,状态码： ", resp.StatusCode, ".")
			//执行重试前的方法
			if r.RetryFunc != nil{
				r.RetryFunc(r)
			}
			return r.Do()
		case "file":
			r.FailFunc()
			return nil
		}
	}

	return nil
}

/*
	Description : 对外提供的方法
	Author : ManGe
	Version : v0.1
	Date : 2021-04-23
*/

package gathertool

import (
	"errors"
	"log"
	"net/http"
	"sync"
)

// Get 请求, 当请求失败或状态码是失败的则会先执行 ff 再回调
//
// @url  请求链接
// @maxTimes  重试次数
// @sf  请求成功后做的事情, 200等
// @ff  请求失败后做的事情, 403等，502等
// @vs  可变参数
// @vs UserAgentType  设置指定类型 user agent 如 AndroidAgent
//
func Get(url string, vs ...interface{}) (*Context,error){
	var (
		client *http.Client
		maxTimes RetryTimes = 10
		task *Task
		start StartFunc
		succeed SucceedFunc
		failed FailedFunc
		retry RetryFunc
		end EndFunc
	)

	if url == "" {
		return nil, errors.New("空的地址")
	}

	//初始化 Request
	req, err := http.NewRequest("GET",url,nil)
	if err != nil{
		return nil,err
	}

	//添加默认的Header
	req.Header.Add("Connection","close")
	req.Header.Add("User-Agent", GetAgent(PCAgent))

	//解析可变参
	for _, v := range vs {
		//log.Println("参数： ", v)
		switch vv := v.(type) {
		// 使用方传入了 header
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
			// 使用方传入了 *http.Client
		case *http.Client:
			client = vv
		case UserAgentType:
			req.Header.Add("User-Agent", GetAgent(vv))
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
		}
	}

	// 如果使用方未传入Client，  初始化 Client
	if client == nil{
		//log.Println("使用方未传入Client， 默认 client")
		client = &http.Client{}
	}

	// 创建对象
	return &Context{
		Client: client,
		Req : req,
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

// StartJobGet 并发执行Get,直到队列任务为空
// @jobNumber 并发数，
// @queue 全局队列，
// @client 单个并发任务的client，
// @SucceedFunc 成功方法，
// @ RetryFunc重试方法，
// @FailedFunc 失败方法
//
func StartJobGet(jobNumber int, queue TodoQueue, client *http.Client,
	SucceedFunc func(ctx *Context),
	RetryFunc func(ctx *Context),
	FailedFunc func(ctx *Context)){
	var wg sync.WaitGroup
	for job:=0;job<jobNumber;job++{
		wg.Add(1)
		go func(i int){
			log.Println("启动第",i ,"个任务")
			defer wg.Done()
			for {
				if queue.IsEmpty(){
					break
				}
				task := queue.Poll()
				log.Println("第",i,"个任务取的值： ", task)
				ctx, err := Get(task.Url, client, task)
				if err != nil {
					log.Println(err)
					return
				}
				ctx.SetSucceedFunc(SucceedFunc)
				ctx.SetRetryFunc(RetryFunc)
				ctx.SetFailedFunc(FailedFunc)
				ctx.Do()
			}
			log.Println("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	log.Println("执行完成！！！")
}
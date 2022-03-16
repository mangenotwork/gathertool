/*
	Description : 并发工作任务
	Author : ManGe
	Version : v0.2
	Date : 2021-08-12
*/

package gathertool

import (
	"net/http"
	"runtime"
	"sync"
)

// StartJob 开始运行并发
func StartJob(jobNumber int, queue TodoQueue,f func(task *Task)){
	var wg sync.WaitGroup
	for job:=0;job<jobNumber;job++{
		wg.Add(1)
		go func(i int){
			Info("启动第",i ,"个任务")
			defer wg.Done()
			for {
				if queue.IsEmpty() {
					break
				}
				task := queue.Poll()
				Info("第", i, "个任务取的值： ", task)
				f(task)
			}
			Info("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	Info("执行完成！！！")
}

// 设置遇到错误执行 Retry 事件
type Err2Retry bool

// StartJobGet 并发执行Get,直到队列任务为空
// jobNumber 并发数，
// queue 全局队列，
func StartJobGet(jobNumber int, queue TodoQueue, vs ...interface{}){

	var (
		client *http.Client
		succeed SucceedFunc
		retry RetryFunc
		failed FailedFunc
		err2Retry Err2Retry
	)

	runtime.GOMAXPROCS(runtime.NumCPU())

	for _,v := range vs{
		switch vv := v.(type) {
		case *http.Client:
			client = vv
		case SucceedFunc:
			succeed = vv
		case FailedFunc:
			failed = vv
		case RetryFunc:
			retry = vv
		case Err2Retry:
			err2Retry = vv
			}
	}

	var wg sync.WaitGroup
	for job:=0;job<jobNumber;job++{
		wg.Add(1)
		go func(i int){
			Info("启动第",i ,"个任务")
			defer wg.Done()
			for {
				if queue.IsEmpty(){
					break
				}
				task := queue.Poll()
				Info("第",i,"个任务取的值： ", task)
				ctx := NewGet(task.Url, task)
				if client != nil {
					ctx.Client = client
				}
				if succeed != nil {
					ctx.SetSucceedFunc(succeed)
				}
				if retry != nil {
					ctx.SetRetryFunc(retry)
				}
				if failed != nil {
					ctx.SetFailedFunc(failed)
				}
				if err2Retry {
					ctx.OpenErr2Retry()
				}

				switch task.Type {
				case "","do":
					ctx.Do()
				case "upload":
					if task.SavePath == ""{
						task.SavePath = task.SaveDir + task.FileName
					}
					ctx.Upload(task.SavePath)
				default:
					ctx.Do()
				}
			}
			Info("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	Info("执行完成！！！")
}


// StartJobPost 开始运行并发Post
func StartJobPost(jobNumber int, queue TodoQueue, vs ...interface{}){
	var (
		client *http.Client
		succeed SucceedFunc
		retry RetryFunc
		failed FailedFunc
		err2Retry Err2Retry
	)

	for _,v := range vs{
		switch vv := v.(type) {
		case *http.Client:
			client = vv
		case SucceedFunc:
			succeed = vv
		case FailedFunc:
			failed = vv
		case RetryFunc:
			retry = vv
		case Err2Retry:
			err2Retry = vv
		}
	}

	var wg sync.WaitGroup
	for job:=0;job<jobNumber;job++{
		wg.Add(1)
		go func(i int){
			Info("启动第",i ,"个任务")
			defer wg.Done()
			for {
				if queue.IsEmpty(){
					break
				}
				task := queue.Poll()
				Info("第",i,"个任务取的值： ", task, task.HeaderMap)
				ctx := NewPost(task.Url, []byte(task.JsonParam), "application/json;", task)
				if client != nil {
					ctx.Client = client
				}
				if succeed != nil {
					ctx.SetSucceedFunc(succeed)
				}
				if retry != nil {
					ctx.SetRetryFunc(retry)
				}
				if failed != nil {
					ctx.SetFailedFunc(failed)
				}
				if err2Retry {
					ctx.OpenErr2Retry()
				}

				switch task.Type {
				case "","do":
					ctx.Do()
				case "upload":
					if task.SavePath == ""{
						task.SavePath = task.SaveDir + task.FileName
					}
					ctx.Upload(task.SavePath)
				default:
					ctx.Do()
				}

			}
			Info("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	Info("执行完成！！！")
}

// CPUMax 多核执行
func CPUMax(){
	runtime.GOMAXPROCS(runtime.NumCPU())
}


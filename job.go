/*
	Description : 并发工作任务
	Author : ManGe
	Version : v0.1
	Date : 2021-04-28
*/

package gathertool

import (
	"log"
	"net/http"
	"sync"
)

//TODO:  StartJob 开始运行并发
func StartJob(){}


// StartJobGet 并发执行Get,直到队列任务为空
// @jobNumber 并发数，
// @queue 全局队列，
// @client 单个并发任务的client，
// @SucceedFunc 成功方法，
// @ RetryFunc重试方法，
// @FailedFunc 失败方法
func StartJobGet(jobNumber int, queue TodoQueue, vs ...interface{}){

	var (
		client *http.Client
		succeed SucceedFunc
		retry RetryFunc
		failed FailedFunc
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
			}
	}

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
				ctx, err := Get(task.Url, task)
				if err != nil {
					log.Println(err)
					return
				}
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
			log.Println("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	log.Println("执行完成！！！")
}


// TODO:  StartJobPost 开始运行并发Post
func StartJobPost(jobNumber int, queue TodoQueue, vs ...interface{}){
	var (
		client *http.Client
		succeed SucceedFunc
		retry RetryFunc
		failed FailedFunc
	)

	for _,v := range vs{
		switch vv := v.(type) {
		case *http.Client:
			loger("have Client")
			client = vv
		case SucceedFunc:
			succeed = vv
		case FailedFunc:
			failed = vv
		case RetryFunc:
			retry = vv
		}
	}

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
				log.Println("第",i,"个任务取的值： ", task, task.HeaderMap)
				ctx, err := PostJson(task.Url, task.JsonParam, task)
				if err != nil {
					log.Println(err)
					return
				}
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
			log.Println("第",i ,"个任务结束！！")
		}(job)
	}
	wg.Wait()
	log.Println("执行完成！！！")
}

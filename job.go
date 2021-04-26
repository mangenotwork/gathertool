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


// TODO:  StartJobPost 开始运行并发Post
func StartJobPost(){}
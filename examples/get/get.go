package main

import (
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"net/http"
	"time"
)

func main(){
	// 设置一个 http.Client 也可以是来自第三方代理的 http.Client
	client := &http.Client{
		Timeout: 5*time.Second,
	}
	// 执行一个get 请求，最多重试10次
	c, err := gt.Get("http://192.168.0.9",client,
		gt.SucceedFunc(succeed),
		gt.FailedFunc(fail),
		gt.RetryFunc(retry))
	if err != nil {
		log.Println(err)
		return
	}
	// 执行
	c.Do()
}

// 成功后的方法
func succeed(c *gt.Context){
	fmt.Printf(string(c.RespBody))
	//处理数据
}

// 设置需要重试状态码， 重试前的方法
func retry(c *gt.Context){
	log.Println(c)
	log.Println(c.MaxTimes)
	c.Client = &http.Client{
		Timeout: 1*time.Second,
	}
	log.Println("休息1s")
	time.Sleep(1*time.Second)
}

// 失败后的方法
func fail(c *gt.Context) {
	fmt.Printf("请求失败")
	time.Sleep(1*time.Second)
}
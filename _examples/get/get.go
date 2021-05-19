package main

import (
	"log"
	"net/http"
	"time"

	gt "github.com/mangenotwork/gathertool"
)

func main(){
	//SimpleGet1()
	//SimpleGet2()
	//SimpleGet3()
	//SimpleGet4()
}


// 简单的get请求实例, 写法一： 方法做为请求函数的参数；
func SimpleGet1(){
	// 创建请求
	c, err := gt.Get("http://192.168.0.1",
		//设置请求成功后的方法： 请求的数据
		gt.SucceedFunc(func(ctx *gt.Context){
			log.Println(string(ctx.RespBody))
		}),

		//设置请求失败后的方法： 打印失败信息
		gt.FailedFunc(func(ctx *gt.Context){
			log.Println(ctx.Err)
		}),

		//设置重试前的方法（遇到403,502 等状态码会重试）： 睡眠1s再重试
		gt.RetryFunc(func(ctx *gt.Context){
			time.Sleep(1*time.Second)
		}),
	)
	if err != nil {
		log.Println("请求创建失败: ", err)
		return
	}
	// 执行创建的请求
	c.Do()
}


// 最简单的get请求实例， 写法二： 上下文处理；
func SimpleGet2(){
	// 创建请求
	c, err := gt.Get("http://192.168.0.1")
	if err != nil {
		log.Println("请求创建失败: ", err)
		return
	}
	// 执行创建的请求
	c.Do()
	// 打印请求结果与请求错误
	log.Println(string(c.RespBody), c.Err)
}


// 简单的get请求实例, 写法三： 给请求设置方法；
func SimpleGet3()  {
	// 创建请求
	c, err := gt.Get("http://192.168.0.1")
	if err != nil {
		log.Println("请求创建失败: ", err)
		return
	}
	//设置请求成功后的方法
	c.SetSucceedFunc(func(ctx *gt.Context){
		log.Println(string(ctx.RespBody))
	})
	//设置请求失败后的方法
	c.SetFailedFunc(func(ctx *gt.Context){
		log.Println(ctx.Err)
	})
	//设置重试次数
	c.SetRetryTimes(5)
	//设置重试前的方法
	c.SetRetryFunc(func(*gt.Context) {
		time.Sleep(1*time.Second)
	})
	// 执行创建的请求
	c.Do()
}

// 简单的get请求实例, 写法四： 外部函数为请求方法；
func SimpleGet4(){
	c, err := gt.Get("http://192.168.0.1",
		gt.SucceedFunc(succeed),
		gt.FailedFunc(fail),
		gt.RetryFunc(retry),
		)
	if err != nil {
		log.Println("请求创建失败: ", err)
		return
	}
	// 执行创建的请求
	c.Do()
}

// 成功后的方法
func succeed(ctx *gt.Context){
	log.Println(string(ctx.RespBody))
	//处理数据
}

// 设置需要重试状态码， 重试前的方法
func retry(ctx *gt.Context){
	ctx.Client = &http.Client{
		Timeout: 1*time.Second,
	}
	log.Println("休息1s")
	time.Sleep(1*time.Second)
}

// 失败后的方法
func fail(ctx *gt.Context) {
	log.Println("请求失败: ", ctx.Err)
}

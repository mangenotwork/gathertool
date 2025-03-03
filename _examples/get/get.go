package main

import (
	"log"
	"net/http"
	"time"

	gt "github.com/mangenotwork/gathertool"
)

func main() {
	SimpleGet1()
	SimpleGet2()
	SimpleGet3()
	SimpleGet4()

	//Case()
}

// SimpleGet1 简单的get请求实例, 写法一： 方法做为请求函数的参数；
func SimpleGet1() {
	// 创建请求
	ctx := gt.NewGet("http://192.168.3.1",
		//设置请求成功后的方法： 请求的数据
		gt.SucceedFunc(func(ctx *gt.Context) {
			log.Println(string(ctx.RespBody))
		}),

		//设置请求失败后的方法： 打印失败信息
		gt.FailedFunc(func(ctx *gt.Context) {
			log.Println(ctx.Err)
		}),

		//设置重试前的方法（遇到403,502 等状态码会重试）： 睡眠1s再重试
		gt.RetryFunc(func(ctx *gt.Context) {
			time.Sleep(1 * time.Second)
		}),
	)
	// 执行创建的请求
	ctx.Do()

}

// SimpleGet2 最简单的get请求实例， 写法二： 上下文处理；
func SimpleGet2() {
	// 创建请求
	ctx, err := gt.Get("http://192.168.3.1")
	// 打印请求结果与请求错误
	log.Println(ctx.RespBodyString(), err)
}

// SimpleGet3 简单的get请求实例, 写法三： 给请求设置方法；
func SimpleGet3() {
	// 创建请求
	gt.NewGet("http://192.168.3.1").SetSucceedFunc(func(ctx *gt.Context) {
		//设置请求成功后的方法
		log.Println(string(ctx.RespBody))
	}).SetFailedFunc(func(ctx *gt.Context) {
		//设置请求失败后的方法
		log.Println(ctx.Err)
	}).SetRetryFunc(func(*gt.Context) {
		//设置重试前的方法
		time.Sleep(1 * time.Second)
	}).Do()
}

// SimpleGet4 简单的get请求实例, 写法四： 外部函数为请求方法；
func SimpleGet4() {
	_, _ = gt.Get("http://192.168.3.1",
		gt.SucceedFunc(succeed),
		gt.FailedFunc(fail),
		gt.RetryFunc(retry),
	)
}

// 成功后的方法
func succeed(ctx *gt.Context) {
	log.Println(string(ctx.RespBody))
	//处理数据
}

// 设置需要重试状态码， 重试前的方法
func retry(ctx *gt.Context) {
	ctx.Client = &http.Client{
		Timeout: 1 * time.Second,
	}
	log.Println("休息1s")
	time.Sleep(1 * time.Second)
}

// 失败后的方法
func fail(ctx *gt.Context) {
	log.Println("请求失败: ", ctx.Err)
}

func Case() {
	a := gt.Any2String("aasddas")
	log.Println(a)
}

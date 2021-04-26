# gathertool
轻量级爬虫，接口测试，压力测试框架, 提高开发对应场景的golang程序。

## 概念
1. 一个请求含有 请求成功后方法， 请求重试前方法，请求失败后方法；举个例子 
 请一个URL  请求遇到状态200或302等算成功，则执行成功后方法。  请求遇到403或502
 等需要执行重试，则执行重试方法，重试方法主要含添加等待时间更换代理IP等。 遇到404
 或500等算失败则执行失败方法。
 
2. 并发执行，第一需要创建TODO队列，等TODO队列加载完后，每个队列对象含有上下文，
在创建时应该富裕上文数据或对象，开始执行并发任务，每个并发任务是一个独立的cilet，
当队列任务取完后则整个并发结束。注意这里的每个并发任务都是独立的，没有chan操作。


## 请求
### Get
> 简单的get请求实例, 写法一： 方法做为请求函数的参数；
```golang
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
```

> 最简单的get请求实例， 写法二： 上下文处理；
```golang
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
```

> 简单的get请求实例, 写法三： 给请求设置方法；
```golang
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
```

> 简单的get请求实例, 写法四： 外部函数为请求方法；
```golang
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
```
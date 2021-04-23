package main

import (
	"fmt"
	"github.com/mangenotwork/gathertool"
	"log"
	"net/http"
	"time"
)

func main(){
	// 设置一个 http.Client 也可以是来自第三方代理的 http.Client
	c := &http.Client{
		Timeout: 5*time.Second,
	}
	// 执行一个get 请求，最多重试10次
	req, err := gathertool.Get("http://192.168.0.1", 10,
		c,
		"adasdas",
		1235646)
	if err != nil{
		log.Println(err)
		return
	}
	// 设置成功执行的回调
	req.Succeed(succeed)
	// 设置这个请求遇到失败状态码的重试前的操作，如402等的处理事件
	// 例如添加 等待时间，更换代理，更换Header等
	req.Failed(failed)
	// 执行
	req.Do()
}

// 成功后的方法
func succeed(b []byte){
	fmt.Printf(string(b))
	//处理数据
}

// 错误状态码重试前的方法
func failed(c *gathertool.Req){
	log.Println(c)
	log.Println(c.MaxTimes)
	c.Client = &http.Client{
		Timeout: 1*time.Second,
	}
	log.Println("休息1s")
	time.Sleep(1*time.Second)
}
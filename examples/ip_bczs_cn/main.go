package main

import (
	"github.com/PuerkitoBio/goquery"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"net/http"
	"time"
)

// 全局变量
var (
	queue = gt.NewQueue()
)

func main(){
	// 抓取国内ip段， 入口 http://ip.bczs.net/country/CN
	// 然后便利抓取所有ip段的信息

	//1. 先请入口获取所有ip段然后创建任务队列
	req, err := gt.Get("http://ip.bczs.net/country/CN",gt.SucceedFunc(IPListSucceed),
		gt.FailedFunc(func(ctx *gt.Context){time.Sleep(1*time.Second)}),
		gt.RetryFunc(func(ctx *gt.Context){log.Println("请求失败")}))
	if err != nil{
		log.Println(err)
		return
	}
	req.Do()

	//创建并发任务请求的client对象
	client := &http.Client{
		// 设置代理
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(uri),
		//},
		Timeout: 5*time.Second,
	}
	//2. 抓取详情数据
	// 当队列加载完毕后，执行并发任务，只有当 queue 完成后结束
	gt.StartJobGet(100,queue, client, GetIPSucceed, GetIPRetry, GetIPFailed)

	queue.Print()
}

// 请求成功执行
func IPListSucceed(cxt *gt.Context){
	html := string(cxt.RespBody)
	dom,err := gt.NewGoquery(html)
	if err != nil{
		log.Println(err)
		return
	}
	result := dom.Find("div[id=result] tbody")
	result.Find("tr").Each(func(i int, tr *goquery.Selection){
		log.Println("第", i+1, "tr：")

		td := tr.Find("td")
		// IP起始
		startIp := td.Eq(0).Text()
		log.Println("IP起始 : ", startIp)
		// 结束
		endIP := td.Eq(1).Text()
		log.Println("结束 : ", endIP)
		// 数量
		number := td.Eq(2).Text()
		log.Println("数量 : ", number)

		// 创建队列 抓取详情信息
		// http://ip.bczs.net/1.0.1.0
		queue.Add(&gt.Task{
			Url: "http://ip.bczs.net/"+startIp,
			Context: map[string]interface{}{
				"start_ip":startIp,
				"end_ip":endIP,
				"number":number,
			},
		},
		)

		log.Println("\n\n")
	})
}

// 获取详情信息成功的处理
func GetIPSucceed(c *gt.Context){
	log.Println(c.Task)

	html := string(c.RespBody)
	dom,err := gt.NewGoquery(html)
	if err != nil{
		log.Println(err)
		return
	}
	result,err := dom.Find("div[id=result] .well").Html()
	if err != nil{
		log.Println(err)
	}
	log.Println(result)
	time.Sleep(3*time.Second)
}

// 获取详情信息重试的处理
func GetIPRetry(c *gt.Context){
	//更换代理
	c.Client = &http.Client{
		// 设置代理
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(uri),
		//},
		Timeout: 5*time.Second,
	}
	log.Println("休息1s")
	time.Sleep(1*time.Second)
}

// 获取详情信息失败执行
func GetIPFailed(c *gt.Context){
	log.Println("请求失败")
}
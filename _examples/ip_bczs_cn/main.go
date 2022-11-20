package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	gt "github.com/mangenotwork/gathertool"
)

func main() {
	//main2()
	main1()
}

var (
	// 全局声明抓取任务队列
	queue = gt.NewQueue()
	// 全局声明数据库客户端对象
	//host   = "192.168.0.192"
	//port    = 3306
	//user      = "root"
	//password  = "root123"
	//dbname = "test"
	//db,_ = gt.NewMysql(host, port, user, password, dbname)
	c, _ = gt.NewCSV("data.xls")
)

func main1() {
	// 1.在页面 http://ip.bczs.net/country/CN 获取所有ip
	_, _ = gt.Get("http://ip.bczs.net/country/CN", gt.SucceedFunc(IPListSucceed))

	// 2. 并发抓取详情数据, 20个并发
	gt.StartJobGet(10, queue,
		gt.SucceedFunc(GetIPSucceed), //请求成功后执行的方法
		gt.RetryFunc(GetIPRetry),     //遇到 502,403 等状态码重试前执行的方法，一般为添加休眠时间或更换代理
		gt.FailedFunc(GetIPFailed),   //请求失败后执行的方法
	)
	//c.Close()

	select {}

}

// IPListSucceed 请求成功执行
func IPListSucceed(ctx *gt.Context) {
	i := 1
	for _, tbody := range gt.RegHtmlTbody(ctx.RespBodyString()) {
		for _, tr := range gt.RegHtmlTr(tbody) {
			if i > 500 {
				goto End
			}
			td := gt.RegHtmlTdTxt(tr)
			log.Println(td)
			if len(td) < 3 {
				gt.Error("异常数据 ： ", td)
				continue
			}
			startIp := gt.Any2String(gt.RegHtmlATxt(td[0])[0]) // IP起始
			endIP := td[1]                                     // 结束
			number := td[2]                                    // 数量
			// 创建队列 抓取详情信息
			queue.Add(&gt.Task{
				Url: "http://ip.bczs.net/" + startIp,
				Data: map[string]interface{}{
					"start_ip": startIp,
					"end_ip":   endIP,
					"number":   number,
				},
			})
			i++
		}
	}
End:
}

// GetIPSucceed 获取详情信息成功的处理
func GetIPSucceed(cxt *gt.Context) {
	// 使用goquery库提取数据
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(cxt.RespBodyString()))
	if err != nil {
		log.Println(err)
		return
	}
	result, err := dom.Find("div[id=result] .well").Html()
	if err != nil {
		log.Println(err)
	}
	// 固定顺序map
	gd := gt.NewGDMap().Add("start_ip", cxt.Task.GetDataStr("start_ip"))
	gd.Add("end_ip", cxt.Task.GetDataStr("end_ip"))
	gd.Add("number", cxt.Task.GetDataStr("number")).Add("result", result)
	log.Println("采集到的数据: ", []string{cxt.Task.GetDataStr("start_ip"), cxt.Task.GetDataStr("end_ip"),
		cxt.Task.GetDataStr("number"), result})
	// 保存抓取数据
	//err = db.InsertAtGd("ip_result", gd)

	if err != nil {
		panic(err)
	}
	err = c.Add([]string{cxt.Task.GetDataStr("start_ip"), cxt.Task.GetDataStr("end_ip"),
		cxt.Task.GetDataStr("number"), result})
	log.Println("写入csv失败: ", err)
	if err != nil {
		panic(err)
	}

}

// GetIPRetry 获取详情信息重试的处理
func GetIPRetry(c *gt.Context) {
	//更换代理
	//c.SetProxy(uri)

	// or
	c.Client = &http.Client{
		// 设置代理
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(uri),
		//},
		Timeout: 5 * time.Second,
	}

	log.Println("休息1s")
	time.Sleep(1 * time.Second)
}

// GetIPFailed 获取详情信息失败执行返还给队列
func GetIPFailed(c *gt.Context) {
	queue.Add(c.Task) //请求失败归还到队列
}

// 第二种实现
func main2() {
	// 连接数据库
	var (
		host     = "192.168.0.192"
		port     = 3306
		user     = "root"
		password = "root123"
		dbname   = "test"
	)
	db, err := gt.NewMysql(host, port, user, password, dbname)
	if err != nil {
		panic(err)
	}
	// 请求数据
	ctx, err := gt.Get("http://ip.bczs.net/country/CN")
	if err != nil {
		log.Println(err)
		return
	}
	// 提取数据并保存
	for _, tbody := range gt.RegHtmlTbody(ctx.Html) {
		for _, tr := range gt.RegHtmlTr(tbody) {
			td := gt.RegHtmlTdTxt(tr)
			if len(td) < 3 {
				gt.Error("异常数据 ： ", td)
				continue
			}
			// 保存抓取数据
			err := db.InsertAt("ip", map[string]interface{}{
				"start": td[0],
				"end":   td[1],
				"num":   td[2],
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

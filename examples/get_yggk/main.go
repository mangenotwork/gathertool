/*
	抓取招生章程
	一级页面 学校列表
	https://gaokao.chsi.com.cn/zsgs/zhangcheng/listVerifedZszc--method-index,lb-1,start-0.dhtml
	二级页 每年的
	https://gaokao.chsi.com.cn/zsgs/zhangcheng/listZszc--schId-5.dhtml
	三级页 招生章程内容
    https://gaokao.chsi.com.cn/zsgs/zhangcheng/listVerifedZszc--infoId-2708104715,method-view,schId-5.dhtml
*/

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/PuerkitoBio/goquery"
	gt "github.com/mangenotwork/gathertool"
)

var (
	PG1Queue = gt.NewQueue() // 一级页面抓取任务队列
	PG2Queue = gt.NewQueue() // 二级页面抓取任务队列
	PG3Queue = gt.NewQueue() // 三级页面抓取任务队列
	host = "192.168.0.197"
	port = 3306
	user = "root"
	password = "root123"
	database = "spider"
	tableName = "yggk_zszc_1"
)

func init(){
	gt.NewMysqlDB(host, port, user, password, database)
	err := gt.MysqlDB.Conn()
	if err != nil{
		log.Panic("数据库初始化失败")
	}
	//初始化表
	gt.MysqlDB.NewTable(tableName, map[string]string{
		"school_name": "varchar(30)",
		"zszcName" : "varchar(50)",
		"fbTime" : "varchar(30)",
		"url" : "varchar(200)",
		"content": "text",
	})
}

// 全量抓取
func main(){
	// 一共 29页
	for p:=0; p<30;p++{
		pg1 := "https://gaokao.chsi.com.cn/zsgs/zhangcheng/listVerifedZszc--method-index,lb-1,start-%d.dhtml"
		PG1Queue.Add(&gt.Task{
			Url: fmt.Sprintf(pg1, p*100),
		})
	}
	//并发执行一级页面
	gt.StartJobGet(50,PG1Queue,
		gt.SucceedFunc(Pg1Succeed),//请求成功后执行的方法
		gt.RetryFunc(Retry),
		gt.FailedFunc(Fail),
	)
	log.Println(" 并发执行一级页面 完成")
	//并发执行二级页面
	gt.StartJobGet(50,PG2Queue,
		gt.SucceedFunc(Pg2Succeed),//请求成功后执行的方法
		gt.RetryFunc(Retry),
		gt.FailedFunc(Fail),
	)
	log.Println(" 并发执行二级页面 完成")
	//并发执行三级页面
	gt.StartJobGet(100,PG3Queue,
		gt.SucceedFunc(Pg3Succeed),//请求成功后执行的方法
		gt.RetryFunc(Retry),
		gt.FailedFunc(Fail),
	)
	log.Println(" 并发执行三级页面 完成")


	////测试第二级页面
	//c,_ := gt.Get("https://gaokao.chsi.com.cn/zsgs/zhangcheng/listZszc--schId-1.dhtml",
	//	gt.SucceedFunc(Pg2Succeed))
	//c.Do()

	//测试第三级页面
	//for i:=0; i<100;i++{
	//	go func(){
	//		c,_ := gt.Get("https://gaokao.chsi.com.cn/zsgs/zhangcheng/listVerifedZszc--infoId-2697675279,method-view,schId-1.dhtml",
	//			gt.SucceedFunc(Pg3Succeed))
	//		c.Do()
	//	}()
	//}
	//time.Sleep(10*time.Second)
}

// 抓取第一级页面成功后
func Pg1Succeed(ctx *gt.Context){
	html := string(ctx.RespBody)
	dom,err := gt.NewGoquery(html)
	if err != nil{
		log.Println(err)
		return
	}
	result := dom.Find("table tbody")
	if len(result.Nodes) < 2{
		log.Println("没有找到table")
	}
	result.Eq(1).Each(func(i int, tr *goquery.Selection){
		tr.Find("td").Each(func(i int, td *goquery.Selection){
			schoolName := td.Text()
			href,_ := td.Find("a").Attr("href")
			log.Println(schoolName, href)
			// 加入二级页面队列
			PG2Queue.Add(&gt.Task{
				Url: "https://gaokao.chsi.com.cn/" + href,
				Data: map[string]interface{}{
					"school_name":schoolName,
				},
			})
		})
	})
}

// 抓取第二级页面成功后
func Pg2Succeed(ctx *gt.Context){
	html := string(ctx.RespBody)
	dom,err := gt.NewGoquery(html)
	if err != nil{
		log.Println(err)
		return
	}
	result := dom.Find(".zszcdel table tbody")
	//log.Println(result.Html())
	result.Find("tr").Each(func(i int, tr *goquery.Selection){
		td := tr.Find("td")
		zszcName := td.Eq(0).Text()
		href,_ := td.Eq(0).Find("a").Attr("href")
		fbTime := td.Eq(1).Text()
		log.Println(zszcName, href, fbTime)
		ctx.Task.Data["zszcName"] = zszcName
		ctx.Task.Data["fbTime"] = fbTime
		ctx.Task.Data["url"] = "https://gaokao.chsi.com.cn/" + href
		ctx.Task.Url = "https://gaokao.chsi.com.cn/" + href
		PG3Queue.Add(ctx.Task)
	})
}

//
func Pg3Succeed(ctx *gt.Context){
	html := string(ctx.RespBody)
	//log.Println(html)
	dom,err := gt.NewGoquery(html)
	if err != nil{
		log.Println(err)
		return
	}
	content,err := dom.Find(".content").Html()
	log.Println(content, err)
	if err != nil || content == ""{
		log.Println("还给队列")
		PG3Queue.Add(ctx.Task)
	}
	// 写入数据库
	schoolName := gt.StringValue(ctx.Task.Data["school_name"]) // 转换成字符串
	zszcName := gt.StringValue(ctx.Task.Data["zszcName"])
	fbTime := gt.StringValue(ctx.Task.Data["fbTime"])
	err = gt.MysqlDB.Insert(tableName, map[string]interface{}{
		"school_name": gt.CleaningStr(schoolName), // 清理字符串前后空格和换行符等
		"zszcName" : gt.CleaningStr(zszcName),
		"fbTime" : gt.CleaningStr(fbTime),
		"url" : ctx.Task.Data["url"],
		"content": content,
	})
	log.Println(err)
}

func Retry(*gt.Context){
	time.Sleep(2*time.Second)
}
func Fail(ctx *gt.Context){
	log.Println(ctx.Err)
}
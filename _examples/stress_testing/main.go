package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
	"net/http"
)

func main() {
	gt.CPUMax()

	// 普通 GET api压测
	url := "http://192.168.0.9:8002/v1/health"
	// 请求10000次 并发数 1000
	test := gt.NewTestUrl(url,"Get",10000,1000)
	test.Run()
	test.Run(gt.SucceedFunc(func(ctx *gt.Context){
		log.Println(ctx.JobNumber, "测试完成!!", ctx.Ms)
	}))

	//// 设置 GET Header 的压测
	////url2 := "http://192.168.0.9:18084/static_service/v1/auth/video/page"
	//url2 := "http://192.168.0.9:18084/static_service/v1/auth/quality_article/list"
	//token := &http.Header{}
	//token.Add("token", tokenStr)
	//token.Add("source", "2")
	//test2 := gt.NewTestUrl(url2,"Get",10000,1000)
	//test2.Run(token)


	//url3 := "https://baidu.com"
	//test3 := gt.NewTestUrl(url3,"Get",100000,10000)
	//test3.Run()

	//gt.NewTestUrl("http://192.168.0.8:18090","Get",100000,10000).Run()
	//test3.Run()

	//// Post
	//url4 := "http://ggzyjy.sc.gov.cn/WebBuilder/frontAppAction.action?cmd=addPageView"
	//test4 := gt.NewTestUrl(url4,"Post",100,10)
	//test4.SetJson(`{
	//			"viewGuid":"cms_002",
	//			"siteGuid":"7eb5f7f1-9041-43ad-8e13-8fcb82ea831a"
	//			}`)
	//test4.Run(gt.SucceedFunc(func(c *gt.Context) {
	//	log.Println(string(c.RespBody))
	//	log.Println(c.Resp.Cookies())
	//}))

	// 含步骤压力测试


	//// 牛票票压力测试
	//nppurl1 := "http://192.168.0.9:8025/v2/index/recommend?index=1&limit=20&uid=8722"
	//token := &http.Header{}
	//token.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjM4OTA1ODQsImlzcyI6Im5pdS5jbiIsIm5iZiI6MTYyMDI5MDU4NCwiSWQiOjg3MjJ9.ij45vZILk9Kr35XoiPyIhjVGmCoERKiEBk6zz6P9P0g")
	//npptest1 := gt.NewTestUrl(nppurl1,"Get",30000,3000)
	//npptest1.Run(token)

	//点赞
	//npplike := "http://192.168.0.9:8025/v2/like/add"
	//npplike()

	//caseurl := "https://www.jy863.com:18443/"
	//caseurl := "https://www.uks678.com:10442/api/site/gdnotice/list"
	//npptest1 := gt.NewTestUrl(caseurl,"Get",10000,1000)
	//hd := &http.Header{}
	//hd.Add("requested-device", "APP")
	//hd.Add("requested-language", "CN")
	//hd.Add("requested-site", "www.uks678.com:10442")
	//npptest1.Run(hd)

	// 发帖
	//nppTopicCreate()

	//评论
	//nppComment()

}

var tokenStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pdF9vcmRlcl90eXBlIjowLCJhdmF0YXIiOiJodHRwczovL3AyLnltenkuY24vYXZlLzIwMjEwNDIyLzZlOTRlMzNiMmNmZWYzYzEuanBnIiwiY2FsX3JhbmsiOjAsImNvbnZlcnRfc2NvcmUiOjAsImV4cCI6MTYxOTY2OTM0MCwiZmlyc3Rfc3ViamVjdCI6MCwiZ2tfeWVhciI6MCwiZ3JhZGUiOjAsImdyb3VwX2V4cF90aW1lIjowLCJpYXQiOjE2MTkwNjQ1MzAsImlzX3ZpcCI6ZmFsc2UsIm1vYmlsZSI6IiIsIm1vYmlsZV9lbmNyeXB0aW9uIjoiIiwibW9kZSI6MCwibmJmIjoxNjE5MDY0NTMwLCJuaWNrbmFtZSI6IuiKseiKseiKseS4gyIsInByb3ZfaWQiOjAsInByb3ZfbW9kZWwiOjEsInByb3ZfbW9kZWxfZXgiOjEsInByb3ZfbmFtZSI6IiIsInJhbmsiOjAsInNjb3JlIjowLCJzY29yZV9kZWYiOjAsInNjb3JlcyI6MCwic3ViamVjdF9pZCI6MCwic3ViamVjdF9uYW1lIjoiIiwidXNlcl9ncm91cF9pZCI6NCwidXNlcl9pZCI6MzA0NzEyfQ.PpSvN0nrTDEOgK5bL0fEEzWDyx7KKXvUIr7PwoBuFdQ"


var (
	nppLikeTestQueue = gt.NewQueue()
	host = "192.168.0.197"
	port = 3306
	user = "root"
	password = "root123"
	database = "niu_pp"
	tableName = "tbl_user_ext"
)

//npp 点赞
func npplike(){
	npplike := "http://192.168.0.9:8025/v2/like/add"
	gt.NewMysqlDB(host, port, user, password, database)
	gt.MysqlDB.CloseLog()
	conn,err := gt.GetMysqlDBConn()
	if err != nil {
		log.Panic("数据库初始化失败")
	}
	allUid, _ := conn.Select("select uid from tbl_user_ext")
	allTid, _ := conn.Select("select tid from tbl_topic")
	for _,v := range allUid{
		token := GetToken(v["uid"])
		h := &http.Header{}
		h.Add("token", token)
		// 所有帖子
		for _, t  := range allTid {
			addnpplike := npplike + "?tid="+t["tid"]+"&uid="+v["uid"]

			gt.Get(addnpplike, h)

			//log.Println(addnpplike)
			//nppLikeTestQueue.Add(gt.CrawlerTask(addnpplike,"", h))
		}
	}
	//gt.StartJobGet(100, nppLikeTestQueue, gt.SucceedFunc(func(c *gt.Context) {
	//	log.Println(string(c.RespBody))
	//}))
}

// npp 获取token - 通用
func GetToken(uid string) string {
	link := "http://192.168.0.9:8025/login?uid="+uid
	token := ""
	c,err := gt.Get(link)
	if err != nil{
		log.Println(err)
	}
	c.Do()
	dataMap := gt.Json2Map(string(c.RespBody))
	body := gt.Any2Map(dataMap["body"])
	data := gt.Any2Map(body["data"])
	//log.Println(data["token"])
	token = gt.Any2String(data["token"])
	return token
}

// 发帖
func nppTopicCreate() {
	caseUrl := `http://192.168.0.9:8025/v2/topic/create`

	gt.NewMysqlDB(host, port, user, password, database)
	gt.MysqlDB.CloseLog()
	conn,err := gt.GetMysqlDBConn()
	if err != nil {
		log.Panic("数据库初始化失败")
	}
	allUid, _ := conn.Select("select uid from tbl_user_ext limit 100,10000")
	for _,v := range allUid{
		token := GetToken(v["uid"])
		h := &http.Header{}
		h.Add("token", token)
		data := `
{
	"uid": `+v["uid"]+`,
	"qi": 2021315,
	"fid": 1,
	"source": 1,
	"content": "我爱牛票票"
}`
		ctx,err := gt.PostJson(caseUrl, data, h)
		log.Println(ctx.RespBodyString(), err)
	}

}

// 评论
func nppComment() {
	caseUrl := "http://192.168.0.9:8025/v1/comment/add"
	gt.NewMysqlDB(host, port, user, password, database)
	gt.MysqlDB.CloseLog()
	conn,err := gt.GetMysqlDBConn()
	if err != nil {
		log.Panic("数据库初始化失败")
	}
	allUid, _ := conn.Select("select uid from tbl_user_ext")
	allTid, _ := conn.Select("select tid from tbl_topic")
	i:=0
	for _,v := range allUid{
		for _, t := range allTid {
			token := GetToken(v["uid"])
			h := &http.Header{}
			h.Add("token", token)
			i++
			if i%3 == 0 {
				continue
			}
			data := `
{
	"content":"我们都爱牛票票",
	"tid": `+t["tid"]+`,
	"cmtid":0,
	"uid": `+v["uid"]+`,
	"source":1
}`
			gt.PostJson(caseUrl, data, h)
		}

	}
}
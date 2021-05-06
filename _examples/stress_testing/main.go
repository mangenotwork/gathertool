package main

import (
	gt "github.com/mangenotwork/gathertool"
	"net/http"
)

func main() {
	// 普通 GET api压测
	//url := "http://192.168.0.9:18084/static_service/v1/allow/school/page"
	//test := gt.NewTestUrl(url,"Get",1000,500)
	//test.Run()
	//test.Run(gt.SucceedFunc(func(ctx *gt.Context){
	//	log.Println(ctx.JobNumber, "测试完成!!", ctx.Ms)
	//}))

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


	// 牛票票压力测试
	nppurl1 := "http://192.168.0.9:8025/v2/index/list"
	token := &http.Header{}
	token.Add("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjM4OTA1ODQsImlzcyI6Im5pdS5jbiIsIm5iZiI6MTYyMDI5MDU4NCwiSWQiOjg3MjJ9.ij45vZILk9Kr35XoiPyIhjVGmCoERKiEBk6zz6P9P0g")
	npptest1 := gt.NewTestUrl(nppurl1,"Get",30000,1000)
	npptest1.Run(token)


	//run2()

	//ZyzMajor2021Test()

}

var tokenStr = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pdF9vcmRlcl90eXBlIjowLCJhdmF0YXIiOiJodHRwczovL3AyLnltenkuY24vYXZlLzIwMjEwNDIyLzZlOTRlMzNiMmNmZWYzYzEuanBnIiwiY2FsX3JhbmsiOjAsImNvbnZlcnRfc2NvcmUiOjAsImV4cCI6MTYxOTY2OTM0MCwiZmlyc3Rfc3ViamVjdCI6MCwiZ2tfeWVhciI6MCwiZ3JhZGUiOjAsImdyb3VwX2V4cF90aW1lIjowLCJpYXQiOjE2MTkwNjQ1MzAsImlzX3ZpcCI6ZmFsc2UsIm1vYmlsZSI6IiIsIm1vYmlsZV9lbmNyeXB0aW9uIjoiIiwibW9kZSI6MCwibmJmIjoxNjE5MDY0NTMwLCJuaWNrbmFtZSI6IuiKseiKseiKseS4gyIsInByb3ZfaWQiOjAsInByb3ZfbW9kZWwiOjEsInByb3ZfbW9kZWxfZXgiOjEsInByb3ZfbmFtZSI6IiIsInJhbmsiOjAsInNjb3JlIjowLCJzY29yZV9kZWYiOjAsInNjb3JlcyI6MCwic3ViamVjdF9pZCI6MCwic3ViamVjdF9uYW1lIjoiIiwidXNlcl9ncm91cF9pZCI6NCwidXNlcl9pZCI6MzA0NzEyfQ.PpSvN0nrTDEOgK5bL0fEEzWDyx7KKXvUIr7PwoBuFdQ"



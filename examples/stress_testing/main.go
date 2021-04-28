package main

import (
	gt "github.com/mangenotwork/gathertool"
)

func main() {
	// 普通 GET api压测
	//url := "http://192.168.0.9:18084/static_service/v1/allow/school/page"
	//test := gt.NewTestUrl(url,"Get",1000,500)
	//test.Run()
	//test.Run(gt.SucceedFunc(func(ctx *gt.Context){
	//	log.Println(ctx.JobNumber, "测试完成!!", ctx.Ms)
	//}))

	// 设置 GET Header 的压测
	//url2 := "http://192.168.0.9:18084/static_service/v1/auth/video/page"
	//test2 := gt.NewTestUrl(url2,"Get",100000,10000)
	//test2.Run()


	url3 := "https://baidu.com"
	test3 := gt.NewTestUrl(url3,"Get",100000,10000)
	test3.Run()
}
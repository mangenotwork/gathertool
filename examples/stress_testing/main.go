package main

import (
	gt "github.com/mangenotwork/gathertool"
)

func main() {
	// 单接口测试 timeOut 是ms单位
	test := gt.NewTestUrl("http://192.168.0.9:9001","Get",100000,2000)
	test.Run()
	//test.Run(gt.SucceedFunc(func(ctx *gt.Context){
	//	log.Println(ctx.JobNumber, "测试完成!!", ctx.Ms)
	//}))

	// 含步骤的测试

}
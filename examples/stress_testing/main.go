package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
)

func main() {
	test := gt.NewTestUrl("http://192.168.0.9:9001","Get",10000,1000)
	test.Run(gt.SucceedFunc(func(ctx *gt.Context){
		log.Println(ctx.JobNumber, "测试完成!!")
	}))
}
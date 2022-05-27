package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
)

func main(){
	ipt := &gt.Intercept{
		Ip : "0.0.0.0:8111",
		HttpPackageFunc: func(pack *gt.HttpPackage){
			// 查看 ContentType
			log.Println("ContentType = ", pack.ContentType)
			// 获取数据包数据为 http,json等 TXT格式的数据
			log.Println("Txt = ", pack.Html())
			// 获取数据包为图片，将图片转为 base64
			log.Println("img base64 = ", pack.Img2Base64())
			// 获取数据包为图片，存储图片
			log.Println(pack.SaveImage(""))
		},
	}
	// 启动服务
	ipt.RunServer()
}

/*
golang实现http&https代理服务器
golang实现http代理服务器并解析包
 */
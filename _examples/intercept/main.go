package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
)

func main(){
	ipt := &gt.Intercept{
		Ip : "0.0.0.0:8111",
		HttpPackageFunc: func(pack *gt.HttpPackage){
			log.Println("ContentType = ", pack.ContentType)
			log.Println("Txt = ", pack.Html())
			log.Println("img base64 = ", pack.Img2Base64())
			log.Println(pack.SaveImage(""))
		},
	}
	ipt.RunServer()
}

/*
	环球医药网 爬虫
	http://data.qgyyzs.net

	药智数据
	https://db.yaozh.com/instruct?p=1&pageSize=20

 */
package main

import (
	"github.com/PuerkitoBio/goquery"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"net/http"
	"strings"
)

func main(){
	//GcypList()

	//HospitalList()

	Instruct()
}

// 抓取国产药品列表
func GcypList(){
	caseUrl := "http://data.qgyyzs.net/gcyp_1.aspx"
	ctx, _ := gt.Get(caseUrl)
	log.Println(ctx.RespBodyString())
}

// 抓取医院列表
func HospitalList(){
	caseUrl := "http://data.qgyyzs.net/hospital_1.aspx"
	head := http.Header{}
	ctx, _ := gt.Get(caseUrl, head)
	htmlStr := gt.ConvertByte2String(ctx.RespBody, gt.GBK)
	//log.Println(htmlStr)
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil{
		log.Println(err)
		return
	}
	result := dom.Find("div[class=content]")
	resultStr, _ := result.Html()
	for _, v := range gt.RegHtmlLi(resultStr) {
		log.Println(v)
	}
}

// 抓取药品说明书列表
func Instruct(){
	caseUrl := "https://db.yaozh.com/instruct?p=1&pageSize=20"
	ctx, _ := gt.Get(caseUrl)
	log.Println(ctx.RespBodyString())
}
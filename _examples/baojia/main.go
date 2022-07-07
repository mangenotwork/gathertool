/*
	url : http://www.100ppi.com/mprice/mlist-1--1.html
	抓取商品报价

 */
package main

import (
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
)

// 数据库对象
var (
	host192  = "192.168.0.192"
	port     = 3306
	user     = "root"
	password = "root123"
	dbName      = "test"
	DB, _ = gt.NewMysql(host192, port, user, password, dbName)
	OutTable = "sp_jiage_20220707"
	Host = "http://www.100ppi.com/mprice/"
)

func main(){
	caseUrl := "http://www.100ppi.com/mprice/mlist-1--%d.html"
	// 50页
	for i:=1; i< 51; i++ {
		// 请求
		ctx, _ := gt.Get(fmt.Sprintf(caseUrl, i))
		// 数据提取
		table := gt.RegHtmlTable(ctx.Html)
		if len(table) > 0 {
			trList := gt.RegHtmlTr(table[0])
			for _, tr := range trList {
				log.Println("tr = ", tr)
				tdList := gt.RegHtmlTd(tr)
				if len(tdList) == 0 {
					continue
				}
				spNameHtml := tdList[0]
				spName := gt.RegHtmlATxt(spNameHtml)[0]
				spUrl := Host + gt.RegHtmlHrefTxt(spNameHtml)[0]
				spNoteHtml := tdList[1]
				spNote := gt.CleaningStr(gt.RegHtmlDivTxt(spNoteHtml)[0])
				pingpaiHtml := tdList[2]
				pingpai := gt.RegHtmlTdTxt(pingpaiHtml)[0]
				baojiaHtml := tdList[3]
				baojia := gt.CleaningStr(gt.RegHtmlTdTxt(baojiaHtml)[0])
				baojiaTypeHtml := tdList[4]
				baojiaType := gt.CleaningStr(gt.RegHtmlTdTxt(baojiaTypeHtml)[0])
				diziHtml := tdList[5]
				dizi := gt.CleaningStr(gt.RegHtmlTdTxt(diziHtml)[0])
				shangjiaHtml := tdList[6]
				shangjia := ""
				shangjiaA :=  gt.RegHtmlATxt(shangjiaHtml)
				if len(shangjiaA) > 0 {
					shangjia = shangjiaA[0]
				}else{
					shangjia = gt.RegHtmlDivTxt(shangjiaHtml)[0]
				}
				shangjiaUrl := ""
				shangjiaUrlList := gt.RegHtmlHrefTxt(spNameHtml)
				if len(shangjiaUrlList) > 0 {
					shangjiaUrl = Host + shangjiaUrlList[0]
				}
				dateHtml := tdList[7]
				date := gt.RegHtmlTdTxt(dateHtml)[0]
				log.Println("商品名称 = ", spName)
				log.Println("商品详情地址 = ", spUrl)
				log.Println("规格 = ", spNote)
				log.Println("品牌/产地 = ", pingpai)
				log.Println("报价 = ", baojia)
				log.Println("报价类型 = ", baojiaType)
				log.Println("交货地 = ", dizi)
				log.Println("交易商 = ", shangjia)
				log.Println("交易商详情地址 = ", shangjiaUrl)
				log.Println("发布时间 = ", date)
				// 固定字段顺序Map, 写入mysql数据库
				gd := gt.NewGDMap().Add("sp_name", spName).Add("sp_url", spUrl).Add("sp_note", spNote)
				gd.Add("pingpai", pingpai).Add("baojia", baojia).Add("baojia_type", baojiaType)
				gd.Add("dizi", dizi).Add("shangjia", shangjia).Add("shangjia_url", shangjiaUrl)
				gd.Add("date", date).Add("req_md5", ctx.CheckMd5())
				err := DB.InsertAtGd(OutTable, gd)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
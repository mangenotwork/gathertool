/*
url : http://www.100ppi.com/mprice/mlist-1--1.html
抓取商品报价
*/
package main

import (
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"strings"
	"time"
)

// 数据库对象
var (
	host192  = "127.0.0.1"
	port     = 3306
	user     = "root"
	password = "123456"
	dbName   = "sp_data"
	DB, _    = gt.NewMysql(host192, port, user, password, dbName)
	Host     = "http://www.100ppi.com/mprice/"
)

func main() {
	//GetMPrice()

	//RunGetMList()

	RunGetDetail()
}

// GetMPrice 抓取商品分类数据处理后存入数据库
func GetMPrice() {
	OutTable := "table_fenlei"
	caseUrl := "https://www.100ppi.com/mprice/plist-1-1-1.html"
	ctx, _ := gt.Get(caseUrl)
	data, _ := gt.GetPointClassHTML(ctx.Html, "div", "p_list2")
	for _, v := range data {
		aDataList := gt.RegHtmlA(v)
		for _, aData := range aDataList {
			gt.Info(aData)
			href := "https://www.100ppi.com/mprice/" + gt.RegHtmlHrefTxt(aData)[0]
			fenlei := gt.RegHtmlATxt(aData)[0]
			fenlei = gt.RegDelHtml(fenlei)
			if strings.Index(href, "plist") == -1 {
				continue
			}
			gt.Info("href = ", href)
			gt.Info("fenlei = ", fenlei)
			gt.Info("数据提取完成，准备写入数据库")
			gd := gt.NewGDMap().Add("fenlei", fenlei).Add("href", href)
			err := DB.InsertAtGd(OutTable, gd)
			if err != nil {
				panic(err)
			}
		}
	}
}

// RunGetMList 从数据库中获取分类链接进行抓取商品列表
func RunGetMList() {
	hrefList, err := DB.Select("select * from table_fenlei")
	if err != nil {
		log.Panic(err)
	}
	for _, hrefData := range hrefList {
		href := hrefData["href"]
		href = href[0:len(href)-6] + "%d.html"
		gt.Info(href)
		GetMList(href)
	}
}

// GetMList 抓取商品列表数据并处理后存储到数据库
func GetMList(caseUrl string) {
	OutTable := "m_list"
	// 50页
	for i := 1; i < 11; i++ {
		time.Sleep(200 * time.Millisecond)
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
				shangjiaA := gt.RegHtmlATxt(shangjiaHtml)
				if len(shangjiaA) > 0 {
					shangjia = shangjiaA[0]
				} else {
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

// RunGetDetail 从数据库中获取商品详细链接进行抓取商品详情和商家信息
func RunGetDetail() {
	caseList, err := DB.Select("select * from m_list")
	if err != nil {
		log.Panic(err)
	}
	for _, caseData := range caseList {
		GetDetail(caseData["sp_url"])
	}
}

// GetDetail 抓取商品详细数据并处理后存储到数据库
func GetDetail(caseUrl string) {
	//caseUrl := "http://www.100ppi.com/mprice/detail-4264843.html"
	time.Sleep(200 * time.Millisecond)
	ctx, _ := gt.Get(caseUrl)
	//gt.Info(ctx.Html)

	// 提取商品详情信息
	infoData, _ := gt.GetPointClassHTML(ctx.Html, "table", "st2-table tac")
	//gt.Info("infoData = ", infoData)

	// 提取商品描述
	noteData, _ := gt.GetPointClassHTML(ctx.Html, "table", "mb20 st2-table tac")
	//gt.Info("noteData = ", noteData[0])

	// 提取报价走势图
	imgData, _ := gt.GetPointClassHTML(ctx.Html, "div", "vain")
	//gt.Info("imgData = ", imgData)

	if len(noteData) < 2 {
		return
	}

	// 商家信息
	shangJiaData := noteData[1]
	//gt.Info("shangJiaData = ", shangJiaData)

	// 存储商品详细信息
	gd1 := gt.NewGDMap().Add("info", infoData[0]).Add("note", noteData[0]).Add("img", imgData[0])
	gd1.Add("shangjia", shangJiaData)
	err := DB.InsertAtGd("sp_detail", gd1)
	if err != nil {
		panic(err)
	}

	// 将商家信息单独提取出来存储数据
	trList := gt.RegHtmlTrTxt(shangJiaData)

	gs := GetSJData(trList[1])
	gt.Info(gs)
	fzr := GetSJData(trList[2])
	gt.Info(fzr)
	lx1 := GetSJData(trList[3])
	gt.Info(lx1)
	lx2 := GetSJData(trList[4])
	gt.Info(lx2)
	gd2 := gt.NewGDMap().Add("gs", gs).Add("fzr", fzr).Add("lx1", lx1).Add("lx2", lx2)
	err = DB.InsertAtGd("sj_info", gd2)
	if err != nil {
		panic(err)
	}

}

func GetSJData(data string) string {
	data = gt.RegDelHtml(data)
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)
	data = strings.Replace(data, "\t", "", -1)
	data = strings.Replace(data, " ", "", -1)
	data = strings.Replace(data, "进入报价通", "", -1)
	data = strings.Replace(data, "联系报价人", "", -1)
	data = strings.Replace(data, "扫码分享", "", -1)
	return data
}

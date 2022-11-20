/*
url: http://futures.100ppi.com/qhb/day-2022-07-05.html
抓取期货数据
*/
package main

import (
	"fmt"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"strings"
)

// 数据库对象
var (
	host192  = "192.168.0.192"
	port     = 3306
	user     = "root"
	password = "root123"
	dbName   = "test"
	DB, _    = gt.NewMysql(host192, port, user, password, dbName)
	OutTable = "qihuo"
)

func main() {
	date := "2022-07-05"
	caseUrl := "http://futures.100ppi.com/qhb/day-%s.html"
	ctx, _ := gt.Get(fmt.Sprintf(caseUrl, date))
	//log.Println(ctx.Html)
	// 数据提取
	datas, _ := gt.GetPointHTML(ctx.Html, "div", "id", "domestic")
	Data(datas, date, "内期表", "备注：内期表=国内期货主力合约表")
	datas, _ = gt.GetPointHTML(ctx.Html, "div", "id", "overseas")
	Data(datas, date, "外期表", "备注：外期表=国外期货主力合约表")
}

func Data(datas []string, date, typeName, note string) {
	for _, data := range datas {
		table, _ := gt.GetPointHTML(data, "table", "id", "fdata")
		if len(table) > 0 {
			trList := gt.RegHtmlTr(table[0])
			jys := ""
			for _, tr := range trList {
				td := gt.RegHtmlTd(tr)
				log.Println("td = ", td, len(td))
				if len(td) == 1 {
					jys = gt.RegHtmlTdTxt(td[0])[0]
					continue
				}
				name := gt.RegHtmlTdTxt(td[0])[0]
				if strings.Index(name, "商品名称") != -1 {
					continue
				}
				zlhy := gt.RegHtmlTdTxt(td[1])[0]
				jsj := gt.RegHtmlTdTxt(td[2])[0]
				zd := gt.RegDelHtml(gt.RegHtmlTdTxt(td[3])[0])
				cjj := gt.RegHtmlTdTxt(td[4])[0]
				ccl := gt.RegHtmlTdTxt(td[5])[0]
				dw := gt.RegHtmlTdTxt(td[6])[0]
				log.Println("日期 = ", date)
				log.Println("机构 = ", jys)
				log.Println("商品名称 = ", name)
				log.Println("主力合约 = ", zlhy)
				log.Println("结算价 = ", jsj)
				log.Println("涨跌 = ", zd)
				log.Println("成交量 = ", cjj)
				log.Println("持仓量 = ", ccl)
				log.Println("单位 = ", dw)
				gd := gt.NewGDMap().Add("date", date).Add("type_name", typeName).Add("note", note)
				gd.Add("jys", jys).Add("name", name).Add("zlhy", zlhy).Add("jsj", jsj)
				gd.Add("zd", zd).Add("cjj", cjj).Add("ccl", ccl).Add("dw", dw)
				err := DB.InsertAtGd(OutTable, gd)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

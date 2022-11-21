package main

import (
	gt "github.com/mangenotwork/gathertool"
)

// 行业信息网 : http://www.cnlinfo.net/

func init() {
	// 读取配置
	err := gt.NewConf("conf.yaml")
	if err != nil {
		panic(err)
	}
}

// Rse 抓取结果保存的数据结构
type Rse struct {
	Broad  string   `json:"broad"`
	Fenlei []string `json:"fenlei"`
}

var (
	rse = make([]*Rse, 0)
)

// 获取行业分类
func main() {
	caseUrl := "http://www.cnlinfo.net/allgongsifenlei/yiqiyibiao.htm"
	ctx, _ := gt.Get(caseUrl)
	hangyeList, err := gt.GetPointClassHTML(ctx.Html, "ul", "hangye-list")
	if err != nil {
		gt.Error("没有获取到 ul class = hangye-list")
		return
	}
	hangyeHtml := ""
	if len(hangyeList) > 0 {
		hangyeHtml = hangyeList[0]
	}
	// 获取并便利分类大类
	for _, li := range gt.RegHtmlA(hangyeHtml) {
		broad := gt.RegHtmlATxt(li)[0]
		href := gt.RegHtmlHrefTxt(li)[0]
		// 抓取小类
		getClassify2(broad, href)
	}
	// 保存并输出到json文件
	err = gt.OutJsonFile(rse, gt.Config.GetStr("out_file"))
	if err != nil {
		gt.Error(err)
	}
}

// 获取小类
func getClassify2(broad, caseUrl string) {
	var (
		notData     = 0
		reqDelayMax = gt.Config.GetInt("req_delay_max")
		reqDelayMin = gt.Config.GetInt("req_delay_min")
	)
	d := &Rse{
		Broad:  broad,
		Fenlei: make([]string, 0),
	}
	gt.Infof("抓取大类 %v : %v", broad, caseUrl)
	ctx, _ := gt.Get(caseUrl, gt.SetSleep(reqDelayMin, reqDelayMax))
	hangyeContent, err := gt.GetPointClassHTML(ctx.Html, "ul", "hangye-fenlei-content")
	if err != nil {
		gt.Error("没有获取到 ul class = hangye-fenlei-content")
		notData++
	}
	for _, ul := range hangyeContent {
		for _, a := range gt.RegHtmlATxt(ul) {
			d.Fenlei = append(d.Fenlei, a)
		}
	}
	// 第二种样式
	hangyeContent1, err := gt.GetPointClassHTML(ctx.Html, "ul", "fenlei_list1")
	if err != nil {
		gt.Error("没有获取到 ul class = fenlei_list1")
		notData++
	}
	for _, ul := range hangyeContent1 {
		for _, a := range gt.RegHtmlATxt(ul) {
			d.Fenlei = append(d.Fenlei, a)
		}
	}
	rse = append(rse, d)
	if notData == 2 {
		panic("该页面没有数据")
	}
}

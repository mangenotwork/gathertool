package main

import (
	"github.com/PuerkitoBio/goquery"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"net/http"
	"strings"
)

func main(){


	/*
		category
					0      热门
					1760   头条
					99991  榜单
					10011  高笑
					7      社会
					12     时尚
					10018  电影
					10007  美女
					3      体育
					10005  动漫
	 */

	url := "https://weibo.com/a/aj/transform/loadingmoreunlogin?ajwvr=6&category=0&page=2&lefnav=0&cursor=&__rnd="+ gt.Timestamp()
	ctx, _ := gt.Get(url, gt.SucceedFunc(succed))
	// TODO 通过主domain得到
	ctx.Req.AddCookie(&http.Cookie{Name: "login_sid_t",Value: "fe4df185ba73522be106feef5bbff035", HttpOnly: true})
	ctx.Req.AddCookie(&http.Cookie{Name: "_s_tentry",Value: "passport.weibo.com", HttpOnly: true})
	ctx.Req.AddCookie(&http.Cookie{Name: "WBStorage",Value: "8daec78e6a891122|undefined", HttpOnly: true})
	ctx.Req.AddCookie(&http.Cookie{Name: "ULV",Value: "1620468432944:1:1:1:9339385764296.797.1620468432934:", HttpOnly: true})
	ctx.Req.AddCookie(&http.Cookie{Name: "SUBP",Value: "0033WrSXqPxfM72-Ws9jqgMF55529P9D9WWmEanFMEIiGoQ6G-cq2.Uz", HttpOnly: true})
	ctx.Req.AddCookie(&http.Cookie{Name: "SUB",Value: "_2AkMXyu34f8NxqwJRmPwczW7iaIlwywzEieKhlhwjJRMxHRl-yj9jqhUutRB6PErDF0hBnN-8Mfy2pdZU1yAQqOPZIudN", HttpOnly: true})
	ctx.Req.AddCookie(&http.Cookie{Name: "SINAGLOBAL",Value: "9339385764296.797.1620468432934", HttpOnly: true})
	ctx.Req.AddCookie(&http.Cookie{Name: "Apache",Value: "9339385764296.797.1620468432934", HttpOnly: true})
	ctx.Do()
}

func succed(ctx *gt.Context) {

	//html := gt.ConvertByte2String(ctx.RespBody, gt.GB2312)
	htmlBody,err := gt.UnescapeUnicode(ctx.RespBody)
	html := string(htmlBody)
	html = strings.Replace(html,"\\r","", -1)
	html = strings.Replace(html,"\\n","", -1)
	html = strings.Replace(html,"\\","", -1)

	dom,err := gt.NewGoquery(html)
	if err != nil{
		log.Println(err)
		return
	}

	dom.Find("div[action-type=feed_list_item]").Each(func(i int, div *goquery.Selection){
		divHtml,err := div.Html()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("\n\n\n\n *******************************  \n ")
		dataList := gt.RegHtmlA(divHtml)
		for _, v := range dataList{

			if strings.Contains(v,`class="subinfo_face`) {
				log.Println("[头像] : ", v)
				scr := gt.RegHtmlSrc(v)
				log.Println(scr)
			}

			if strings.Contains(v,`class="subinfo S_txt2`) {
				log.Println("[昵称] : ", v)
				href := gt.RegHtmlHref(v)
				log.Println(href)
				name := gt.RegHtmlSpan(v)
				log.Println(name)

			}
		}
		log.Println( "\n===================================\n\n")
	})
}


// cookie 持久化
func case2(){
	////https://passport.weibo.com/visitor/visitor?a=restore&cb=restore_back&from=weibo&_rand=0.5904772874565642
	//u1 := "https://weibo.com"
	//ctx, _ := gt.Get(u1)
	//ctx.Do()
	//log.Println(ctx.Resp.Cookies())
}
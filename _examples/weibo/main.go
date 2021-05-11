package main

import (
	"github.com/PuerkitoBio/goquery"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"net/http"
	"strings"
)

func main(){

	//case1()

	case2()

}

func case1(){


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
	//ctx.Req.AddCookie(&http.Cookie{Name: "login_sid_t",Value: "fe4df185ba73522be106feef5bbff035", HttpOnly: true})
	//ctx.Req.AddCookie(&http.Cookie{Name: "_s_tentry",Value: "passport.weibo.com", HttpOnly: true})
	//ctx.Req.AddCookie(&http.Cookie{Name: "WBStorage",Value: "8daec78e6a891122|undefined", HttpOnly: true})
	//ctx.Req.AddCookie(&http.Cookie{Name: "ULV",Value: "1620468432944:1:1:1:9339385764296.797.1620468432934:", HttpOnly: true})
	ctx.Req.AddCookie(&http.Cookie{Name: "SUBP",Value: "0033WrSXqPxfM72-Ws9jqgMF55529P9D9WWENAjmKyIZz1AWjDi68mRw", HttpOnly: true})
	ctx.Req.AddCookie(&http.Cookie{Name: "SUB",Value: "_2AkMXxWiSf8NxqwFRmPoWz2nlbop1zwvEieKhmZlJJRMxHRl-yT9jqlAItRB6PEVGfTP09XmsX_7CR2H1OUv6b-f-1bJl", HttpOnly: true})
	//ctx.Req.AddCookie(&http.Cookie{Name: "SINAGLOBAL",Value: "9339385764296.797.1620468432934", HttpOnly: true})
	//ctx.Req.AddCookie(&http.Cookie{Name: "Apache",Value: "9339385764296.797.1620468432934", HttpOnly: true})
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
	// ==========   获取 tid
	//由  https://passport.weibo.com/js/visitor/mini_original.js?v=20161116  //main 执行入口 js代码 找到tid的方式
	// 也找到了  fpPostInterface: "visitor/genvisitor"
	// fpCollectInterface: "visitor/record"

	// cb 固定值 gen_callback
	//getTidUrl := "https://passport.weibo.com/visitor/genvisitor?cb=gen_callback&fp={\"os\":\"1\",\"browser\":\"Chrome70,0,3538,25\",\"fonts\":\"undefined\",\"screenInfo\":\"1920*1080*24\",\"plugins\":\"\"}"
	//ctx, _ := gt.Get(getTidUrl)
	//ctx.Do()
	//log.Println(string(ctx.RespBody))



	// ============   获取sub subp
	// F12观察到 https://passport.weibo.com/visitor/visitor?a=restore&cb=restore_back&from=weibo&_rand=0.8977533486179248
	// 获取的cookie
	// a 固定 incarnate
	// t 上一步得到的tid
	// w 上一步 new_tid = true 就是3  否则为2
	// cb 固定 cross_domain
	// from  固定  weibo

	getSubUrl := "https://passport.weibo.com/visitor/visitor?a=incarnate&t=hWJuIDjcJt14rjdJBmelhsQ0ReEl6ATZZnf2EQbrBQM=&w=3&c&cb=restore_back&from=weibo"
	ctx,_ := gt.Get(getSubUrl)
	ctx.Do()
	log.Println(ctx.Resp)
	log.Println(ctx.Resp.Cookies())
	log.Println(string(ctx.RespBody))

	//getSubUrl := "https://passport.weibo.com/visitor/visitor?a=restore&cb=restore_back&from=weibo&_rand=0.8977533486179248"
	//ctx,_ := gt.Get(getSubUrl)
	//ctx.Req.AddCookie(&http.Cookie{Name: "tid",Value: `HVyCvxI2gJHD6a6/hDOaWo3MNgwy9X8LwZFUG7o3vNA=__095`, HttpOnly: true})
	////ctx.Req.AddCookie(&http.Cookie{Name: "SRT",Value: "D.QqHBTrsR5QPiUcRtOeYoWr9NUPB3R39Qi-bYNdo35QWwMdbbN-YjTmntNbHi5mYNUCsuTZbgVdYC43MNAZSAMQHK549Q4qiaK4S1VFM6R4YbVP9GUqYYT3AqW-kmdA77%2AB.vAflW-P9Rc0lR-ykKDvnJqiQVbiRVPBtS%21r3J8sQVqbgVdWiMZ4siOzu4DbmKPWQU4PYU%21SiM4b9M-yMi%21VkR3mpIbPw", HttpOnly: true})
	////ctx.Req.AddCookie(&http.Cookie{Name: "SRF",Value: "1620470800", HttpOnly: true})
	//ctx.Do()
	//log.Println(ctx.Resp)
	//log.Println(ctx.Resp.Cookies())
	//log.Println(ctx.Resp.TransferEncoding)
	//log.Println(string(ctx.RespBody))
}
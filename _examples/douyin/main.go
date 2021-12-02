package main

import (
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	gt "github.com/mangenotwork/gathertool"
	"log"
	"net/url"
	"strings"
)

func main(){
	caseUrl := "https://www.douyin.com/search/美女"
	ctx, _ := gt.Get(caseUrl)
	html := ctx.RespBodyString()
	log.Println("state = ", ctx.Resp.Status)
	log.Println("html = ", html)

	//https://www.douyin.com/aweme/v1/web/search/item/?device_platform=webapp&aid=6383&channel=channel_pc_web&search_channel=aweme_video_web&sort_type=0&publish_time=0&keyword=%E7%BE%8E%E5%A5%B3&search_source=normal_search&query_correct_type=1&is_filter_search=0&offset=0&count=24&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1920&screen_height=1080&browser_language=zh-CN&browser_platform=Linux+x86_64&browser_name=Mozilla&browser_version=5.0+(X11%3B+Linux+x86_64)+AppleWebKit%2F537.36+(KHTML,+like+Gecko)+Chrome%2F88.0.4324.182+Safari%2F537.36&browser_online=true&msToken=hYnIG63GJJkmsYDWKUjWsvEDcDPmZEuy8MYgsjXxSlU1ZN6MC8hm8GhE_9a3ZyzPYfjdiNH1NAyJ3f9kVKYoWZKhLA97yJEFpZE0iuAmTv-d_ANDzS_XIdM=&X-Bogus=DFSzswSO6iXANti9Sm4Fs7kkNL-I&_signature=_02B4Z6wo000012Jp3cwAAIDAlRxGJKBU6ztibdlAALnad5Pv.LnK808o1d1sNqHrn0TuRCRBKdT7SBouK177aJxXYha8PtRCONuke.0SGNVqrwEC0oaLm3b1csTU2JaAVmPgWXxliTxTbadT4d
	//https://www.douyin.com/aweme/v1/web/search/item/?device_platform=webapp&aid=6383&channel=channel_pc_web&search_channel=aweme_video_web&sort_type=0&publish_time=0&keyword=%E7%BE%8E%E5%A5%B3&search_source=normal_search&query_correct_type=1&is_filter_search=0&offset=0&count=24&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1920&screen_height=1080&browser_language=zh-CN&browser_platform=Linux+x86_64&browser_name=Mozilla&browser_version=5.0+(X11%3B+Linux+x86_64)+AppleWebKit%2F537.36+(KHTML,+like+Gecko)+Chrome%2F88.0.4324.182+Safari%2F537.36&browser_online=true&msToken=V87HYI3JZjqkHIzHUXeCX-ydzD-tg0a-3XnKe4ve8CKjU3aTgNooos1-YkuX4vqKdeh8M-AxIGQ7Se9ZYuaY9hMQ_DMOusY6fvd8-8m-s_PB18E3jWDonM4x&X-Bogus=DFSzswSO6EGANti9Sm4FddkkNL-N&_signature=_02B4Z6wo00001reRrigAAIDBQOQ1w1S1fya3laqAAMy9pcfZE.rxv6pEdxi2yA2oLpfeiTpV7qKmXzx8P8z6GWU76dRxCzblqXzukwCUP3ENdbAGBoU090l4J8MOqXBZgVDV19sE4LsODg5b07

	/*
		H = 6383
		, q = "douyin_website"
		, X = ""
		, W = "cbae2a42b075f1dfc39b9e7df764637c821e7bf2";

	https://lf1-cdn-tos.bytegoofy.com/goofy/ies/douyin/search/index.50db626a.js:formatted  --> 11452


	*/


}

func main1(){
	//ctx, err := gt.Upload("https://v26-web.douyinvod.com/1671afda348c88fb13af3d73bc727002/6155656b/video/tos/cn/tos-cn-ve-15/1b4fef278b124e659a69c648793dd278/?a=6383&br=2249&bt=2249&cd=0%7C0%7C0&ch=26&cr=0&cs=0&cv=1&dr=0&ds=4&er=&ft=~MeSY~88-o10DN7nh7TQPqeUfTusSYjVWy&l=2021093014201201021205308106205B1D&lr=all&mime_type=video_mp4&net=0&pl=0&qs=0&rc=amlybzk6ZnF2ODMzNGkzM0ApZjQ2M2g2OTw0Nzo7OzxmZmdrbS1scjRnaC9gLS1kLTBzczYyMmMzYGJiYGAwYzQ1YDU6Yw%3D%3D&vl=&vr=",
	//	"/home/mange/Desktop/dy1.mp4")
	//log.Println(ctx.Resp.StatusCode, err)

	caseUrl := "https://www.douyin.com/video/7000282511469268264"
	ctx, _ := gt.Get(caseUrl)
	html := ctx.RespBodyString()
	log.Println("state = ", ctx.Resp.Status)
	log.Println("html = ", html)
	dom,err := goquery.NewDocumentFromReader(strings.NewReader(ctx.RespBodyString()))
	if err != nil{
		log.Println(err)
		return
	}
	result := dom.Find("script[id=RENDER_DATA]")
	log.Println(result.Html())
	res,_ := result.Html()
	unescape, _ := url.QueryUnescape(res)
	log.Println(unescape)

	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(unescape), &m)
	log.Println("json err = ", err)
	log.Println("m = ", m, "\n\n\n\n")

	c19 := m["C_19"].(map[string]interface{})
	//log.Println("c19 = ", c19)

	aweme := c19["aweme"].(map[string]interface{})
	//log.Println("aweme = ", aweme)

	detail := aweme["detail"].(map[string]interface{})
	log.Println("detail = ", detail)

	awemeId :=  detail["awemeId"].(string)
	log.Println("awemeId = ", awemeId)

	awemeType := detail["awemeType"]
	log.Println("awemeType = ", awemeType)

	groupId := detail["groupId"]
	log.Println("groupId = ", groupId)

	authorInfo := detail["authorInfo"].(map[string]interface{})
	log.Println("authorInfo = ", authorInfo)

	desc := detail["desc"]
	log.Println("desc = ", desc)

	authorUserId := detail["authorUserId"]
	log.Println("authorUserId = ", authorUserId)

	createTime := detail["createTime"]
	log.Println("createTime = ", createTime)

	video := detail["video"]
	log.Println("video = ", video)

	download := detail["download"].(map[string]interface{})
	log.Println("download = ", download)

	url := download["url"].(string)
	ctx1, err := gt.Upload(url,"/home/mange/Desktop/"+awemeId+".mp4")
	log.Println(ctx1.Resp.StatusCode, err)

}

/*
 */
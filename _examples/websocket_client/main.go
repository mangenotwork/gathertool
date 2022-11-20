package main

import (
	gt "github.com/mangenotwork/gathertool"
	"log"
	"math/rand"
	"time"
)

type Token struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	Time int64  `json:"time"`
	B    Body   `json:"body"`
}

type Body struct {
	P Pager `json:"pager"`
	D Data  `json:"data"`
}

type Pager struct {
	Page      int64 `json:"page"`
	Count     int64 `json:"count"`
	PageCount int64 `json:"page_count"`
}

type Data struct {
	Userid      int64  `json:"userid"`
	Appid       int64  `json:"appid"`
	TokenString string `json:"token_string"`
	Password    string `json:"password"`
}

func main() {
	GetAllUser()
}

func main1() {
	idStr := "104"
	//tokenUrl := "https://test.api.ymzy.cn/kimi/logic/allow/init/user"
	//par := `{
	//            "nick_name":"test-` + idStr + `",
	//            "user_id":104,
	//            "password":"1",
	//            "app_id":1
	//        }`
	//ctx, _ := gt.PostJson(tokenUrl, par)
	//
	//m := Token{}
	//
	//err := json.Unmarshal(ctx.RespBody, &m)
	//if err != nil {
	//	log.Println("Umarshal failed:", err)
	//	return
	//}
	//log.Println("Token = ", m.B.D.TokenString)

	/*
		/v1/ws?
		user_id=304930&app_id=1&name=%E6%88%91%E8%BF%98%E5%93%88%E5%93%88%E5%93%88%E5%93%88%E5%93%88%E5%93%88%E5%93%88%E5%93%88%E5%93%88%E5%93%88
		&avatar=https://p2.ymzy.cn/ave/20220215/7ebc2087f6aba85c.jpg
		&token=bmJmIjoxNjI2MzMxNzYzLCJuaWNrbmFtZSI6IjEyMzAwMDAwMDIzIiwicHJvdl9pZCI6NTAsInByb3ZfbW9kZWwiOjIsIn
		&device_id=304930-e9f44268-b747-4c2e-971c-1c15a855a15a
		&source=1

	*/

	token := "bmJmIjoxNjI2MzMxNzYzLCJuaWNrbmFtZSI6IjEyMzAwMDAwMDIzIiwicHJvdl9pZCI6NTAsInByb3ZfbW9kZWwiOjIsIn"

	// 连接
	host := "wss://test.api.ymzy.cn"
	path := "/kimi/v1/ws?user_id=" + idStr + "&device_id=" + idStr + "&token=" + token + "&app_id=1&name=104&avatar=&source=1"
	wsc, err := gt.WsClient(host, path, false)
	log.Println(wsc, err)
	rand.Seed(time.Now().UnixNano())

	//// 发消息
	//gt.PostForm("http://192.168.0.9:8001/auth/send", map[string]string{
	//	"conversation_type":"0",
	//	"sender_type":"2",
	//	"app_id":"1",
	//	"from":idStr,
	//	"from_device_id":"test_"+idStr,
	//	"from_avatar": userAvatar,
	//	"from_nickname": userNickName,
	//	"to": "-1",
	//	"content_type":"1",
	//	"content":"你好我找一下客服",
	//
	//})

	for {
		// CHANGE
		// QUESTION
		// ARTIFICIAL
		//		commands := []string{"CHANGE", "QUESTION", "ARTIFICIAL"}
		//
		//    	command := `{
		//    "type":12,
		//    "request_id":"",
		//    "command": "`+commands[rand.Intn(len(commands)-1)]+`",
		//    "device_id":"123123213123213213",
		//    "app_id":1,
		//    "user_id":104,
		//    "receiver_id":1,
		//    "seq":12321321321323
		//}`
		//
		//		err = wsc.Send([]byte(command))
		//		log.Println(err)
		data := make([]byte, 100)
		err = wsc.Read(data)
		log.Println(err)
		log.Println("data = ", string(data))
		time.Sleep(5 * time.Second)
		//time.Sleep(10*time.Microsecond)
	}

	wsc.Close()
}

var (
	host197       = "192.168.0.197"
	port          = 3306
	user          = "root"
	password      = "root123"
	gkzyUser      = "gkzy-user"
	GkzyUserDB, _ = gt.NewMysql(host197, port, user, password, gkzyUser)
)

func GetAllUser() {
	data, _ := GkzyUserDB.Select("select * from user limit 300, 100")
	for _, v := range data {
		log.Println(v["user_id"])
		go Work(v["user_id"], v["avatar"], v["nick_name"])
	}
	select {}
}

func Work(idStr, userAvatar, userNickName string) {
	token := "bmJmIjoxNjI2MzMxNzYzLCJuaWNrbmFtZSI6IjEyMzAwMDAwMDIzIiwicHJvdl9pZCI6NTAsInByb3ZfbW9kZWwiOjIsIn"

	// 连接
	host := "wss://test.api.ymzy.cn"
	path := "/kimi/v1/ws?user_id=" + idStr + "&device_id=" + idStr + "&token=" + token + "&app_id=1&name=104&avatar=&source=1"
	wsc, err := gt.WsClient(host, path, false)
	log.Println(wsc, err)
	rand.Seed(time.Now().UnixNano())

	for {
		// 发消息
		header := gt.NewHeader(map[string]string{
			"kimi-token": token,
		})
		ctx, _ := gt.PostJson("https://test.api.ymzy.cn/kimi/logic/auth/send", `{
		"conversation_type":0,
		"sender_type":2,
		"app_id":1,
		"from":`+idStr+`,
		"from_device_id":"test_`+idStr+`",
		"from_avatar": "`+userAvatar+`",
		"from_nickname": "`+userNickName+`",
		"to": -1,
		"content_type":1,
		"content":"[自动化压力测试] 你好我找一下客服"
	}`, header)

		log.Println(ctx.Json)

		data := make([]byte, 100)
		if wsc != nil {
			//err = wsc.Read(data)
			//log.Println(err)
			log.Println("data = ", string(data))
		}

		time.Sleep(30 * time.Second)
		//time.Sleep(10*time.Microsecond)
	}

	wsc.Close()
}

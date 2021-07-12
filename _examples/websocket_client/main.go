package main

import (
	"encoding/json"
	gt "github.com/mangenotwork/gathertool"
	"log"
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

func main(){
	idStr := "1"
	tokenUrl := "https://test.api.ymzy.cn/kimi/logic/allow/init/user"
	par := `{
                "nick_name":"test-` + idStr + `",
                "user_id":` + idStr + `,
                "password":"1",
                "app_id":3
            }`
	ctx, _ := gt.PostJson(tokenUrl, par)

	m := Token{}

	err := json.Unmarshal(ctx.RespBody, &m)
	if err != nil {
		log.Println("Umarshal failed:", err)
		return
	}
	log.Println("Token = ", m.B.D.TokenString)

	host := "wss://test.api.ymzy.cn"
	path := "/kimi/v1/ws?user_id=" + idStr + "&device_id=" + idStr + "&token=" + m.B.D.TokenString + "&app_id=3"
    wsc, err := gt.WsClient(host, path)
    log.Println(wsc, err)
    for {
		err = wsc.Send([]byte(`{"type":3,"request_id":"1"}`))
		log.Println(err)
		data := make([]byte,100)
		err = wsc.Read(data)
		log.Println(err)
		log.Println("data = ", string(data))
		time.Sleep(1*time.Second)
	}

	wsc.Close()
}
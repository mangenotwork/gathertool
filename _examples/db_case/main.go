package main

import (
	gt "github.com/mangenotwork/gathertool"
)

var (
	host192  = "192.168.0.192"
	port     = 3306
	user     = "root"
	password = "root123"
	spiderBase      = "spider-base"
	SpiderBaseDB, _ = gt.NewMysql(host192, port, user, password, spiderBase)

)

func main(){
	//case1()

	case2()
}

func case1(){
	data := map[string]interface{}{
		"num":123,
		"txt":"aaaaa'aaaaa",
		"txt2":`bbbb"bbb'bbbbbb'bbb"bbbbb`,
	}
	err := SpiderBaseDB.InsertAt("test", data)
	if err != nil {
		panic(err)
	}
}

func case2(){
	gd := gt.NewGDMap().Add("num", 123).Add("txt", "aaaaa'aaaaa")
	gd.Add("txt2", `bbbb"bbb'bbbbbb'bbb"bbbbb`)
	err := SpiderBaseDB.InsertAtGd("test2", gd)
	if err != nil {
		panic(err)
	}
}

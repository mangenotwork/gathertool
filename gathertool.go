package gathertool

import (
	"log"
	"net/http"
)

// Get 请求, 当请求失败或状态码是失败的则会先执行 ff 再回调
//
// @url  请求链接
// @maxTimes  重试次数
// @sf  请求成功后做的事情, 200等
// @ff  请求失败后做的事情, 403等，502等
// @vs  可变参数
//
func Get(url string, maxTimes int64, vs ...interface{}) (*Req,error){
	var (
		client *http.Client
	)

	//初始化 Request
	req, err := http.NewRequest("GET",url,nil)
	if err != nil{
		return nil,err
	}

	//添加默认的Header
	req.Header.Add("Connection","close")
	req.Header.Add("User-Agent","zhaofan")

	//解析可变参
	for _, v := range vs {
		log.Println("参数： ", v)
		switch vv := v.(type) {
		// 使用方传入了 header
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
			// 使用方传入了 *http.Client
		case *http.Client:
			client = vv
		}
	}

	// 如果使用方未传入Client，  初始化 Client
	if client == nil{
		log.Println("使用方未传入Client， 默认 client")
		client = &http.Client{}
	}

	// 创建对象
	return &Req{
		Client: client,
		Req : req,
		times : 0,
		MaxTimes : maxTimes,
	},nil
}
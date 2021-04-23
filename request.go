package gathertool

import (
	"io/ioutil"
	"log"
	"net/http"
)

type Req struct {

	// client
	Client *http.Client

	// 请求
	Req *http.Request


	// 执行的次数 初始化都是0
	times int64

	// 最大允许重试次数
	MaxTimes int64

	// 请求成功了需要处理的事件
	SuccessFunc func([]byte)

	// 请求失败了需要做的事， 如休息等待或设置代理
	FailFunc func(req *Req)
}


// Get 请求, 当请求失败或状态码是失败的则会先执行 ff 再回调
//
// @url  请求链接
// @maxTimes  重试次数
// @sf  请求成功后做的事情, 200等
// @ff  请求失败后做的事情, 403等，502等
// @vs  可变参数
//
func Get(url string, maxTimes int64, sf func([]byte), vs ...interface{}) (*Req,error){
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
		SuccessFunc : sf,
	},nil
}

// SetFailed 设置错误处理
func (r *Req) SetFailed(failedFunc func(c *Req)) {
	r.FailFunc = failedFunc
}

// Do 执行请求
func (r *Req) Do() func(){

	//空验证
	if r == nil{
		log.Println("空对象")
		return nil
	}

	//重试验证
	r.times++
	if r.times > r.MaxTimes{
		log.Println("请求失败操过", r.MaxTimes, "次了")
		return nil
	}

	//执行请求
	resp,err := r.Client.Do(r.Req)
	if err != nil || resp.StatusCode != 200{
		log.Println("第", r.times, "请求失败.")

		//执行失败方法
		if r.FailFunc != nil{
			r.FailFunc(r)
		}

		//回调
		return r.Do()
	}
	defer resp.Body.Close()

	//请求后的结果
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Println(err)
		return nil
	}

	//执行成功方法
	r.SuccessFunc(body)
	return nil
}

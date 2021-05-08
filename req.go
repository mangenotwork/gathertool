/*
	Description : 对外提供的方法
	Author : ManGe
	Version : v0.1
	Date : 2021-04-23
*/

package gathertool

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	UrlBad error = errors.New("url is bad.") // 错误的url
)

type ReqTimeOut int
type ReqTimeOutMs int


// Get 请求, 当请求失败或状态码是失败的则会先执行 ff 再回调
func Get(url string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	return	Req(request, vs...)
}


// Get 请求直接执行
func GetRun(url string, vs ...interface{}) (ctx *Context,err error) {
	if !isUrl(url) {
		err = errors.New("请求 url 为空.")
		return
	}
	request, err := http.NewRequest("GET", url, nil)
	ctx, err = Req(request, vs...)
	if ctx != nil {
		ctx.Do()
	}
	return
}


// POST 请求
func Post(url string, data []byte, contentType string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return	Req(request, vs...)
}


// POST json 请求
func PostJson(url string, jsonStr string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	return	Req(request, vs...)
}


// Put
func Put(url string, data []byte, contentType string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("PUT", url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return	Req(request, vs...)
}


// Delete
func Delete(url string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	return	Req(request, vs...)
}


// Options
func Options(url string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	return	Req(request, vs...)
}


// Request 请求
func Request(url, method string, data []byte, contentType string, vs ...interface{}) (*Context,error){
	if !isUrl(url) {
		return nil, UrlBad
	}
	request, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(data)))
	if err != nil{
		log.Println("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return	Req(request, vs...)
}


// Upload
func Upload(url, savePath string, vs ...interface{})  error {
	if !isUrl(url) {
		return UrlBad
	}
	c,err := Get(url,vs)
	if err != nil{
		return err
	}
	c.Upload(savePath)
	return nil
}


// isUrl 验证是否是有效的 url
func isUrl(url string) bool {
	if url == ""{
		return false
	}
	return true
}


type Header map[string]string


// Req 初始化请求
// @url  请求链接
// @maxTimes  重试次数
// @sf  请求成功后做的事情, 200等
// @ff  请求失败后做的事情, 403等，502等
// @vs  可变参数
// @vs UserAgentType  设置指定类型 user agent 如 AndroidAgent
func Req(request *http.Request, vs ...interface{}) (*Context,error){
	var (
		client *http.Client
		maxTimes RetryTimes = 10
		task *Task
		start StartFunc
		succeed SucceedFunc
		failed FailedFunc
		retry RetryFunc
		end EndFunc
		reqTimeOut ReqTimeOut
		reqTimeOutMs ReqTimeOutMs
	)

	//添加默认的Header
	request.Header.Set("Connection","close")
	request.Header.Set("User-Agent", GetAgent(PCAgent))

	//解析可变参
	for _, v := range vs {
		switch vv := v.(type) {
		case http.Header:
			for key, values := range vv {
				for _, value := range values {
					request.Header.Add(key, value)
				}
			}
		case *http.Header:
			for key, values := range *vv {
				for _, value := range values {
					request.Header.Add(key, value)
				}
			}
		case Header:
			for key, value := range vv {
				request.Header.Add(key, value)
			}
		case *http.Client:
			client = vv
		case UserAgentType:
			request.Header.Set("User-Agent", GetAgent(vv))
		case *http.Cookie:
			request.AddCookie(vv)
		case RetryTimes:
			maxTimes = vv
		case *Task:
			task = vv
		case StartFunc:
			start = vv
		case SucceedFunc:
			succeed = vv
		case FailedFunc:
			failed = vv
		case RetryFunc:
			retry = vv
		case EndFunc:
			end = vv
		case ReqTimeOut:
			reqTimeOut = vv
		case ReqTimeOutMs:
			reqTimeOutMs = vv
		}
	}

	// task Header
	if task != nil && task.HeaderMap != nil {
		for k,v := range *task.HeaderMap {
			for _, value := range v {
				request.Header.Add(k, value)
			}
		}
	}

	// 如果使用方未传入Client，  初始化 Client
	if client == nil{
		//log.Println("使用方未传入Client， 默认 client")
		client = &http.Client{
			Timeout: 60*time.Second,
		}
	}

	if reqTimeOut > 0 {
		client.Timeout =  time.Duration(reqTimeOut) * time.Second
	}

	if reqTimeOutMs > 0 {
		client.Timeout =  time.Duration(reqTimeOutMs) * time.Millisecond
	}

	// 创建对象
	return &Context{
		Client: client,
		Req : request,
		times : 0,
		MaxTimes : maxTimes,
		Task: task,
		StartFunc: start,
		SucceedFunc: succeed,
		FailedFunc: failed,
		RetryFunc: retry,
		EndFunc: end,
	},nil
}

func SearchDomain(ip string){
	addr, err := net.LookupTXT(ip)
	log.Println(addr, err)
}

// 扫描ip的端口
func SearchPort(ipStr string, vs ...interface{}) {

	timeOut := 4*time.Second

	for _, v := range vs {
		switch vv := v.(type) {
		case time.Duration:
			timeOut = vv
		}
	}

	queue := NewQueue()

	for i:=0;i<65536;i++{
		buf := &bytes.Buffer{}
		buf.WriteString(ipStr)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(i))
		queue.Add(&Task{
			Url: buf.String(),
		})
	}

	var wg sync.WaitGroup

	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	for i:=0;i<65536;i++{
	//		ip := net.ParseIP(ipStr)
	//
	//		tcpAddr := &net.TCPAddr{
	//			IP:ip,
	//			Port:i,
	//		}
	//
	//		queue.Add(&Task{
	//			Url: tcpAddr.String(),
	//		})
	//	}
	//}()


	for job:=0;job<65536;job++{
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			for {

				if queue.IsEmpty(){
					break
				}

				task := queue.Poll()
				//log.Println(task.Url)
				if task == nil {
					continue
				}

				conn, err := net.DialTimeout("tcp", task.Url, timeOut)
				if err == nil {
					log.Println(task.Url, "开放")
					conn.Close()
				}
			}
		}(job)
	}

	wg.Wait()




	log.Println("执行完成！！！")
}


func Ping(ip string){

	before := time.Now()
	defer func(tStart time.Time){
		dur := time.Now().Sub(before)
		log.Println("来自 ",ip," 的回复: 时间 = ", dur)
	}(before)

	c, err := net.Dial("ip4:icmp", ip)
	if err != nil {
		return
	}
	log.Println(c)
	c.SetDeadline(time.Now().Add(1 * time.Second))
	defer c.Close()

	var msg [512]byte
	msg[0] = 8
	msg[1] = 0
	msg[2] = 0
	msg[3] = 0
	msg[4] = 0
	msg[5] = 13
	msg[6] = 0
	msg[7] = 37
	msg[8] = 99
	len := 9
	check := checkSum(msg[0:len])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 0xff)
	fmt.Println(msg[0:len])

	_, err = c.Write(msg[0:len])
	if err!= nil{
		log.Println(ip," -> ping err : ", err)
		return
	}

	c.Write(msg[0:len])


	_, err = c.Read(msg[0:])
	if err!= nil{
		log.Println(ip," -> ping err : ", err)
		return
	}

	//log.Println(string(msg[0 : 20+len]))

	//log.Println("Got response")
	if msg[20+5] != 13 {
		log.Println(ip," -> ping err : Identifier not matches")
		return
	}
	if msg[20+7] != 37 {
		log.Println(ip," -> ping err : Sequence not matches")
		return
	}
	if msg[20+8] != 99 {
		log.Println(ip," -> ping err : Custom data not matches")
		return
	}

	log.Println("ping ok : ", ip)
}

func checkSum(msg []byte) uint16 {
	sum := 0

	len := len(msg)
	for i := 0; i < len-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if len%2 == 1 {
		sum += int(msg[len-1]) * 256 // notice here, why *256?
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}
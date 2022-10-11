/*
	Description : 请求相关
	Author : ManGe
			2912882908@qq.com
			https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"bytes"
	"crypto/tls"
	"errors"
	"golang.org/x/net/publicsuffix"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	POST    = "POST"
	GET     = "GET"
	HEAD    = "HEAD"
	PUT     = "PUT"
	DELETE  = "DELETE"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
	ANY     = ""
)

var (
	UrlBad = errors.New("url is bad.") // 错误的url
	UrlNil = errors.New("url is null.") // 空的url
	randEr = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type ReqTimeOut int
type ReqTimeOutMs int
type Sleep time.Duration

// Request 请求
func Request(url, method string, data []byte, contentType string, vs ...interface{}) (*Context, error) {
	request, err := http.NewRequest(method, urlStr(url), bytes.NewBuffer([]byte(data)))
	if err != nil{
		Error("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	cxt :=	Req(request, vs...)
	cxt.Do()
	return cxt, nil
}

// NewRequest 请求
func NewRequest(url, method string, data []byte, contentType string, vs ...interface{}) *Context {
	request, err := http.NewRequest(method, urlStr(url), bytes.NewBuffer([]byte(data)))
	if err != nil{
		Error("err->", err)
		return nil
	}
	request.Header.Set("Content-Type", contentType)
	return	Req(request, vs...)
}

// isUrl 验证是否是有效的 url
func isUrl(url string) bool {
	if url == ""{
		return false
	}
	return true
}

func urlStr(url string) string {
	l := len(url)
	if l < 1 {
		panic("url is null")
	}
	if l > 8 && (url[:7] == "http://" || url[:8] == "https://") {
		return url
	}
	return "http://" + url
}

// SetSleep 设置请求随机休眠时间， 单位秒
func SetSleep(min, max int) Sleep {
	r := randEr.Intn(max) + min
	return Sleep(time.Duration(r)*time.Second)
}

// SetSleepMs 设置请求随机休眠时间， 单位毫秒
func SetSleepMs(min, max int) Sleep {
	r := randEr.Intn(max) + min
	return Sleep(time.Duration(r)*time.Millisecond)
}

// Header
type Header map[string]string

// NewHeader 新建Header
func NewHeader(data map[string]string) Header {
	return data
}

func (h Header) haveObj() {
	if h == nil {
		h = Header{}
	}
}

// Set Header Set
func (h Header) Set(key, value string) Header {
	h.haveObj()
	h[key] = value
	return h
}

// Delete Header Delete
func (h Header) Delete(key string) Header {
	h.haveObj()
	delete(h, key)
	return h
}

// Cookie
type Cookie map[string]string

// NewCookie 新建Cookie
func NewCookie(data map[string]string) Cookie {
	return data
}

func (c Cookie) haveObj() {
	if c == nil {
		c = Cookie{}
	}
}

// Set Cookie Set
func (c Cookie) Set(key, value string) Cookie {
	c.haveObj()
	c[key] = value
	return c
}

// Delete Cookie Delete
func (c Cookie) Delete(key string) Cookie {
	c.haveObj()
	delete(c, key)
	return c
}

// Req 初始化请求
func Req(request *http.Request, vs ...interface{}) *Context {
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
		islog IsLog
		proxyUrl string
		sleep Sleep
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
		case Cookie:
			for key, value := range vv {
				request.AddCookie(&http.Cookie{Name: key, Value: value, HttpOnly: true})
			}
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
		case IsLog:
			islog = vv
		case ProxyUrl:
			proxyUrl = string(vv)
		case Sleep:
			sleep = vv
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
	if task == nil {
		task = NewTask()
	}
	// 如果使用方未传入Client，  初始化 Client
	if client == nil{
		client = &http.Client{
			Timeout: 60*time.Second,
		}
		// Transport 设置
		client.Transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				FallbackDelay: -1 * time.Nanosecond,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			// gathertool默认每个请求实例都创建一个独立的client，
			// 不复用client，这样设计是在高并发中，每个请求都是独立的
			MaxIdleConns:          10,
			MaxIdleConnsPerHost:   5, // 默认是 2
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives: true,//DisableKeepAlives这个字段可以用来关闭长连接，默认值为false
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	if reqTimeOut > 0 {
		client.Timeout =  time.Duration(reqTimeOut) * time.Second
	}

	if reqTimeOutMs > 0 {
		client.Timeout =  time.Duration(reqTimeOutMs) * time.Millisecond
	}

	if proxyUrl != "" {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			Error("设置代理失败:", err)
		}else{
			client.Transport =  &http.Transport{Proxy: http.ProxyURL(proxy)}
		}
	}

	// CookieJar管理
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil {
		client.Jar = jar
	}

	if l := request.Header.Get("Content-Length"); l != "" {
		request.ContentLength = Str2Int64(l)
	}

	if retry == nil {
		retry = defaultRetry
	}

	if succeed == nil {
		succeed = defaultSucceed
	}

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
		IsLog: islog,
		sleep: sleep,
		Param: make(map[string]interface{}),
	}
}

// SearchDomain 搜索host
func SearchDomain(ip string){
	addr, err := net.LookupTXT(ip)
	Info(addr, err)
}

// SearchPort 扫描ip的端口
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
	for job:=0;job<65536;job++{
		wg.Add(1)
		go func(i int){
			defer wg.Done()
			for {
				if queue.IsEmpty(){
					break
				}
				task := queue.Poll()
				if task == nil {
					continue
				}
				conn, err := net.DialTimeout("tcp", task.Url, timeOut)
				if err == nil {
					Error(task.Url, "开放")
					conn.Close()
				}
			}
		}(job)
	}
	wg.Wait()
	Info("执行完成！！！")
}

// Ping ping IP
func Ping(ip string){
	before := time.Now()
	defer func(tStart time.Time){
		dur := time.Now().Sub(before)
		Info("来自 ",ip," 的回复: 时间 = ", dur)
	}(before)
	c, err := net.Dial("ip4:icmp", ip)
	if err != nil {
		return
	}
	Info(c)
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
	l := 9
	check := checkSum(msg[0:l])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 0xff)
	Info(msg[0:l])
	_, err = c.Write(msg[0:l])
	if err!= nil{
		log.Println(ip," -> ping err : ", err)
		return
	}
	c.Write(msg[0:l])
	_, err = c.Read(msg[0:])
	if err!= nil{
		log.Println(ip," -> ping err : ", err)
		return
	}
	if msg[20+5] != 13 {
		Error(ip," -> ping err : Identifier not matches")
		return
	}
	if msg[20+7] != 37 {
		Error(ip," -> ping err : Sequence not matches")
		return
	}
	if msg[20+8] != 99 {
		Error(ip," -> ping err : Custom data not matches")
		return
	}
	Info("ping ok : ", ip)
}

func checkSum(msg []byte) uint16 {
	sum := 0
	l := len(msg)
	for i := 0; i < l-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if l%2 == 1 {
		sum += int(msg[l-1]) * 256 // notice here, why *256?
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}

func defaultRetry(ctx *Context) {
	Info("2s后准备重试......")
	time.Sleep(2*time.Second)
}

func defaultSucceed(ctx *Context) {
	Info("请求成功")
	//Info(ctx.RespBodyString())
}

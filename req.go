/*
*	Description : 请求相关
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

func init() {
	rand.NewSource(time.Now().UnixNano())
}

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
	randEr = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type ReqTimeOut int
type ReqTimeOutMs int
type Sleep time.Duration

// Request 请求
func Request(url, method string, data []byte, contentType string, vs ...interface{}) (*Context, error) {
	request, err := http.NewRequest(method, urlStr(url), bytes.NewBuffer(data))
	if err != nil {
		Error("err->", err)
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	cxt := Req(request, vs...)
	cxt.Do()
	return cxt, nil
}

// NewRequest 请求
func NewRequest(url, method string, data []byte, contentType string, vs ...interface{}) *Context {
	request, err := http.NewRequest(method, urlStr(url), bytes.NewBuffer(data))
	if err != nil {
		Error("err->", err)
		return nil
	}
	request.Header.Set("Content-Type", contentType)
	return Req(request, vs...)
}

// isUrl 验证是否是有效的 url
func isUrl(url string) bool {
	if url == "" {
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
	return "https://" + url
}

func IsUrl(url string) bool {
	return isUrl(url)
}

func UrlStr(url string) string {
	return urlStr(url)
}

// SetSleep 设置请求随机休眠时间， 单位秒
func SetSleep(min, max int) Sleep {
	r := randEr.Intn(max) + min
	return Sleep(time.Duration(r) * time.Second)
}

// SetSleepMs 设置请求随机休眠时间， 单位毫秒
func SetSleepMs(min, max int) Sleep {
	r := randEr.Intn(max) + min
	return Sleep(time.Duration(r) * time.Millisecond)
}

type Header map[string]string

// NewHeader 新建Header
func NewHeader(data map[string]string) Header {
	return data
}

// Set Header Set
func (h Header) Set(key, value string) Header {
	h[key] = value
	return h
}

// Delete Header Delete
func (h Header) Delete(key string) Header {
	delete(h, key)
	return h
}

type Cookie map[string]string

// NewCookie 新建Cookie
func NewCookie(data map[string]string) Cookie {
	return data
}

// Set Cookie Set
func (c Cookie) Set(key, value string) Cookie {
	c[key] = value
	return c
}

// Delete Cookie Delete
func (c Cookie) Delete(key string) Cookie {
	delete(c, key)
	return c
}

// Req 初始化请求
func Req(request *http.Request, vs ...interface{}) *Context {
	var (
		client       *http.Client
		maxTimes     RetryTimes = 10
		task         *Task
		start        StartFunc
		succeed      SucceedFunc
		failed       FailedFunc
		retry        RetryFunc
		end          EndFunc
		reqTimeOut   ReqTimeOut
		reqTimeOutMs ReqTimeOutMs
		isLog        IsLog
		isRetry      IsRetry
		proxyUrl     string
		sleep        Sleep
	)

	//添加默认的Header
	request.Header.Set("Connection", "close")
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
			isLog = vv
		case IsRetry:
			isRetry = vv
		case ProxyUrl:
			proxyUrl = string(vv)
		case Sleep:
			sleep = vv
		}
	}

	// task Header
	if task != nil && task.HeaderMap != nil {
		for k, v := range *task.HeaderMap {
			for _, value := range v {
				request.Header.Add(k, value)
			}
		}
	}
	if task == nil {
		task = NewTask()
	}

	// 如果使用方未传入Client，  初始化 Client
	if client == nil {
		client = &http.Client{
			Timeout: 30 * time.Second,
		}
		// Transport 设置
		client.Transport = &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				n := &net.Dialer{
					Timeout:       30 * time.Second,
					KeepAlive:     30 * time.Second,
					FallbackDelay: -1 * time.Nanosecond,
				}
				conn, err := n.DialContext(ctx, network, addr)
				if err == nil && conn != nil {
					request.RemoteAddr = conn.RemoteAddr().String()
				}
				return conn, err
			},
			ForceAttemptHTTP2: true,
			// gatherTool 默认每个请求实例都创建一个独立的client，
			// 不复用client，这样设计是在高并发中，每个请求都是独立的
			MaxIdleConns:          10,
			MaxIdleConnsPerHost:   5, // 默认是 2
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     true, // 这个字段可以用来关闭长连接，默认值为false
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		}
	}

	if reqTimeOut > 0 {
		client.Timeout = time.Duration(reqTimeOut) * time.Second
	}

	if reqTimeOutMs > 0 {
		client.Timeout = time.Duration(reqTimeOutMs) * time.Millisecond
	}

	if proxyUrl != "" {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			Error("设置代理失败:", err)
		} else {
			client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
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
		Client:      client,
		Req:         request,
		times:       0,
		MaxTimes:    maxTimes,
		Task:        task,
		StartFunc:   start,
		SucceedFunc: succeed,
		FailedFunc:  failed,
		RetryFunc:   retry,
		EndFunc:     end,
		IsLog:       isLog,
		isRetry:     isRetry,
		sleep:       sleep,
		Param:       make(map[string]interface{}),
	}
}

// SearchDomain 搜索host
func SearchDomain(ip string) {
	addr, err := net.LookupTXT(ip)
	Info(addr, err)
}

// SearchPort 扫描ip的端口
func SearchPort(ipStr string, vs ...interface{}) {
	timeOut := 4 * time.Second
	for _, v := range vs {
		switch vv := v.(type) {
		case time.Duration:
			timeOut = vv
		}
	}

	queue := NewQueue()
	for i := 0; i < 65536; i++ {
		buf := &bytes.Buffer{}
		buf.WriteString(ipStr)
		buf.WriteString(":")
		buf.WriteString(strconv.Itoa(i))
		_ = queue.Add(&Task{
			Url: buf.String(),
		})
	}

	var wg sync.WaitGroup
	for job := 0; job < 65536; job++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for {
				if queue.IsEmpty() {
					break
				}
				task := queue.Poll()
				if task == nil {
					continue
				}
				conn, err := net.DialTimeout("tcp", task.Url, timeOut)
				if err == nil {
					Error(task.Url, "开放")
					_ = conn.Close()
				}
			}
		}(job)
	}
	wg.Wait()
	Info("执行完成！！！")
}

var pingTerminalPrint = true

func ClosePingTerminalPrint() {
	pingTerminalPrint = false
}

// Ping ping IP
func Ping(ip string) (time.Duration, error) {
	before := time.Now()
	c, err := net.Dial("ip4:icmp", ip)
	if err != nil {
		return time.Duration(0), err
	}
	_ = c.SetDeadline(time.Now().Add(1 * time.Second))
	defer func() {
		_ = c.Close()
	}()
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
	_, err = c.Write(msg[0:l])
	if err != nil {
		if pingTerminalPrint {
			Info(ip, " -> ping err : ", err)
		}
		return time.Duration(0), err
	}
	_, _ = c.Write(msg[0:l])
	_, err = c.Read(msg[0:])
	if err != nil {
		if pingTerminalPrint {
			Info(ip, " -> ping err : ", err)
		}
		return time.Duration(0), err
	}
	if msg[20+5] != 13 {
		if pingTerminalPrint {
			Error(ip, " -> ping err : Identifier not matches")
		}
		return time.Duration(0), fmt.Errorf(ip, " -> ping err : Identifier not matches")
	}
	if msg[20+7] != 37 {
		if pingTerminalPrint {
			Error(ip, " -> ping err : Sequence not matches")
		}
		return time.Duration(0), fmt.Errorf(ip, " -> ping err : Sequence not matches")
	}
	if msg[20+8] != 99 {
		if pingTerminalPrint {
			Error(ip, " -> ping err : Custom data not matches")
		}
		return time.Duration(0), fmt.Errorf(ip, " -> ping err : Custom data not matches")
	}
	dur := time.Now().Sub(before)
	if pingTerminalPrint {
		Info("来自 ", ip, " 的回复: 时间 = ", dur)
	}
	return dur, nil
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
	sum += sum >> 16
	var answer = uint16(^sum)
	return answer
}

func defaultRetry(ctx *Context) {
	Info("2s后准备重试......")
	time.Sleep(2 * time.Second)
}

func defaultSucceed(ctx *Context) {
	//Info("请求成功")
}

type SSLCertificateInfo struct {
	Url                   string `json:"Url"`                   // url
	EffectiveTime         string `json:"EffectiveTime"`         // 有效时间
	NotBefore             int64  `json:"NotBefore"`             // 起始
	NotAfter              int64  `json:"NotAfter"`              // 结束
	DNSName               string `json:"DNSName"`               // DNSName
	OCSPServer            string `json:"OCSPServer"`            // OCSPServer
	CRLDistributionPoints string `json:"CRLDistributionPoints"` // CRL分发点
	Issuer                string `json:"Issuer"`                // 颁发者
	IssuingCertificateURL string `json:"IssuingCertificateURL"` // 颁发证书URL
	PublicKeyAlgorithm    string `json:"PublicKeyAlgorithm"`    // 公钥算法
	Subject               string `json:"Subject"`               // 颁发对象
	Version               string `json:"Version"`               // 版本
	SignatureAlgorithm    string `json:"SignatureAlgorithm"`    // 证书算法
}

func (s SSLCertificateInfo) Echo() {
	txt := `
Url: %s 
有效时间: %s 
颁发对象: %s 
颁发者: %s 
颁发证书URL: %s 
公钥算法: %s 
证书算法: %s 
版本: %s 
DNSName: %s 
CRL分发点: %s 
OCSPServer: %s `
	Info(fmt.Sprintf(txt, s.Url, s.EffectiveTime, s.Subject, s.Issuer, s.IssuingCertificateURL, s.PublicKeyAlgorithm,
		s.SignatureAlgorithm, s.Version, s.DNSName, s.CRLDistributionPoints, s.OCSPServer))
}
func (s SSLCertificateInfo) Expire() int64 {
	return s.NotAfter - time.Now().Unix()
}

// GetCertificateInfo 获取SSL证书信息
func GetCertificateInfo(caseUrl string) (SSLCertificateInfo, bool) {
	var info = SSLCertificateInfo{
		Url: caseUrl,
	}
	var cert *x509.Certificate
	var err error
	client := http.Client{}
	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
				if len(rawCerts) < 1 {
					return nil
				}
				cert, err = x509.ParseCertificate(rawCerts[0])
				return err
			},
		},
	}
	_, err = client.Get(caseUrl)
	if err != nil || cert == nil {
		return info, false
	}
	info.NotBefore = cert.NotBefore.Unix()
	info.NotAfter = cert.NotAfter.Unix()
	info.EffectiveTime = fmt.Sprintf("%s 到 %s", Timestamp2Date(info.NotBefore), Timestamp2Date(info.NotAfter))
	info.DNSName = strings.Join(cert.DNSNames, ";")
	info.OCSPServer = strings.Join(cert.OCSPServer, ";")
	info.CRLDistributionPoints = strings.Join(cert.CRLDistributionPoints, ";")
	info.Issuer = cert.Issuer.String()
	info.IssuingCertificateURL = strings.Join(cert.IssuingCertificateURL, ";")
	info.PublicKeyAlgorithm = cert.PublicKeyAlgorithm.String()
	info.Subject = cert.Subject.String()
	info.Version = strconv.Itoa(cert.Version)
	info.SignatureAlgorithm = cert.SignatureAlgorithm.String()
	return info, true
}

// UsefulUrl 判断url是否是有效的
func UsefulUrl(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	u, err := url.Parse(str)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	// Check if the URL has a valid scheme (http or https)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	return true
}

// StatusCodeMap 状态码处理映射
// success 该状态码对应执行成功函数
// fail    该状态码对应执行失败函数
// retry   该状态码对应需要重试前执行的函数
var StatusCodeMap = map[int]string{
	200: "success",
	201: "success",
	202: "success",
	203: "success",
	204: "fail",
	300: "success",
	301: "success",
	302: "success",
	400: "fail",
	401: "retry",
	402: "retry",
	403: "retry",
	404: "fail",
	405: "retry",
	406: "retry",
	407: "retry",
	408: "retry",
	412: "success",
	500: "fail",
	501: "fail",
	502: "retry",
	503: "retry",
	504: "retry",
}

// SetStatusCodeSuccessEvent 将指定状态码设置为执行成功事件
func SetStatusCodeSuccessEvent(code int) {
	StatusCodeMap[code] = "success"
}

// SetStatusCodeRetryEvent 将指定状态码设置为执行重试事件
func SetStatusCodeRetryEvent(code int) {
	StatusCodeMap[code] = "retry"
}

// SetStatusCodeFailEvent 将指定状态码设置为执行失败事件
func SetStatusCodeFailEvent(code int) {
	StatusCodeMap[code] = "fail"
}

//// 设置为最大 socket open file ; windows can't!!!
//func SetMaxOpenFile() error {
//	var rLimit syscall.Rlimit
//
//	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
//	if err != nil {
//		return err
//	}
//
//	rLimit.Cur = rLimit.Max
//	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
//}
//
//// 设置指定 socket open file
//func SetRLimit(number int) error {
//	var rLimit syscall.Rlimit
//	var n = uint64(number)
//
//	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
//	if err != nil {
//		return err
//	}
//	if n > rLimit.Max {
//		return errors.New("设置失败：操过最大 Rlimit")
//	}
//	rLimit.Cur = n
//	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
//}

type UserAgentType int

const (
	PCAgent UserAgentType = iota + 1
	WindowsAgent
	LinuxAgent
	MacAgent
	AndroidAgent
	IosAgent
	PhoneAgent
	WindowsPhoneAgent
	UCAgent
)

var UserAgentMap = map[int]string{
	1:  "Mozilla/5.0 (Windows NT 6.2; WOW64; rv:21.0) Gecko/20100101 Firefox/21.0",                                                                                                                                          //Firefox on Windows
	2:  "Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.94 Safari/537.36",                                                                                                      //Chrome on Windows
	3:  "Mozilla/5.0 (compatible; WOW64; MSIE 10.0; Windows NT 6.2)",                                                                                                                                                        //Internet Explorer 10
	4:  "Opera/9.80 (Windows NT 6.1; WOW64; U; en) Presto/2.10.229 Version/11.62",                                                                                                                                           //Opera on Windows
	5:  "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27",                                                                                          //Safari on Windows
	6:  "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:21.0) Gecko/20130331 Firefox/21.0",                                                                                                                                      //Firefox on Ubuntu
	7:  "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/535.11 (KHTML, like Gecko) Ubuntu/11.10 Chromium/27.0.1453.93 Chrome/27.0.1453.93 Safari/537.36",                                                                       //Chrome on Ubuntu
	8:  "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_6; en-US) AppleWebKit/533.20.25 (KHTML, like Gecko) Version/5.0.4 Safari/533.20.27",                                                                                 //Safari on Mac
	9:  "Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.9.168 Version/11.52",                                                                                                                                 //Opera on Mac
	10: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/27.0.1453.93 Safari/537.36",                                                                                           //Chrome on Mac
	11: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.8; rv:21.0) Gecko/20100101 Firefox/21.0",                                                                                                                                 //Firefox on Mac
	12: "Mozilla/5.0 (Linux; Android 4.1.1; Nexus 7 Build/JRO03D) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166  Safari/535.19",                                                                               //Nexus 7 (Tablet)
	13: "Mozilla/5.0 (Linux; U; Android 4.0.4; en-gb; GT-I9300 Build/IMM76D) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",                                                                       //Samsung Galaxy S3 (Handset)
	14: "Mozilla/5.0 (Linux; U; Android 2.2; en-gb; GT-P1000 Build/FROYO) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",                                                                            //Samsung Galaxy Tab (Tablet)
	15: "Mozilla/5.0 (Android; Mobile; rv:14.0) Gecko/14.0 Firefox/14.0",                                                                                                                                                    //Firefox on Android Mobile
	16: "Mozilla/5.0 (Android; Tablet; rv:14.0) Gecko/14.0 Firefox/14.0",                                                                                                                                                    //Firefox on Android Tablet
	17: "Mozilla/5.0 (Linux; Android 4.0.4; Galaxy Nexus Build/IMM76B) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.133 Mobile Safari/535.19",                                                                    //Chrome on Android Mobile
	18: "Mozilla/5.0 (Linux; Android 4.1.2; Nexus 7 Build/JZ054K) AppleWebKit/535.19 (KHTML, like Gecko) Chrome/18.0.1025.166 Safari/535.19",                                                                                //Chrome on Android Tablet
	19: "Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",                                                                            //iPhone
	20: "Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",                                                                                     //iPad
	21: "Mozilla/5.0 (iPhone; CPU iPhone OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",                                                                            //Safari on iPhone
	22: "Mozilla/5.0 (iPad; CPU OS 5_0 like Mac OS X) AppleWebKit/534.46 (KHTML, like Gecko) Version/5.1 Mobile/9A334 Safari/7534.48.3",                                                                                     //Safari on iPad
	23: "Mozilla/5.0 (iPhone; CPU iPhone OS 6_1_4 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) CriOS/27.0.1453.10 Mobile/10B350 Safari/8536.25",                                                                    //Chrome on iPhone
	24: "Mozilla/5.0 (compatible; MSIE 10.0; Windows Phone 8.0; Trident/6.0; IEMobile/10.0; ARM; Touch; NOKIA; Lumia 920)",                                                                                                  //Windows Phone 8
	25: "Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0; SAMSUNG; SGH-i917)",                                                                                                            //Windows Phone 7.5
	26: "User-Agent, UCWEB7.0.2.37/28/999",                                                                                                                                                                                  //UC无
	27: "User-Agent, NOKIA5700/ UCWEB7.0.2.37/28/999",                                                                                                                                                                       //UC标准
	28: "User-Agent, Openwave/ UCWEB7.0.2.37/28/999",                                                                                                                                                                        //UCOpenwave
	29: "User-Agent, Mozilla/4.0 (compatible; MSIE 6.0; ) Opera/UCWEB7.0.2.37/28/999",                                                                                                                                       //UC Opera
	30: "Mozilla/5.0 (Windows; U; Windows NT 6.1; ) AppleWebKit/534.12 (KHTML, like Gecko) Maxthon/3.0 Safari/534.12",                                                                                                       //傲游3.1.7在Win7+ie9,高速模式
	31: "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E)",                    //傲游3.1.7在Win7+ie9,IE内核兼容模式:
	32: "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)", //搜狗3.0在Win7+ie9,IE内核兼容模式
	33: "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.3 (KHTML, like Gecko) Chrome/6.0.472.33 Safari/534.3 SE 2.X MetaSr 1.0",                                                                            //搜狗3.0在Win7+ie9,高速模式
	34: "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; InfoPath.3; .NET4.0C; .NET4.0E)",                    //360浏览器3.0在Win7+ie9
	35: "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/535.1 (KHTML, like Gecko) Chrome/13.0.782.41 Safari/535.1 QQBrowser/6.9.11079.201",                                                                                        //QQ浏览器6.9(11079)在Win7+ie9,极速模式
}

// 每种类型设备的useragent的列表
var (
	listPCAgent           = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35}
	listWindowsAgent      = []int{1, 2, 3, 4, 5, 30, 31, 32, 33, 34, 35}
	listLinuxAgent        = []int{6, 7}
	listMacAgent          = []int{8, 9, 10, 11}
	listAndroidAgent      = []int{12, 13, 14, 15, 16, 17, 18}
	listIosAgent          = []int{19, 20, 21, 22, 23}
	listPhoneAgent        = []int{12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	listWindowsPhoneAgent = []int{24, 25}
	listUCAgent           = []int{26, 27, 28, 29}
)

var AgentType = map[UserAgentType][]int{
	PCAgent:           listPCAgent,
	WindowsAgent:      listWindowsAgent,
	LinuxAgent:        listLinuxAgent,
	MacAgent:          listMacAgent,
	AndroidAgent:      listAndroidAgent,
	IosAgent:          listIosAgent,
	PhoneAgent:        listPhoneAgent,
	WindowsPhoneAgent: listWindowsPhoneAgent,
	UCAgent:           listUCAgent,
}

// GetAgent 随机获取那种类型的 user-agent
func GetAgent(agentType UserAgentType) string {
	switch agentType {
	case PCAgent:
		if v, ok := UserAgentMap[listPCAgent[rand.Intn(len(listPCAgent))]]; ok {
			return v
		}
	case WindowsAgent:
		if v, ok := UserAgentMap[listWindowsAgent[rand.Intn(len(listWindowsAgent))]]; ok {
			return v
		}
	case LinuxAgent:
		if v, ok := UserAgentMap[listLinuxAgent[rand.Intn(len(listLinuxAgent))]]; ok {
			return v
		}
	case MacAgent:
		if v, ok := UserAgentMap[listMacAgent[rand.Intn(len(listMacAgent))]]; ok {
			return v
		}
	case AndroidAgent:
		if v, ok := UserAgentMap[listAndroidAgent[rand.Intn(len(listAndroidAgent))]]; ok {
			return v
		}
	case IosAgent:
		if v, ok := UserAgentMap[listIosAgent[rand.Intn(len(listIosAgent))]]; ok {
			return v
		}
	case PhoneAgent:
		if v, ok := UserAgentMap[listPhoneAgent[rand.Intn(len(listPhoneAgent))]]; ok {
			return v
		}
	case WindowsPhoneAgent:
		if v, ok := UserAgentMap[listWindowsPhoneAgent[rand.Intn(len(listWindowsPhoneAgent))]]; ok {
			return v
		}
	case UCAgent:
		if v, ok := UserAgentMap[listUCAgent[rand.Intn(len(listUCAgent))]]; ok {
			return v
		}
	default:
		if v, ok := UserAgentMap[rand.Intn(len(UserAgentMap))]; ok {
			return v
		}
	}
	return ""
}

func SetAgent(agentType UserAgentType, agent string) {
	UserAgentMap[len(UserAgentMap)] = agent
	AgentType[agentType] = append(AgentType[agentType], len(UserAgentMap))
}

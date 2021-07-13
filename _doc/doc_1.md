# gathertool 开发文档

## 1. HTTP/S 请求所有参数说明与上下文结构
> gathertool 是一个请求工具，用于数据抓取，接口测试； 可在请求阶段执行各个事件，根据状态码自定义事件等； 可扩展性，
>传入任意自定义http请求参数与方法， 可以适用于各种代理； 还有很多创新的地方文档会根据函数与参数的说明来介绍 gathertool
>的创新；
- StartFunc(func (ctx *Context)) ： 请求前执行的事件函数类型；
- SucceedFunc(func (ctx *Context)) ： 请求成功后的事件函数类型；
- FailedFunc(func (ctx *Context))  ： 请求失败后的事件函数类型, 请求错误与默认状态码（参见默认状态码事件）会触发； 
- RetryFunc(func (ctx *Context))  ： 请求重试前的事件函数类型, 默认状态码（参见默认状态码事件）会触发， 可以在此事件更换代理，添加等待时间等等， 重试次数默认是10次，可自行设置； 
- EndFunc(func (ctx *Context)) ： 请求结束后的事件函数类型； 
```、go
import gt "github.com/mangenotwork/gathertool"

gt.Get(`http://192.168.0.1`, 
        gt.StartFunc(func(ctx *Context){
            log.Println("请求前： 添加代理等等操作")
            ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }),
        gt.SucceedFunc(func(ctx *Context){
            log.Println("请求成功： 处理数据或存储等等")
            log.Println(ctx.RespBodyString())
        }),
        gt.FailedFunc(func(ctx *Context){
            log.Println("请求失败： 记录失败等等")
            log.Println(ctx.Err)
        }),
        gt.RetryFunc(func(ctx *Context){
             log.Println("请求重试： 更换代理或添加等待时间等等")
             ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }),
        gt.RetryFunc(func(ctx *Context){
             log.Println("请求结束： 记录结束，处理数据等等")
             log.Println(ctx)
        }),
)
```
- *http.Client 
- *http.Request 
- *http.Response
-  Err ： 记录的错误信息
- context.Context
- RetryTimes(n int) : 设置重试次数
```go
gt.Get(`http://192.168.0.1`, RetryTimes(50))
```
- *Task : 请求的任务对象， 这个任务可被多个请求使用,（见Task任务对象）
```go
gt.Get(`http://192.168.0.1`, &gt.Task{})
```
- IsLog : 是否打印日志, 默认打开
- ProxyUrl(url string) : 设置代理
```go
gt.Get(`http://192.168.0.1`, gt.ProxyUrl("http://192.168.0.2:8888"))
```

## 2. 状态码事件 与 UserAgent
> 状态码对应事件的全局的，可初始化设置，也可随时重置
### 默认状态码事件
|状态码|事件类型|事件描述|
| :-----| ----: | :----: |
|200|success|请求成功后事件|
|201|success|请求成功后事件|
|202|success|请求成功后事件|
|203|success|请求成功后事件|
|204|fail|请求失败后事件|
|300|success|请求成功后事件|
|301|success|请求成功后事件|
|302|success|请求成功后事件|
|400|fail|请求失败后事件|
|401|retry|请求重试前的事件|
|402|retry|请求重试前的事件|
|403|retry|请求重试前的事件|
|404|fail|请求失败后事件|
|405|retry|请求重试前的事件|
|406|retry|请求重试前的事件|
|407|retry|请求重试前的事件|
|408|retry|请求重试前的事件|
|500|fail|请求失败后事件|
|501|fail|请求失败后事件|
|502|retry|请求重试前的事件|
|503|retry|请求重试前的事件|
|504|retry|请求重试前的事件|

### 设置状态码事件
- func SetStatusCodeSuccessEvent(code int) ： 将指定状态码设置为执行成功事件
- func SetStatusCodeRetryEvent(code int) ： 将指定状态码设置为执行重试事件
- func SetStatusCodeFailEvent(code int) ： 将指定状态码设置为执行失败事件

### UserAgent
> gathertool 有UserAgent Map 是全局的，可自定义，可扩展等
> 所有请求默认使用 PCAgent随机的一个

|类型|值|描述|
| :-----| ----: | :----: |
| PCAgent| 1 | pc端的useragent|
|WindowsAgent|2|Windows useragent|
|LinuxAgent|3|Linux useragent|
|MacAgent|4|Mac useragent|
|AndroidAgent|5|Android useragent|
|IosAgent|6|Ios useragent|
|PhoneAgent|7|Phone useragent|
|WindowsPhoneAgent|8|WindowsPhone useragent|
|UCAgent|9|UC useragent|

#### func GetAgent(agentType UserAgentType) string ：  随机获取 user-agent
#### func SetAgent(agentType UserAgentType, agent string) :  设置 user-agent

## 3. 基础请求
#### func Get(url string) (*Context, error)  ： get请求

```go
// 无事件参数的请求
ctx, err := gathertool.Get(`http://192.168.0.1`)
log.Println(ctx.RespBodyString(), err)

// 有事件参数的请求
gathertool.Get(`http://192.168.0.1`,SucceedFunc(succeed))
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### func NewGet(url string) *Context ： 新建一个get请求上下文
```go
gt.NewGet(`http://192.168.0.1`).SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### func Post(url string, data []byte, contentType string) (*Context, error) ： Post请求
```go
ctx, err := gt.Post(`https://httpbin.org/post`, []byte(`{"a":"a"}`), "application/json;")
log.Println(ctx.RespBodyString(), err)
```

#### func PostJson(url string, jsonStr string) (*Context, error) : Post请求 - json参数
```go
ctx, err := gt.PostJson(`https://httpbin.org/post`, `{"a":"a"}`)
log.Println(ctx.RespBodyString(), err)
```

#### func PostForm(url string, data url.Values) (*Context, error)  : Post请求 - Form 
```go
ctx, err := gt.PostForm(`https://httpbin.org/post`, url.Values(`a=123`))
log.Println(ctx.RespBodyString(), err)
```


#### func NewPost(url string, data []byte, contentType string) *Context  ： 新建一个post请求上下文
```go
gt.NewPost(`https://httpbin.org/post`, []byte(`{"a":"a"}`), "application/json;").SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### func NewPut(url string, data []byte, contentType string) *Context  : 新建一个put请求上下文
```go
gt.NewPut(`https://httpbin.org/put`, []byte(`{"a":"a"}`), "application/json;").SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### func Put(url string, data []byte, contentType string) (*Context, error)  : put请求
```go
ctx, err := gt.Put(`https://httpbin.org/put`, url.Values(`a=123`))
log.Println(ctx.RespBodyString(), err)
```

#### func NewDelete(url string) *Context   :  新建一个delete请求上下文
```go
gt.NewDelete(`https://httpbin.org/delete`).SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### func Delete(url string) (*Context, error)  : delete 请求
```go
ctx, err := gt.Delete(`https://httpbin.org/delete`)
log.Println(ctx.RespBodyString(), err)
```

#### func NewOptions(url string) *Context  :  新建一个options请求上下文
```go
gt.NewOptions(`https://httpbin.org/options`).SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### func Options(url string) (*Context, error)   ： options请求
```go
ctx, err := gt.Options(`https://httpbin.org/options`)
log.Println(ctx.RespBodyString(), err)
```

#### func Request(url, method string, data []byte, contentType string) (*Context, error)  ：  一个请求

#### func NewRequest(url, method string, data []byte, contentType string) *Context  ：  新建一个请求上下文

#### func Upload(url, savePath string) (*Context, error) : 下载
```go
ctx, err := gt.Upload("https://jyzd.bfsu.edu.cn/uploadfile/bfsu/front/default/upload_file_35256.pdf", "/home/mange/Desktop/upload_file_3.pdf")
log.Println(ctx.Resp.StatusCode, err)
```

#### func Req(request *http.Request) *Context  ：  自定义  *http.Request 的请求

#### func WsClient(host, path string) (WSClient, error)  :  websocket请求
```go
host := "wss://xxx.api.xxx.cn"
path := "/v1/ws?user_id=xxxx&device_id=xxxx&token=xxxxx&app_id=xxxx"
wsc, err := gt.WsClient(host, path)
if err != nil {
    log.Println(err)
    return
}
for {
	err = wsc.Send([]byte(`{"type":3,"request_id":"1"}`))
	if err != nil {
            log.Println(err)
            break
        }
	data := make([]byte,100)
	err = wsc.Read(data)
	if err != nil {
		log.Println(err)
		break
             }
	log.Println("data = ", string(data))
	time.Sleep(1*time.Second)
}
```

#### func SSHClient(user string, pass string, addr string) (*ssh.Client, error)  :  ssh 请求
```go
user = "root"
password = "root123"
addr = "192.168.0.2:22" // ip+port
sshc, err := gt.SSHClient(user, password, addr)
```

#### [TODO]  Tcp 请求

#### [TODO]  Udf 请求

#### [TODO]  FTP 请求

#### [TODO]  Smtp | Pop3 请求

#### [TODO]  MQTT 请求

#### [TODO]  grpc 等等

#### type  RetryTimes int : 重试次数
     
#### type StartFunc func(c *Context)    : 请求开始前的方法类型
     
#### type SucceedFunc func(c *Context)    :  成功后的方法类型
     
#### type FailedFunc func(c *Context)   : 失败后的方法类型  

#### type RetryFunc func(c *Context)   :  重试前的方法类型

#### type EndFunc func(c *Context)   :  请求结束后的方法类型
     
#### type IsLog bool   :   是否开启日志
     
#### type ProxyUrl string  :   代理地址

#### type Context struct :  请求上下文

#### func (*Context) SetSucceedFunc(successFunc func(c *Context)) *Context : 设置成功后的执行方法
```go
gt.NewGet(`http://192.168.0.1`).SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### func (*Context) SetFailedFunc(failedFunc func(c *Context)) *Context : 设置错误后的方法
```go
gt.NewGet(`http://192.168.0.1`).SetFailedFunc(failed).Do()
func failed(ctx *gt.Context) {
    log.Println(ctx.Err)
}
```

#### func (*Context) SetRetryFunc(retryFunc func(c *Context)) *Context : 设置重试，在重试前的方法
```go
gt.NewGet(`http://192.168.0.1`).SetRetryFunc(retry).Do()
func retry(ctx *gt.Context) {
    //更换代理
}
```

#### func (*Context) SetRetryTimes(times int) *Context  ：  设置重试次数
```go
gt.NewGet(`http://192.168.0.1`).SetRetryFunc(retry).SetRetryTimes(3).Do()
func retry(ctx *gt.Context) {
    //更换代理
}
```

#### func (*Context) Do()  ： 执行请求

#### func (*Context) RespBodyString() string  :  请求结果输出字符串
```go
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString())
}
```

#### func (*Context) CheckReqMd5() string :  请求的唯一md5, 请求Url+参数+methd； 使用场景确认请求的唯一性

#### func (*Context) CheckMd5() string  ：  请求的唯一md5, 请求Url+参数+methd+返回结果； 使用场景确认请求的唯一性

#### func (*Context) AddHeader(k,v string) :  add header
```go
gt.NewGet(`http://192.168.0.1`).AddHeader("token","gathertool").Do()
```

#### func (*Context) AddCookie(k, v string) :  add Cookie
```go
gt.NewGet(`http://192.168.0.1`).AddCookie("SUBP", "0033WrSXqPxfM72-Ws9jqgMF55529P9D9WWENAjmKyIZz1AWjDi68mRw").Do()
```

#### func (*Context) SetProxy(proxyUrl string) :  set proxy
```go
gt.NewGet(`http://192.168.0.1`).SetProxy("http://192.168.0.2:8888").Do()
```

#### func (*Context) Upload(filePath string)  :  Upload 下载
```go
gt.NewGet(`http://192.168.0.3:8081/img.png`).Upload("/home/Desktop/img.png").Do()
```   

#### func (*Context) CloseLog() : close log 关闭日志
```go
gt.NewGet(`http://192.168.0.1`).CloseLog().Do()
``` 

#### func (*Context) OpenErr2Retry() : 开启请求失败执行retry
```go
gt.NewGet(`http://192.168.0.1`).OpenErr2Retry().SetRetryFunc(retry).Do()
``` 

#### func (*Context) CloseRetry() : 关闭重试
     
#### var CookiePool  cookie 池

#### CookiePool.Add(cookie *http.Cookie)  添加

#### CookiePool.Get() *http.Cookie   随机获取

#### type Task struct   任务对象
    
#### func CrawlerTask(url, jsonParam string, vs ...interface{}) *Task  ： 创建一个发送任务


## 4. 队列与高并发请求
> gathertool 采用本地队列加上多个goroutine实现高并发， 每个并发函数默认开启了runtime.GOMAXPROCS(runtime.NumCPU())会充分利用每个cpu核心；

#### func NewQueue() TodoQueue : 新建一个队列

#### func NewUploadQueue() TodoQueue  新建一个下载队列

#### type TodoQueue interface 队列接口

#### func (*Queue) Add(task *Task) error   添加队列

#### func (*Queue) Poll() *Task  移除队列中最前面的额元素

#### func (*Queue) Clear() bool   清除队列

```go
var queue = gt.NewQueue() //全局声明抓取任务队列
html := string(cxt.RespBody)
	dom,err := gt.NewGoquery(html)
	if err != nil{
		log.Println(err)
		return
	}
	result := dom.Find("div[id=result] tbody")
	result.Find("tr").Each(func(i int, tr *goquery.Selection){
		td := tr.Find("td")
		startIp := td.Eq(0).Text()// IP起始
		endIP := td.Eq(1).Text()// 结束
		number := td.Eq(2).Text()// 数量
		// 创建队列 抓取详情信息
		queue.Add(&gt.Task{
			Url: "http://ip.bczs.net/"+startIp,
			Data: map[string]interface{}{
				"start_ip":startIp,
				"end_ip":endIP,
				"number":number,
			},
		})
	})
```

#### func StartJobGet(jobNumber int, queue TodoQueue, vs ...interface{})  并发执行Get,直到队列任务为空
     // @jobNumber 并发数，
     // @queue 全局队列，
     // @client 单个并发任务的client，
     // @SucceedFunc 成功方法，
     // @ RetryFunc重试方法，
     // @FailedFunc 失败方法
```go
gt.StartJobGet(100,queue,
	gt.SucceedFunc(GetIPSucceed),//请求成功后执行的方法
	gt.RetryFunc(GetIPRetry),//遇到 502,403 等状态码重试前执行的方法，一般为添加休眠时间或更换代理
	gt.FailedFunc(GetIPFailed),//请求失败后执行的方法
)
```

#### func StartJobPost(jobNumber int, queue TodoQueue, vs ...interface{}) 开始运行并发Post,直到队列任务为空

## 5. 压力测试

#### func NewTestUrl(url, method string, sum int64, total int) *StressUrl ： 实例化一个新的url压测
    // @method :请求类型
    // @sum  : 总请求次数
    // @total : 并发数
```go
url := "http://192.168.0.9:18084/static_service/v1/allow/school/page"
test := gt.NewTestUrl(url,"Get",1000,500)
test.Run()
```

#### type StressUrl struct  ： 压力测试一个url

#### func (*StressUrl) Run(vs ...interface{})  ： 运行压测， vs接收请求的参数，如header 等等
```go
url4 := "http://ggzyjy.sc.gov.cn/WebBuilder/frontAppAction.action?cmd=addPageView"
test4 := gt.NewTestUrl(url4,"Post",100,10)
test4.SetJson(`{
				"viewGuid":"cms_002",
				"siteGuid":"7eb5f7f1-9041-43ad-8e13-8fcb82ea831a"
				}`)
test4.Run(gt.SucceedFunc(func(c *gt.Context) {
	log.Println(string(c.RespBody))
	log.Println(c.Resp.Cookies())
}))
```
     
#### func (*StressUrl) SetJson(str string)  ： 设置压测请求的json参数

     
## 6. 正则提取

#### func RegHtmlA(str string)[]string 

#### func RegHtmlTitle(str string)[]string

#### func RegHtmlTr(str string)[]string

#### func RegHtmlInput(str string) []string

#### func RegHtmlTd(str string) []string

#### func RegHtmlP(str string) []string

#### func RegHtmlSpan(str string) []string

#### func RegHtmlSrc(str string) []string

#### func RegHtmlHref(str string) []string

#### func RegHtmlATxt(str string)[]string 

#### func RegHtmlTitleTxt(str string)[]string 

#### func RegHtmlTrTxt(str string)[]string 

#### func RegHtmlInputTxt(str string) []string 

#### func RegHtmlTdTxt(str string) []string 

#### func RegHtmlPTxt(str string) []string 

#### func RegHtmlSpanTxt(str string) []string 

#### func RegHtmlSrcTxt(str string) []string 

#### func RegHtmlHrefTxt(str string) []string 


## 7. 常用函数

#### func StringValue(i interface{}) string 任何类型返回值字符串形式

#### func MD5(str string) string 

#### func Json2Map(str string) map[string]interface{}  json转map函数，通用

#### func Any2Map(data interface{}) map[string]interface{}  ：  interface{} -> map[string]interface{} 

#### func Any2String(data interface{}) string  ：  interface{} -> string
     
#### func Any2Int(data interface{}) int  ：  interface{} -> int

#### func Any2int64(data interface{}) int64  ：  interface{} -> int64
     
#### func Any2AnyList(data interface{}) []interface{}   ：  interface{} -> []interface{}

#### func Any2Float64(data interface{}) float64  ：  interface{} -> float64
     
#### func Any2Strings(data interface{}) []string    ：  interface{} -> []string
         
#### func CleaningStr(str string) string  ：  清理字符串前后空白 和回车 换行符号
     
#### func StrDeleteSpace(str string) string   ：  StrDeleteSpace 删除字符串前后的空格
          
#### func Str2Int64(str string) int64     ：   string -> int64
     
#### func Str2Float64(str string) float64  ：   string -> float64
              
#### func Uint82Str(bs []uint8) string    ：   []uint8 -> string

#### func FileSizeFormat(fileSize int64) (size string)  ：  字节的单位转换 保留两位小数
     
#### func Struct2Map(obj interface{}, hasValue bool) (map[string]interface{}, error)   ：   Struct  ->  map
     // hasValue=true表示字段值不管是否存在都转换成map
     // hasValue=false表示字段为空或者不为0则转换成map
     
#### func PanicToError(fn func()) (err error)    ：   panic -> error
     
#### func ConvertByte2String(byte []byte, charset Charset) string   ：   编码转换
     charset:
        UTF8 
        GB18030 
        GB18030 
        GBK 
        GB2312 

#### func UnescapeUnicode(raw []byte) ([]byte, error)   :  Unicode 转码

#### func Base64Encode(str string) string   :   base64 编码
     
#### func Base64Decode(str string) (string,error)  :  base64 解码

#### func Base64UrlEncode(str string) string  :  base64 url 编码
     
#### func Base64UrlDecode(str string) (string,error)  :  base64 url 解码

#### func IsContainStr(items []string, item string) bool  字符串是否等于items中的某个元素
     
#### func FileMd5(path string) (string, error)   文件md5

#### func Byte2Str(b []byte) string ：  []byte -> string
     
#### func ByteToBinaryString(data byte) (str string)  ：   字节 -> 二进制字符串
     
#### func StrsDuplicates(a []string) []string  ：  数组，切片去重和去空串

#### func IsElementStr(list []string, element string) bool  ：  判断字符串是否与数组里的某个字符串相同
     
#### func Timestamp() string    ：  当前 Timestamp

#### func BeginDayUnix() int64   :  获取当天 0点时间戳

#### func EndDayUnix() int64  ：  获取当天 24点时间戳

#### func Daydiff(beginDay string, endDay string) int  ： 两个时间字符串的日期差时间戳

#### func TickerRun(t time.Duration, runFirst bool, f func())  ：  间隔运行
     // t: 间隔时间，  f: 运行的方法

## 8. redis
> redis 的方法使用 github.com/garyburd/redigo/redis

#### func NewRedis(host, port, password string, db int, vs ...interface{}) (*Rds)   ： 实例化redis对象

#### func NewRedisPool(host, port, password string, db, maxIdle, maxActive, idleTimeoutSec int, vs ...interface{}) (*Rds)  ： 实例化redis对象

#### type Rds struct   ：  redis对象

#### func (*Rds) RedisConn() (err error)  ： redis 连接

#### func (*Rds) RedisPool() error  :  redis 连接池连接 *redis.Pool.Get() 获取redis连接
     
#### func (*Rds) SelectDB(dbNumber int) error  ：  redis 切换db

#### func RedisDELKeys(rds *Rds, keys string, jobNumber int)  : 高并发删除key
    jobNumber 并发数
```go
rds := gt.NewRedisPool(redis_host, redis_port, redis_password, dbnumber, 5, 10, 10,
		gt.NewSSHInfo(ssh_addr, ssh_user, ssh_password))                                                       
gt.RedisDELKeys(rds, "in:*", 1000)

```

## 9. mysql
> mysql 的方法使用 github.com/go-sql-driver/mysql


## 10. mongo
> mongo 的方法使用 go.mongodb.org/mongo-driver

## 11. github.com/PuerkitoBio/goquery
> 提取数据推荐使用 github.com/PuerkitoBio/goquery

## 
    
     


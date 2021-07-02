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
- SetStatusCodeSuccessEvent(code int) ： 将指定状态码设置为执行成功事件
- SetStatusCodeRetryEvent(code int) ： 将指定状态码设置为执行重试事件
- SetStatusCodeFailEvent(code int) ： 将指定状态码设置为执行失败事件

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

#### GetAgent(agentType UserAgentType) string ：  随机获取 user-agent
#### SetAgent(agentType UserAgentType, agent string) :  设置 user-agent

## 3. 基础请求
#### Get(url string) (*Context, error)  ： get请求

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

#### NewGet(url string) *Context ： 新建一个get请求上下文
```go
gt.NewGet(`http://192.168.0.1`).SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### Post(url string, data []byte, contentType string) (*Context, error) ： Post请求
```go
ctx, err := gt.Post(`https://httpbin.org/post`, []byte(`{"a":"a"}`), "application/json;")
log.Println(ctx.RespBodyString(), err)
```

#### PostJson(url string, jsonStr string) (*Context, error) : Post请求 - json参数
```go
ctx, err := gt.PostJson(`https://httpbin.org/post`, `{"a":"a"}`)
log.Println(ctx.RespBodyString(), err)
```

#### PostForm(url string, data url.Values) (*Context, error)  : Post请求 - Form 
```go
ctx, err := gt.PostForm(`https://httpbin.org/post`, url.Values(`a=123`))
log.Println(ctx.RespBodyString(), err)
```


#### NewPost(url string, data []byte, contentType string) *Context  ： 新建一个post请求上下文
```go
gt.NewPost(`https://httpbin.org/post`, []byte(`{"a":"a"}`), "application/json;").SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### NewPut(url string, data []byte, contentType string) *Context  : 新建一个put请求上下文
```go
gt.NewPut(`https://httpbin.org/put`, []byte(`{"a":"a"}`), "application/json;").SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### Put(url string, data []byte, contentType string) (*Context, error)  : put请求
```go
ctx, err := gt.Put(`https://httpbin.org/put`, url.Values(`a=123`))
log.Println(ctx.RespBodyString(), err)
```

#### NewDelete(url string) *Context   :  新建一个delete请求上下文
```go
gt.NewDelete(`https://httpbin.org/delete`).SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### Delete(url string) (*Context, error)  : delete 请求
```go
ctx, err := gt.Delete(`https://httpbin.org/delete`)
log.Println(ctx.RespBodyString(), err)
```

#### NewOptions(url string) *Context  :  新建一个options请求上下文
```go
gt.NewOptions(`https://httpbin.org/options`).SetSucceedFunc(succeed).Do()
func succeed(ctx *gt.Context) {
    log.Println(ctx.RespBodyString(), ctx)
}
```

#### Options(url string) (*Context, error)   ： options请求
```go
ctx, err := gt.Options(`https://httpbin.org/options`)
log.Println(ctx.RespBodyString(), err)
```

#### Request(url, method string, data []byte, contentType string) (*Context, error)  ：  一个请求

#### NewRequest(url, method string, data []byte, contentType string) *Context  ：  新建一个请求上下文

#### Upload(url, savePath string) (*Context, error) : 下载
```go
ctx, err := gt.Upload("https://jyzd.bfsu.edu.cn/uploadfile/bfsu/front/default/upload_file_35256.pdf", "/home/mange/Desktop/upload_file_3.pdf")
log.Println(ctx.Resp.StatusCode, err)
```

#### Req(request *http.Request) *Context  ：  自定义  *http.Request 的请求

#### WsClient(host, path string) (WSClient, error)  :  websocket请求
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

#### SSHClient(user string, pass string, addr string) (*ssh.Client, error)  :  ssh 请求
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

## 4. 队列与高并发请求
> gathertool 采用本地队列加上多个goroutine实现高并发， 每个并发函数默认开启了runtime.GOMAXPROCS(runtime.NumCPU())会充分利用每个cpu核心；


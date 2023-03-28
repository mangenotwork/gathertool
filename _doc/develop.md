# 《gathertool开发使用文档》

Date : 2023-03-28

Author : ManGe

Mail : 2912882908@qq.com

Github : https://github.com/mangenotwork/gathertool

## 一、介绍

### 1.1 简介
> gathertool是golang脚本化开发集成库，目的是提高对应场景脚本程序开发的效率；
> 
> gathertool也是一款轻量级爬虫库，特色是分离了请求事件，通俗点理解就是对请求过程状态进行事件处理。
> 
> gathertool也是接口测试&压力测试库，在接口测试脚本开发上有明显的效率优势，
> 
> gathertool还集成了对三方中间件的操作，DB操作等。

---

### 1.2 使用
> import "github.com/mangenotwork/gathertool"
> 
> go get github.com/mangenotwork/gathertool

---

### 1.3 介绍
> gathertool是一个高度封装工具库，包含了http/s的请求，Mysql数据库方法，数据类型处理方法，数据提取方法，websocket相关方法，
> TCP|UDP相关方法，NoSql相关方法，开发常用方法等;  可以用于爬虫程序，接口&压力测试程序，常见网络协议调试程序，数据提取与存储程序等；
> gathertool的请求特点: 会在请求阶段执行各个事件如请求失败后的重试事件,请求前后的事件，请求成功事件等等, 可以根据请求状态码自定义这些事件；
> gathertool还拥有很好的可扩展性， 适配传入任意自定义http请求对象， 能适配各种代理对象等等；
> gathertool还拥有抓取数据存储功能, 比如存储到mysql, redis, mongo, pgsql等等; 还有很多创新的地方文档会根据具体方法进行介绍；
> gathertool还封装了消息队列接口，支持Nsq,Kafka,rabbitmq,redis等消息队列中间件

---

### 1.4 使用场景
1. 爬虫脚本开发
2. 接口测试&压力测试脚本开发
3. http/s代理服务器, socket5代理服务器
4. mysql相关操作方法
5. redis相关操作方法
6. mongo相关操作方法
7. 数据提取&清洗相关操作
8. Websocket相关操作方法
9. TCP客户端
10. UDP客户端
11. SSH客户端
12. 加密解密脚本开发
13. ip扫描，端口扫描脚本开发
14. 消息队列

---

### 1.5 简单例子

#### 简单的get请求

```golang
import gt "github.com/mangenotwork/gathertool"

func main(){
    ctx, err := gt.Get("https://www.baidu.com")
    if err != nil {
        log.Println(err)
    }
    log.Println(ctx.Html)
}
```

#### 含请求事件请求

```golang
import gt "github.com/mangenotwork/gathertool"

func main(){
   ctx := gt.NewGet(`http://192.168.0.1`)

   ctx.SetStartFunc(func(ctx *gt.Context){
            log.Println("请求前： 添加代理等等操作")
            ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }
    )

   ctx.SetSucceedFunc(func(ctx *gt.Context){
            log.Println("请求成功： 处理数据或存储等等")
            log.Println(ctx.RespBodyString())
        }
    )

    ctx.SetFailedFunc(func(ctx *gt.Context){
            log.Println("请求失败： 记录失败等等")
            log.Println(ctx.Err)
        }
    )

    ctx.SetRetryFunc(func(ctx *gt.Context){
             log.Println("请求重试： 更换代理或添加等待时间等等")
             ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }
    )

    ctx.SetEndFunc(func(ctx *gt.Context){
             log.Println("请求结束： 记录结束，处理数据等等")
             log.Println(ctx)
        }
    )

    ctx.Do()
    log.Println(ctx.Html)
}
```

#### 事件方法复用

```golang
func main(){
    gt.NewGet(`http://192.168.0.1`).SetSucceedFunc(succeed).SetFailedFunc(failed).SetRetryFunc(retry).Do()
    gt.NewGet(`http://www.baidu.com`).SetSucceedFunc(baiduSucceed).SetFailedFunc(failed).SetRetryFunc(retry).Do()
}

// 请求成功： 处理数据或存储等等
func succeed(ctx *gt.Context){
    log.Println(ctx.RespBodyString())
}

// 请求失败： 记录失败等等
func failed(ctx *gt.Context){
    log.Println(ctx.Err)
}

// 请求重试： 更换代理或添加等待时间等等
func retry(ctx *gt.Context){
    ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
}

// 百度请求成功后
func baiduSucceed(ctx *gt.Context){
    log.Println(ctx.RespBodyString())
}
```

#### post请求

```golang
 // FormData
    postData := gt.FormData{
        "aa":"aa",	
    }
    
    // header
    header := gt.NewHeader(map[string]string{
        "Accept":"*/*",
        "X-MicrosoftAjax":"Delta=true",
        "Accept-Encoding":"gzip, deflate",
        "XMLHttpRequest":"XMLHttpRequest",
        "Content-Type":"application/x-www-form-urlencoded; charset=UTF-8",
    })
    
    // cookie
    cookie := gt.Cookie{
        "aa":"a"
    }
    
    // 随机休眠 2~6秒
    sleep := gt.SetSleep(2,6)
    c := gt.NewPostForm(caseUrl, postData, header, cookie, sleep)
    c.Do()
    html := c.RespBodyString()
    log.Print(html)

```


#### 数据存储到mysql

```golang
var (
    host   = "192.168.0.100"
    port      = 3306
    user      = "root"
    password  = "root123"
    dbName  = "dbName"
    db,_ = gt.NewMysql(host, port, user, password, dbName)
)

//.... 执行抓取
data1 := "data1"
data2 := "data2"

inputdata := map[string]interface{} {
    "data1" : data1,
    "data2" : data2,
}

tableName := "data"
db.Spider2022DB.InsertAt(tableName, inputdata)
```


#### HTML数据提取
```
func main(){
	date := "2022-07-05"
	caseUrl := "***"
	ctx, _ := gt.Get(fmt.Sprintf(caseUrl, date))
    datas, _ := gt.GetPointHTML(ctx.Html, "div", "id", "domestic")
	Data(datas, date, "内期表", "备注：内期表=国内期货主力合约表")
	datas, _ = gt.GetPointHTML(ctx.Html, "div", "id", "overseas")
	Data(datas, date, "外期表", "备注：外期表=国外期货主力合约表")
}

func Data(datas []string, date, typeName, note string) {
	for _, data := range datas {
		table, _ := gt.GetPointHTML(data, "table", "id", "fdata")
		if len(table) > 0 {
			trList := gt.RegHtmlTr(table[0])
			jys := ""
			for _, tr := range trList {
				td := gt.RegHtmlTd(tr)
				log.Println("td = ", td, len(td))
				if len(td) == 1 {
					jys = gt.RegHtmlTdTxt(td[0])[0]
					continue
				}
				name := gt.RegHtmlTdTxt(td[0])[0]
				if strings.Index(name, "商品名称") != -1 {
					continue
				}
				zlhy := gt.RegHtmlTdTxt(td[1])[0]
				jsj := gt.RegHtmlTdTxt(td[2])[0]
				zd := gt.RegDelHtml(gt.RegHtmlTdTxt(td[3])[0])
				cjj := gt.RegHtmlTdTxt(td[4])[0]
				ccl := gt.RegHtmlTdTxt(td[5])[0]
				dw := gt.RegHtmlTdTxt(td[6])[0]
				log.Println("日期 = ", date)
				log.Println("机构 = ", jys)
				log.Println("商品名称 = ", name)
				log.Println("主力合约 = ", zlhy)
				log.Println("结算价 = ", jsj)
				log.Println("涨跌 = ", zd)
				log.Println("成交量 = ", cjj)
				log.Println("持仓量 = ", ccl)
				log.Println("单位 = ", dw)
			}
		}
	}
}
```

#### Json数据提取
```
func main(){
	txt := `{
    "reason":"查询成功!",
    "result":{
        "city":"苏州",
        "realtime":{
            "temperature":"17",
            "humidity":"69",
            "info":"阴",
            "wid":"02",
            "direct":"东风",
            "power":"2级",
            "aqi":"30"
        },
        "future":[
            {
                "date":"2021-10-25",
                "temperature":"12\/21℃",
                "weather":"多云",
                "wid":{
                    "day":"01",
                    "night":"01"
                },
                "direct":"东风"
            },
            {
                "date":"2021-10-26",
                "temperature":"13\/21℃",
                "weather":"多云",
                "wid":{
                    "day":"01",
                    "night":"01"
                },
                "direct":"东风转东北风"
            },
            {
                "date":"2021-10-27",
                "temperature":"13\/22℃",
                "weather":"多云",
                "wid":{
                    "day":"01",
                    "night":"01"
                },
                "direct":"东北风"
            }
        ]
    },
    "error_code":0
}`

	jx1 := "/result/future/[0]/date"
	jx2 := "/result/future/[0]"
	jx3 := "/result/future"

	log.Println(gt.JsonFind(txt, jx1))
	log.Println(gt.JsonFind2Json(txt, jx2))
	log.Println(gt.JsonFind2Json(txt, jx3))
	log.Println(gt.JsonFind2Map(txt, jx2))
	log.Println(gt.JsonFind2Arr(txt, jx3))

}
```

---

### 1.6 实例
-  [Get请求](https://github.com/mangenotwork/gathertool/tree/main/_examples/get)
-  [阳光高考招生章程抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/get_yggk)
-  [ip地址信息抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/ip_bczs_cn)
-  [压力测试](https://github.com/mangenotwork/gathertool/tree/main/_examples/stress_testing)
-  [文件下载](https://github.com/mangenotwork/gathertool/tree/main/_examples/upload_file)
-  [无登录微博抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/weibo)
-  [百度题库抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/baidu_tk)
-  [搭建http/s代理与抓包](https://github.com/mangenotwork/gathertool/tree/main/_examples/intercept)
-  [搭建socket5代理](https://github.com/mangenotwork/gathertool/tree/main/_examples/socket5)
-  [商品报价信息抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/baojia)
-  [期货信息抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/qihuo)
-  [行业信息网行业分类](https://github.com/mangenotwork/gathertool/tree/main/_examples/cnlinfo)

---

## 二、请求

### 2.1 请求事件

gathertool的特色之一就是纳入了请求事件，一个请求拥有请求前，请求后，根据指定状态码触发请求成功事件，请求失败事件，请求重试事件；每一种事件都可以去完成自定义的业务，可以进行随意搭配使用。

---
#### StartFunc(func (ctx \*Context))  
请求前执行的事件函数类型；

---
#### SucceedFunc(func (ctx \*Context)) 
请求成功后的事件函数类型；请求错误与默认状态码（参见默认状态码事件）会触发；

---
#### FailedFunc(func (ctx \*Context))  
请求失败后的事件函数类型, 请求错误与默认状态码（参见默认状态码事件）会触发；

---
#### RetryFunc(func (ctx \*Context))   
请求重试前的事件函数类型, 默认状态码（参见默认状态码事件）会触发， 可以在此事件更换代理，添加等待时间等等， 重试次数默认是10次，可自行设置；

---
#### EndFunc(func (ctx \*Context)) 
请求结束后的事件函数类型；

---
#### func (c \*Context) SetSucceedFunc(successFunc func(c \*Context)) \*Context 
设置成功后的方法

---
#### func (c \*Context) SetFailedFunc(failedFunc func(c \*Context)) \*Context 
设置错误后的方法

---
#### func (c \*Context) SetRetryFunc(retryFunc func(c \*Context)) \*Context 
设置重试，在重试前的方法

---

例子1：

```golang
import gt "github.com/mangenotwork/gathertool"

func main(){
   ctx := gt.NewGet(`http://192.168.0.1`)

   ctx.SetStartFunc(func(ctx *gt.Context){
            log.Println("请求前： 添加代理等等操作")
            ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }
    )

   ctx.SetSucceedFunc(func(ctx *gt.Context){
            log.Println("请求成功： 处理数据或存储等等")
            log.Println(ctx.RespBodyString())
        }
    )

    ctx.SetFailedFunc(func(ctx *gt.Context){
            log.Println("请求失败： 记录失败等等")
            log.Println(ctx.Err)
        }
    )

    ctx.SetRetryFunc(func(ctx *gt.Context){
             log.Println("请求重试： 更换代理或添加等待时间等等")
             ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }
    )

    ctx.SetEndFunc(func(ctx *gt.Context){
             log.Println("请求结束： 记录结束，处理数据等等")
             log.Println(ctx)
        }
    )

    ctx.Do()
    log.Println(ctx.Html)
}
```

例子2：

```
import gt "github.com/mangenotwork/gathertool"

func main(){
   ctx,_ := gt.Get(`http://192.168.0.1`, gt.StartFunc(Start), 
                           gt.SucceedFunc(Succeed), gt.FailedFunc(Failed), 
                           gt.RetryFunc(Retry), gt.EndFunc(End) )
    log.Println(ctx.Html)
}

func Start(ctx *gt.Context){
     log.Println("请求前： 添加代理等等操作")
     ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
}

func Succeed(ctx *gt.Context){
      log.Println("请求成功： 处理数据或存储等等")
      log.Println(ctx.RespBodyString())
}

func Failed(ctx *gt.Context){
     log.Println("请求失败： 记录失败等等")
     log.Println(ctx.Err)
}

func Retry(ctx *gt.Context){
     log.Println("请求重试： 更换代理或添加等待时间等等")
     ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
}

func End(ctx *gt.Context){
      log.Println("请求结束： 记录结束，处理数据等等")
      log.Println(ctx)
}
```

---

### 2.2 默认状态码事件
状态码对应事件的全局的，可以进行随意修改

---

#### func SetStatusCodeSuccessEvent(code int) 
将指定状态码设置为执行成功事件

---

#### func SetStatusCodeRetryEvent(code int) 
将指定状态码设置为执行重试事件

---
#### func SetStatusCodeFailEvent(code int) 
将指定状态码设置为执行失败事件

---

默认状态码事件表：

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


列子：
```golang
import gt "github.com/mangenotwork/gathertool"
gt.SetStatusCodeSuccessEvent(412) // 指定412状态码执行成功请求事件
gt.SetStatusCodeRetryEvent(413) // 指定413状态码执行重试请求事件
```

---

### 2.3 事件转换
事件是可以进行强制转换的；比较常用的操作，遇到错误就重试请求。

---

#### func (c \*Context) OpenErr2Retry()   
开启请求失败都执行retry请求事件

```
import gt "github.com/mangenotwork/gathertool"

ctx := gt.NewGet(`http://192.168.0.1`)
ctx.OpenErr2Retry() 
ctx.Do()
```

---

### 2.4 请求头
gathertool 有UserAgent Map 是全局的，可自定义，可扩展等; 在不设置的情况下请求默认会使用 PCAgent随机的一个

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


#### func GetAgent(agentType UserAgentType) string 
随机获取 user-agent

---
#### func SetAgent(agentType UserAgentType, agent string) 
设置 user-agent

```
import gt "github.com/mangenotwork/gathertool"
userAgent := gt.GetAgent(gt.WindowsAgent)
gt.SetAgent(LinuxAgent, userAgent)
```

---

### 2.5 重试
默认重试是10次，也可以修改重试次数；也可以及时关闭重试；

#### func (c \*Context) GetRetryTimes() int  
获取当前重试次数

---

#### func (c \*Context) CloseRetry()   
关闭重试

```
import gt "github.com/mangenotwork/gathertool"

func mian(){
    gt.Get(`http://192.168.0.1`,gt.SetRetryFunc(succeed))
}

func retry(ctx *gt.Context) {
    // 超过3次重试就关闭
    if ctx.GetRetryTimes > 3 {
        ctx.CloseRetry()
    }
}
```

---
#### func (c \*Context) SetRetryTimes(times int) \*Context   
设置重试次数

```
import gt "github.com/mangenotwork/gathertool"

func mian(){
    ctx := gt.NewGet(`http://192.168.0.1`)
    ctx.SetRetryTimes(5) // 指定当前请求最多重试5次
    ctx.Do()
}
```

---

### 2.6 请求上下文 Context
这就是一个请求的主体，一个请求可以是三连续，也可以是一次性的，下面是一个请求的方法介绍

#### func (c \*Context) SetSucceedFunc(successFunc func(c \*Context)) \*Context 
设置成功后的方法

```
ctx.SetStartFunc(func(ctx *gt.Context){
            log.Println("请求前： 添加代理等等操作")
            ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }
    )
```

---
#### func (c \*Context) SetFailedFunc(failedFunc func(c \*Context)) \*Context 
设置错误后的方法

```
ctx.SetFailedFunc(func(ctx *gt.Context){
           log.Println("请求成功： 处理数据或存储等等")
            log.Println(ctx.RespBodyString())
        }
    )
```

---
#### func (c \*Context) SetRetryFunc(retryFunc func(c \*Context)) \*Context 
设置重试，在重试前的方法

```
ctx.SetRetryFunc(func(ctx *gt.Context){
             log.Println("请求重试： 更换代理或添加等待时间等等")
             ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }
    )
```

---
#### func (c \*Context) SetRetryTimes(times int) \*Context 
设置重试次数

```
ctx.SetRetryTimes(3)
```

---
#### func (c \*Context) GetRetryTimes() int  
获取当前重试次数

```
times := ctx.GetRetryTimes()
```

---
#### func (c \*Context) Do() func() 
执行请求

```golang
ctx := gt.NewGet(`http://192.168.0.1?a=aaa&b=bbb`)
ctx.Do()
log.Println(ctx.RespBodyString())
```

---
#### func (c \*Context) RespBodyString() string 
Body -> String

```
log.Println(ctx.RespBodyString())
```

---
#### func (c \*Context) RespBodyHtml() string 
Body html string

```golang
log.Println(ctx.RespBodyHtml())
```

---
#### func (c \*Context) RespBodyMap() map[string]interface{} 
Body -> Map

```
log.Println(ctx.RespBodyMap())
```

---
#### func (c \*Context) RespBodyArr() []interface{}  
Body -> Arr

```
log.Println(ctx.RespBodyArr())
```

---
#### func (c \*Context) GetRespHeader() string  
获取 Response 的header

---
#### func (c \*Context) RespContentLength() int64 
获取Response 的 ContentLength

---
#### func (c \*Context) CheckReqMd5() string 
这次请求Request的md5值，是唯一的，应用场景: 记录每次请求的唯一性排除重复请求等

---
#### func (c \*Context) CheckMd5() string 
这次请求的md5值，是唯一的

---
#### func (c \*Context) AddHeader(k,v string) *Context   
添加请求 header

---
#### func (c \*Context) AddCookie(k, v string) *Context  
添加请求cookie

---
#### func (c \*Context) SetProxy(proxyUrl string) *Context 
设置代理

---
#### func (c \*Context) SetProxyFunc(f func() *http.Transport) *Context  
设置代理函数

---
#### func (c \*Context) SetProxyPool(pool *ProxyPool) *Context  
设置代理池，会依次获取代理

--- 
#### func (c \*Context) Upload(filePath string) func()  
下载

---
#### func (c \*Context) CookieNext() error  
使用上次请求返回的cookie

---
#### func (c \*Context) CloseLog()  
关闭打印日志

---
#### func (c \*Context) OpenErr2Retry()  
开启请求失败都执行retry

---
#### func (c \*Context) CloseRetry()  
关闭重试

---
#### func (c \*Context) GetParam(key string) interface{}  
获取上下文参数

---
#### func (c \*Context) AddParam(key string, val interface{})  
添加上下文参数

---
#### func (c \*Context) DelParam(key string)  
删除上下文参数

---
#### func (c \*Context) Upload(filePath string) func()  
下载

---

### 2.7 Context 的成员
```golang
// Context 请求上下文
type Context struct {
	// Token
	Token string

	// http client
	Client *http.Client

	// http Request
	Req *http.Request

	// http Response
	Resp *http.Response

	// Error
	Err error

	// Ctx context.Context
	Ctx context.Context

	// 执行的次数 初始化都是0
	times RetryTimes

	// 最大允许重试次数
	MaxTimes RetryTimes

	// 请求成功了需要处理的事件
	SucceedFunc SucceedFunc

	// 请求失败了需要做的事
	FailedFunc FailedFunc

	// 请求状态码设置了重试，在重试前的事件
	RetryFunc RetryFunc

	// 请求开始前的方法
	StartFunc StartFunc

	// 请求完成后的方法
	EndFunc EndFunc

	// 本次请求的任务
	// 用于有步骤的请求和并发执行请求
	Task *Task

	// 请求返回的结果
	RespBody []byte

	// job编号
	// 在执行多并发执行抓取任务，每个并发都有一个编号
	// 这个编号是递增分配的
	JobNumber int

	// 请求的响应时间 单位ms
	Ms time.Duration

	// 是否显示日志, 默认是显示的
	IsLog IsLog

	// 指定失败执行重试事件
	err2retry bool

	// 是否关闭重试
	isRetry bool

	// 休眠时间
	sleep Sleep

	// 输出字符串
	Text string

	// 输出Json
	Json string

	// 输出xml
	Xml string

	// 输出HTML
	Html string

	// 请求上下文参数
	Param map[string]interface{}
}
```

## 三、 请求使用

### 3.1 Get

#### func Get(url string) (\*Context, error)

```golang
import gt "github.com/mangenotwork/gathertool"

ctx, err := gt.Get(`http://192.168.0.1`)
log.Println(ctx.RespBodyString(), err)
```
---

#### func NewGet(url string) \*Context 
新建一个get请求

```golang
ctx := gt.NewGet(`http://192.168.0.1?a=aaa&b=bbb`)
ctx.Do()
log.Println(ctx.RespBodyString())
```

---

### 3.2 Post

#### func Post(url string, data []byte, contentType string) (\*Context, error)

```golang
ctx, err := gt.Post(`https://httpbin.org/post`, []byte(`{"a":"a"}`), "application/json;")
log.Println(ctx.RespBodyString(), err)
```
---

#### func NewPost(url string, data []byte, contentType string) \*Context

```golang
ctx := gt.NewPost(`https://httpbin.org/post`, []byte(`{"a":"a"}`), "application/json;")
ctx.Do()
log.Println(ctx.RespBodyString())
```
---

#### func PostJson(url string, jsonStr string) (\*Context, error)

```golang
ctx, err := gt.PostJson(`https://httpbin.org/post`, `{"a":"a"}`)
log.Println(ctx.RespBodyString(), err)
```

---

#### func PostForm(url string, data url.Values) (\*Context, error)

```golang
formData := map[string]string{
    "a":"a",
}
ctx, err := gt.PostForm(`https://httpbin.org/post`, formData)
log.Println(ctx.RespBodyString(), err)
```
---

#### func NewPostForm(u string, data map[string]string, vs ...interface{}) \*Context

---

#### func PostFile(url, paramName, filePath string, vs ...interface{}) \*Context
```
ctx := gt.PostFile(`http:/192.168.0.9:8888`, "file", "/home/test.txt")
ctx.Do
log.Println(ctx.RespBodyString(), err)
```

---


### 3.3 Put

#### func Put(url string, data []byte, contentType string, vs ...interface{}) (\*Context, error)

---

#### func NewPut(url string, data []byte, contentType string, vs ...interface{}) \*Context

### 3.4 Delete

#### func Delete(url string, vs ...interface{}) (\*Context, error)

---

#### func NewDelete(url string, vs ...interface{}) \*Context

---

### 3.5 Options

#### func Options(url string, vs ...interface{}) (\*Context, error)

---

#### func NewOptions(url string, vs ...interface{}) \*Context

---

### 3.6 Upload

#### func Upload(url, savePath string, vs ...interface{}) (\*Context, error)
```
ctx, err := gt.PostFile(`http:/192.168.0.9:8888/file/txt1.txt`, "/home/txt1.txt")
log.Println(ctx.RespBodyString(), err)
```

### 3.7 代理

代理ip

```
type ProxyIP struct {
	IP string
	Post int
	User string
	Pass string
	IsTLS bool
}
```
---
#### func NewProxyIP(ip string, port int, user, pass string, isTls bool) \*ProxyIP

```
ip := gt.NewProxyIP("127.0.0.1",1981, "", "", false)
```

---
#### func (p \*ProxyIP) String() string

---
####  type ProxyPool struct
代理池

---
#### func NewProxyPool() \*ProxyPool  
新建代理池

```
ProxyPool = gt.NewProxyPool()
ProxyPool.Add(gt.NewProxyIP("120.26.170.171",1981, "wahaha", "993126", false))
ProxyPool.Add(gt.NewProxyIP("8.134.60.130",1981, "wahaha", "993126", false))
```

--- 
#### func (p \*ProxyPool) Add(proxyIP \*ProxyIP) 
添加代理到代理池

---
#### func (p \*ProxyPool) Get() (string, int) 
获取一个代理，是依次获取。 返回的第二个参数是这个代理在代理池的位置

---
#### func (p \*ProxyPool) Del(n int)  
移除一个代理，传人参数的代理的位置


### 3.8 Cookie
#### var CookiePool \*cookiePool  
cookie池

---
#### func NewCookiePool() \*cookiePool

---
#### func (c \*cookiePool) Add(cookie \*http.Cookie)

代理池添加代理

---
#### func (c \*cookiePool) Get() \*http.Cookie
代理池按获取顺序获取一个代理

---
#### func (p *ProxyPool) Del(n int)
代理池删除代理

---
#### func (c *Context) CookieNext() error
CookieNext 使用上下文cookies

---

### 3.9 Header 

#### type Header map[string]string

---
#### func NewHeader(data map[string]string) Header
NewHeader 新建Header

---
#### func (h Header) Set(key, value string) Header
Set Header Set

---
#### func (h Header) Delete(key string) Header
Delete Header Delete

---

#### func (c *Context) AddHeader(k,v string) *Context
AddHeader 给请求添加Header

---
#### func (c *Context) GetRespHeader() string
GetRespHeader 获取Response的Header

---

### 3.12 Body & Log & Arg

#### func (c *Context) RespBodyString() string
RespBodyString 请求到的Body转换为string类型

---
#### func (c *Context) RespBodyHtml() string
RespBodyHtml 请求到的Body转换为string类型的html格式

---
#### func (c *Context) RespBodyMap() map[string]interface{}
RespBodyMap 请求到的Body转换为map[string]interface{}类型

---
#### func (c *Context) RespBodyArr() []interface{}
RespBodyArr 请求到的Body转换为[]interface{}类型

---
#### func (c *Context) RespContentLength() int64
RespContentLength 获取Response ContentLength的值

---
#### func (c *Context) CheckReqMd5() string
CheckReqMd5 本次请求的md5值， url+reqBodyBytes+method

---
#### func (c *Context) CheckMd5() string
CheckMd5 本次请求上下文包含输出结果的md5值

---
#### func (c *Context) CloseLog()
CloseLog 关闭日志打印

---
#### func (c *Context) GetParam(key string) interface{}
GetParam 获取上下文参数

---
#### func (c *Context) AddParam(key string, val interface{})
AddParam 添加上下文参数

---
#### func (c *Context) DelParam(key string)
DelParam 删除上下文参数

---

3.11 其他

### func (c *Context) OpenErr2Retry()
OpenErr2Retry 开启请求失败都执行retry

---
### func (c *Context) CloseRetry()
CloseRetry 关闭重试

---
### func SearchPort(ipStr string, vs ...interface{})
SearchPort 扫描ip的端口

---
### func Ping(ip string)
Ping ping IP

---

## 四、常用方法

### 4.1 类型转换

#### func Any2String(data interface{}) string
任何类型返回值字符串类型

---
#### func Json2Map(str string) map[string]interface{}
json -> map

---
#### func Map2Json(m map[string]interface{}) (string, error)
map -> json

---
#### func Any2Map(data interface{}) map[string]interface{}
interface{} -> map[string]interface{}

---
#### func Any2Int(data interface{}) int
interface{} -> int

---
#### func Any2int64(data interface{}) int64
interface{} -> int64

---
#### func Any2Arr(data interface{}) []interface{}
interface{} -> []interface{}

---
#### func Any2Float64(data interface{}) float64
interface{} -> float64

---
#### func Any2Strings(data interface{}) []string
interface{} -> []string

---
#### func Any2Json(data interface{}) (string, error)
interface{} -> json string

---
#### func Int2Hex(i int64) string
int -> hex string

---
#### func Int642Hex(i int64) string
int64 -> hex string

---
#### func Hex2Int(s string) int
hex string -> int

---
#### func Hex2Int64(s string) int64
hex string -> int

---
#### func Str2Int64(str string) int64
string -> int64

---
#### func Str2Int(str string) int
string -> int

---
#### func Str2Int32(str string) int32
string -> int32

---
#### func Str2Float64(str string) float64
string -> float64

---
#### func Str2Float32(str string) float32
string -> float32

---
#### func Uint82Str(bs []uint8) string
[]uint8 -> string

---
#### func Str2Byte(s string) []byte
string -> []byte

---
#### func Byte2Str(b []byte) string
[]byte -> string

---
#### func Bool2Byte(b bool) []byte
bool -> []byte

---
#### func Byte2Bool(b []byte) bool
[]byte -> bool

---
#### func Int2Byte(i int) []byte
int -> []byte

---
#### func Byte2Int(b []byte) int
[]byte -> int

---
#### func Int642Byte(i int64) []byte
int64 -> []byte

---
#### func Byte2Int64(b []byte) int64
[]byte -> int64

---
#### func Float322Byte(f float32) []byte
float32 -> []byte

---
#### func Float322Uint32(f float32) uint32
float32 -> uint32

---
#### func Byte2Float32(b []byte) float32
[]byte -> float32

---
#### func Float642Byte(f float64) []byte
float64 -> []byte

---
#### func Float642Uint64(f float64) uint64
float64 -> uint64

---
#### func Byte2Float64(b []byte) float64
[]byte -> float64

---
#### func Byte2Bit(b []byte) []uint8
[]byte -> []uint8 (bit)

---
#### func Bit2Byte(b []uint8) []byte
[]uint8 -> []byte

---
#### func Struct2Map(obj interface{}, hasValue bool) (map[string]interface{}, error)
Struct  ->  map


### 4.2 字符串相关

#### func MD5(str string) string
MD5

---
#### func CleaningStr(str string) string
清理字符串前后空白 和回车 换行符号

---
#### func StrDeleteSpace(str string) string
删除字符串前后的空格

---
#### func IsContainStr(items []string, item string) bool
字符串是否等于items中的某个元素

---

## 4.3 其他

#### func OSLine() string
系统对应换行符

---
#### func EncodeByte(v interface{}) []byte
encode byte

---
#### func DecodeByte(b []byte) (interface{}, error)
decode byte

---
#### func FileSizeFormat(fileSize int64) (size string)
字节的单位转换 保留两位小数

---
#### func DeepCopy(dst, src interface{}) error
深copy

---
#### func FileMd5(path string) (string, error)
文件md5

---
#### func PathExists(path string)
目录不存在则创建

---
#### func StrDuplicates(a []string) []string
数组，切片去重和去空串

---
#### func windowsPath(path string) string
windows平台需要转一下

---
#### func GetNowPath() string
获取当前运行路径

---
#### func FileMd5sum(fileName string) string
文件 Md5

---
#### func SearchBytesIndex(bSrc []byte, b byte) int
[]byte 字节切片 循环查找

---
#### func IF(condition bool, a, b interface{}) interface{}
三元表达式  use: IF(a>b, a, b).(int)

---
#### func CopySlice(s []interface{}) []interface{}
Copy slice

---
#### func CopySliceStr(s []string) []string
Copy string slice

---
#### func CopySliceInt(s []int) []int
Copy int slice

---
#### func CopySliceInt64(s []int64) []int64
Copy int64 slice

---
#### func IsInSlice(s []interface{}, v interface{})  bool
is in slice

---
#### func ReplaceAllToOne(str string, from []string, to string) string
批量统一替换字符串

---
#### func Exists(path string) bool

---
#### func IsDir(path string) bool

---
#### func Pwd() string

---
#### func Chdir(dir string) error

---
#### func HumanFriendlyTraffic(bytes uint64) string

#### func StrToSize(sizeStr string) int64

字节换算

---

### 4.4 字符转码,编码解码

#### func ConvertByte2String(byte []byte, charset Charset) string

编码转换
- UTF8    = Charset("UTF-8")
- GB18030 = Charset("GB18030")
- GBK = Charset("GBK")
- GB2312 = Charset("GB2312")

---
#### func UnicodeDec(raw string) string
#### func UnicodeDecByte(raw []byte) []byte

Unicode 解码

---
#### func UnescapeUnicode(raw []byte) ([]byte, error)
Unicode 转码

---
#### func Base64Encode(str string) string
base64 编码

---
#### func Base64Decode(str string) (string,error)
base64 解码

---
#### func Base64UrlEncode(str string) string
base64 url 编码

---
#### func Base64UrlDecode(str string) (string,error)
base64 url 解码

---
#### func ToUTF8(srcCharset string, src string) (dst string, err error)
其他编码转为UTF8
```
ToUTF8("GB2312", "你好")
```

---
#### func UTF8To(dstCharset string, src string) (dst string, err error)
UTF8转其他编码
```
UTF8To("GB2312", "assdsdfdsf")
```

---
#### func ToUTF16(srcCharset string, src string) (dst string, err error)
ToUTF16

---
#### func UTF16To(dstCharset string, src string) (dst string, err error)
UTF16To

---
#### func ToBIG5(srcCharset string, src string) (dst string, err error)
ToBIG5

---
#### func BIG5To(dstCharset string, src string) (dst string, err error)
BIG5To

---
#### func ToGDK(srcCharset string, src string) (dst string, err error)
ToGDK

---
#### func GDKTo(dstCharset string, src string) (dst string, err error)
GDKTo

---
#### func ToGB18030(srcCharset string, src string) (dst string, err error)
ToGB18030

---
#### func GB18030To(dstCharset string, src string) (dst string, err error)
GB18030To

---
#### func ToGB2312(srcCharset string, src string) (dst string, err error)
ToGB2312

---
#### func GB2312To(dstCharset string, src string) (dst string, err error)
GB2312To

---
#### func ToHZGB2312(srcCharset string, src string) (dst string, err error)
ToHZGB2312

---
#### func HZGB2312To(dstCharset string, src string) (dst string, err error)
HZGB2312To

---

### 4.5 集合

#### type Set map[string]struct{}

---
#### func (s Set) Has(key string) bool

---
#### func (s Set) Add(key string)

---
#### func (s Set) Delete(key string)

### 4.6 栈

#### type Stack struct

---
#### func New() *Stack

---
#### func (s *Stack) Push(data interface{})

---
#### func (s *Stack) Pop()

---
#### func (s *Stack) String() string

---


### 4.7 Map


#### func MapCopy(data map[string]interface{}) (copy map[string]interface{})
map 深copy

---
#### func MapMergeCopy(src ...map[string]interface{}) (copy map[string]interface{})
多个map 合并为一个新的map

---
#### func Map2Slice(data interface{}) []interface{}
Map2Slice Eg: {"K1": "v1", "K2": "v2"} => ["K1", "v1", "K2", "v2"]


### 4.8 固定顺序Map

GDMapApi 固定顺序 Map 接口
```go
type GDMapApi interface {
	Add(key string, value interface{}) *gDMap
	Get(key string) interface{}
	Del(key string) *gDMap
	Len() int
	KeyList() []string
	AddMap(data map[string]interface{}) *gDMap
	Range(f func(k string, v interface{})) *gDMap
	RangeAt(f func(id int, k string, v interface{})) *gDMap
	CheckValue(value interface{}) bool // 检查是否存在某个值
	Reverse()                          //反序
}
```
---
#### func NewGDMap() *gDMap
NewGDMap ues: NewGDMap().Add(k,v)

---
#### func (m *gDMap) Add(key string, value interface{}) *gDMap
 Add  添加kv

```go
NewGDMap().Add("a", 1).Add("b",2)
```

---
#### func (m *gDMap) Get(key string) interface{}
Get 通过key获取值

---
#### func (m *gDMap) Del(key string) *gDMap
Del 删除指定key的值
```go
NewGDMap().Del("a", 1).Del("b",2)
```

---
#### func (m *gDMap) Len() int
 Len map的长度

---
#### func (m *gDMap) KeyList() []string
KeyList 打印map所有的key

---
#### func (m *gDMap) AddMap(data map[string]interface{}) *gDMap
AddMap 写入map

---
#### func (m *gDMap) Range(f func(k string, v interface{})) *gDMap
Range 遍历map

---
#### func (m *gDMap) RangeAt(f func(id int, k string, v interface{})) *gDMap
RangeAt Range 遍历map含顺序id

---
#### func (m *gDMap) CheckValue(value interface{}) bool
CheckValue 查看map是否存在指定的值

---
#### func (m *gDMap) Reverse()
Reverse map反序

---

### 4.9 Slice

#### func SliceCopy(data []interface{}) []interface{}
slice 深copy

---
#### func Slice2Map(slice interface{}) map[string]interface{}
```
["K1", "v1", "K2", "v2"] => {"K1": "v1", "K2": "v2"}
["K1", "v1", "K2"]       => nil
```

---
#### func SliceTool() *sliceTool
SliceTool use : SliceTool().CopyInt64(a)

---
#### func (sliceTool) CopyInt64(a []int64) []int64
CopyInt64 copy int64

---
#### func (sliceTool) CopyStr(a []string) []string
 CopyStr copy string

---
#### func (sliceTool) CopyInt(a []int) []int
CopyInt copy int

---
#### func (sliceTool) ContainsByte(a []byte, x byte) bool
 ContainsByte contains byte

---
#### func (sliceTool) ContainsInt(a []int, x int) bool
ContainsInt contains int

---
#### func (sliceTool) ContainsInt64(a []int64, x int64) bool
ContainsInt64  contains int64

---
#### func (sliceTool) ContainsStr(a []string, x string) bool
ContainsStr contains str

---
#### func (sliceTool) DeduplicateInt(a []int) []int

---
#### func (sliceTool) DeduplicateInt64(a []int64) []int64
 DeduplicateInt64 deduplicate int64

---
#### func (sliceTool) DeduplicateStr(a []string) []string
DeduplicateStr  deduplicate string

---
#### func (sliceTool) DelInt(a []int, i int) []int
DelInt del int

---
#### func (sliceTool) DelInt64(a []int64, i int) []int64
DelInt64 del int64

---
#### func (sliceTool) DelStr(a []string, i int) []string
DelStr delete str

---
#### func (sliceTool) MaxInt(a []int) int

---
#### func (sliceTool) MaxInt64(a []int64) int64

---
#### func (sliceTool) MinInt(a []int) int

---
#### func (sliceTool) MinInt64(a []int64) int64

---
#### func (sliceTool) PopInt(a []int) (int, []int)

---
#### func (sliceTool) PopInt64(a []int64) (int64, []int64)

---
#### func (sliceTool) PopStr(a []string) (string, []string)

---
#### func (sliceTool) ReverseInt(a []int) []int
ReverseInt  反转

---
#### func (sliceTool) ReverseInt64(a []int64) []int64
ReverseInt64 reverse int64

---
#### func (sliceTool) ReverseStr(a []string) []string
ReverseStr  reverseStr

---
#### func (sliceTool) ShuffleInt(a []int) []int
ShuffleInt 洗牌

---

## 4.10 时间相关方法

> func Timestamp() string

时间戳

---
> func BeginDayUnix() int64

获取当天 0点

---
> func EndDayUnix() int64

获取当天 24点

---
> func MinuteAgo(i int) int64

获取多少分钟前的时间戳

---
> func HourAgo(i int) int64

获取多少小时前的时间戳

---
> func DayAgo(i int) int64

获取多少天前的时间戳

---
> func Daydiff(beginDay string, endDay string) int

两个时间字符串的日期差

---
> func TickerRun(t time.Duration, runFirst bool, f func())

间隔运行

---

func Timestamp2Date(timestamp int64) string


// GetChineseMonthDay 获取农历
func GetChineseMonthDay(date string) (rmonth, rday int64)


// NowToEnd 计算当前时间到这天结束还有多久
func NowToEnd() (int64, error)

// IsLeap 是否是闰年
func IsLeap(year int) bool

// IsToday 判断是否是今天   "2006-01-02 15:04:05"
// timestamp 需要判断的时间
func IsToday(timestamp int64) string

// IsTodayList 列表页的时间显示  "01-02 15:04"
func IsTodayList(timestamp int64) string

func Timestamp2Week(timestamp int64) string

func Timestamp2WeekXinQi(timestamp int64) string





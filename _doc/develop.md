# 《gathertool开发使用文档》

Date : 2023-03-28

Author : ManGe

Mail : 2912882908@qq.com

Github : https://github.com/mangenotwork/gathertool

[toc]

[一、介绍](#一介绍)

## 一、介绍

### 1.1 简介
- gathertool是golang脚本化开发集成库，目的是提高对应场景脚本程序开发的效率；
- gathertool也是一款轻量级爬虫库，特色是分离了请求事件，通俗点理解就是对请求过程状态进行事件处理。 
- gathertool也是接口测试&压力测试库，在接口测试脚本开发上有明显的效率优势，
- gathertool还集成了对三方中间件的操作，DB操作等。

---

### 1.2 使用
go get github.com/mangenotwork/gathertool

---

### 1.3 介绍
gathertool是一个高度封装工具库，包含了http/s的请求，Mysql数据库方法，数据类型处理方法，数据提取方法，websocket相关方法，
TCP|UDP相关方法，NoSql相关方法，开发常用方法等;  可以用于爬虫程序，接口&压力测试程序，常见网络协议调试程序，数据提取与存储程序等；
gathertool的请求特点: 会在请求阶段执行各个事件如请求失败后的重试事件,请求前后的事件，请求成功事件等等, 可以根据请求状态码自定义这些事件；
gathertool还拥有很好的可扩展性， 适配传入任意自定义http请求对象， 能适配各种代理对象等等；
gathertool还拥有抓取数据存储功能, 比如存储到mysql, redis, mongo, pgsql等等; 还有很多创新的地方文档会根据具体方法进行介绍；
gathertool还封装了消息队列接口，支持Nsq,Kafka,rabbitmq,redis等消息队列中间件

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

#### func Timestamp() string
时间戳

---
#### func BeginDayUnix() int64
获取当天 0点

---
#### func EndDayUnix() int64
获取当天 24点

---
#### func MinuteAgo(i int) int64
获取多少分钟前的时间戳

---
> func HourAgo(i int) int64

获取多少小时前的时间戳

---
#### func DayAgo(i int) int64
获取多少天前的时间戳

---
#### func Daydiff(beginDay string, endDay string) int
两个时间字符串的日期差

---
#### func TickerRun(t time.Duration, runFirst bool, f func())
间隔运行

---
#### func Timestamp2Date(timestamp int64) string

---
#### func GetChineseMonthDay(date string) (rMonth, rDay int64)
GetChineseMonthDay 获取农历

---
#### func NowToEnd() (int64, error)
 NowToEnd 计算当前时间到这天结束还有多久

---
#### func IsLeap(year int) bool
IsLeap 是否是闰年

---
#### func IsToday(timestamp int64) string
IsToday 判断是否是今天   "2006-01-02 15:04:05"

---
#### func IsTodayList(timestamp int64) string
IsTodayList 列表页的时间显示  "01-02 15:04"

---
#### func Timestamp2Week(timestamp int64) string

---
#### func Timestamp2WeekXinQi(timestamp int64) string

---

## 五、数据提取

### 5.1 正则

#### func RegFindAll(regStr, rest string) [][]string
正则提取所有
```
list := RegFindAll(`<option(.*?)</option>`, txt)
```

---

#### func RegHtmlA(str string, property ...string) []string
提取a标签
```
`(?is:<a.*?</a>)`
```

---
#### func RegHtmlTitle(str string, property ...string) []string

提取title标签
```
`(?is:<title.*?</title>)`
```

---
#### func RegHtmlTr(str string, property ...string) []string
提取tr标签
```
`(?is:<tr.*?</tr>)`
```

---

#### func RegHtmlInput(str string, property ...string) []string
提取input标签
```
`(?is:<input.*?>)`
```

---

#### func RegHtmlTd(str string, property ...string) []string
提取td标签
```
`(?is:<td.*?</td>)`
```

---

#### func RegHtmlP(str string, property ...string) []string
提取P标签
```
`(?is:<p.*?</p>)`
```

---

#### func RegHtmlSpan(str string, property ...string) []string
提取span标签
```
`(?is:<span.*?</span>)`
```

---

#### func RegHtmlSrc(str string, property ...string) []string
提取src内容
```
`(?is:src=\".*?\")`
```

---

#### func RegHtmlHref(str string, property ...string) []string
提取href内容
```
`(?is:href=\".*?\")`
```

---

#### func RegHtmlVideo(str string, property ...string) []string
提取video内容
```
`(?is:<video.*?</video>)`
```

---

#### func RegHtmlCanvas(str string, property ...string) []string
提取canvas
```
`(?is:<canvas.*?</canvas>)`
```

---

#### func RegHtmlCode(str string, property ...string) []string
提取code标签
```
`(?is:<code.*?</code>)`
```

---

#### func RegHtmlImg(str string, property ...string) []string
提取img标签
```
`(?is:<img.*?/>)`
```

---

#### func RegHtmlUl(str string, property ...string) []string
提取ul标签
```
`(?is:<ul.*?</ul>)`
```

---

#### func RegHtmlLi(str string, property ...string) []string
提取li标签
```
`(?is:<li.*?</li>)`
```

---

#### func RegHtmlMeta(str string, property ...string) []string
提取meta标签
```
`(?is:<meta.*?>)`
```

---

#### func RegHtmlSelect(str string, property ...string) []string
提取select标签
```
`(?is:<select.*?</select>)`
```

---

#### func RegHtmlTable(str string, property ...string) []string
提取table标签
```
`(?is:<table.*?</table>)`
```

---

#### func RegHtmlButton(str string, property ...string) []string
提取button标签
```
`(?is:<button.*?</button>)`
```

---

#### func RegHtmlH(str, typeH string, property ...string) []string
提取h标签
```
`(?is:<h1.*?</h1>)`
```

---

#### func RegHtmlTbody(str string, property ...string) []string
提取Tbody标签
```
`(?is:<table.*?</table>)`
```

---

#### func RegHtmlATxt(str string, property ...string) []string
提取a标签内容非标签部分
```
`(?is:<a.*?>(.*?)</a>)`
```

---

#### func RegHtmlTitleTxt(str string, property ...string) []string
提取title标签内容非标签部分
```
`(?is:<title.*?>(.*?)</title>)`
```

---

#### func RegHtmlTrTxt(str string, property ...string) []string
提取tr标签内容非标签部分
```
`(?is:<tr.*?>(.*?)</tr>)`
```

---

#### func RegHtmlInputTxt(str string, property ...string) []string
提取input标签内容非标签部分
```
`(?is:<input(.*?)>)`
```

---

#### func RegHtmlTdTxt(str string, property ...string) []string
提取td标签内容非标签部分
```
`(?is:<td.*?>(.*?)</td>)`
```

---

#### func RegHtmlPTxt(str string, property ...string) []string
提取p标签内容非标签部分
```
`(?is:<p.*?>(.*?)</p>)`
```

---

#### func RegHtmlSpanTxt(str string, property ...string) []string
提取span标签内容非标签部分
```
`(?is:<span.*?>(.*?)</span>)`
```

---

#### func RegHtmlSrcTxt(str string, property ...string) []string
提取Src内容
```
`(?is:src=\"(.*?)\")`
```

---

#### func RegHtmlHrefTxt(str string, property ...string) []string
提取href内容
```
`(?is:href=\"(.*?)\")`
```

---

#### func RegHtmlCodeTxt(str string, property ...string) []string
提取code标签内容非标签部分
```
`(?is:<code.*?>(.*?)</code>)`
```

---

#### func RegHtmlUlTxt(str string, property ...string) []string
提取ul标签内容非标签部分
```
`(?is:<ul.*?>(.*?)</ul>)`
```

---

#### func RegHtmlLiTxt(str string, property ...string) []string
提取li标签内容非标签部分
```
`(?is:<li.*?>(.*?)</li>)`
```

---

#### func RegHtmlSelectTxt(str string, property ...string) []string
提取select标签内容非标签部分
```
`(?is:<select.*?>(.*?)</select>)`
```

---

#### func RegHtmlTableTxt(str string, property ...string) []string
提取table标签内容非标签部分
```
`(?is:<table.*?>(.*?)</table>)`
```

---

#### func RegHtmlButtonTxt(str string, property ...string) []string
提取button标签内容非标签部分
```
`(?is:<button.*?>(.*?)</button>)`
```

---

#### func RegHtmlHTxt(str, typeH string, property ...string) []string
提取h标签内容非标签部分
```
`(?is:<h1.*?>(.*?)</h1>)`
```

---

#### func RegDelHtml(str string) string
删除所有标签

---

#### func RegDelNumber(str string) string
删除所有数字

---

#### func RegDelHtmlA(str string) string
删除所有a标签

---

#### func RegDelHtmlTitle(str string) string
删除所有html标签

---

#### func RegDelHtmlTr(str string) string
删除所有tr标签

---

#### func RegDelHtmlInput(str string, property ...string) string
删除所有input标签

---

#### func RegDelHtmlTd(str string, property ...string) string
删除所有td标签

---

#### func RegDelHtmlP(str string, property ...string) string
删除所有p标签

---

#### func RegDelHtmlSpan(str string, property ...string) string
删除所有span标签

---

#### func RegDelHtmlSrc(str string, property ...string) string
删除所有src

---

#### func RegDelHtmlHref(str string, property ...string) string
删除所有href

---

#### func RegDelHtmlVideo(str string, property ...string) string
删除所有video标签

---

#### func RegDelHtmlCanvas(str string, property ...string) string
删除所有canvas标签

---

#### func RegDelHtmlCode(str string, property ...string) string
删除所有code标签

---

#### func RegDelHtmlImg(str string, property ...string) string
删除所有img标签

---

#### func RegDelHtmlUl(str string, property ...string) string
删除所有ul标签

---

#### func RegDelHtmlLi(str string, property ...string) string
删除所有li标签

---

#### func RegDelHtmlMeta(str string, property ...string) string
删除所有meta标签

---

#### func RegDelHtmlSelect(str string, property ...string) string
删除所有select标签

---

#### func RegDelHtmlTable(str string, property ...string) string
删除所有table标签

---

#### func RegDelHtmlButton(str string, property ...string) string
删除所有button标签

---

#### func RegDelHtmlH(str, typeH string, property ...string) string
删除所有h标签

---

#### func RegDelHtmlTbody(str string, property ...string) string
删除所有body标签

---

#### func IsNumber(str string) bool
验证是否含有number
```
`^[0-9]*$`
```

---

#### func IsNumber2Len(str string, l int) bool
验证是否含有连续长度不超过长度l的number
```
`[0-9]{%d}`
```

---

#### func IsNumber2Heard(str string, n int) bool
验证是否含有n开头的number
```
`^(%d)[0-9]*$`
```

---

#### func IsFloat(str string) bool
验证是否是标准正负小数(123. 不是小数)
```
 `^(-?\d+\.\d+)?$`
```

---

#### func IsFloat2Len(str string, l int) bool
验证是否含有带不超过len个小数的小数
```
`^(-?\d+\.\d{%d})?$`
```

---

#### func IsChineseAll(str string) bool
验证是否是全汉字

---

#### func IsChinese(str string) bool
验证是否含有汉字

---

#### func IsChineseN(str string, number int) bool
验证是否含有number个汉字

---

#### func IsChineseNumber(str string) bool
验证是否全是汉字数字

---

#### func IsChineseMoney(str string) bool
验证是否是中文钱大写

---

#### func IsEngAll(str string) bool
验证是否是全英文
```
`^[A-Za-z]*$`
```

---

#### func IsEngLen(str string, l int) bool
验证是否含不超过len个英文字符
```
`^[A-Za-z]{%d}$`
```

---

#### func IsEngNumber(str string) bool
验证是否是英文和数字
```
`^[A-Za-z0-9]*$`
```

---

#### func IsAllCapital(str string) bool
验证是否全大写

---

#### func IsHaveCapital(str string) bool
验证是否有大写

---

#### func IsAllLower(str string) bool
验证是否全小写

---

#### func IsHaveLower(str string) bool
验证是否有小写

---

#### func IsLeastNumber(str string, n int) bool
验证不低于n个数字
```
`[0-9]{%d,}?`
```

---

#### func IsLeastCapital(str string, n int) bool
验证不低于n个大写字母
```
`[A-Z]{%d,}?`
```

---

#### func IsLeastLower(str string, n int) bool
验证不低于n个小写字母
```
`[a-z]{%d,}?`
```

---

#### func IsLeastSpecial(str string, n int) bool
验证不低于n特殊字符
```
`[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\&\$\.\*\+\?\{\}\(\)\[\]\|\!\_\@\#\%\-\=]{%d,}?`
```

---

#### func IsDomain(str string) bool
验证域名
```
`[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(/.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+/.?`
```

---

#### func IsURL(str string) bool
验证URL
```
`//([\w-]+\.)+[\w-]+(/[\w-./?%&=]*)?$`
```

---

#### func IsPhone(str string) bool
验证手机号码
```
`^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`
```

---

#### func IsLandline(str string) bool
验证电话号码("XXX-XXXXXXX"、"XXXX-XXXXXXXX"、"XXX-XXXXXXX"、"XXX-XXXXXXXX"、"XXXXXXX"和"XXXXXXXX)
```
`^(\(\d{3,4}-)|\d{3.4}-)?\d{7,8}$`
```

---

#### func IsIP(str string) bool
IP地址：((?:(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d)\\.){3}(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d))

---

#### func AccountRational(str string) bool
帐号合理性验证
```
`^[a-zA-Z][a-zA-Z0-9_]{4,15}$`
```

---

#### func IsXMLFile(str string) bool
是否三xml文件
```
`^*+\\.[x|X][m|M][l|L]$`
```

---

#### func IsUUID3(str string) bool
是否是uuid
```
`^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$`
```

---

#### func IsUUID4(str string) bool
是否是uuid
```
`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
```

---

#### func IsUUID5(str string) bool
是否是uuid
```
`^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
```

---

#### func IsRGB(str string) bool
是否是 rgb
```
`^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$`
```

---

#### func IsFullWidth(str string) bool
是否是全角字符
```
`[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]`
```

---

#### func IsHalfWidth(str string) bool
是否是半角字符
```
`[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]`
```

---

#### func IsBase64(str string) bool
是否是base64
```
`^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`
```

---

#### func IsLatitude(str string) bool
是否是纬度
```
`^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$`
```

---

#### func IsLongitude(str string) bool
是否是经度
```
`^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$`
```

---

#### func IsDNSName(str string) bool
是否是dns 名称
```
`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
```

---

#### func IsIPv4(str string) bool
是否是ipv4
```
`([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
```

---

#### func IsWindowsPath(str string) bool
是否是windos路径
```
`^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
```

---

#### func IsUnixPath(str string) bool
是否是unix路径
```
`^(/[^/\x00]*)+/?$`
```

---

#### func RegTime(str string, property ...string) []string
提取时间
```
`(?i)\d{1,2}:\d{2} ?(?:[ap]\.?m\.?)?|\d[ap]\.?m\.?`
```

---

#### func RegLink(str string, property ...string) []string
提取链接
```
`(?:(?:https?:\/\/)?(?:[a-z0-9.\-]+|www|[a-z0-9.\-])[.](?:[^\s()<>]+|\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\))+(?:\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:\'".,<>?]))`
```

---

#### func RegEmail(str string, property ...string) []string
提取邮件
```
`(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`
```

---

#### func RegIPv4(str string, property ...string) []string
提取ipv4
```
`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
```

---

#### func RegIPv6(str string, property ...string) []string
提取ipv6
```
`(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`
```

---

#### func RegIP(str string, property ...string) []string
提取ip
```
`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)|(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`
```

---

#### func RegMD5Hex(str string, property ...string) []string
提取md5
```
`[0-9a-fA-F]{32}`
```

---

#### func RegSHA1Hex(str string, property ...string) []string
提取sha1
```
`[0-9a-fA-F]{40}`
```

---

#### func RegSHA256Hex(str string, property ...string) []string
提取sha256
```
`[0-9a-fA-F]{64}`
```

---

#### func RegGUID(str string, property ...string) []string
提取guid
```
`[0-9a-fA-F]{8}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{12}`
```

---

#### func RegMACAddress(str string, property ...string) []string
提取MACAddress
```
`(([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2}))`
```

---

#### func RegEmail2(str string, property ...string) []string
提取邮件
```
"^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
```

---

#### func RegUUID3(str string, property ...string) []string
提取uuid
```
"^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
```

---

#### func RegUUID4(str string, property ...string) []string
提取uuid
```
"^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
```

---

#### func RegUUID5(str string, property ...string) []string
提取uuid
```
"^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
```

---

#### func RegUUID(str string, property ...string) []string
提取uuid
```
"^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
```

---

#### func RegInt(str string, property ...string) []string
提取整形
```
"^(?:[-+]?(?:0|[1-9][0-9]*))$"
```

---

#### func RegFloat(str string, property ...string) []string
提取浮点型
```
"^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
```

---

#### func RegRGBColor(str string, property ...string) []string
提取RGB值
```
"^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"
```

---

#### func RegFullWidth(str string, property ...string) []string
提取全角字符
```
"[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
```

---

#### func RegHalfWidth(str string, property ...string) []string
提取半角字符
```
"[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
```

---

#### func RegBase64(str string, property ...string) []string
提取base64字符串
```
"^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
```

---

#### func RegLatitude(str string, property ...string) []string
提取纬度
```
"^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
```

---

#### func RegLongitude(str string, property ...string) []string
提取经度
```
"^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
```

---

#### func RegDNSName(str string, property ...string) []string
提取dns
```
`^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
```

---

#### func RegFullURL(str string, property ...string) []string
提取url
```
`^(?:ftp|tcp|udp|wss?|https?):\/\/[\w\.\/#=jQuery1124048736386703191026_1648193326187&]+$`
```

---

#### func RegURLSchema(str string, property ...string) []string
提取url schema
```
`((ftp|tcp|udp|wss?|https?):\/\/)`
```

---

#### func RegURLUsername(str string, property ...string) []string
提取url username
```
`(\S+(:\S*)?@)`
```

---

#### func RegURLPath(str string, property ...string) []string
提取url path
```
`((\/|\?|#)[^\s]*)`
```

---

#### func RegURLPort(str string, property ...string) []string
提取url port
```
`(:(\d{1,5}))`
```

---

#### func RegURLIP(str string, property ...string) []string
提取 url ip
```
`([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
```

---

#### func RegURLSubdomain(str string, property ...string) []string
提取 url sub domain
```
`((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
```

---

#### func RegWinPath(str string, property ...string) []string
提取 windows路径
```
`^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
```

---

#### func RegUnixPath(str string, property ...string) []string
提取 unix路径

```
`^(/[^/\x00]*)+/?$`
```

---

### 4.2 Html提取

#### func GetPointHTML(htmlStr, label, attr, val string) ([]string, error) 
获取指定位置的HTML， 用标签， 标签属性， 属性值来定位
```go
    date := "2022-07-05"
	caseUrl := "http://futures.100ppi.com/qhb/day-%s.html"
	ctx, _ := gt.Get(fmt.Sprintf(caseUrl, date))
	//log.Println(ctx.Html)
	// 数据提取
	datas, _ := gt.GetPointHTML(ctx.Html, "div", "id", "domestic")
```

---
#### func GetPointIDHTML(htmlStr, label, val string) ([]string, error)
获取指定标签id属性的html

---
#### func GetPointClassHTML(htmlStr, label, val string) ([]string, error)
获取指定标签class属性的html

### 4.3 Json提取

#### func JsonFind(jsonStr, find string) (interface{}, error)
JsonFind 按路径寻找指定json值
```
用法参考  ./_examples/json/main.go
@find : 寻找路径，与目录的url类似， 下面是一个例子：
json:  {a:[{b:1},{b:2}]}
find=/a/[0]  =>   {b:1}
find=a/[0]/b  =>   1
```


实例：
```
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
            },
            {
                "date":"2021-10-28",
                "temperature":"13\/21℃",
                "weather":"多云转晴",
                "wid":{
                    "day":"01",
                    "night":"00"
                },
                "direct":"东北风"
            },
            {
                "date":"2021-10-29",
                "temperature":"14\/21℃",
                "weather":"多云转小雨",
                "wid":{
                    "day":"01",
                    "night":"07"
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
```

---
#### func JsonFind2Json(jsonStr, find string) (string, error)
寻找json,输出 json格式字符串

---
#### func JsonFind2Map(jsonStr, find string) (map[string]interface{}, error)
寻找json,输出 map[string]interface{}

---
#### func JsonFind2Arr(jsonStr, find string) ([]interface{}, error)
JsonFind2Arr 寻找json,输出 []interface{}

---
#### func IsJson(str string) bool
IsJson 是否是json格式

## 五、爬虫篇

## 爬虫篇
现在各个数据源日益月薪，变化非常快，爬虫程序无法做到固定，对端源一变化爬虫程序就要跟着变，有时候与其在原有基础上迭代还不如重新编写。重新编写就会出现一个问题，这个问题就是时间成本增加，为了降低时间成本所以我编写并开源了gathertool。
gathertool作为一款轻量级爬虫框架，其解决核心问题是提高编写爬虫程序效率降低时间成本。gathertool主要使用请求部分+提取部分+数据处理部分+存储部分作为爬虫框架的方法集，使用者可以根据具体业务进行灵活使用，下面先通过一个例子来介绍：

---

### 5.1 例子
已抓取 http://ip.bczs.net/country/CN ，将页面内容保存到数据为例
```
import gt "github.com/mangenotwork/gathertool"
func main(){
	db,err := gt.NewMysql("192.168.0.192", 3306,  "root", "root123", "test")  // 连接数据库
	if err != nil { panic(err) }
	ctx, err := gt.Get("http://ip.bczs.net/country/CN") // 请求数据
	if err != nil { panic(err) }
	for _, tbody := range gt.RegHtmlTbody(ctx.Html) { // 提取数据并保存
		for _, tr := range gt.RegHtmlTr(tbody) {
			td := gt.RegHtmlTdTxt(tr)
			if len(td) < 3 {
				gt.Error("异常数据 ： ", td)
				continue
			}
			err := db.InsertAt("ip", map[string]interface{}{"start": td[0], "end": td[1], "num": td[2]})
			if err != nil { panic(err) }
		}
	}
}
```

抓取结果如图：

![](http://mange1.oss-cn-beijing.aliyuncs.com/test/31d42875c8c6a41b6ee8f04d87f0deb5.png "")

从上面例子就可以看到使用gathertool 18行代码就能完成请求，提取，存储到数据库的全过程; 编码时间只有1分钟。

实例代码位置:  [ip地址信息抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/ip_bczs_cn)

---

### 5.2 并发抓取
gathertool的请求对象是独立的实例，这样做的目的也是为了并发请求时都是独立的，不使用重复对象; gathertool的并发抓取采用本地队列+任务对象;下面将具体介绍:


队列方法

```
type TodoQueue interface {
	Add(task *Task) error  //向队列中添加元素
	Poll()   *Task  //移除队列中最前面的元素
	Clear()  bool   //清空队列
	Size()  int     //获取队列的元素个数
	IsEmpty() bool  //判断队列是否是空
	Print() // 打印
}
```

---
 队列

```
type Queue struct {
	mux *sync.Mutex
	list []*Task
}
```

---
#### func NewQueue() TodoQueue

新建一个队列

---
#### func (q \*Queue) Add(task \*Task) error

向队列中添加元素

---
#### func (q \*Queue) Poll() \*Task

移除队列中最前面的额元素

---
#### func (q \*Queue) Clear() bool

清空队列

---
#### func (q \*Queue) Size() int

队列长度

---
#### func (q \*Queue) IsEmpty() bool

队列是否空

---
#### func (q *Queue) Print()

队列打印

---
#### type UploadQueue struct

下载队列

---
#### func NewUploadQueue() TodoQueue

新建一个下载队列

---
#### func (q *UploadQueue) Add(task *Task) error

向队列中添加元素

---
#### func (q *UploadQueue) Poll() *Task

移除队列中最前面的额元素

---
#### 任务对象

```
type Task struct {
	Url string
	JsonParam string
	HeaderMap *http.Header
	Data map[string]interface{} // 上下文传递的数据
	Urls []*ReqUrl // 多步骤使用
	Type string // "", "upload", "do"
	SavePath string
	SaveDir string
	FileName string
	once *sync.Once
}
```

---
#### func (task Task) GetDataStr(key string) string
获取上下文数据

---
#### func (task Task) AddData(key string, value interface{}) Task
添加上下问数据

---
#### func CrawlerTask(url, jsonParam string, vs ...interface{}) *Task
爬虫任务

---
#### func StartJob(jobNumber int, queue TodoQueue,f func(task *Task))
启动并发

---
#### func StartJobGet(jobNumber int, queue TodoQueue, vs ...interface{})
并发执行Get,直到队列任务为空;  jobNumber 并发数，queue 全局队列，

---
#### func StartJobPost(jobNumber int, queue TodoQueue, vs ...interface{})
开始运行并发Post,直到队列任务为空;  jobNumber 并发数，queue 全局队列，

---
#### func CPUMax()
启用多核执行

---

### 5.3 并发抓取列子
已抓取 http://ip.bczs.net/country/CN ，并发抓取每个ip的二级页面数据并存入数据库
```
var (
	// 全局声明抓取任务队列
	queue = gt.NewQueue()
	// 全局声明数据库客户端对象
	host   = "192.168.0.192"
	port    = 3306
	user      = "root"
	password  = "root123"
	dbname = "test"
	db,_ = gt.NewMysql(host, port, user, password, dbname)
)

func main(){
	// 1.在页面 http://ip.bczs.net/country/CN 获取所有ip
	_,_=gt.Get("http://ip.bczs.net/country/CN",gt.SucceedFunc(IPListSucceed))

	// 2. 并发抓取详情数据, 20个并发
	gt.StartJobGet(20,queue,
		gt.SucceedFunc(GetIPSucceed),//请求成功后执行的方法
		gt.RetryFunc(GetIPRetry),//遇到 502,403 等状态码重试前执行的方法，一般为添加休眠时间或更换代理
		gt.FailedFunc(GetIPFailed),//请求失败后执行的方法
	)
}

// 请求成功执行
func IPListSucceed(ctx *gt.Context){
	for _, tbody := range gt.RegHtmlTbody(ctx.RespBodyString()) {
		for _, tr := range gt.RegHtmlTr(tbody) {
			td := gt.RegHtmlTdTxt(tr)
			log.Println(td)
			if len(td) < 3 {
				gt.Error("异常数据 ： ", td)
				continue
			}
			startIp := gt.Any2String(gt.RegHtmlATxt(td[0])[0])// IP起始
			endIP := td[1]// 结束
			number := td[2]// 数量
			// 创建队列 抓取详情信息
			queue.Add(&gt.Task{
				Url: "http://ip.bczs.net/"+startIp,
				Data: map[string]interface{}{
					"start_ip":startIp,
					"end_ip":endIP,
					"number":number,
				},
			})
		}
	}
}


// 获取详情信息成功的处理
func GetIPSucceed(cxt *gt.Context){
	// 使用goquery库提取数据
	dom,err := goquery.NewDocumentFromReader(strings.NewReader(cxt.RespBodyString()))
	if err != nil{
		log.Println(err)
		return
	}
	result, err := dom.Find("div[id=result] .well").Html()
	if err != nil{
		log.Println(err)
	}
	// 固定顺序map
	gd := gt.NewGDMap().Add("start_ip", cxt.Task.GetDataStr("start_ip"))
	gd.Add("end_ip", cxt.Task.GetDataStr("end_ip"))
	gd.Add("number", cxt.Task.GetDataStr("number")).Add("result", result)

	// 保存抓取数据
	err = db.InsertAtGd("ip_result", gd)
	if err != nil {
		panic(err)
	}
}

// 获取详情信息重试的处理
func GetIPRetry(c *gt.Context){
	//更换代理
	//c.SetProxy(uri)

	// or
	c.Client = &http.Client{
		// 设置代理
		//Transport: &http.Transport{
		//	Proxy: http.ProxyURL(uri),
		//},
		Timeout: 5*time.Second,
	}

	log.Println("休息1s")
	time.Sleep(1*time.Second)
}

// 获取详情信息失败执行返还给队列
func GetIPFailed(c *gt.Context){
	queue.Add(c.Task)//请求失败归还到队列
}

```

抓取结果如图:

![](http://mange1.oss-cn-beijing.aliyuncs.com/test/effd5a6d3e8b26ee994f3c4e76efebed.png "")

---

### 5.4 更多实例子  [见](https://github.com/mangenotwork/gathertool/tree/main/_examples/) 
-  [阳光高考招生章程抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/get_yggk)
-  [ip地址信息抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/ip_bczs_cn)
-  [文件下载](https://github.com/mangenotwork/gathertool/tree/main/_examples/upload_file)
-  [无登录微博抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/weibo)
-  [百度题库抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/baidu_tk)
- ...


## 六、Mysql存储篇
gathertool 基于 "github.com/go-sql-driver/mysql" + "database/sql" 封装了数据操作的方法；

### 6.1 例子
```
// =========================== 数据库初始化
var (
	host2   = "192.168.0.2"
	host3   = "192.168.0.3"
	port    = 3306
	user2      = "root"
	password2  = "root123"
    user3      = "root3"
	password3  = "root333"

	db1  = "db1"
	DB1,_ = gt.NewMysql(host2, port, user2, password2, db1)

	db2  = "db2"
	DB2,_ = gt.NewMysql(host2, port, user2, password2, db2)

        db3  = "db3"
	DB3,_ = gt.NewMysql(host3, port, user3, password3, db3)

)

// 不存在表table1 存入json key为字段，value为值， 非固定顺序字段
jsonStr := `{"a":"a","b":"b"}`
db1.InsertAtJson("table1", jsonStr)

// 不存在表table2 存入非固定顺序字段的数据(map是无序的)
data1 := map[string]interface{}{"a":"a", "b":1}
db2.InsertAt("table2", data1)

// 不存在表table3 存入固定顺序字段的数据  (a,b)
gd := gt.NewGDMap().Add("a", 1).Add("b", 2)
db3.InsertAt("table3", gd)

// 存在表table2 存入数据
db2.Insert("table2", data1)
......
```

#### var MysqlDB
全局mysql客户端对象

---
#### type Mysql struct
mysql客户端结构体

---
#### func NewMysqlDB(host string, port int, user, password, database string) (err error)
初始化mysql客户端对象并进行连接
```
gt.NewMysqlDB(host2, port, user2, password2, db1)
gt.MysqlDB.Insert("table2", data1)
```

---
#### func NewMysql(host string, port int, user, password, database string) (\*Mysql, error)
创建一个mysql客户端对象

---
#### func GetMysqlDBConn() (\*Mysql,error)
获取mysql 连接, 前置需要初始化 NewMysqlDB

---
#### func (m \*Mysql) CloseLog()
关闭终端日志打印

---
#### func (m \*Mysql) SetMaxOpenConn(number int)
设置最大连接数

---
#### func (m \*Mysql) SetMaxIdleConn(number int)
最大等待连接中的数量

---
#### func (m \*Mysql) Conn() (err error)
连接mysql

---
#### func (m \*Mysql) IsHaveTable(table string) bool
表是否存在

---
#### type TableInfo struct
表信息

---
#### func (m \*Mysql) Describe(table string) ( \*tableDescribe, error)
查看表结构

---
#### func (m \*Mysql) Select(sql string) ([]map[string]string, error)
查询语句 返回 map

---
#### func (m \*Mysql) NewTable(table string, fields map[string]string) error
创建表; 字段顺序不固定;
fields  字段:类型； name:varchar(10);

---
#### func (m \*Mysql) NewTableGd(table string, fields \*gDMap) error
创建新的固定map顺序为字段的表(见下文固定map)

---
#### func (m \*Mysql) Insert(table string, fieldData map[string]interface{}) error
新增数据; table是存在的

---
#### func (m \*Mysql) InsertAt(table string, fieldData map[string]interface{}) error
新增数据; 如果没有表则先创建表; 字段顺序不是固定的;

---
#### func (m \*Mysql) InsertAtGd(table string, fieldData \*gDMap) error
新增数据; 如果没有表则先创建表; 字段顺序是固定的;

---
#### func (m \*Mysql) InsertAtJson(table, jsonStr string) error
json字符串存入数据库; 如果没有表则先创建表; 字段顺序不是固定的;

---
#### func (m \*Mysql) Update(sql string) error
Update

---
#### func (m \*Mysql) Exec(sql string) error
Exec

---
#### func (m \*Mysql) Query(sql string) ([]map[string]string, error)
Query

---
#### func (m \*Mysql) Delete(sql string) error
Delete

---
#### func (m \*Mysql) ToVarChar(data interface{}) string
写入mysql 的字符类型

---
#### func (m \*Mysql) DeleteTable(tableName string) error
删除表

---
#### func (m \*Mysql) HasTable(tableName string) bool
判断表是否存

---
#### func (m \*Mysql) ToXls(sql, outPath string)
查询数据导出到excel
```
db.ToXls("selecr * from tabe1", "table1.xls")
```

## 七、接口测试篇

gathertool可以用于接口测试和压力测试，极大的提升了测试编码效率。

#### type StressUrl struct

压力测试一个url

```
type StressUrl struct {
	Url string
	Method string
	Sum int64
	Total int
	TQueue TodoQueue

	// 请求时间累加
	sumReqTime int64

	// 测试结果
	avgReqTime time.Duration

	// 接口传入的json
	JsonData string

	// 接口传入类型
	ContentType string

	stateCodeList []*stateCodeData
	stateCodeListMux *sync.Mutex
}
```

---
#### func NewTestUrl(url, method string, sum int64, total int) \*StressUrl
实例化一个新的url压测

---
#### func (s \*StressUrl) SetJson(str string)
设置json数据

---
#### func (s \*StressUrl) Run(vs ...interface{})
运行压测


#### 例子：

```
import (
	gt "github.com/mangenotwork/gathertool"
	"log"
)

func main() {
	gt.CPUMax()
    // 普通 GET api压测
	url := "http://192.168.0.9:8002/v1/health"
	// 请求10000次 并发数 1000
	test := gt.NewTestUrl(url,"Get",10000,1000)
	test.Run()
	test.Run(gt.SucceedFunc(func(ctx *gt.Context){
		log.Println(ctx.JobNumber, "测试完成!!", ctx.Ms)
	}))
}
```

结果:
```
2022-03-28 16:17:09 |Info  | 执行次数 : 10000
2022-03-28 16:17:09 |Info  | 状态码分布: map[200:10000]
2022-03-28 16:17:09 |Info  | 平均用时: 35.555388957699996ms
2022-03-28 16:17:09 |Info  | 最高用时: 235.268968ms
2022-03-28 16:17:09 |Info  | 最低用时: 0.123935ms
2022-03-28 16:17:09 |Info  | 执行完成！！！

```

其他例子，碎片话代码，仅供参考:
```
/// 设置 GET Header 的压测
	////url2 := "http://192.168.0.9:18084/static_service/v1/auth/video/page"
	//url2 := "http://192.168.0.9:18084/static_service/v1/auth/quality_article/list"
	//token := &http.Header{}
	//token.Add("token", tokenStr)
	//token.Add("source", "2")
	//test2 := gt.NewTestUrl(url2,"Get",10000,1000)
	//test2.Run(token)


//// Post
	//url4 := "http://ggzyjy.sc.gov.cn/WebBuilder/frontAppAction.action?cmd=addPageView"
	//test4 := gt.NewTestUrl(url4,"Post",100,10)
	//test4.SetJson(`{
	//			"viewGuid":"cms_002",
	//			"siteGuid":"7eb5f7f1-9041-43ad-8e13-8fcb82ea831a"
	//			}`)
	//test4.Run(gt.SucceedFunc(func(c *gt.Context) {
	//	log.Println(string(c.RespBody))
	//	log.Println(c.Resp.Cookies())
	//}))


//点赞
	//npplike := "http://192.168.0.9:8025/v2/like/add"
	//npplike()

	//caseurl := "https://www.jy863.com:18443/"
	//caseurl := "https://www.uks678.com:10442/api/site/gdnotice/list"
	//npptest1 := gt.NewTestUrl(caseurl,"Get",10000,1000)
	//hd := &http.Header{}
	//hd.Add("requested-device", "APP")
	//hd.Add("requested-language", "CN")
	//hd.Add("requested-site", "www.uks678.com:10442")
	//npptest1.Run(hd)


caseUrl := `http://192.168.0.9:8025/v2/topic/create`

	gt.NewMysqlDB(host, port, user, password, database)
	gt.MysqlDB.CloseLog()
	conn,err := gt.GetMysqlDBConn()
	if err != nil {
		log.Panic("数据库初始化失败")
	}
	allUid, _ := conn.Select("select uid from tbl_user_ext limit 100,10000")
	for _,v := range allUid{
		token := GetToken(v["uid"])
		h := &http.Header{}
		h.Add("token", token)
		data := `
{
	"uid": `+v["uid"]+`,
	"qi": 2021315,
	"fid": 1,
	"source": 1,
	"content": "我爱牛票票"
}`
		ctx,err := gt.PostJson(caseUrl, data, h)
		log.Println(ctx.RespBodyString(), err)
	}


// 评论
func nppComment() {
	caseUrl := "http://192.168.0.9:8025/v1/comment/add"
	gt.NewMysqlDB(host, port, user, password, database)
	gt.MysqlDB.CloseLog()
	conn,err := gt.GetMysqlDBConn()
	if err != nil {
		log.Panic("数据库初始化失败")
	}
	allUid, _ := conn.Select("select uid from tbl_user_ext")
	allTid, _ := conn.Select("select tid from tbl_topic")
	i:=0
	for _,v := range allUid{
		for _, t := range allTid {
			token := GetToken(v["uid"])
			h := &http.Header{}
			h.Add("token", token)
			i++
			if i%3 == 0 {
				continue
			}
			data := `
{
	"content":"我们都爱牛票票",
	"tid": `+t["tid"]+`,
	"cmtid":0,
	"uid": `+v["uid"]+`,
	"source":1
}`
			gt.PostJson(caseUrl, data, h)
		}

	}
}

```

---

## 八、加密解密篇

加密解码相关封装方法

```
const (
	CBC = "CBC"
	ECB = "ECB"
	CFB = "CFB"
	CTR = "CTR"
)
```

---
#### type AES interface
```
type AES interface {
	Encrypt(str, key []byte) ([]byte, error)
	Decrypt(str, key []byte) ([]byte, error)
}
```

---
#### type DES interface
```
type DES interface {
	Encrypt(str, key []byte) ([]byte, error)
	Decrypt(str, key []byte) ([]byte, error)
}
```

---
#### func NewAES(typeName string, arg ...[]byte) AES
use NewAES("cbc")
typeName:
- "cbc", "Cbc","CBC"
- "ecb", "Ecb", ECB
- "cfb", "Cfb", CFB
- "ctr", "Ctr", CTR

---
#### func NewDES(typeName string, arg ...[]byte) DES
use NewAES("cbc")
typeName:
- "cbc", "Cbc","CBC"
- "ecb", "Ecb", ECB
- "cfb", "Cfb", CFB
- "ctr", "Ctr", CTR

---
#### func (cbc *cbcObj) Encrypt(str, key []byte) ([]byte, error)
AES CBC Encrypt

---
#### func (cbc *cbcObj) Decrypt(str, key []byte) ([]byte, error)
AES CBC Decrypt

---
#### func (ecb *ecbObj) Encrypt(str, key []byte) ([]byte, error)
AES ECB Encrypt

---
#### func (ecb *ecbObj) Decrypt(str, key []byte) ([]byte, error)
AES ECB Decrypt

---
#### func (cfb *cfbObj) Encrypt(str, key []byte) ([]byte, error)
AES CFB Encrypt

---
#### func (cfb *cfbObj) Decrypt(str, key []byte) ([]byte, error)
AES CFB Decrypt

---
#### func (ctr *ctrObj) Encrypt(str, key []byte) ([]byte, error)
AES CTR Encrypt

---
#### func (ctr *ctrObj) Decrypt(str, key []byte) ([]byte, error)
AES CTR Decrypt

---
#### func HmacMD5(str, key string) string
HmacMD5

---
#### func HmacSHA1(str, key string) string
HmacSHA1

---
#### func HmacSHA256(str, key string) string
HmacSHA256

---
#### func HmacSHA512(str, key string) string
HmacSHA512

---
#### func PBKDF2(str, salt []byte, iterations, keySize int) ([]byte)
PBKDF2

---
#### func JwtEncrypt(data map[string]interface{}, secret, method string) (string, error)
Jwt Encrypt

---
#### func JwtEncrypt256(data map[string]interface{}, secret string) (string, error)
Jwt Encrypt 256

---
#### func JwtEncrypt384(data map[string]interface{}, secret string) (string, error)
Jwt Encrypt 384

---
#### func JwtEncrypt512(data map[string]interface{}, secret string) (string, error)
Jwt Encrypt 512

---
#### func JwtDecrypt(tokenString, secret string) (data map[string]interface{}, err error)
Jwt Decrypt

---

## 九、Redis篇
gathertool的redis方法是基于"github.com/garyburd/redigo/redis"再次封装的,再加上了ssh连接通道，实现了对云端redis的连接; 主要常用如：快速删除大量指定key等。

### 9.1 连接

#### type Rds struct

```
type Rds struct {
	SSHUser string
	SSHPassword string
	SSHAddr string
	RedisHost string
	RedisPost string
	RedisPassword string

	// redis DB
	RedisDB int

	// 单个连接
	Conn redis.Conn

	//	最大闲置数，用于redis连接池
	RedisMaxIdle int

	//	最大连接数
	RedisMaxActive int

	//	单条连接Timeout
	RedisIdleTimeoutSec int

	// 连接池
	Pool *redis.Pool
}
```

---
#### func (r \*Rds) RedisConn() (err error)
redis连接

---
#### func (r \*Rds) RedisPool() error
连接池连接, 返回redis连接池  *redis.Pool.Get() 获取redis连接

---
#### func (r \*Rds) GetConn() redis.Conn
获取连接

---
#### func (r \*Rds) SelectDB(dbNumber int) error
切换db

---
#### func NewSSHInfo( addr, user, password string) \*SSHConnInfo
连接ssh

---
#### func NewRedis(host, port, password string, db int, vs ...interface{}) (\*Rds)
实例化redis对象

---
#### func NewRedisPool(host, port, password string, db, maxIdle, maxActive, idleTimeoutSec int, vs ...interface{}) (\*Rds)
实例化redis连接池对象

---
#### func RedisDELKeys(rds \*Rds, keys string, jobNumber int)
并发删除key；keys 模糊key; jobNumber 并发数;
```
rds := gt.NewRedisPool(redis_host, redis_port, redis_password, dbnumber, 5, 10, 10,
		gt.NewSSHInfo(ssh_addr, ssh_user, ssh_password))

	gt.RedisDELKeys(rds, "user:*", 500)
```

### 9.2 keys

### 9.3 string

### 9.4 list

### 9.5 hash

### 9.6 set

### 9.7 zset

## 十、
















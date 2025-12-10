## 简介
- gathertool是golang脚本化开发集成库，目的是提高爬虫，数据处理，接口测试，模拟请求等对应场景脚本程序开发的效率；
- gathertool也是一款轻量级爬虫库，特色是抽象了请求事件，通俗点理解就是对请求过程状态进行事件处理。
- gathertool也是接口测试&压力测试库，在接口测试脚本开发上有明显的效率优势，
- gathertool还集成了对三方中间件的操作，DB操作等。

---

## 使用
- go get github.com/mangenotwork/gathertool

---

## 介绍
 gathertool是一个高度封装工具库，包含了http/s的请求，Mysql数据库方法，数据类型处理方法，数据提取方法，websocket相关方法，
TCP|UDP相关方法，NoSql相关方法，开发常用方法等;  可以用于爬虫程序，接口&压力测试程序，常见网络协议调试程序，数据提取与存储程序等； 

- gathertool的请求特点: 会在请求阶段执行各个事件如请求失败后的重试事件,请求前后的事件，请求成功事件等等, 可以根据请求状态码自定义这些事件；

- gathertool还拥有很好的可扩展性， 适配传入任意自定义http请求对象， 能适配各种代理对象等等；

- gathertool还拥有抓取数据存储功能, 存储到mysql；配合 https://github.com/mangenotwork/dbHelper 可存储到更多的数据存储上; 还有很多创新的地方文档会根据具体方法进行介绍；

- gathertool还封装了消息队列接口，支持Nsq,Kafka,rabbitmq,redis等消息队列中间件

## 文档： 
[开发文档](https://github.com/mangenotwork/gathertool/blob/main/_doc/develop.md)

[pkg.go.dev](https://pkg.go.dev/github.com/mangenotwork/gathertool)


## 开始使用
> go get github.com/mangenotwork/gathertool


#### 简单的get请求
```go
import gt "github.com/mangenotwork/gathertool"

func main(){
    ctx, err := gt.Get("https://www.baidu.com")
    if err != nil {
        log.Println(err)
    }
    log.Println(ctx.RespBodyString)
}
```

#### 含请求事件请求

- StartFunc 开始事件，请求开始前执行的方法类型

- SucceedFunc 成功事件，请求成功后执行的方法类型

- FailedFunc 失败事件，请求失败执行的方法类型

- RetryFunc 重试事件，重试请求前执行的方法类型，是否重试来自上次请求的状态码来确定，见StatusCodeMap;

- EndFunc 结束事件，请求流程结束后执行的方法类型

```go
import gt "github.com/mangenotwork/gathertool"

func main(){

    gt.NewGet(`http://192.168.0.1`).SetStartFunc(func(ctx *gt.Context){
            log.Println("请求前： 添加代理等等操作")
            ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }
    ).SetSucceedFunc(func(ctx *gt.Context){
            log.Println("请求成功： 处理数据或存储等等")
            log.Println(ctx.RespBodyString())
        }
    ).SetFailedFunc(func(ctx *gt.Context){
            log.Println("请求失败： 记录失败等等")
            log.Println(ctx.Err)
        }
    ).SetRetryFunc(func(ctx *gt.Context){
             log.Println("请求重试： 更换代理或添加等待时间等等")
             ctx.Client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
        }
    ).SetEndFunc(func(ctx *gt.Context){
             log.Println("请求结束： 记录结束，处理数据等等")
             log.Println(ctx)
        }
    ).Do()
    
}
```

#### 事件方法复用
```go
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
```go
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
```go
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
```go
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
```go
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

......

## 实例
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

## JetBrains 开源证书支持

`gathertool` 项目一直以来都是在 JetBrains 公司旗下的 `GoLand` 集成开发环境中进行开发，基于 **free JetBrains Open Source license(s)** 正版免费授权，在此表达我的谢意。

<a href="https://www.jetbrains.com/?from=gathertool" target="_blank"><img src="https://raw.githubusercontent.com/moonD4rk/staticfiles/master/picture/jetbrains-variant-4.png" width="256" align="middle"/></a>


## 使用注意&不足之处
1. Mysql相关方法封装使用的字符串拼接，如抓取无需关心Sql注入这类场景放心使用，生产环境使用会出现安全风险；
2. 变量类型转换使用反射，性能方面会差一些，追求性能请使用其他；


## 经验积累
从2022年v0.3.4版本到至今已经应用在我的实际工作中，给边边角角的工作带来了巨大的脚步编写效率，经过不断的积累和总结质量也在不断的提示；感兴趣的伙伴可以fork代码，
将你的宝贵经验提交进来。

## 配合使用
- 数据存储与读取 https://github.com/mangenotwork/dbHelper
- 图片处理  https://github.com/mangenotwork/imgHelper

## TODO
- Redis连接方法改为连接池, 干掉redis队列
- 关闭重试
- 提供通过html定位加正则获取内容的方法
- 引入请求工作与队列的功能，在广度爬取场景使用

## BUG
- 固定map调用不了  .\main.go:124:16: undefined: gt.OrderMap
- panic
```azure
panic: runtime error: index out of range [1] with length 1

goroutine 1 [running]:
github.com/mangenotwork/gathertool.RegFindAllTxt({0xb042f8?, 0x1?}, {0xc00020c000, 0x3fe5})
        D:/go/pkg/mod/github.com/mangenotwork/gathertool@v0.4.7/reg.go:32 +0x135
```


## 三方引用 感谢这些开源项目
- github.com/Shopify/sarama 
- github.com/dgrijalva/jwt-go 
- github.com/go-sql-driver/mysql 
- github.com/nsqio/go-nsq 
- github.com/streadway/amqp 
- github.com/xuri/excelize/v2 
- golang.org/x/crypto 
- golang.org/x/net 
- golang.org/x/text 
- gopkg.in/gomail.v2
- gopkg.in/yaml.v3

## 里程碑

#### v0.4.8
```
1. 升级到go版本1.19并进行优化
2. 多类型重复的方法改为泛型，减轻代码量  <todo>
3. 增加测试代码和大量测试与整理实例 <todo>
4. 优化项目文件结构，规范库的注释
5. 更新和优化文档  <todo>
6. 修复 panic <todo>
7. 修复问题: <todo>
  - Any2Float64 不支持字符串转float64 <todo>
  - 下载的函数名错误 Upload 执行GET下载请求 <todo>
  - 缺少url提取文件名，文件后缀的函数 <todo>
  - 类型转换缺少空值判断 <todo>
8. 新增:  <todo>
  - 正则提取增加  section, dl, dt, dd 标签 <todo>
  - 去掉 gorm.io/gorm, mongo-driver, github.com/gomodule/redigo
  - html解析需要引入属性值包含关系 <todo>
  - 下载资源需要打印下载url, 以及检查下载大小，如果是0则抛出错误 <todo>
```

#### v0.4.7
```
1. Redis连接方法改为连接池
2. 增加关闭重试方法
3. 增加Whois查询
4. 测试与优化
```

#### v0.4.4 ~ v0.4.6
```
1. 移除 grpc相关方法
2. 新增DNS查询
3. 新增证书信息获取
4. 新增Url扫描
5. 新增邮件发送
6. 优化ICMP Ping的方法
```

#### v0.4.3
```
1. 新增网站链接采集应用场景的方法 
2. 修复v0.4.2的json导入bug 
3. 修复redis封装方法的bug 
4. 请求日志添加请求地址信息
5. 优化抓取实例
```

#### v0.4.2
```
新增:
1. redis, nsq, rabbitmq, kafka 消息队列方法
2. 新增开发文档
3. 新增redis相关方法
```

#### v0.4.1
```
新增:
1. 文件相关处理
2. 文件压缩解压
3. 新增抓取实列 _examples/cnlinfo
```

#### v0.2.1 ~ v0.3.9
```
旧的迭代版本，不稳定，请使用 > v0.3.4的版本
- 新增代理与抓包
- 启动一个代理服务，并抓包对数据包进行处理
- 新增socket5代理

v0.3.6 
新增:
1. func SimpleServer(addr string, g *grpc.Server) : grpc服务端
2. func NewTask() *Task : task的实例方法
3. func (task *Task) SetUrl(urlStr string) *Task & func (task *Task) SetJsonParam(jsonStr string) *Task

修复:
1. Context对象，Task为空的问题

v0.3.7
新增:
1. 新增html解析，指定html内容提取
2. 新增抓取实例
3. 优化部分方法

v0.3.9
新增:
1. 新增配置，支持yaml
2. 优化部分方法
```




# gathertool
gathertool是golang脚本化开发库，目的是提高对应场景程序开发的效率；
轻量级爬虫库，接口测试&压力测试库，DB操作库。
请使用最新版本!!!

## 使用场景
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


## 文档： 
[https://pkg.go.dev 点击开始](https://pkg.go.dev/github.com/mangenotwork/gathertool)
[http://mange.work/doc?id=1 点击开始](http://mange.work/doc?id=1)

## 开始使用
> go get github.com/mangenotwork/gathertool


简单的get请求
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

含请求事件请求
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

事件方法复用
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

post请求
```
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

数据存储到mysql
```
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

HTML数据提取
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

更多方法见 [文档](https://github.com/mangenotwork/gathertool/tree/main/_doc/doc_1.md)

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

## 里程碑

#### v0.2.1 ~ v0.2.9
```
2021.10 ~ 2022.3
库的基本完成与完善，用于正式工作环境
```

#### v0.3.4
```
新增代理与抓包
启动一个代理服务，并抓包对数据包进行处理
```

#### v0.3.5
```
新增socket5代理
```

#### v0.3.6
```
新增:
1. func SimpleServer(addr string, g *grpc.Server) : grpc服务端
2. func NewTask() *Task : task的实例方法
3. func (task *Task) SetUrl(urlStr string) *Task & func (task *Task) SetJsonParam(jsonStr string) *Task

修复:
1. Context对象，Task为空的问题

```

#### v0.3.7
```
新增:
1. 新增html解析，指定html内容提取
2. 新增抓取实例
3. 优化部分方法
4. 添加注释
```

#### v0.3.9
```
新增:
1. 新增配置，支持yaml

优化部分方法
代码审查
```

#### v0.4.1
```
新增:
1. 文件相关处理
2. 文件压缩解压
3. 新增抓取实列 _examples/cnlinfo

优化部分方法
代码审查

```

#### v0.4.2
```shell
新增:
1. 并发下载
2. Rabbit, RC4, RIPEMD-160 加密解密
3. 日志路径打印两个段
4. 暴力登录 web admin 场景
5. ssh暴力登录

```

## TODO List
- 三方队列（redis, mysql, nsq）

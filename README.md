# gathertool
轻量级爬虫，接口测试，压力测试框架, 提高开发对应场景的golang程序的效率。

## 文档： [点击开始](https://github.com/mangenotwork/gathertool/tree/main/_doc/doc_1.md)

## 开始
> go get github.com/mangenotwork/gathertool

```go
import gt "github.com/mangenotwork/gathertool"

func main(){
    caseUrl := "https://www.baidu.com"
    ctx, err := gt.Get(caseUrl)
    if err != nil {
        log.Println(err)
    }
    log.Println(ctx.RespBodyString)
}
```

> 含请求事件请求
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

> 优雅的写法
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

## 实例
-  [Get请求](https://github.com/mangenotwork/gathertool/tree/main/_examples/get)
-  [阳光高考招生章程抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/get_yggk)
-  [ip地址信息抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/ip_bczs_cn)
-  [压力测试](https://github.com/mangenotwork/gathertool/tree/main/_examples/stress_testing)
-  [文件下载](https://github.com/mangenotwork/gathertool/tree/main/_examples/upload_file)
-  [无登录微博抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/weibo)
-  [百度题库抓取](https://github.com/mangenotwork/gathertool/tree/main/_examples/baidu_tk)

### BUG
- MysqlDB.NewTable() 字段参数是map, 创建的表会乱序
-  client.Transport = &http.Transport{
  	 	MaxIdleConnsPerHost: 1000,
      }
  解决连接数高的问题
  
###  常见的反爬虫策略
- User-Agent反爬
- IP访问频率限制
- 强制登录
- 漏桶、令牌桶之类的算法限制接口访问频率
- 数据接口加上庞大的加密解密和混淆算法
- JS代码执行后，会发送一个带参数key的请求，后台通过判断key的值来决定是响应真实的页面，还是响应伪造或错误的页面。因为key参数是动态生成的，每次都不一样，难以分析出其生成方法，使得无法构造对应的http请求
-


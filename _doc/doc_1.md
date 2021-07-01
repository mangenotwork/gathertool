# gathertool 开发文档

## HTTP/S 请求所有参数说明与上下文结构
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
- IsLog : 是否打印日志

## 状态码事件
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

## 一、 get请求
### 1. Get(url string) get请求，返回上下文；

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


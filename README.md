# gathertool
轻量级爬虫，接口测试，压力测试框架, 提高开发对应场景的golang程序的效率。
文档： [点击开始](https://github.com/mangenotwork/gathertool/tree/main/_doc/doc_1.md)

## 概念
1. 一个请求含有 请求成功后方法， 请求重试前方法，请求失败后方法；举个例子 
 请一个URL  请求遇到状态200或302等算成功，则执行成功后方法。  请求遇到403或502
 等需要执行重试，则执行重试方法，重试方法主要含添加等待时间更换代理IP等。 遇到404
 或500等算失败则执行失败方法。
 
2. 并发执行，第一需要创建TODO队列，等TODO队列加载完后，每个队列对象含有上下文，
在创建时应该富裕上文数据或对象，开始执行并发任务，每个并发任务是一个独立的cilet，
当队列任务取完后则整个并发结束。注意这里的每个并发任务都是独立的，没有chan操作。


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



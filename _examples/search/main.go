/*
	1. ip 扫描
	2. DOS 攻击
	3. DDOS 攻击
	4. arp
	5. ip的端口扫描

 */
package main

import (
	gt "github.com/mangenotwork/gathertool"
	"time"
)

func main(){
	//gt.SearchDomain("110.242.68.3")
	//gt.SearchPort("120.79.88.59", 5*time.Second)
	//gt.Ping("192.168.0.23")
	//
	//for i:=0; i<256; i++{
	//	gt.Ping(fmt.Sprintf("192.168.0.%d", i))
	//}

	for {
		gt.Ping("120.79.88.59")
		time.Sleep(1*time.Second)
	}

}



/*
	搜索引擎采集策略
	- 种子url 扩增扫描
	- ip 扫描到站点
	- 各大搜索引擎抓取

	内容识别字段:  title, description, <h1>...<h6>, keywords, site_name,  property="og:...,
 */

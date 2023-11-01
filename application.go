/*
	Description : 一些常用的应用场景函数
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"fmt"
	"strings"
	"time"
)

var ApplicationTerminalOut = true

func ApplicationTerminalOutClose() {
	ApplicationTerminalOut = false
}

// HostToolEr TODO HostToolEr
type HostToolEr interface {
	Run() ([]string, int)
}

// HostScanUrl ========================================================================================================
// Host站点下 A标签 Url扫描， 从更目录开始扫描指定深度 get Url 应用函数
type HostScanUrl struct {
	Host     string
	Depth    int // 页面深度
	UrlSet   map[string]struct{}
	Count    int
	MaxCount int64
}

func NewHostScanUrl(host string, depth int) *HostScanUrl {
	return &HostScanUrl{
		Host:     host,
		Depth:    depth,
		UrlSet:   make(map[string]struct{}),
		Count:    0,
		MaxCount: 0,
	}
}

func (scan *HostScanUrl) Run() ([]string, int) {
	CloseLog()
	scan.do(scan.Host, 0)
	urls := make([]string, 0)
	for k, _ := range scan.UrlSet {
		urls = append(urls, urlStr(k))
	}
	if ApplicationTerminalOut {
		fmt.Printf("\n")
	}
	return urls, scan.Count
}

func (scan *HostScanUrl) do(caseUrl string, df int) {
	if len(caseUrl) < 1 {
		return
	}
	if df > scan.Depth {
		return
	}
	// 如果不是host下的域名
	if strings.Index(caseUrl, scan.Host) == -1 {
		if string(caseUrl[0]) == "/" {
			caseUrl = scan.Host + caseUrl
			goto G
		}
		return
	}
G:
	if _, ok := scan.UrlSet[caseUrl]; ok {
		return
	}
	ctx, err := Get(caseUrl)
	if err != nil {
		Error(err)
		return
	}
	if ApplicationTerminalOut {
		fmt.Print(".")
	}
	df++
	scan.UrlSet[caseUrl] = struct{}{}
	scan.Count++
	a := RegHtmlA(ctx.Html)
	for _, v := range a {
		links := RegHtmlHrefTxt(v)
		if len(links) < 1 {
			continue
		}
		link := links[0]
		// 请求并验证
		scan.do(link, df)
	}
}

// HostScanExtLinks ===================================================================================================
// Host站点下的外链采集 应用函数
type HostScanExtLinks struct {
	Host string
}

func NewHostScanExtLinks(host string) *HostScanExtLinks {
	return &HostScanExtLinks{
		Host: host,
	}
}

func (scan *HostScanExtLinks) Run() ([]string, int) {
	var (
		urlSet = make(map[string]struct{})
		count  int
	)
	CloseLog()
	ctx, err := Get(scan.Host)
	if err != nil {
		Error(err)
		return []string{}, 0
	}
	a := RegHtmlA(ctx.Html)
	for _, v := range a {
		links := RegHtmlHrefTxt(v)
		if len(links) < 1 {
			continue
		}
		link := links[0]
		if strings.Index(link, scan.Host) == -1 && string(link[0]) != "/" && string(link[0]) != "#" {
			if _, ok := urlSet[link]; ok {
				continue
			}
			if UsefulUrl(link) {
				urlSet[link] = struct{}{}
				count++
			}
		}
	}
	urls := make([]string, 0)
	for k, _ := range urlSet {
		urls = append(urls, k)
	}
	return urls, count
}

// HostScanBadLink ====================================================================================================
// Host站点下 HTML Get Url 死链接扫描 应用函数
type HostScanBadLink struct {
	Host      string
	Depth     int // 页面深度
	PageState map[string]int
	UrlSet    map[string]struct{} // 检查页面重复
}

func NewHostScanBadLink(host string, depth int) *HostScanBadLink {
	return &HostScanBadLink{
		Host:      host,
		Depth:     depth,
		PageState: make(map[string]int),
		UrlSet:    make(map[string]struct{}),
	}
}

func (scan *HostScanBadLink) Run() ([]string, int) {
	CloseLog()
	scan.do(scan.Host, 0)
	badUrls := make([]string, 0)
	for k, v := range scan.PageState {
		if v == 404 || v == 410 {
			badUrls = append(badUrls, k)
		}
	}
	if ApplicationTerminalOut {
		fmt.Printf("\n")
	}
	return badUrls, len(badUrls)
}

func (scan *HostScanBadLink) Result() map[string]int {
	res := make(map[string]int)
	for k, v := range scan.PageState {
		if v == 404 || v == 410 {
			res[k] = v
		}
	}
	return res
}

func (scan *HostScanBadLink) Report() map[string]int {
	return scan.PageState
}

func (scan *HostScanBadLink) do(caseUrl string, df int) {
	if df > scan.Depth {
		return
	}
	if len(caseUrl) < 1 {
		return
	}
	if string(caseUrl[0]) == "/" {
		caseUrl = scan.Host + caseUrl
	}
	if _, ok := scan.UrlSet[caseUrl]; ok {
		return
	}
	ctx, err := Get(caseUrl)
	if err != nil {
		ctx.StateCode = 404
	}
	if ApplicationTerminalOut {
		fmt.Print(".")
	}
	df++
	scan.UrlSet[caseUrl] = struct{}{}
	scan.PageState[caseUrl] = ctx.StateCode
	a := RegHtmlA(ctx.Html)
	for _, v := range a {
		links := RegHtmlHrefTxt(v)
		if len(links) < 1 {
			continue
		}
		link := links[0]
		// 请求并验证
		scan.do(link, df)
	}
}

// HostPageSpeedCheck =================================================================================================
// Host站点下 HTML Get 测速 应用函数
type HostPageSpeedCheck struct {
	Host      string
	Depth     int // 页面深度
	PageSpeed map[string]time.Duration
	UrlSet    map[string]struct{} // 检查页面重复
}

func NewHostPageSpeedCheck(host string, depth int) *HostPageSpeedCheck {
	return &HostPageSpeedCheck{
		Host:      host,
		Depth:     depth,
		PageSpeed: make(map[string]time.Duration),
		UrlSet:    make(map[string]struct{}),
	}
}

// Run int 单位 ms
func (scan *HostPageSpeedCheck) Run() ([]string, int) {
	CloseLog()
	scan.do(scan.Host, 0)
	urls := make([]string, 0)
	for k, v := range scan.PageSpeed {
		urls = append(urls, fmt.Sprintf("%s:%v", k, v))
	}
	if ApplicationTerminalOut {
		fmt.Printf("\n")
	}
	return urls, len(urls)
}

func (scan *HostPageSpeedCheck) Result() map[string]time.Duration {
	return scan.PageSpeed
}

func (scan *HostPageSpeedCheck) do(caseUrl string, df int) {
	if len(caseUrl) < 1 {
		return
	}
	if df > scan.Depth {
		return
	}
	if string(caseUrl[0]) == "/" {
		caseUrl = scan.Host + caseUrl
	}
	if _, ok := scan.UrlSet[caseUrl]; ok {
		return
	}
	if strings.Index(caseUrl, scan.Host) == -1 {
		return
	}
	ctx, err := Get(caseUrl)
	if err != nil {
		Error(err)
		return
	}
	if ApplicationTerminalOut {
		fmt.Print(".")
	}
	df++
	scan.UrlSet[caseUrl] = struct{}{}
	scan.PageSpeed[caseUrl] = ctx.Ms
	a := RegHtmlA(ctx.Html)
	for _, v := range a {
		links := RegHtmlHrefTxt(v)
		if len(links) < 1 {
			continue
		}
		link := links[0]
		// 请求并验证
		scan.do(link, df)
	}
}

// AverageSpeed 平均时间
func (scan *HostPageSpeedCheck) AverageSpeed() float64 {
	var (
		n int64 = 0
		t int64 = 0
	)
	for _, v := range scan.PageSpeed {
		n++
		t += v.Milliseconds()
	}
	return float64(t) / float64(n)
}

// MaxSpeed 最高用时
func (scan *HostPageSpeedCheck) MaxSpeed() int64 {
	var max int64 = 0
	for _, v := range scan.PageSpeed {
		if v.Milliseconds() > max {
			max = v.Milliseconds()
		}
	}
	return max
}

// MinSpeed 最低用时
func (scan *HostPageSpeedCheck) MinSpeed() int64 {
	var min int64 = 0
	for _, v := range scan.PageSpeed {
		if v.Milliseconds() < min {
			min = v.Milliseconds()
		}
	}
	return min
}

// Report 报告
func (scan *HostPageSpeedCheck) Report() map[string]string {
	report := make(map[string]string)
	report["avg"] = fmt.Sprintf("%f", scan.AverageSpeed())
	report["max"] = fmt.Sprintf("%d", scan.MaxSpeed())
	report["min"] = fmt.Sprintf("%d", scan.MinSpeed())
	return report
}

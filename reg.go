/*
	Description : 正则
	Author : ManGe
	Version : v0.1
	Date : 2021-04-30
*/

package gathertool

import (
	"regexp"
	"runtime"
	"strings"
)

// td:=`<td>(.*?)</td>`
// tdreg := regexp.MustCompile(td)
// tdList := tdreg.FindAllStringSubmatch(rest,-1)

func RegFindAll(regStr, rest string) [][]string {
	reg := regexp.MustCompile(regStr)
	List := reg.FindAllStringSubmatch(rest, -1)
	reg.FindStringSubmatch(rest)
	return List
}

var regMap = map[string]string{
	"RegHtmlA": `(?is:<a.*?</a>)`,
	"RegHtmlTitle": `(?is:<title.*?</title>)`,
	"RegHtmlTr": `(?is:<tr.*?</tr>)`,
	"RegHtmlInput": `(?is:<input.*?>)`,
	"RegHtmlTd": `(?is:<td.*?</td>)`,
	"RegHtmlP": `(?is:<p.*?</p>)`,
	"RegHtmlSpan": `(?is:<span.*?</span>)`,
	"RegHtmlSrc": `(?is:src=\".*?\")`,
	"RegHtmlHref": `(?is:href=\".*?\")`,
}

// 获取正在运行的函数名
func runFuncName()string{
	pc := make([]uintptr,1)
	runtime.Callers(2,pc)
	f := runtime.FuncForPC(pc[0])
	fName := f.Name()
	fList := strings.Split(fName,".")
	return fList[len(fList)-1]
}

// 执行正则提取
func regFind(funcName, txt string) (dataList []string) {
	regStr, ok := regMap[funcName]
	if !ok{
		loger("reg func is not")
		return
	}
	reg := regexp.MustCompile(regStr)
	resList := reg.FindAllStringSubmatch(txt, -1)
	for _,v := range resList{
		dataList = append(dataList, v[0])
	}
	return
}

func RegHtmlA(str string)[]string { return regFind(runFuncName(), str) }

func RegHtmlTitle(str string)[]string { return regFind(runFuncName(), str) }

func RegHtmlTr(str string)[]string { return regFind(runFuncName(), str) }

func RegHtmlInput(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlTd(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlP(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlSpan(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlSrc(str string) []string { return regFind(runFuncName(), str) }

func RegHtmlHref(str string) []string { return regFind(runFuncName(), str) }

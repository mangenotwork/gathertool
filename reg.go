/*
	Description : 正则
	Author : ManGe
	Version : v0.4
	Date : 2021-12-03
*/

package gathertool

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"unicode"
)

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
	"RegHtmlH1": `(?is:<h1.*?</h1>)`,
	"RegHtmlH2": `(?is:<h2.*?</h2>)`,
	"RegHtmlH3": `(?is:<h3.*?</h3>)`,
	"RegHtmlH4": `(?is:<h4.*?</h4>)`,
	"RegHtmlH5": `(?is:<h5.*?</h5>)`,
	"RegHtmlH6": `(?is:<h6.*?</h6>)`,
	"RegHtmlTbody": `(?is:<tbody.*?</tbody>)`,
	"RegHtmlVideo": `(?is:<video.*?</video>)`,
	"RegHtmlCanvas": `(?is:<canvas.*?</canvas>)`,
	"RegHtmlCode": `(?is:<code.*?</code>)`,
	"RegHtmlImg": `(?is:<img.*?/>)`,
	"RegHtmlUl": `(?is:<ul.*?</ul>)`,
	"RegHtmlLi": `(?is:<li.*?</li>)`,
	"RegHtmlMeta": `(?is:<meta.*?>)`,
	"RegHtmlSelect": `(?is:<select.*?</select>)`,
	"RegHtmlTable": `(?is:<table.*?</table>)`,
	"RegHtmlButton": `(?is:<button.*?</button>)`,

	// 取标签内容
	"RegHtmlATxt" : `(?is:<a.*?>(.*?)</a>)`,
	"RegHtmlTitleTxt" : `(?is:<title.*?>(.*?)</title>)`,
	"RegHtmlTrTxt": `(?is:<tr.*?>(.*?)</tr>)`,
	"RegHtmlTdTxt": `(?is:<td.*?>(.*?)</td>)`,
	"RegHtmlPTxt": `(?is:<p.*?>(.*?)</p>)`,
	"RegHtmlSpanTxt": `(?is:<span.*?>(.*?)</span>)`,
	"RegHtmlSrcTxt": `(?is:src=\"(.*?)\")`,
	"RegHtmlHrefTxt": `(?is:href=\"(.*?)\")`,
	"RegHtmlHTxt1": `(?is:<h1.*?>(.*?)</h1>)`,
	"RegHtmlHTxt2": `(?is:<h2.*?>(.*?)</h2>)`,
	"RegHtmlHTxt3": `(?is:<h3.*?>(.*?)</h3>)`,
	"RegHtmlHTxt4": `(?is:<h4.*?>(.*?)</h4>)`,
	"RegHtmlHTxt5": `(?is:<h5.*?>(.*?)</h5>)`,
	"RegHtmlHTxt6": `(?is:<h6.*?>(.*?)</h6>)`,
	"RegHtmlCodeTxt": `(?is:<code.*?>(.*?)</code>)`,
	"RegHtmlUlTxt": `(?is:<ul.*?>(.*?)</ul>)`,
	"RegHtmlLiTxt": `(?is:<li.*?>(.*?)</li>)`,
	"RegHtmlSelectTxt": `(?is:<select.*?>(.*?)</select>)`,
	"RegHtmlTableTxt": `(?is:<table.*?>(.*?)</table>)`,
	"RegHtmlButtonTxt": `(?is:<button.*?>(.*?)</button>)`,

	// 删除
	"RegDelHtml": `\<[\S\s]+?\>`,
	"RegDelNumber": `^[0-9]*$`,

	// 是否含有
	"IsNumber":        `^[0-9]*$`,
	"IsNumber2Len":    `[0-9]{%d}`,
	"IsNumber2Heard":  `^(%d)[0-9]*$`,
	"IsFloat":         `^(-?\d+\.\d+)?$`,
	"IsFloat2Len":     `^(-?\d+\.\d{%d})?$`,
	"IsEngAll":        `^[A-Za-z]*$`,
	"IsEngLen":        `^[A-Za-z]{%d}$`,
	"IsEngNumber":     `^[A-Za-z0-9]*$`,
	"IsLeastNumber":   `[0-9]{%d,}?`,
	"IsLeastCapital":  `[A-Z]{%d,}?`,
	"IsLeastLower":    `[a-z]{%d,}?`,
	"IsLeastSpecial":  `[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\&\$\.\*\+\?\{\}\(\)\[\]\|\!\_\@\#\%\-\=]{%d,}?`,
	"HaveNumber":      `[0-9]+`,
	"HaveSpecial":     `[\f\t\n\r\v\123\x7F\x{10FFFF}\\\^\&\$\.\*\+\?\{\}\(\)\[\]\|\!\_\@\#\%\-\=]+`,
	"IsEmail":         `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`,
	"IsDomain":        `[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(/.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+/.?`,
	"IsURL":           `//([\w-]+\.)+[\w-]+(/[\w-./?%&=]*)?$`,
	"IsPhone":         `^(13[0-9]|14[5|7]|15[0|1|2|3|5|6|7|8|9]|18[0|1|2|3|5|6|7|8|9])\d{8}$`,
	"IsLandline":      `^(\(\d{3,4}-)|\d{3.4}-)?\d{7,8}$`,
	"AccountRational": `^[a-zA-Z][a-zA-Z0-9_]{4,15}$`,
	"IsXMLFile":       `^*+\\.[x|X][m|M][l|L]$`,
	"IsUUID3":         `^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$`,
	"IsUUID4":         `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
	"IsUUID5":         `^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
	"IsRGB":           `^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$`,
	"IsFullWidth":     `[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]`,
	"IsHalfWidth":     `[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]`,
	"IsBase64":        `^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$`,
	"IsLatitude":      `^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$`,
	"IsLongitude":     `^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$`,
	"IsDNSName":       `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`,
	"IsIPv4":          `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`,
	"IsWindowsPath":   `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`,
	"IsUnixPath":      `^(/[^/\x00]*)+/?$`,
}

// 常用正则
const (
	RegTime = `(?i)\d{1,2}:\d{2} ?(?:[ap]\.?m\.?)?|\d[ap]\.?m\.?`
	RegLink = `(?:(?:https?:\/\/)?(?:[a-z0-9.\-]+|www|[a-z0-9.\-])[.](?:[^\s()<>]+|\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\))+(?:\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:\'".,<>?]))`
	RegEmail = `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`
	RegIPv4 = `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	RegIPv6 = `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`
	RegIP = RegIPv4 + `|` + RegIPv6
	RegMD5Hex = `[0-9a-fA-F]{32}`
	RegSHA1Hex = `[0-9a-fA-F]{40}`
	RegSHA256Hex = `[0-9a-fA-F]{64}`
	RegGUID = `[0-9a-fA-F]{8}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{12}`
	RegMACAddress = `(([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2}))`
	Email        = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	UUID3        = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	UUID4        = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID5        = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID         = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	Int          = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	Float        = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	RGBColor     = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"
	FullWidth    = "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	HalfWidth    = "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	Base64       = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	Latitude     = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	Longitude    = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	DNSName      = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	FullURL      = `^(?:ftp|tcp|udp|wss?|https?):\/\/[\w\.\/#=?&]+$`
	URLSchema    = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername  = `(\S+(:\S*)?@)`
	URLPath      = `((\/|\?|#)[^\s]*)`
	URLPort      = `(:(\d{1,5}))`
	URLIP        = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	URLSubdomain = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	WinPath      = `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
	UnixPath     = `^(/[^/\x00]*)+/?$`
)

var (
	ChineseNumber = []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "零"}
	ChineseMoney = []string{"壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
	ChineseMoneyAll = []string{"壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖", "拾", "佰", "仟", "万", "亿", "元", "角", "分", "零", "整", "正", "貳", "陸", "億", "萬", "圓"}
)

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

func RegHtmlVideo(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlCanvas(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlCode(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlImg(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlUl(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlLi(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlMeta(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlSelect(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlTable(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlButton(str string) []string { return regFind(runFuncName(), str)}

func RegHtmlH(str, typeH string) []string {
	funcName := runFuncName()
	return regFind(funcName+typeH, str)
}

func RegHtmlTbody(str string) []string { return regFind(runFuncName(), str) }


// 执行正则提取 只取内容
func regFindTxt(funcName, txt string) (dataList []string) {
	regStr, ok := regMap[funcName]
	if !ok{
		loger("reg func is not")
		return
	}
	reg := regexp.MustCompile(regStr)
	resList := reg.FindAllStringSubmatch(txt, -1)
	for _,v := range resList{
		if len(v) > 1{
			dataList = append(dataList, v[1])
		}
	}
	return
}

func RegHtmlATxt(str string)[]string { return regFindTxt(runFuncName(), str) }

func RegHtmlTitleTxt(str string)[]string { return regFindTxt(runFuncName(), str) }

func RegHtmlTrTxt(str string)[]string { return regFindTxt(runFuncName(), str) }

func RegHtmlInputTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlTdTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlPTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlSpanTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlSrcTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlHrefTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlCodeTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlUlTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlLiTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlSelectTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlTableTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlButtonTxt(str string) []string { return regFindTxt(runFuncName(), str) }

func RegHtmlHTxt(str, typeH string) []string {
	funcName := runFuncName()
	return regFindTxt(funcName+typeH, str)
}

// 删除正则匹配的字符
func replace(funcName, rest string) string {
	regStr, ok := regMap[funcName]
	if !ok{
		loger("reg func is not")
		return ""
	}
	re, err := regexp.Compile(regStr)
	if err != nil {
		loger(err)
		return ""
	}
	return re.ReplaceAllString(rest, "")
}

// 删除所有标签
func RegDelHtml(str string) string { return replace(runFuncName(), str) }

// 删除所有数字
func RegDelNumber(str string) string { return replace(runFuncName(), str) }


// 是否含有正则匹配的字符
func isHaveStr(regStr, rest string) bool {
	isHave, err := regexp.MatchString(regStr, rest)
	if err != nil {
		loger(err)
		return false
	}
	return isHave
}

// 是否含有正则匹配的字符
func isHave(funcName, rest string) bool {
	regStr, ok := regMap[funcName]
	if !ok{
		loger("reg func is not")
		return false
	}
	return isHaveStr(regStr, rest)
}



// 验证是否含有number
func IsNumber(str string) bool { return isHave(runFuncName(), str) }

// 验证是否含有连续长度不超过长度l的number
func IsNumber2Len(str string, l int) bool {
	regStr, ok := regMap[runFuncName()]
	if !ok{
		loger("reg func is not")
		return false
	}
	return isHaveStr(fmt.Sprintf(regStr, l), str)
}

// 验证是否含有n开头的number
func IsNumber2Heard(str string, n int) bool {
	regStr, ok := regMap[runFuncName()]
	if !ok{
		loger("reg func is not")
		return false
	}
	return isHaveStr(fmt.Sprintf(regStr, n), str)
}

// 验证是否是标准正负小数(123. 不是小数)
func IsFloat(str string) bool { return isHave(runFuncName(), str) }

//	验证是否含有带不超过len个小数的小数
func IsFloat2Len(str string, l int) bool {
	regStr, ok := regMap[runFuncName()]
	if !ok{
		loger("reg func is not")
		return false
	}
	return isHaveStr(fmt.Sprintf(regStr, l), str)
}

//	验证是否是全汉字
func IsChineseAll(str string) bool {
	if str == "" {
		return false
	}
	for _, v := range str {
		if !unicode.Is(unicode.Han, v) {
			return false
		}
	}
	return true
}

//	验证是否含有汉字
func IsChinese(str string) bool {
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

//	验证是否含有number个汉字
func IsChineseN(str string, number int) bool {
	count := 0
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			count++
		}
	}
	if count == number {
		return true
	}
	return false
}

//	验证是否全是汉字数字
func IsChineseNumber(str string) bool {
	if str == "" {
		return false
	}
	for _, value := range str {
		if !isArrayStr(string(value), ChineseNumber) {
			return false
		}
	}
	return true
}

//	验证是否是中文钱大写
func IsChineseMoney(str string) bool {
	if str == "" {
		return false
	}
	if !isArrayStr(string(str[0]), ChineseMoney) {
		return false
	}
	for _, value := range str {
		if !isArrayStr(string(value), ChineseMoneyAll) {
			return false
		}
	}
	return true
}

//	验证是否是全英文
func IsEngAll(str string) bool { return isHave(runFuncName(), str) }

//	验证是否含不超过len个英文字符
func IsEngLen(str string, l int) bool {
	regStr, ok := regMap[runFuncName()]
	if !ok{
		loger("reg func is not")
		return false
	}
	return isHaveStr(fmt.Sprintf(regStr, l), str)
}

//	验证是否是英文和数字
func IsEngNumber(str string) bool { return isHave(runFuncName(), str) }

//	验证是否全大写
func IsAllCapital(str string) bool {
	for _, value := range str {
		if value > 91 || value < 64 {
			return false
		}
	}
	return true
}

//	验证是否有大写
func IsHaveCapital(str string) bool {
	for _, value := range str {
		if value < 91 && value > 64 {
			return true
		}
	}
	return false
}

//	验证是否全小写
func IsAllLower(str string) bool {
	for _, value := range str {
		if value > 123 || value < 96 {
			return false
		}
	}
	return true
}

//	验证是否有小写
func IsHaveLower(str string) bool {
	for _, value := range str {
		if value < 123 && value > 96 {
			return true
		}
	}
	return false
}

//	验证不低于n个数字
func IsLeastNumber(str string, n int) bool { return isHave(runFuncName(), str) }

//	验证不低于n个大写字母
func IsLeastCapital(str string, n int) bool { return isHave(runFuncName(), str) }

//	验证不低于n个小写字母
func IsLeastLower(str string, n int) bool { return isHave(runFuncName(), str) }

//	验证不低于n特殊字符
func IsLeastSpecial(str string, n int) bool { return isHave(runFuncName(), str) }

//	验证域名
func IsDomain(str string) bool { return isHave(runFuncName(), str) }

//	验证URL
func IsURL(str string) bool { return isHave(runFuncName(), str) }

//	验证手机号码
func IsPhone(str string) bool { return isHave(runFuncName(), str) }

//	验证电话号码("XXX-XXXXXXX"、"XXXX-XXXXXXXX"、"XXX-XXXXXXX"、"XXX-XXXXXXXX"、"XXXXXXX"和"XXXXXXXX)：
func IsLandline(str string) bool { return isHave(runFuncName(), str) }

//	IP地址：((?:(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d)\\.){3}(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d))
func IsIP(str string) bool { return isHave(runFuncName(), str) }


func isArrayStr(s string, slist []string) bool {
	for _, value := range slist {
		if s == value {
			return true
		}
	}
	return false
}

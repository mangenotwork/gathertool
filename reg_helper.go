/*
*	Description : 正则   TODO 测试
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

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

func RegFindAllTxt(regStr, rest string) (dataList []string) {
	reg := regexp.MustCompile(regStr)
	resList := reg.FindAllStringSubmatch(rest, -1)
	for _, v := range resList {
		if len(v) < 1 {
			continue
		}
		dataList = append(dataList, v[1])
	}
	return
}

var regMap = map[string]string{
	"RegHtmlA":           `(?is:<a.*?</a>)`,
	"RegHtmlTitle":       `(?is:<title.*?</title>)`,
	"RegHtmlKeyword":     `(?is:<meta name="keywords".*?>)`,
	"RegHtmlDescription": `(?is:<meta name="description".*?>)`,
	"RegHtmlTr":          `(?is:<tr.*?</tr>)`,
	"RegHtmlInput":       `(?is:<input.*?>)`,
	"RegHtmlTd":          `(?is:<td.*?</td>)`,
	"RegHtmlP":           `(?is:<p.*?</p>)`,
	"RegHtmlSpan":        `(?is:<span.*?</span>)`,
	"RegHtmlSrc":         `(?is:src=\".*?\")`,
	"RegHtmlHref":        `(?is:href=\".*?\")`,
	"RegHtmlH1":          `(?is:<h1.*?</h1>)`,
	"RegHtmlH2":          `(?is:<h2.*?</h2>)`,
	"RegHtmlH3":          `(?is:<h3.*?</h3>)`,
	"RegHtmlH4":          `(?is:<h4.*?</h4>)`,
	"RegHtmlH5":          `(?is:<h5.*?</h5>)`,
	"RegHtmlH6":          `(?is:<h6.*?</h6>)`,
	"RegHtmlTbody":       `(?is:<tbody.*?</tbody>)`,
	"RegHtmlVideo":       `(?is:<video.*?</video>)`,
	"RegHtmlCanvas":      `(?is:<canvas.*?</canvas>)`,
	"RegHtmlCode":        `(?is:<code.*?</code>)`,
	"RegHtmlImg":         `(?is:<img.*?/>)`,
	"RegHtmlUl":          `(?is:<ul.*?</ul>)`,
	"RegHtmlLi":          `(?is:<li.*?</li>)`,
	"RegHtmlMeta":        `(?is:<meta.*?>)`,
	"RegHtmlSelect":      `(?is:<select.*?</select>)`,
	"RegHtmlTable":       `(?is:<table.*?</table>)`,
	"RegHtmlButton":      `(?is:<button.*?</button>)`,
	"RegHtmlTableOlny":   `(?is:<table>.*?</table>)`,
	"RegHtmlDiv":         `(?is:<div.*?</div>)`,
	"RegHtmlOption":      `(?is:<option.*?</option>)`,

	// 取标签内容
	"RegHtmlATxt":           `(?is:<a.*?>(.*?)</a>)`,
	"RegHtmlTitleTxt":       `(?is:<title.*?>(.*?)</title>)`,
	"RegHtmlKeywordTxt":     `(?is:<meta name="keywords".*?content="(.*?)".*?>)`,
	"RegHtmlDescriptionTxt": `(?is:<meta name="description".*?content="(.*?)".*?.*?>)`,
	"RegHtmlTrTxt":          `(?is:<tr.*?>(.*?)</tr>)`,
	"RegHtmlTdTxt":          `(?is:<td.*?>(.*?)</td>)`,
	"RegHtmlPTxt":           `(?is:<p.*?>(.*?)</p>)`,
	"RegHtmlSpanTxt":        `(?is:<span.*?>(.*?)</span>)`,
	"RegHtmlSrcTxt":         `(?is:src=\"(.*?)\")`,
	"RegHtmlHrefTxt":        `(?is:href=\"(.*?)\")`,
	"RegHtmlHTxt1":          `(?is:<h1.*?>(.*?)</h1>)`,
	"RegHtmlHTxt2":          `(?is:<h2.*?>(.*?)</h2>)`,
	"RegHtmlHTxt3":          `(?is:<h3.*?>(.*?)</h3>)`,
	"RegHtmlHTxt4":          `(?is:<h4.*?>(.*?)</h4>)`,
	"RegHtmlHTxt5":          `(?is:<h5.*?>(.*?)</h5>)`,
	"RegHtmlHTxt6":          `(?is:<h6.*?>(.*?)</h6>)`,
	"RegHtmlCodeTxt":        `(?is:<code.*?>(.*?)</code>)`,
	"RegHtmlUlTxt":          `(?is:<ul.*?>(.*?)</ul>)`,
	"RegHtmlLiTxt":          `(?is:<li.*?>(.*?)</li>)`,
	"RegHtmlSelectTxt":      `(?is:<select.*?>(.*?)</select>)`,
	"RegHtmlTableTxt":       `(?is:<table.*?>(.*?)</table>)`,
	"RegHtmlButtonTxt":      `(?is:<button.*?>(.*?)</button>)`,
	"RegHtmlDivTxt":         `(?is:<div.*?>(.*?)</div>)`,
	"RegHtmlOptionTxt":      `(?is:<option.*?>(.*?)</option>)`,
	"RegValue":              `(?is:value=\"(.*?)\")`,

	// 删除
	"RegDelHtml":   `\<[\S\s]+?\>`,
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

	//常用正则
	"RegTime":         `(?i)\d{1,2}:\d{2} ?(?:[ap]\.?m\.?)?|\d[ap]\.?m\.?`,
	"RegLink":         `(?:(?:https?:\/\/)?(?:[a-z0-9.\-]+|www|[a-z0-9.\-])[.](?:[^\s()<>]+|\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\))+(?:\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:\'".,<>?]))`,
	"RegEmail":        `(?i)([A-Za-z0-9!#$%&'*+\/=?^_{|.}~-]+@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)`,
	"RegIPv4":         `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`,
	"RegIPv6":         `(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`,
	"RegIP":           `(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)|(?:(?:(?:[0-9A-Fa-f]{1,4}:){7}(?:[0-9A-Fa-f]{1,4}|:))|(?:(?:[0-9A-Fa-f]{1,4}:){6}(?::[0-9A-Fa-f]{1,4}|(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){5}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,2})|:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(?:(?:[0-9A-Fa-f]{1,4}:){4}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,3})|(?:(?::[0-9A-Fa-f]{1,4})?:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){3}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,4})|(?:(?::[0-9A-Fa-f]{1,4}){0,2}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){2}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,5})|(?:(?::[0-9A-Fa-f]{1,4}){0,3}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?:(?:[0-9A-Fa-f]{1,4}:){1}(?:(?:(?::[0-9A-Fa-f]{1,4}){1,6})|(?:(?::[0-9A-Fa-f]{1,4}){0,4}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(?::(?:(?:(?::[0-9A-Fa-f]{1,4}){1,7})|(?:(?::[0-9A-Fa-f]{1,4}){0,5}:(?:(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(?:\.(?:25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(?:%.+)?\s*`,
	"RegMD5Hex":       `[0-9a-fA-F]{32}`,
	"RegSHA1Hex":      `[0-9a-fA-F]{40}`,
	"RegSHA256Hex":    `[0-9a-fA-F]{64}`,
	"RegGUID":         `[0-9a-fA-F]{8}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{4}-?[a-fA-F0-9]{12}`,
	"RegMACAddress":   `(([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2}))`,
	"RegEmail2":       "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$",
	"RegUUID3":        "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$",
	"RegUUID4":        "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
	"RegUUID5":        "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$",
	"RegUUID":         "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$",
	"RegInt":          "^(?:[-+]?(?:0|[1-9][0-9]*))$",
	"RegFloat":        "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$",
	"RegRGBColor":     "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$",
	"RegFullWidth":    "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]",
	"RegHalfWidth":    "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]",
	"RegBase64":       "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$",
	"RegLatitude":     "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$",
	"RegLongitude":    "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$",
	"RegDNSName":      `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`,
	"RegFullURL":      `^(?:ftp|tcp|udp|wss?|https?):\/\/[\w\.\/#=?&]+$`,
	"RegURLSchema":    `((ftp|tcp|udp|wss?|https?):\/\/)`,
	"RegURLUsername":  `(\S+(:\S*)?@)`,
	"RegURLPath":      `((\/|\?|#)[^\s]*)`,
	"RegURLPort":      `(:(\d{1,5}))`,
	"RegURLIP":        `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`,
	"RegURLSubdomain": `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`,
	"RegWinPath":      `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`,
	"RegUnixPath":     `^(/[^/\x00]*)+/?$`,
}

// 常用正则
const ()

var (
	ChineseNumber   = []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "零"}
	ChineseMoney    = []string{"壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖"}
	ChineseMoneyAll = []string{"壹", "贰", "叁", "肆", "伍", "陆", "柒", "捌", "玖", "拾", "佰", "仟", "万", "亿", "元", "角", "分", "零", "整", "正", "貳", "陸", "億", "萬", "圓"}
)

// runFuncName 获取正在运行的函数名
func runFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	fName := f.Name()
	fList := strings.Split(fName, ".")
	return fList[len(fList)-1]
}

// regFind 执行正则提取
func regFind(funcName, txt string, property ...string) (dataList []string) {
	regStr, ok := regMap[funcName]
	if !ok {
		return
	}
	reg := regexp.MustCompile(regStr)
	resList := reg.FindAllStringSubmatch(txt, -1)
	for _, v := range resList {
		if len(v) < 1 {
			continue
		}
		if len(property) == 0 || strings.Count(v[0], strings.Join(property, " ")) > 0 {
			dataList = append(dataList, v[0])
		}
	}
	return
}

func RegHtmlA(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlTitle(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlKeyword(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlDescription(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlTr(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlInput(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlTd(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlP(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlSpan(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlSrc(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlHref(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlVideo(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlCanvas(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlCode(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlImg(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlUl(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlLi(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlMeta(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlSelect(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlTable(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlTableOlny(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlButton(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlH(str, typeH string, property ...string) []string {
	funcName := runFuncName()
	return regFind(funcName+typeH, str, property...)
}

func RegHtmlTbody(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlDiv(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

func RegHtmlOption(str string, property ...string) []string {
	return regFind(runFuncName(), str, property...)
}

// regFindTxt 执行正则提取 只取内容
func regFindTxt(funcName, txt string, property ...string) (dataList []string) {
	regStr, ok := regMap[funcName]
	if !ok {
		return
	}
	reg := regexp.MustCompile(regStr)
	resList := reg.FindAllStringSubmatch(txt, -1)
	for _, v := range resList {
		if len(v) < 1 {
			continue
		}
		if len(property) == 0 || strings.Count(v[0], strings.Join(property, " ")) > 0 {
			dataList = append(dataList, v[1])
		}
	}
	return
}

func RegHtmlATxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlTitleTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlKeywordTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlDescriptionTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlTrTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlInputTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlTdTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlPTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlSpanTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlSrcTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlHrefTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlCodeTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlUlTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlLiTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlSelectTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlTableTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlButtonTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlDivTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlOptionTxt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegValue(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

func RegHtmlHTxt(str, typeH string, property ...string) []string {
	funcName := runFuncName()
	return regFindTxt(funcName+typeH, str, property...)
}

// replace 删除正则匹配的字符
func replace(funcName, rest string) string {
	regStr, ok := regMap[funcName]
	if !ok {
		return ""
	}
	re, err := regexp.Compile(regStr)
	if err != nil {
		Error(err)
		return ""
	}
	return re.ReplaceAllString(rest, "")
}

// RegDelHtml 删除所有标签
func RegDelHtml(str string) string { return replace(runFuncName(), str) }

// RegDelNumber 删除所有数字
func RegDelNumber(str string) string { return replace(runFuncName(), str) }

func RegDelHtmlA(str string) string { return replace("RegHtmlA", str) }

func RegDelHtmlTitle(str string) string { return replace("RegHtmlTitle", str) }

func RegDelHtmlTr(str string) string { return replace("RegHtmlTr", str) }

func RegDelHtmlInput(str string, property ...string) string { return replace("RegHtmlInput", str) }

func RegDelHtmlTd(str string, property ...string) string { return replace("RegHtmlTd", str) }

func RegDelHtmlP(str string, property ...string) string { return replace("RegHtmlP", str) }

func RegDelHtmlSpan(str string, property ...string) string { return replace("RegHtmlSpan", str) }

func RegDelHtmlSrc(str string, property ...string) string { return replace("RegHtmlSrc", str) }

func RegDelHtmlHref(str string, property ...string) string { return replace("RegHtmlHref", str) }

func RegDelHtmlVideo(str string, property ...string) string { return replace("RegHtmlVideo", str) }

func RegDelHtmlCanvas(str string, property ...string) string { return replace("RegHtmlCanvas", str) }

func RegDelHtmlCode(str string, property ...string) string { return replace("RegHtmlCode", str) }

func RegDelHtmlImg(str string, property ...string) string { return replace("RegHtmlImg", str) }

func RegDelHtmlUl(str string, property ...string) string { return replace("RegHtmlUl", str) }

func RegDelHtmlLi(str string, property ...string) string { return replace("RegHtmlLi", str) }

func RegDelHtmlMeta(str string, property ...string) string { return replace("RegHtmlMeta", str) }

func RegDelHtmlSelect(str string, property ...string) string { return replace("RegHtmlSelect", str) }

func RegDelHtmlTable(str string, property ...string) string { return replace("RegHtmlTable", str) }

func RegDelHtmlButton(str string, property ...string) string { return replace("RegHtmlButton", str) }

func RegDelHtmlH(str, typeH string, property ...string) string { return replace("RegHtmlH"+typeH, str) }

func RegDelHtmlTbody(str string, property ...string) string { return replace("RegHtmlTbody", str) }

// isHaveStr 是否含有正则匹配的字符
func isHaveStr(regStr, rest string) bool {
	isHave, err := regexp.MatchString(regStr, rest)
	if err != nil {
		Error(err)
		return false
	}
	return isHave
}

// isHave 是否含有正则匹配的字符
func isHave(funcName, rest string) bool {
	regStr, ok := regMap[funcName]
	if !ok {
		Error("reg func is not")
		return false
	}
	return isHaveStr(regStr, rest)
}

// IsNumber 验证是否含有number
func IsNumber(str string) bool { return isHave(runFuncName(), str) }

// IsNumber2Len 验证是否含有连续长度不超过长度l的number
func IsNumber2Len(str string, l int) bool {
	regStr, ok := regMap[runFuncName()]
	if !ok {
		Error("reg func is not")
		return false
	}
	return isHaveStr(fmt.Sprintf(regStr, l), str)
}

// IsNumber2Heard 验证是否含有n开头的number
func IsNumber2Heard(str string, n int) bool {
	regStr, ok := regMap[runFuncName()]
	if !ok {
		Error("reg func is not")
		return false
	}
	return isHaveStr(fmt.Sprintf(regStr, n), str)
}

// IsFloat 验证是否是标准正负小数(123. 不是小数)
func IsFloat(str string) bool { return isHave(runFuncName(), str) }

// IsFloat2Len 验证是否含有带不超过len个小数的小数
func IsFloat2Len(str string, l int) bool {
	regStr, ok := regMap[runFuncName()]
	if !ok {
		Error("reg func is not")
		return false
	}
	return isHaveStr(fmt.Sprintf(regStr, l), str)
}

// IsChineseAll 验证是否是全汉字
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

// IsChinese 验证是否含有汉字
func IsChinese(str string) bool {
	for _, v := range str {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

// IsChineseN 验证是否含有number个汉字
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

// IsChineseNumber 验证是否全是汉字数字
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

// IsChineseMoney 验证是否是中文钱大写
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

// IsEngAll 验证是否是全英文
func IsEngAll(str string) bool { return isHave(runFuncName(), str) }

// IsEngLen 验证是否含不超过len个英文字符
func IsEngLen(str string, l int) bool {
	regStr, ok := regMap[runFuncName()]
	if !ok {
		Error("reg func is not")
		return false
	}
	return isHaveStr(fmt.Sprintf(regStr, l), str)
}

// IsEngNumber 验证是否是英文和数字
func IsEngNumber(str string) bool { return isHave(runFuncName(), str) }

// IsAllCapital 验证是否全大写
func IsAllCapital(str string) bool {
	for _, value := range str {
		if value > 91 || value < 64 {
			return false
		}
	}
	return true
}

// IsHaveCapital 验证是否有大写
func IsHaveCapital(str string) bool {
	for _, value := range str {
		if value < 91 && value > 64 {
			return true
		}
	}
	return false
}

// IsAllLower 验证是否全小写
func IsAllLower(str string) bool {
	for _, value := range str {
		if value > 123 || value < 96 {
			return false
		}
	}
	return true
}

// IsHaveLower 验证是否有小写
func IsHaveLower(str string) bool {
	for _, value := range str {
		if value < 123 && value > 96 {
			return true
		}
	}
	return false
}

// IsLeastNumber 验证不低于n个数字
func IsLeastNumber(str string, n int) bool { return isHave(runFuncName(), str) }

// IsLeastCapital 验证不低于n个大写字母
func IsLeastCapital(str string, n int) bool { return isHave(runFuncName(), str) }

// IsLeastLower 验证不低于n个小写字母
func IsLeastLower(str string, n int) bool { return isHave(runFuncName(), str) }

// IsLeastSpecial 验证不低于n特殊字符
func IsLeastSpecial(str string, n int) bool { return isHave(runFuncName(), str) }

// IsDomain 验证域名
func IsDomain(str string) bool { return isHave(runFuncName(), str) }

// IsURL 验证URL
func IsURL(str string) bool { return isHave(runFuncName(), str) }

// IsPhone 验证手机号码
func IsPhone(str string) bool { return isHave(runFuncName(), str) }

// IsLandline 验证电话号码("XXX-XXXXXXX"、"XXXX-XXXXXXXX"、"XXX-XXXXXXX"、"XXX-XXXXXXXX"、"XXXXXXX"和"XXXXXXXX)：
func IsLandline(str string) bool { return isHave(runFuncName(), str) }

// IsIP IP地址：((?:(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d)\\.){3}(?:25[0-5]|2[0-4]\\d|[01]?\\d?\\d))
func IsIP(str string) bool { return isHave(runFuncName(), str) }

// AccountRational 帐号合理性
func AccountRational(str string) bool { return isHave(runFuncName(), str) }

// IsXMLFile 是否三xml文件
func IsXMLFile(str string) bool { return isHave(runFuncName(), str) }

// IsUUID3 是否是uuid
func IsUUID3(str string) bool { return isHave(runFuncName(), str) }

// IsUUID4 是否是uuid
func IsUUID4(str string) bool { return isHave(runFuncName(), str) }

// IsUUID5 是否是uuid
func IsUUID5(str string) bool { return isHave(runFuncName(), str) }

// IsRGB 是否是 rgb
func IsRGB(str string) bool { return isHave(runFuncName(), str) }

// IsFullWidth 是否是全角字符
func IsFullWidth(str string) bool { return isHave(runFuncName(), str) }

// IsHalfWidth 是否是半角字符
func IsHalfWidth(str string) bool { return isHave(runFuncName(), str) }

// IsBase64 是否是base64
func IsBase64(str string) bool { return isHave(runFuncName(), str) }

// IsLatitude 是否是纬度
func IsLatitude(str string) bool { return isHave(runFuncName(), str) }

// IsLongitude 是否是经度
func IsLongitude(str string) bool { return isHave(runFuncName(), str) }

// IsDNSName 是否是dns 名称
func IsDNSName(str string) bool { return isHave(runFuncName(), str) }

// IsIPv4 是否是ipv4
func IsIPv4(str string) bool { return isHave(runFuncName(), str) }

// IsWindowsPath 是否是windows路径
func IsWindowsPath(str string) bool { return isHave(runFuncName(), str) }

// IsUnixPath 是否是unix路径
func IsUnixPath(str string) bool { return isHave(runFuncName(), str) }

func isArrayStr(s string, sList []string) bool {
	for _, value := range sList {
		if s == value {
			return true
		}
	}
	return false
}

// RegTime 提取时间
func RegTime(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegLink 提取链接
func RegLink(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegEmail 提取邮件
func RegEmail(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegIPv4 提取ipv4
func RegIPv4(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegIPv6 提取ipv6
func RegIPv6(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegIP 提取ip
func RegIP(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegMD5Hex 提取md5
func RegMD5Hex(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegSHA1Hex 提取sha1
func RegSHA1Hex(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegSHA256Hex 提取sha256
func RegSHA256Hex(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegGUID 提取guid
func RegGUID(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegMACAddress 提取MACAddress
func RegMACAddress(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegEmail2 提取邮件
func RegEmail2(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegUUID3 提取uuid
func RegUUID3(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegUUID4 提取uuid
func RegUUID4(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegUUID5 提取uuid
func RegUUID5(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegUUID 提取uuid
func RegUUID(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegInt 提取整形
func RegInt(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegFloat 提取浮点型
func RegFloat(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegRGBColor 提取RGB值
func RegRGBColor(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegFullWidth 提取全角字符
func RegFullWidth(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegHalfWidth 提取半角字符
func RegHalfWidth(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegBase64 提取base64字符串
func RegBase64(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegLatitude 提取纬度
func RegLatitude(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegLongitude 提取经度
func RegLongitude(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegDNSName 提取dns
func RegDNSName(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegFullURL 提取url
func RegFullURL(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegURLSchema  提取url schema
func RegURLSchema(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegURLUsername  提取url username
func RegURLUsername(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegURLPath  提取url path
func RegURLPath(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegURLPort  提取url port
func RegURLPort(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegURLIP  提取 url ip
func RegURLIP(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegURLSubdomain  提取 url sub domain
func RegURLSubdomain(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegWinPath  提取 windows路径
func RegWinPath(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

// RegUnixPath  提取 unix路径
func RegUnixPath(str string, property ...string) []string {
	return regFindTxt(runFuncName(), str, property...)
}

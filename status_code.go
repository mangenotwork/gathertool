/*
	Description : 请求状态码对应事件配置
	Author : ManGe
	Version : v0.1
	Date : 2021-04-23
*/

package gathertool

// StatusCodeMap 状态码处理映射
// success 该状态码对应执行成功函数
// fail    该状态码对应执行失败函数
// retry   该状态码对应需要重试前执行的函数
var StatusCodeMap map[int]string = map[int]string{
	200:"success",
	201:"success",
	202:"success",
	203:"success",
	204:"fail",
	300:"success",
	301:"success",
	302:"success",
	400:"fail",
	401:"retry",
	402:"retry",
	403:"retry",
	404:"fail",
	405:"retry",
	406:"retry",
	407:"retry",
	408:"retry",
	500:"fail",
	501:"fail",
	502:"retry",
	503:"retry",
	504:"retry",
}

// 将指定状态码设置为执行成功事件
func SetStatusCodeSuccessEvent(code int){
	StatusCodeMap[code] = "success"
}

// 将指定状态码设置为执行重试事件
func SetStatusCodeRetryEvent(code int){
	StatusCodeMap[code] = "retry"
}

// 将指定状态码设置为执行失败事件
func SetStatusCodeFailEvent(code int){
	StatusCodeMap[code] = "fail"
}

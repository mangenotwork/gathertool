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
// null    该状态码没有可执行，直接返回
var StatusCodeMap map[int]string = map[int]string{
	200:"success",
	201:"success",
	202:"success",
	203:"success",
	204:"null",
	300:"success",
	301:"success",
	302:"success",
	400:"null",
	401:"fail",
	402:"fail",
	403:"fail",
	404:"null",
	405:"fail",
	406:"fail",
	407:"fail",
	408:"fail",
	500:"null",
	501:"null",
	502:"fail",
	503:"fail",
	504:"fail",
}
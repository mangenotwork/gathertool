/*
	Description : 请求状态码对应事件配置
	Author : ManGe
	Version : v0.1
	Date : 2021-04-23
*/

package gathertool

import (
	"errors"
	"syscall"
)

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

// 设置为最大 socket open file
func SetMaxOpenFile() error {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err
	}

	rLimit.Cur = rLimit.Max
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}

// 设置指定 socket open file
func SetRLimit(number int) error {
	var rLimit syscall.Rlimit
	var n = uint64(number)

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err
	}
	if n > rLimit.Max {
		return errors.New("设置失败：操过最大 Rlimit")
	}
	rLimit.Cur = n
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}
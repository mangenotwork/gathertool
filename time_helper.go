/*
	Description : 时间相关的操作
	Author : ManGe
	Version : v0.1
	Date : 2021-04-26
*/

package gathertool

import (
	"strconv"
	"time"
)

func Timestamp() string {
	return strconv.FormatInt(time.Now().Unix(),10)
}

// BeginDayUnix 获取当天 0点
func BeginDayUnix() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t.Unix()
}

// EndDayUnix 获取当天 24点
func EndDayUnix() int64 {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t.Unix() + 86400
}

// TODO 获取多少分钟前的时间戳

// TODO 获取多少小时前的时间戳

// TODO 获取多少天前的时间戳

// 两个时间字符串的日期差
func Daydiff(beginDay string, endDay string) int {
	begin, _ := time.Parse("2006-01-02 15:04:05", beginDay+" 00:00:00")
	end, _ := time.Parse("2006-01-02 15:04:05", endDay+" 00:00:00")

	diff := end.Unix() - begin.Unix()
	return int(diff / (24 * 60 * 60))
}


// 间隔运行
// t: 间隔时间，  f: 运行的方法
func TickerRun(t time.Duration, runFirst bool, f func()){
	if runFirst {
		f()
	}
	tick := time.NewTicker(t)
	for range tick.C {
		f()
	}
}


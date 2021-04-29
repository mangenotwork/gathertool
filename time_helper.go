package gathertool

import "time"

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

// 获取多少分钟前的时间戳

// 获取多少小时前的时间戳

// 获取多少天前的时间戳



/*
*	Description : 时间相关的操作
*	Author 		: ManGe
*	Mail 		: 2912882908@qq.com
**/

package gathertool

import (
	"fmt"
	"strconv"
	"time"
)

const TimeTemplate = "2006-01-02 15:04:05"
const TimeMilliTemplate = "2006-01-02 15:04:05.000"
const TimeTemplateYMD = "2006-01-02"
const OneDay = 24
const OneHours = 60
const OneMinutes = 60
const ZeroHMS = " 00:00:00"
const TimeFormatList = "01-02 15:04"

// Timestamp 获取时间戳
func Timestamp() int64 {
	return time.Now().Unix()
}

func TimestampStr() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// BeginDayUnix 获取当天 0点
func BeginDayUnix() int64 {
	timeStr := time.Now().Format(TimeTemplateYMD)
	t, _ := time.ParseInLocation(TimeTemplateYMD, timeStr, time.Local)
	return t.Unix()
}

// EndDayUnix 获取当天 24点
func EndDayUnix() int64 {
	timeStr := time.Now().Format(TimeTemplateYMD)
	t, _ := time.ParseInLocation(TimeTemplateYMD, timeStr, time.Local)
	return t.Unix() + OneDay*OneHours*OneMinutes
}

// MinuteAgo 获取多少分钟前的时间戳
func MinuteAgo(i int) int64 {
	return time.Now().Unix() - int64(i*OneMinutes)
}

// HourAgo 获取多少小时前的时间戳
func HourAgo(i int) int64 {
	return time.Now().Unix() - int64(i*OneHours*OneMinutes)
}

// DayAgo 获取多少天前的时间戳
func DayAgo(i int) int64 {
	return time.Now().Unix() - int64(i*OneHours*OneMinutes*OneDay)
}

// DayDiff 两个时间字符串的日期差, 返回的是天
func DayDiff(beginDay string, endDay string) int {
	begin, _ := time.Parse(TimeTemplate, beginDay+ZeroHMS)
	end, _ := time.Parse(TimeTemplate, endDay+ZeroHMS)

	diff := end.Unix() - begin.Unix()
	return int(diff / (OneDay * OneHours * OneMinutes))
}

func DayDiffAtUnix(s, e int64) int {
	diff := e - s
	return int(diff / (OneDay * OneHours * OneMinutes))
}

// TickerRun 间隔运行
// t: 间隔时间， runFirst: 间隔前或者后执行  f: 运行的方法
func TickerRun(t time.Duration, runFirst bool, f func()) {
	if runFirst {
		f()
	}
	tick := time.NewTicker(t)
	for range tick.C {
		f()
	}
}

func Timestamp2Date(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format(TimeTemplate)
}

// NowToEnd 计算当前时间到这天结束还有多久
func NowToEnd() (int64, error) {
	now := time.Now()
	nextDay := now.Add(24 * time.Hour)
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return -1, err
	}
	endTime := fmt.Sprintf("%d-%d-%d 00:00:00", nextDay.Year(), nextDay.Month(), nextDay.Day())
	times, err := time.ParseInLocation(TimeTemplate, endTime, loc)
	if err != nil {
		return -1, err
	}
	distanceSec := times.Unix() - now.Unix()
	return distanceSec, nil
}

// IsToday 判断是否是今天   "2006-01-02 15:04:05"
// timestamp 需要判断的时间
func IsToday(timestamp int64) string {
	nowTime := time.Now()
	t := time.Unix(timestamp, 0)
	if t.Day() == nowTime.Day() && t.Month() == nowTime.Month() && t.Year() == nowTime.Year() {
		return "今天 " + t.Format("15:04:05")
	}
	return t.Format(TimeTemplate)
}

// IsTodayList 列表页的时间显示  "01-02 15:04"
func IsTodayList(timestamp int64) string {
	nowTime := time.Now()
	t := time.Unix(timestamp, 0)
	if t.Day() == nowTime.Day() && t.Month() == nowTime.Month() && t.Year() == nowTime.Year() {
		return "今天 " + t.Format("15:04")
	}
	return t.Format(TimeFormatList)
}

func Timestamp2Week(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	tm.Weekday()
	switch tm.Weekday() {
	case time.Sunday:
		return "周天"
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	}
	return "周一"
}

func Timestamp2WeekXinQi(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	tm.Weekday()
	switch tm.Weekday() {
	case time.Sunday:
		return "星期天"
	case time.Monday:
		return "星期一"
	case time.Tuesday:
		return "星期二"
	case time.Wednesday:
		return "星期三"
	case time.Thursday:
		return "星期四"
	case time.Friday:
		return "星期五"
	case time.Saturday:
		return "星期六"
	}
	return "星期一"
}

// LatestDate 返回最近好多天的日期
func LatestDate(date int) []string {
	outList := make([]string, 0)
	now := time.Now()
	outList = append(outList, now.Format("2006-01-02"))
	nowInt := now.Unix()
	for i := 0; i < date; i++ {
		nowInt -= 86400
		outList = append(outList, time.Unix(nowInt, 0).Format("2006-01-02"))
	}
	return outList
}

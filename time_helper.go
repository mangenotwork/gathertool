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

// Timestamp2Date timestamp -> "2006-01-02 15:04:05"
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

// Timestamp2Week 获取timestamp是周几
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

// Timestamp2WeekXinQi 获取timestamp是星期几
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

// GetCurrentMonthRange 获取当前月的时间范围（月初1号0点 ~ 月末最后一刻）
// 返回：
//
//	firstDay：当月1号 00:00:00
//	lastDay：当月最后一天 23:59:59.999999999
func GetCurrentMonthRange() (firstDay, lastDay time.Time) {
	now := time.Now() // 当前时间（当地时区）

	// 年、月不变，日设为1，时/分/秒/纳秒设为0
	firstDay = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 先获取下个月1号0点，再减1天得到当月最后一天0点，最后减1纳秒得到23:59:59.999999999
	nextMonthFirstDay := firstDay.AddDate(0, 1, 0)             // 下个月1号0点
	lastDayStart := nextMonthFirstDay.AddDate(0, 0, -1)        // 当月最后一天0点
	lastDay = lastDayStart.Add(24*time.Hour - time.Nanosecond) // 最后一天23:59:59.999999999

	return firstDay, lastDay
}

// GetCurrentWeekRange 获取当前周的时间范围（周一凌晨 ~ 周日23:59:59）
// 返回：
//
//	monday ：本周一 00:00:00
//	sundayEnd：本周日 23:59:59.999999999
func GetCurrentWeekRange() (monday, sundayEnd time.Time) {
	now := time.Now() // 当前时间（当地时区）

	// 计算当前是星期几（转换为：周一=1，周二=2，...，周日=7）
	currentWeekday := int(now.Weekday())
	if currentWeekday == 0 { // 若为周日（0），转为7
		currentWeekday = 7
	}

	// 计算距离本周一的天数差
	daysToMonday := currentWeekday - 1

	// 计算本周一的凌晨（00:00:00）
	monday = now.AddDate(0, 0, -daysToMonday).Truncate(24 * time.Hour)

	// 计算本周日的23:59:59.999999999（当天最后一刻）
	sundayDate := monday.AddDate(0, 0, 6) // 周一加6天 = 周日（日期）
	sundayEnd = time.Date(
		sundayDate.Year(),
		sundayDate.Month(),
		sundayDate.Day(),
		23, 59, 59, 999999999, // 时:分:秒:纳秒（当天最后一刻）
		now.Location(), // 保持与当前时间相同的时区
	)

	return monday, sundayEnd
}

// GetTodayRange 按东八区（Asia/Shanghai）获取当天自然日范围（00:00:00 至 23:59:59.999999999）
// 返回：
//
//	start：当天00:00:00（东八区）
//	end：当天23:59:59.999999999（东八区）
func GetTodayRange() (start, end time.Time) {
	// 指定时区为东八区（Asia/Shanghai）
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)

	// 计算东八当天凌晨（00:00:00）
	start = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		loc, // 明确使用东八
	)

	// 计算东八当天结束（23:59:59.999999999）
	end = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		23, 59, 59, 999999999,
		loc, // 明确使用东八
	)

	return start, end
}

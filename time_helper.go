/*
	Description : 时间相关的操作
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"fmt"
	"strconv"
	"time"
)

const TimeTemplate = "2006-01-02 15:04:05"

// Timestamp 获取时间戳
func Timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
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

// MinuteAgo 获取多少分钟前的时间戳
func MinuteAgo(i int) int64 {
	return time.Now().Unix() - int64(i*60)
}

// HourAgo 获取多少小时前的时间戳
func HourAgo(i int) int64 {
	return time.Now().Unix() - int64(i*3600)
}

// DayAgo 获取多少天前的时间戳
func DayAgo(i int) int64 {
	return time.Now().Unix() - int64(i*3600*24)
}

// DayDiff 两个时间字符串的日期差
func DayDiff(beginDay string, endDay string) int {
	begin, _ := time.Parse("2006-01-02 15:04:05", beginDay+" 00:00:00")
	end, _ := time.Parse("2006-01-02 15:04:05", endDay+" 00:00:00")

	diff := end.Unix() - begin.Unix()
	return int(diff / (24 * 60 * 60))
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

const TimeMilliTemplate = "2006-01-02 15:04:05.000"

var (
	chineseMonth     = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
	chineseNumber    = []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"}
	gan              = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	zhi              = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
	chineseTen       = []string{"初", "十", "廿", "卅"}
	animals          = []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
	year, month, day int
	leap             bool

	timeFormat         = "2006-01-02 15:04:05"
	timeFormatList     = "01-02 15:04"
	timeFormatYYYYMMDD = "20060102"
	lunarInfo          = []int{0x04bd8, 0x04ae0, 0x0a570,
		0x054d5, 0x0d260, 0x0d950, 0x16554, 0x056a0, 0x09ad0, 0x055d2,
		0x04ae0, 0x0a5b6, 0x0a4d0, 0x0d250, 0x1d255, 0x0b540, 0x0d6a0,
		0x0ada2, 0x095b0, 0x14977, 0x04970, 0x0a4b0, 0x0b4b5, 0x06a50,
		0x06d40, 0x1ab54, 0x02b60, 0x09570, 0x052f2, 0x04970, 0x06566,
		0x0d4a0, 0x0ea50, 0x06e95, 0x05ad0, 0x02b60, 0x186e3, 0x092e0,
		0x1c8d7, 0x0c950, 0x0d4a0, 0x1d8a6, 0x0b550, 0x056a0, 0x1a5b4,
		0x025d0, 0x092d0, 0x0d2b2, 0x0a950, 0x0b557, 0x06ca0, 0x0b550,
		0x15355, 0x04da0, 0x0a5d0, 0x14573, 0x052d0, 0x0a9a8, 0x0e950,
		0x06aa0, 0x0aea6, 0x0ab50, 0x04b60, 0x0aae4, 0x0a570, 0x05260,
		0x0f263, 0x0d950, 0x05b57, 0x056a0, 0x096d0, 0x04dd5, 0x04ad0,
		0x0a4d0, 0x0d4d4, 0x0d250, 0x0d558, 0x0b540, 0x0b5a0, 0x195a6,
		0x095b0, 0x049b0, 0x0a974, 0x0a4b0, 0x0b27a, 0x06a50, 0x06d40,
		0x0af46, 0x0ab60, 0x09570, 0x04af5, 0x04970, 0x064b0, 0x074a3,
		0x0ea50, 0x06b58, 0x055c0, 0x0ab60, 0x096d5, 0x092e0, 0x0c960,
		0x0d954, 0x0d4a0, 0x0da50, 0x07552, 0x056a0, 0x0abb7, 0x025d0,
		0x092d0, 0x0cab5, 0x0a950, 0x0b4a0, 0x0baa4, 0x0ad50, 0x055d9,
		0x04ba0, 0x0a5b0, 0x15176, 0x052b0, 0x0a930, 0x07954, 0x06aa0,
		0x0ad50, 0x05b52, 0x04b60, 0x0a6e6, 0x0a4e0, 0x0d260, 0x0ea65,
		0x0d530, 0x05aa0, 0x076a3, 0x096d0, 0x04bd7, 0x04ad0, 0x0a4d0,
		0x1d0b6, 0x0d250, 0x0d520, 0x0dd45, 0x0b5a0, 0x056d0, 0x055b2,
		0x049b0, 0x0a577, 0x0a4b0, 0x0aa50, 0x1b255, 0x06d20, 0x0ada0}
)

// GetChineseMonthDay 获取农历
func GetChineseMonthDay(date string) (rMonth, rDay int64) {
	var monCyl, leapMonth int = 0, 0
	t1, _ := time.Parse(timeFormat, "1900-01-31 00:00:00")
	t2, err := time.Parse(timeFormatYYYYMMDD, date)
	if err != nil {
		return 0, 0
	}
	offset := int((t2.UnixNano() - t1.UnixNano()) / 1000000 / 86400000)
	monCyl = 14
	var iYear, daysOfYear int = 0, 0

	for iYear = 1900; iYear < 2050 && offset > 0; iYear++ {
		daysOfYear = yearDays(iYear)
		offset -= daysOfYear
		monCyl += 12
	}

	if offset < 0 {
		offset += daysOfYear
		iYear--
		monCyl -= 12
	}
	year = iYear
	leapMonth = leapMonthMethod(iYear)
	leap = false

	var iMonth, daysOfMonth int = 0, 0

	for iMonth = 1; iMonth < 13 && offset > 0; iMonth++ {
		if leapMonth > 0 && iMonth == (leapMonth+1) && !leap {
			iMonth--
			leap = true
			daysOfMonth = leapDays(year)
		} else {
			daysOfMonth = monthDays(year, iMonth)
		}
		offset -= daysOfMonth
		if leap && iMonth == (leapMonth+1) {
			leap = false
		}
		if !leap {
			monCyl++
		}
	}

	if offset == 0 && leapMonth > 0 && iMonth == leapMonth+1 {
		if leap {
			leap = false
		} else {
			leap = true
			iMonth--
			monCyl--
		}
	}

	if offset < 0 {
		offset += daysOfMonth
		iMonth--
		monCyl--
	}
	month = iMonth
	day = offset + 1

	// doubleMonth := ""
	// if leap {
	// 	doubleMonth = "闰"
	// }
	//return cyclical() + animalsYear() + "年" + doubleMonth + chineseMonth[month-1] + "月" + getChinaDayString(day)
	rMonth, err = strconv.ParseInt(chineseMonth[month-1], 10, 64)
	if err != nil {
		Error(err)
	}
	rDay, err = strconv.ParseInt(getChinaDayString(day), 10, 64)
	if err != nil {
		Error(err)
	}
	return rMonth, rDay
}

func yearDays(y int) int {
	var i, sum int = 348, 348
	for i = 0x8000; i > 0x8; i >>= 1 {
		if (lunarInfo[y-1900] & i) != 0 {
			sum++
		}
	}
	return (sum + leapDays(y))
}

func getChinaDayString(day int) string {
	// n := day
	if day > 30 {
		return ""
	}
	return fmt.Sprintf("%d", day)
	// if n%10 == 0 {
	// 	n = 9
	// } else {
	// 	n = day%10 - 1
	// }
	// if day > 30 {
	// 	return ""
	// } else if day == 10 {
	// 	return "初十"
	// } else {
	// 	return chineseTen[day/10] + chineseNumber[n]
	// }
}

func leapMonthMethod(y int) int {
	return (int)(lunarInfo[y-1900] & 0xf)
}

func monthDays(y, m int) int {
	if (lunarInfo[y-1900] & (0x10000 >> uint(m))) == 0 {
		return 29
	}
	return 30
}

func leapDays(y int) int {
	if leapMonthMethod(y) != 0 {
		if (lunarInfo[y-1900] & 0x10000) != 0 {
			return 30
		}
		return 29
	}
	return 0
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
	times, err := time.ParseInLocation("2006-1-2 15:04:05", endTime, loc)
	if err != nil {
		return -1, err
	}
	distanceSec := times.Unix() - now.Unix()
	return distanceSec, nil
}

// Leaps 闰年的天数
var Leaps = []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// Pyears 平年天数
var Pyears = []int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

// IsLeap 是否是闰年
func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// IsToday 判断是否是今天   "2006-01-02 15:04:05"
// timestamp 需要判断的时间
func IsToday(timestamp int64) string {
	nowTime := time.Now()
	t := time.Unix(timestamp, 0)
	if t.Day() == nowTime.Day() && t.Month() == nowTime.Month() && t.Year() == nowTime.Year() {
		return "今天 " + t.Format("15:04:05")
	}
	return t.Format(timeFormat)
}

// IsTodayList 列表页的时间显示  "01-02 15:04"
func IsTodayList(timestamp int64) string {
	nowTime := time.Now()
	t := time.Unix(timestamp, 0)
	if t.Day() == nowTime.Day() && t.Month() == nowTime.Month() && t.Year() == nowTime.Year() {
		return "今天 " + t.Format("15:04")
	}
	return t.Format(timeFormatList)
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

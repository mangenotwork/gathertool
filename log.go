/*
	Description : 内部日志打印
	Author : ManGe
	Mail : 2912882908@qq.com
	Github : https://github.com/mangenotwork/gathertool
*/

package gathertool

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// LogClose 是否关闭日志
var LogClose bool = true
var std = newStd()

// CloseLog 关闭日志
func CloseLog() {
	LogClose = false
}

type logger struct {
	outFile       bool
	outFileWriter *os.File
}

func newStd() *logger {
	return &logger{}
}

// SetLogFile 设置日志输出到的指定文件
func SetLogFile(name string) {
	std.outFile = true
	std.outFileWriter, _ = os.OpenFile(name+time.Now().Format("-20060102")+".log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// Level 日志等级
type Level int

var LevelMap = map[Level]string{
	1: "[Info]  ",
	2: "[Debug] ",
	3: "[Warn]  ",
	4: "[Error] ",
	5: "[HTTP/s]",
}

// Log 日志
func (l *logger) Log(level Level, args string, times int) {
	var buffer bytes.Buffer
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05.000 "))
	buffer.WriteString(LevelMap[level])
	_, file, line, _ := runtime.Caller(times)
	fileList := strings.Split(file, "/")
	// 最多显示两级路径
	if len(fileList) > 3 {
		fileList = fileList[len(fileList)-3 : len(fileList)]
	}
	buffer.WriteString(strings.Join(fileList, "/"))
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(line))
	buffer.WriteString(" \t| ")
	buffer.WriteString(args)
	buffer.WriteString("\n")
	out := buffer.Bytes()
	if LogClose {
		_, _ = buffer.WriteTo(os.Stdout)
	}

	// 输出到文件或远程日志服务
	if l.outFile {
		_, _ = l.outFileWriter.Write(out)
	}
}

// Info 日志-信息
func Info(args ...interface{}) {
	std.Log(1, fmt.Sprint(args...), 2)
}

// Infof 日志-信息
func Infof(format string, args ...interface{}) {
	std.Log(1, fmt.Sprintf(format, args...), 2)
}

// InfoTimes 日志-信息, 指定日志代码位置的定位调用层级
func InfoTimes(times int, args ...interface{}) {
	std.Log(1, fmt.Sprint(args...), times)
}

// InfofTimes 日志-信息, 指定日志代码位置的定位调用层级
func InfofTimes(format string, times int, args ...interface{}) {
	std.Log(1, fmt.Sprintf(format, args...), times)
}

// Debug 日志-调试
func Debug(args ...interface{}) {
	std.Log(2, fmt.Sprint(args...), 2)
}

// Debugf 日志-调试
func Debugf(format string, args ...interface{}) {
	std.Log(2, fmt.Sprintf(format, args...), 2)
}

// DebugTimes 日志-调试, 指定日志代码位置的定位调用层级
func DebugTimes(times int, args ...interface{}) {
	std.Log(1, fmt.Sprint(args...), times)
}

// DebugfTimes 日志-调试, 指定日志代码位置的定位调用层级
func DebugfTimes(format string, times int, args ...interface{}) {
	std.Log(1, fmt.Sprintf(format, args...), times)
}

// Warn 日志-警告
func Warn(args ...interface{}) {
	std.Log(3, fmt.Sprint(args...), 2)
}

// Warnf 日志-警告
func Warnf(format string, args ...interface{}) {
	std.Log(3, fmt.Sprintf(format, args...), 2)
}

// WarnTimes 日志-警告, 指定日志代码位置的定位调用层级
func WarnTimes(times int, args ...interface{}) {
	std.Log(1, fmt.Sprint(args...), times)
}

// WarnfTimes 日志-警告, 指定日志代码位置的定位调用层级
func WarnfTimes(format string, times int, args ...interface{}) {
	std.Log(1, fmt.Sprintf(format, args...), times)
}

// Error 日志-错误
func Error(args ...interface{}) {
	std.Log(4, fmt.Sprint(args...), 2)
}

// Errorf 日志-错误
func Errorf(format string, args ...interface{}) {
	std.Log(4, fmt.Sprintf(format, args...), 2)
}

// ErrorTimes 日志-错误, 指定日志代码位置的定位调用层级
func ErrorTimes(times int, args ...interface{}) {
	std.Log(4, fmt.Sprint(args...), times)
}

// ErrorfTimes 日志-错误, 指定日志代码位置的定位调用层级
func ErrorfTimes(format string, times int, args ...interface{}) {
	std.Log(4, fmt.Sprintf(format, args...), times)
}

// HTTPTimes 日志-信息, 指定日志代码位置的定位调用层级
func HTTPTimes(times int, args ...interface{}) {
	std.Log(5, fmt.Sprint(args...), times)
}

// Bar 终端显示的进度条
type Bar struct {
	percent int64  //百分比
	cur     int64  //当前进度位置
	total   int64  //总进度
	rate    string //进度条
	graph   string //显示符号
}

func (bar *Bar) NewOption(start, total int64) {
	bar.cur = start
	bar.total = total
	if bar.graph == "" {
		bar.graph = "█"
	}
	bar.percent = bar.getPercent()
	for i := 0; i < int(bar.percent); i += 2 {
		bar.rate += bar.graph //初始化进度条位置
	}
}

func (bar *Bar) getPercent() int64 {
	return int64(float32(bar.cur) / float32(bar.total) * 100)
}

func (bar *Bar) NewOptionWithGraph(start, total int64, graph string) {
	bar.graph = graph
	bar.NewOption(start, total)
}

func (bar *Bar) Play(cur int64) {
	bar.cur = cur
	last := bar.percent
	bar.percent = bar.getPercent()
	if bar.percent != last && bar.percent%2 == 0 {
		bar.rate += bar.graph
	}
	fmt.Printf("\r[%-50s]%3d%% %8d/%d", bar.rate, bar.percent, bar.cur, bar.total)
}

func (bar *Bar) Finish() {
	fmt.Println()
}

/*
	Description : 内部日志打印
	Author : ManGe
	Version : v0.2
	Date : 2021-12-27
*/

package gathertool

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"
)

var std = newStd()

type logger struct {
	outFile bool
	outFileWriter *os.File
}

func newStd() *logger {
	return &logger{}
}

func SetLogFile(name string) {
	std.outFile = true
	std.outFileWriter, _ = os.OpenFile( name+time.Now().Format("-20060102")+".log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

type Level int

var LevelMap = map[Level]string {
	1 : "Info  ",
	2 : "Debug ",
	3 : "Warn  ",
	4 : "Error ",
}

func (l *logger) Log(level Level, args string, times int) {
	var buffer bytes.Buffer
	buffer.WriteString(time.Now().Format("2006-01-02 15:04:05 |"))
	buffer.WriteString(LevelMap[level])
	_, file, line, _ := runtime.Caller(times)
	buffer.WriteString("|")
	buffer.WriteString(file)
	buffer.WriteString(":")
	buffer.WriteString(strconv.Itoa(line))
	buffer.WriteString(" : ")
	buffer.WriteString(args)
	buffer.WriteString("\n")
	outString := buffer.Bytes()
	_,_ = buffer.WriteTo(os.Stdout)
	if l.outFile {
		_,_ = l.outFileWriter.Write(outString)
	}
}

func Info(args ...interface{}) {
	std.Log(1, fmt.Sprint(args...), 2)
}

func Infof(format string, args ...interface{}) {
	std.Log(1, fmt.Sprintf(format, args...), 2)
}

func InfoTimes(times int, args ...interface{}) {
	std.Log(1, fmt.Sprint(args...), times)
}

func Debug(args ...interface{}) {
	std.Log(2, fmt.Sprint(args...), 2)
}

func Debugf(format string, args ...interface{}) {
	std.Log(2, fmt.Sprintf(format, args...), 2)
}

func Warn(args ...interface{}) {
	std.Log(3, fmt.Sprint(args...), 2)
}

func Warnf(format string, args ...interface{}) {
	std.Log(3, fmt.Sprintf(format, args...), 2)
}

func Error(args ...interface{}) {
	std.Log(4, fmt.Sprint(args...), 2)
}

func Errorf(format string, args ...interface{}) {
	std.Log(4, fmt.Sprintf(format, args...), 2)
}

// cli 进度条
type Bar struct {
	percent int64 //百分比
	cur   int64 //当前进度位置
	total  int64 //总进度
	rate  string //进度条
	graph  string //显示符号
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

func (bar *Bar) Finish(){
	fmt.Println()
}

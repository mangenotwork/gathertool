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
	std.Log(3, fmt.Sprintf(format, args...), 2)
}


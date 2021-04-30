/*
	Description : 内部日志打印
	Author : ManGe
	Version : v0.1
	Date : 2021-04-27
*/

package gathertool

import (
	"fmt"
	"log"
	"runtime"
)

func loger(v ...interface{}){
	_, file, line, _ := runtime.Caller(2)
	//fun := runtime.FuncForPC(pc)
	//funName := fun.Name()
	log.Println(fmt.Sprintf("%v:%v", file, line), v)
}
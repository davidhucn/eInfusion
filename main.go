package main

import (
	//	. "eInfusion/comm"
	eh "eInfusion/httpOperate"
	"eInfusion/logs"
	et "eInfusion/tcpOperate"
	"runtime"
)

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func tt(){

	for i:=0;i<2;i++{
		println(i)
	}
}

func main() {
	tt()
	go eh.StartHttpServer()
	et.StartTcpServer(7778)
}

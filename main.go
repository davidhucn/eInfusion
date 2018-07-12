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

func main() {
	go eh.StartHttpServer(7779)
	et.StartTcpServer(7778)
	
}

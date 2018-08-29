package main

import (
	"eInfusion/logs"
	eh "eInfusion/thttp"
	eq "eInfusion/tqueue"
	et "eInfusion/ttcp"
	"runtime"
)

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {
	go eq.StartSendQueueListener()
	go et.StartTCPServer(7778)
	eh.StartHTTPServer(7779)
}

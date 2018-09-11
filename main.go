package main

import (
	eh "eInfusion/thttp"
	logs "eInfusion/tlogs"
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
	go et.StartTCPServer(7778)
	go eh.StartSendQueueListener()
	eh.StartHTTPServer(7779)
}

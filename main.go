package main

import (
	eh "eInfusion/http"
	et "eInfusion/tcp"
	logs "eInfusion/tlogs"
	"runtime"
	// eh "eInfusion/thttp"
	// et "eInfusion/ttcp"
)

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main() {
	// go et.StartTCPServer(7778)
	// go eh.StartSendQueueListener()
	// eh.StartHTTPServer(7779)
	wc := eh.NewWebClients()
	go eh.StartHTTPServer(wc, 7778)
	dt := et.NewDevices()
	et.StartTCPService(dt, 7779)

}

package main

import (
	eh "eInfusion/http"
	et "eInfusion/tcp"
	logs "eInfusion/tlogs"
	"runtime"
)

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	wc := eh.NewWebClients()
	go eh.StartHTTPServer(wc, 7779)
	tcpSer := et.NewTCPServer(300, 45)
	et.RunTCPService(tcpSer, 7778)

}

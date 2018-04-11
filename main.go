package main

import (
	//	eh "eInfusion/httpOperate"
	"eInfusion/logs"
	et "eInfusion/tcpOperate"
)

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()
}

func main() {
	//	go eh.StartHttpServer()
	//	et.StartTcpServer()
	et.TryTcpServer()
}

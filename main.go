package main

import (
	//	eh "eInfusion/httpOperate"
	"eInfusion/logs"
	//	et "eInfusion/tcpOperate"
	. "eInfusion/comm"
	//	ec "strconv"
	//	"bytes"
)

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()
}

func main() {
	//	go eh.StartHttpServer()
	//	et.StartTcpServer()
	test()

}

func test() {
	//	str2 := "0xao"
	//	comm.Msg("str2:", comm.ConvertBasStrToInt(10, str2))
	//	comm.Msg("typeof", comm.GetVarType(comm.ConvertBasStrToInt(10, str2)))
	//	data2 := []byte(str2)

	var s []byte

	v := "c000000"
	d := string(v[0]) + string(v[1])
	Msg(d)
	s = append(s, ConvertBasStrToUint(10, d))
	Msg(s)
	//	Msg(string(v[0]))
	//	Msg(GetPartOfStringToBytes(v, 1, 3))

}

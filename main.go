package main

import (
	"net"
	"tcp/comm"
	"tcp/logs"

	tp "tcp/protocol"
)

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()

}

func main() {
	netListen, err := net.Listen("tcp", ":7778")
	defer netListen.Close()
	comm.ScrPrint("TCP数据接收平台开始运行...")
	if err != nil {
		logs.LogMain.Critical("监听TCP出错", err)
		return
	}
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		logs.LogMain.Info("客户端：" + conn.RemoteAddr().String() + "连接成功")
		//	comm.ScrPrint("tcp客户端连接成功，ip地址:", conn.RemoteAddr().String())
		go receiveData(conn)
		//	time.Sleep(time.Second * 2)
	}
}

func receiveData(conn net.Conn) {
	for {
		//	指定接收数据包头的帧长
		recDataHeader := make([]byte, tp.GetDataHeaderLength())
		_, err := conn.Read(recDataHeader)
		logs.LogMain.Error("接收Tcp包头失败", err)
		// 数据包长度记录变量
		var intPckContentLength int
		// 判断包头是否正确，如果不正确直接退回
		if !tp.DecodeHeader(recDataHeader, &intPckContentLength) {
			return
		}
		// 如果包头接收完整
		recDataContent := make([]byte, intPckContentLength)
		_, err = conn.Read(recDataContent)
		logs.LogMain.Error("接收包数据出错", err)
		// 处理数据包内容
		tp.DecodeData(recDataContent)

	}
}

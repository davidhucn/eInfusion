package tcpOperate

import (
	"eInfusion/comm"
	"eInfusion/logs"
	ep "eInfusion/protocol"
	"net"
	//	"sync"
	//	"reflect"
)

const c_Msg_ServerStart = "Transfusion平台运行中 …… "
const c_TcpServer_Port = "7778"

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()
}

func StartTcpServer() {
	netListen, err := net.Listen("tcp", ":"+c_TcpServer_Port)
	defer netListen.Close()
	//	系统开始运行时log记录时间
	logs.LogMain.Info(c_Msg_ServerStart + "（" + comm.GetCurrentDate() + "）")
	if err != nil {
		logs.LogMain.Critical("监听TCP出错", err)
		panic(err)
	}
	comm.ShowScreen("TCP Port:" + c_TcpServer_Port)
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		comm.ShowScreen("----------------------------------------------------------")
		logs.LogMain.Info("客户端：" + conn.RemoteAddr().String() + " 连接!")
		go receiveData(conn)
		//	time.Sleep(time.Second * 2)
	}
}

func receiveData(conn net.Conn) {
	for {
		//	指定接收数据包头的帧长
		recDataHeader := make([]byte, ep.GetDataHeaderLength())
		_, err := conn.Read(recDataHeader)
		if err != nil {
			//	logs.LogMain.Error("接收Tcp包头失败", err)
			comm.ShowScreen(conn.RemoteAddr(), " 客户端连接丢失!")
			return
		}
		// 数据包长度记录变量
		var intPckContentLength int
		// 判断包头是否正确，如果正确，获取长度
		if !ep.DecodeHeader(recDataHeader, &intPckContentLength) {
			return
		}
		// 如果包头接收完整
		recDataContent := make([]byte, intPckContentLength)
		_, err = conn.Read(recDataContent)
		if err != nil {
			logs.LogMain.Error("接收包数据出错", err)
		}
		// 处理数据包内容
		ep.DecodeRcvData(recDataContent)
	}
}

func sendData(conn net.Conn, bData []byte) {

	conn.Write([]byte("hello")) // don't care about return value
	conn.Close()
}

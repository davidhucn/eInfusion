package tcpBackup

import (
	"eInfusion/comm"
	"eInfusion/logs"
	ep "eInfusion/protocol"
	"net"
	//	"sync"
)

func StartTcpServer_backup() {
	//	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	//	checkError(err)
	netListen, err := net.Listen("tcp", ":"+c_TcpServer_Port)
	defer netListen.Close()
	//	系统开始运行时log记录时间
	logs.LogMain.Info(c_Msg_Info_ServerStart + "（" + comm.GetCurrentDate() + "）")
	if err != nil {
		logs.LogMain.Critical("监听TCP出错", err)
		panic(err)
	}
	comm.SepLi(60)
	comm.Msg("TCP Port:" + c_TcpServer_Port)
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		comm.SepLi(60)
		logs.LogMain.Info("客户端：" + conn.RemoteAddr().String() + " 连接!")
		go receiveData_backup(conn)
		//	time.Sleep(time.Second * 2)
	}
}

func receiveData_backup(conn net.Conn) {
	for {
		//	指定接收数据包头的帧长
		recDataHeader := make([]byte, ep.G_TsCmd.HeaderLength)
		_, err := conn.Read(recDataHeader)
		if err != nil {
			comm.SepLi(60)
			comm.Msg(conn.RemoteAddr(), " 客户端连接丢失!")
			return
		}
		// 数据包数据内容长度记录变量
		var intPckContentLength int
		// 判断包头是否正确，如果正确，获取长度
		if !ep.DecodeHeader(recDataHeader, &intPckContentLength) {
			comm.Msg("调试信息：数据包头不正确")
			continue
		}
		// 如果包头接收
		recDataContent := make([]byte, intPckContentLength)
		_, err = conn.Read(recDataContent)
		if !comm.CkErr("接收报文出错", err) {
			// 处理报文数据内容
			ep.DecodeRcvData(recDataContent, conn.RemoteAddr().(*net.TCPAddr).IP.String())
		}
	}
}

func sendData_backup(conn net.Conn, packetData []byte) {
	_, err := conn.Write(packetData) // don't care about return value
	defer conn.Close()
	if err != nil {
		comm.Msg(c_Msg_Err_SendData, err)
		logs.LogMain.Critical(c_Msg_Err_SendData, err)
		return
	}
}

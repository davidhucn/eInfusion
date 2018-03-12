package tcpOperate

import (
	"eInfusion/comm"
	//	"eInfusion/db"
	"eInfusion/logs"
	ep "eInfusion/protocol"
	"net"
)

func init() {
	// 初始化日志
	logs.LogDisable()
	logs.LogConfigLoad()

}

func StartTcpServer() {
	//	db.TestDb()
	netListen, err := net.Listen("tcp", ":7778")
	defer netListen.Close()
	comm.ShowScreen("["+comm.GetCurrentTime()+"]", "Transfusion数据接收平台开始运行...")
	//	系统开始运行时log记录时间
	logs.LogMain.Info("[" + comm.GetCurrentTime() + "]" + "Transfusion数据接收平台开始运行... ")
	if err != nil {
		logs.LogMain.Critical("监听TCP出错", err)
		panic(err)
	}
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		logs.LogMain.Info("客户端：" + conn.RemoteAddr().String() + "连接成功!")
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
			logs.LogMain.Error("接收Tcp包头失败", err)
			return
		}
		// 数据包长度记录变量
		var intPckContentLength int
		// 判断包头是否正确，如果不正确直接退回
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
		//		tp.DecodeToOrderData(recDataContent)

	}
}

func sendData(conn net.Conn) {

	conn.Write([]byte("hello")) // don't care about return value
	conn.Close()
}

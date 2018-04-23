package tcpOperate

import (
	"eInfusion/comm"
	"eInfusion/logs"
	ep "eInfusion/protocol"
	"net"
	//	"strconv"
	//	"strings"
)

func init() {
	//	初始化连接对象集
	G_tConns = make(map[string]TcpConn)
}

func StartTcpServer() {
	//	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":8989")
	//	checkError(err)
	netListen, err := net.Listen("tcp", ":"+c_TcpServer_Port)
	defer netListen.Close()
	//	系统开始运行时log记录时间
	logs.LogMain.Info(c_Msg_ServerStart + "（" + comm.GetCurrentDate() + "）")
	if err != nil {
		logs.LogMain.Critical("监听TCP出错", err)
		panic(err)
	}
	comm.SepLi(60)
	comm.Msg("TCP Port:" + c_TcpServer_Port)
	//	最大连接数不能超过规定数
	if len(G_tConns) <= c_MaxConnectionAmount {
		for {
			conn, err := netListen.Accept()
			if err != nil {
				continue
			}
			////////////////临时Tconn连接对象////////////////////////////////
			var c TcpConn
			c.ID = conn.RemoteAddr().String()
			c.IsAlive = true
			c.Conn = conn
			c.IPAddr = conn.RemoteAddr().(*net.TCPAddr).IP.String()
			G_tConns[c.ID] = c
			///////////////////////////////////////////////////////////////
			comm.SepLi(60)
			logs.LogMain.Info("客户端：" + c.ID + " 连接!")
			PreSendOrders()
			go receiveData(c)
			//	time.Sleep(time.Second * 2)
			///////////////////////////////////////////////////////////////
		}
	} else {
		//超出连接数则不再接收连接
		logs.LogMain.Warn(c_Msg_OutOfMaximumConnection)
	}

}

func receiveData(c TcpConn) {
	for {
		//	指定接收数据包头的帧长
		recDataHeader := make([]byte, ep.G_TsCmd.HeaderLength)
		_, err := c.Conn.Read(recDataHeader)
		if err != nil {
			comm.SepLi(60)
			c.Mutex.Lock()
			delete(G_tConns, c.ID)
			c.Mutex.Unlock()
			return
		}
		// 数据包数据内容长度记录变量
		var intPckContentLength int
		// 判断包头是否正确，如果正确，获取长度
		if !ep.DecodeHeader(recDataHeader, &intPckContentLength) {
			comm.Msg("调试信息：数据包头不正确")
			continue
		}
		// 如果包头接收正确
		recDataContent := make([]byte, intPckContentLength)
		_, err = c.Conn.Read(recDataContent)
		if !comm.CkErr("接收报文出错", err) {
			// 处理报文数据内容
			ep.DecodeRcvData(recDataContent, c.IPAddr)
		}
	}
}

// 整合信息发送至指定客户端
func PreSendOrders() {
	//遍历所有连接结点，发送命令
	for connsID, _ := range G_tConns {
		var orders []byte
		var RcvID []byte
		s := "A0000000"
		RcvID = comm.ConvertPerTwoOxCharOfStrToBytes(s)
		orders = ep.CmdGetRcvStatus(RcvID)
		SendData(G_tConns[connsID].Conn, orders)
	}
}

func SendData(conn net.Conn, packetData []byte) {
	_, err := conn.Write(packetData) // don't care about return value
	//	defer conn.Close()
	if err != nil {
		comm.Msg(c_Msg_SendDataErr, err)
		logs.LogMain.Error(c_Msg_SendDataErr, err)
		return
	}
}

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
	l, err := net.Listen("tcp", ":"+c_TcpServer_Port)
	defer l.Close()
	//	系统开始运行时log记录时间
	logs.LogMain.Info(c_Msg_Info_ServerStart + "（" + comm.GetCurrentDate() + "）")
	if err != nil {
		logs.LogMain.Critical("监听TCP出错", err)
		panic(err)
	}
	comm.SepLi(60)
	comm.Msg("TCP Port:" + c_TcpServer_Port)
	newConn(l)
}

//生成新tcp连接
func newConn(lser net.Listener) {
	//	最大连接数不能超过规定数
	if len(G_tConns) <= c_MaxConnectionAmount {
		for {
			var c TcpConn
			var err error
			c.Conn, err = lser.Accept()
			if err != nil {
				continue
				logs.LogMain.Error(c_Msg_Err_AcceptConnection)
			}
			c.Flag = c.Conn.RemoteAddr().String()
			c.IPAddr = c.Conn.RemoteAddr().(*net.TCPAddr).IP.String()
			G_tConns[c.ID] = c
			comm.SepLi(60)
			logs.LogMain.Info("客户端：" + c.ID + " 连接!")
			go receiveData(c)
			//	time.Sleep(time.Second * 2)
		}
	} else {
		//超出连接数则不再接收连接
		logs.LogMain.Warn(c_Msg_Warn_OutOfMaxConnection)
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
			delete(G_tConns, c.Flag)
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
		r := make([]byte, intPckContentLength)
		_, err = c.Conn.Read(r)
		if !comm.CkErr("接收报文出错", err) {
			// 处理报文数据内容
			ep.DecodeRcvData(r, c.IPAddr)
		}
	}
}

// 整合信息发送至指定客户端
func SendOrders(connID string, packetData []byte) {

	//如果连接存在
	if _, ok := G_tConns[connID]; ok {
		writeToConn(G_tConns[connID].Conn, packetData)
	} else { /*如果连接断开，则保留2小时*/

	}
}

func writeToConn(conn net.Conn, packetData []byte) {
	_, err := conn.Write(packetData) // don't care about return value
	//	defer conn.Close()
	if err != nil {
		comm.Msg(c_Msg_Err_SendData, err)
		logs.LogMain.Error(c_Msg_Err_SendData, err)
		return
	}
}

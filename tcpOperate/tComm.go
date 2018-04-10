package tcpOperate

import "net"

const c_TcpServer_Port = "7778"

const (
	c_Msg_SendDataErr = "发送数据错误！"
	c_Msg_ServerStart = "Transfusion平台运行中 …… "
)

type TConn struct {
	IPAddr   string
	Port     int
	Conn     *net.TCPConn
	IsConned bool
	ConnID   int
}

// tcp全局连接slice
var g_Conns []TConn

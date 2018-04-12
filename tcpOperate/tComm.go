package tcpOperate

import "net"

const c_TcpServer_Port = "7778"

const (
	c_Msg_SendDataErr            = "发送数据错误！"
	c_Msg_ServerStart            = "Transfusion平台运行中 …… "
	c_Msg_OutOfMaximumConnection = "超出设定连接数!"
)

//tcp最大连接数
const c_MaxConnectionAmount = 30

type TConn struct {
	IPAddr  string
	Conn    net.Conn
	IsAlive bool
	ID      string
}

// tcp全局连接map
var g_Conns map[string]TConn

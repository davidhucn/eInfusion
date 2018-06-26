package tcpOperate

import (
	"net"
	//	"sync"
)

const c_TcpServer_Port = "7778"

const (
	c_Msg_Err_SendData            = "发送数据错误！"
	c_Msg_Info_ServerStart        = "Transfusion平台运行中 …… "
	c_Msg_Err_AcceptConnection    = "错误，Tcp接收连接错误！"
	c_Msg_Warn_OutOfMaxConnection = "警告,超出设定连接数!"
)

//tcp最大连接数
const c_MaxConnectionAmount = 2

//定义锁
var (
	connMkMutex  sync.Mutex
	connDelMutex sync.Mutex
)

type Clienter struct {
	IPAddr string
	Conn   net.Conn
	Flag   string /*临时标记，用于索引,暂定时间戳*/
	ID     string /*设备编号唯一标识*/
	//	sync.Mutex
}

//指令消息结构
type OrdersQueue struct {
	OrderContent []byte
	CreateTime   string
	TargetID     string
}

//全局tcp连接对象
//var G_tConns map[string]TcpConn

//消息
var G_cOrders chan OrdersQueue

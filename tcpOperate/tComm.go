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
const c_MaxConnectionAmount = 30

//定义锁
var (
	connMkMutex  sync.Mutex
	connDelMutex sync.Mutex
)

<<<<<<< HEAD
//全局tcp连接对象
var clisConnMap map[string]*net.TCPConn
=======
type EndPointer struct {
	IPAddr string
	Conn   net.Conn
	Flag   string /*临时标记，用于索引,暂定时间戳*/
	ID     string /*设备编号唯一标识*/
}
>>>>>>> 7d5589a2bc1ca06c570d39bf6c057c1a6ce15959

//指令消息结构
type OrdersQueue struct {
	OrderContent []byte
	CreateTime   string
	TargetID     string
}
<<<<<<< HEAD
=======

//全局连接对象集
var G_epMap map[string]EndPointer

//消息
var G_cOrders chan OrdersQueue
>>>>>>> 7d5589a2bc1ca06c570d39bf6c057c1a6ce15959

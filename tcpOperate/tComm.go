package tcpOperate

import (
	"net"
	"sync"
)

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

//全局tcp连接对象
var clisConnMap map[string]*net.TCPConn

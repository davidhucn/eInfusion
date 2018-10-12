package tcp

import (
	cm "eInfusion/comm"
	"net"
	"sync"
)

type tcpMsg struct {
	SendError              string
	ReceiveError           string
	StartServiceMsg        string
	HeaderDataError        string
	CanNotFindConnection   string
	OutOfMaxConnAmount     string
	SendSuccess            string
	SourceError            string
	SendFailureForLongTime string
}

func init() {
	TCPMsg.SendError = "错误，发送数据错误！"
	TCPMsg.ReceiveError = "错误，Tcp接收数据错误！"
	TCPMsg.OutOfMaxConnAmount = "提示,超出设定连接数！"
	TCPMsg.HeaderDataError = "错误，数据包头错误！"
	TCPMsg.CanNotFindConnection = "错误，未找到TCP连接！"
	TCPMsg.SourceError = "错误，TCP服务资源错误！"
	TCPMsg.SendSuccess = "提示,发送指令和数据成功！"
	TCPMsg.SendFailureForLongTime = "错误，由于设备长时间断线或者故障，发送指令和设备失败！"
	TCPMsg.StartServiceMsg = "提示，Transfusion平台运行中 ……"
}

// MaxTCPConnectLimit :TCP最大连接数
// const MaxTCPConnectLimit int = 3

// TCPMsg :消息对象
var TCPMsg tcpMsg

// Devices :设备对象，用于衔接TCP操用
type Devices struct {
	// 连接集合
	Connections   map[string]*net.TCPConn
	Orders        chan *cm.Cmd
	MaxTCPConnect int
	sync.Mutex
}

// NewDevices :创建设备对象
func NewDevices(rMaxTCPConn int) *Devices {
	return &Devices{
		Connections:   make(map[string]*net.TCPConn),
		Orders:        make(chan *cm.Cmd, 1024),
		MaxTCPConnect: rMaxTCPConn,
	}
}

package tcp

import (
	cm "eInfusion/comm"
	"net"
	"sync"
	"time"
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

// WaitOrder ：TCP发送数据对象（指令 + 数据）
type WaitOrder struct {
	CreateTime time.Time
	SendData   *cm.Cmd
}

func newWaitOrder(rTime time.Time, rOrder *cm.Cmd) *WaitOrder {
	return &WaitOrder{
		CreateTime: rTime,
		SendData:   rOrder,
	}
}

// TClient : TCP客户端对象
type TClient struct {
	Connection *net.TCPConn
	SendData   chan *cm.Cmd
}

// newTClient :创建新的 TClient
func newTClient(rConn *net.TCPConn) *TClient {
	return &TClient{
		Connection: rConn,
		SendData:   make(chan *cm.Cmd, 1024),
	}
}

// TServer :TCP服务对象
type TServer struct {
	// 连接集合
	Connections         map[string]*TClient
	WaitOrders          []WaitOrder
	ExpireTimeByMinutes int
	MaxTCPConnect       int
	sync.Mutex
}

// NewTCPServer :创建设备对象
func NewTCPServer(rMaxTCPConnectAmount int, rExpireTimeByMinutes int) *TServer {
	return &TServer{
		Connections:         make(map[string]*TClient),
		WaitOrders:          make([]WaitOrder, 0),
		ExpireTimeByMinutes: rExpireTimeByMinutes,
		MaxTCPConnect:       rMaxTCPConnectAmount,
	}
}

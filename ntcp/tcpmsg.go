package ntcp

type tcpMsg struct {
	SendError              string
	ReceiveError           string
	ReceiveDataOutOfRange  string
	StartServiceMsg        string
	HeaderDataError        string
	CanNotFindConnection   string
	OutOfMaxConnAmount     string
	SendSuccess            string
	SourceError            string
	SendFailureForLongTime string
}

// TCPMsg :消息对象
var TCPMsg tcpMsg

func init() {
	TCPMsg.SendError = "错误，发送数据错误！"
	TCPMsg.ReceiveError = "错误，TCP接收数据错误！"
	TCPMsg.OutOfMaxConnAmount = "提示,超出设定连接数！"
	TCPMsg.HeaderDataError = "错误，数据包头错误！"
	TCPMsg.CanNotFindConnection = "错误，未找到TCP连接！"
	TCPMsg.SourceError = "错误，TCP服务资源错误！"
	TCPMsg.SendSuccess = "提示,发送指令和数据成功！"
	TCPMsg.SendFailureForLongTime = "错误，由于设备长时间断线或者故障，发送指令和设备失败！"
	TCPMsg.ReceiveDataOutOfRange = "接收的数据包超出限定长度"
	TCPMsg.StartServiceMsg = "提示，TCP平台开始运行 ……"
}

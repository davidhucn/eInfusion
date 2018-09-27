package tcp

type tcpMsg struct {
	SendError          string
	ReceiveError       string
	StartServiceMsg    string
	OutOfMaxConnAmount string
}

// TCPMsg :消息对象
var TCPMsg tcpMsg

func init() {
	TCPMsg.SendError = "错误，发送数据错误！"
	TCPMsg.ReceiveError = "错误，Tcp接收连接错误！"
	TCPMsg.OutOfMaxConnAmount = "提示,超出设定连接数!"
	TCPMsg.StartServiceMsg = "Transfusion平台运行中 ……"
}

//tcp最大连接数
const maxOfConnAmount = 3

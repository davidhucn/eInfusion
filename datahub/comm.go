//Package datahub :此包主要用于并行的数据交互(tcp,websocket)
package datahub

import (
	cm "eInfusion/comm"
	"sync"
)

// TCPOrderQueue :TCP指令队列
var TCPOrderQueue chan *cm.Cmd

// WebMsgQueue :回写到web的消息发送队列
var WebMsgQueue chan *cm.Cmd

// RequestOrder : 客户端请求对象
type RequestOrder struct {
	TargetID  string
	CmdType   uint8
	Args      string
	RequestID string //可用于wsID
}

// NewReqestOrder ：新建请求指令对象
func NewReqestOrder(rTargetID string, rCmdType uint8, rArgs string, rRequestID string) *RequestOrder {
	return &RequestOrder{
		TargetID:  rTargetID,
		CmdType:   rCmdType,
		Args:      rArgs,
		RequestID: rRequestID,
	}
}

type dataHubMsg struct {
	GetServerDataErr string
	CmdInvaildErr    string
	CmdRepeatNotice  string
}

// DataHubMsg :Datahub消息提示对象
var DataHubMsg dataHubMsg

type reqOrdersUnion struct {
	RequestOrders map[string]*RequestOrder //map[targetID + randstring]*RequestOrder
	sync.Mutex
}

// ReqOrdersUnion :客户端请求指令记录池,记录TCP指令ID，通过ID匹配，便于回写到前端web
var ReqOrdersUnion reqOrdersUnion

func init() {
	TCPOrderQueue = make(chan *cm.Cmd, 1024)
	WebMsgQueue = make(chan *cm.Cmd, 1024)
	ReqOrdersUnion.RequestOrders = make(map[string]*RequestOrder)

	DataHubMsg.GetServerDataErr = "错误，获取服务器数据出错！"
	DataHubMsg.CmdInvaildErr = "错误，非法或不可识别指令！"
	DataHubMsg.CmdRepeatNotice = "提示，申请指令重复"
}

//Package datahub :此包主要用于并行的数据交互(tcp,websocket)
package datahub

import (
	cm "eInfusion/comm"
	"strings"
	"sync"
)

// TCPOrderQueue :TCP指令队列
var tcpOrderQueue chan *cm.Cmd

// WebMsgQueue :回写到web的消息发送队列
var WebMsgQueue chan *cm.Cmd

// RequestOrder : 客户端请求对象
type RequestOrder struct {
	TargetID  string
	CmdType   uint8
	Args      string
	RequestID string //可用于wsID、MQTTID
}

// NewReqestOrder ：新建请求指令对象
func NewReqestOrder(rRequestID string, rTargetID string, rCmdType uint8, rArgs string) *RequestOrder {
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

// ReqOrdersUnion :客户端请求指令记录池,记录TCP指令ID，
// 通过ID匹配，便于回写到前端web，避免记录web重复操作
var ReqOrdersUnion reqOrdersUnion

func init() {
	tcpOrderQueue = make(chan *cm.Cmd, 1024)
	WebMsgQueue = make(chan *cm.Cmd, 1024)
	ReqOrdersUnion.RequestOrders = make(map[string]*RequestOrder)

	DataHubMsg.GetServerDataErr = "错误，获取服务器数据出错！"
	DataHubMsg.CmdInvaildErr = "错误，非法或不可识别指令！"
	DataHubMsg.CmdRepeatNotice = "提示，申请指令重复"
}

// NewWSOrderID :生成新的websocket消息编号
func NewWSOrderID(rWSConnectionID string) string {
	return rWSConnectionID + "#" + cm.GetRandString(8)
}

// DecodeToWSConnID :解析生成websocket连接序号
func DecodeToWSConnID(rWSOrderID string) string {
	return strings.Split(rWSOrderID, "#")[0]
}

// NewTCPOrderID :生成TCP包约定指令序号
func NewTCPOrderID(rStrCnt string, rTCPConnectionID string) string {
	return strings.Split(rStrCnt, "@")[0] + "@" + rTCPConnectionID
}

// DecodeToTCPConnID :解析指令ID为TCP连接序号
func DecodeToTCPConnID(rStrCnt string) string {
	return strings.Split(rStrCnt, "@")[1]
}

// NewReqOrdersUnionID : 生成新的RequestOrdersUnion ID
func NewReqOrdersUnionID(rTargetID string) string {
	return rTargetID + "~" + cm.GetRandString(8)
}

// DecodeToReqOrderID : 解析为RequestOrder的ID
func DecodeToReqOrderID(rStrCnt string) string {
	return strings.Split(rStrCnt, "~")[0]
}

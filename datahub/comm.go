//Package datahub :此包主要用于并行的数据交互(tcp,websocket)
package datahub

import (
	cm "eInfusion/comm"
	wk "eInfusion/dbwork"
	ep "eInfusion/protocol"
	logs "eInfusion/tlogs"
	"strings"
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
	RequestID string
}

// ReqOrder :客户端请求指令对象，便于WS回写定位
// var ReqOrder requestOrder

// ReqOrdersUnion :指令ID记录池,记录TCP指令ID，通过ID匹配，便于回写到前端web
// var ReqOrdersUnion map[string][]*requestOrder
var ReqOrdersUnion sync.Map //map[tcpConnectionID][]*requestOrder

func init() {
	TCPOrderQueue = make(chan *cm.Cmd, 1024)
	WebMsgQueue = make(chan *cm.Cmd, 1024)
}

// AddToTCPQueue ：通过TCP协议发送指令至设备
func addToTCPSendQueue(rCmd *cm.Cmd, rTargetID string) {
	// 保存cmdID,便于回写到web前端时可以对应
	tcpID := strings.Split(rCmd.CmdID, "@")[1]
	wsID := strings.Split(rCmd.CmdID, "@")[0]
	// 指令池内存放wsID + 设备号
	val := wsID + "~" + rTargetID
	ReqOrdersUnion.Store(tcpID, val)
	TCPOrderQueue <- rCmd
}

// SendMsgToWeb :回写到web前端
func SendMsgToWeb(rCmd *cm.Cmd) {
	WebMsgQueue <- rCmd
}

// NewTCPOrderID :生成TCP包约定指令序号
func NewTCPOrderID(rWSOrderID string, rTCPConnectionID string) string {
	return rWSOrderID + "@" + rTCPConnectionID
}

// DecodeToTCPConnID :解析指令ID为TCP连接序号
func DecodeToTCPConnID(rStrCnt string) string {
	return strings.Split(rStrCnt, "@")[1]
}

// SendOrderToDeviceByTCP :添加到TCP指令发送队列
func SendOrderToDeviceByTCP(rWSOrderID string, rTargetID string, rCmdType uint8, rArgs string) bool {
	// 如果不是接收器, TODO:稍候改成直接判断是否为检测器
	if !wk.IsReceiver(rTargetID) {
		rcvID := cm.ConvertStrToBytesByPerTwoChar(wk.GetRcvID(rTargetID))
		ipAddr := wk.GetRcvIP(wk.GetRcvID(rTargetID))
		if rcvID == "" || ipAddr == "" {
			// 根据DetID获取RcvID和IP地址失败时，返回错误
			return false
		}
		addDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.AddDetect)
		delDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.DelDetect)
		if rCmdType == addDet || rCmdType == delDet {
			// TODO:存入指令集合ReqOrdersUnion
			reqOrder := &RequestOrder{TargetID: rTargetID, CmdType: rCmdType, Args: rArgs, RequestID: rWSOrderID}
			reqOrderID := NewTCPOrderID(cm.GetRandString(8), ipAddr)
			ReqOrdersUnion.Store(reqOrderID, reqOrder)
			// TCP指令标识:wsOrderID + 随机字符 + IP地址组成
			tcpOrderID := NewTCPOrderID(rWSOrderID, ipAddr)
			detID := cm.ConvertStrToBytesByPerTwoChar(rTargetID)
			od := cm.NewOrder(tcpOrderID, ep.CmdOperateDetect(rCmdType, rcvID, 1, detID))
			addToTCPSendQueue(od, rTargetID)
		} else {
			logs.LogMain.Error("错误：没有目标设备编码或者无法获取相关设备编码！")
		}
	}
	return true
}

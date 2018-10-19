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

// OrdersIDUnion :指令ID记录池,记录TCP指令ID，通过ID匹配，便于回写到前端web
var OrdersIDUnion sync.Map //map[tcpConnectionID]wsID

func init() {
	TCPOrderQueue = make(chan *cm.Cmd, 1024)
	WebMsgQueue = make(chan *cm.Cmd, 1024)
}

// AddToTCPQueue ：通过TCP协议发送指令至设备
func addToTCPSendQueue(rCmd *cm.Cmd, rDeviceID string) {
	// 保存cmdID,便于回写到web前端时可以对应
	tcpID := strings.Split(rCmd.CmdID, "@")[1]
	wsID := strings.Split(rCmd.CmdID, "@")[0]
	// 指令池内存放wsID + 设备号
	val := wsID + "~" + rDeviceID
	OrdersIDUnion.Store(tcpID, val)
	TCPOrderQueue <- rCmd
}

// SendMsgToWeb :回写到web前端
func SendMsgToWeb(rCmd *cm.Cmd) {
	WebMsgQueue <- rCmd
}

// NewTCPOrderID :生成TCP包约定指令序号
func NewTCPOrderID(rOrderID string, rTCPConnectionID string) string {
	return rOrderID + "@" + rTCPConnectionID
}

// DecodeToTCPConnID :解析指令ID为TCP连接序号
func DecodeToTCPConnID(rStrCnt string) string {
	return strings.Split(rStrCnt, "@")[1]
}

// SendOrderToDeviceByTCP :添加到TCP指令发送队列
func SendOrderToDeviceByTCP(rOrderID string, rDeviceID string, rCmdType uint8, rArgs string) {
	// 如果不是接收器
	if !wk.IsReceiver(rDeviceID) {
		addDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.AddDetect)
		delDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.DelDetect)
		if rCmdType == addDet || rCmdType == delDet {
			// 根据DetID获取RcvID,如果rcvID不为空
			if wk.GetRcvID(rDeviceID) != "" {
				rcvID := cm.ConvertStrToBytesByPerTwoChar(wk.GetRcvID(rDeviceID))
				detID := cm.ConvertStrToBytesByPerTwoChar(rDeviceID)
				// 获取rcv相关的IP
				ipAddr := wk.GetRcvIP(wk.GetRcvID(rDeviceID))
				// TCP指令标识:wsID + 随机字符 + IP地址组成
				tcpOrderID := NewTCPOrderID(rOrderID, ipAddr)
				od := cm.NewOrder(tcpOrderID, ep.CmdOperateDetect(rCmdType, rcvID, 1, detID))
				addToTCPSendQueue(od, rDeviceID)
			} else {
				logs.LogMain.Error("错误：没有目标设备编码或者无法获取相关设备编码！")
			}
		}
	} else {
		// TODO:如果是接收器
	}

}

// func AddToWSDataQueue(rID, rOrder []byte) {

// }

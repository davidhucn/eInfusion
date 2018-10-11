//Package datahub :此包主要用于并行的数据交互(tcp,websocket)
package datahub

import (
	cm "eInfusion/comm"
	wk "eInfusion/dbwork"
	ep "eInfusion/protocol"
	logs "eInfusion/tlogs"
)

// DeviceTCPOrderQueue :设备命令对象队列
var DeviceTCPOrderQueue chan *cm.Cmd

// WebMsgQueue :回写到web的消息发送队列
var WebMsgQueue chan *cm.Cmd

func init() {
	DeviceTCPOrderQueue = make(chan *cm.Cmd, 1024)
	WebMsgQueue = make(chan *cm.Cmd, 1024)
}

// AddToTCPQueue ：通过TCP协议发送指令至设备
func addToTCPSendQueue(rCmd *cm.Cmd) {
	DeviceTCPOrderQueue <- rCmd
}

// SendMsgToWeb :回写到web前端
func SendMsgToWeb(rCmd *cm.Cmd) {
	WebMsgQueue <- rCmd
}

// SendOrderToDeviceByTCP :添加到TCP数据发送队列
func SendOrderToDeviceByTCP(rOrderID string, rDeviceID string, rCmdType uint8, rArgs string) {
	addDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.AddDetect)
	delDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.DelDetect)
	if rCmdType == addDet || rCmdType == delDet {
		// 根据DetID获取RcvID,如果rcvID不为空
		if wk.GetRcvID(rDeviceID) != "" {
			rcvID := cm.ConvertStrToBytesByPerTwoChar(wk.GetRcvID(rDeviceID))
			detID := cm.ConvertStrToBytesByPerTwoChar(rDeviceID)
			// 获取rcv相关的IP
			ipAddr := wk.GetRcvIP(wk.GetRcvID(rDeviceID))
			// 重组指令标识:由时间戳+IP地址组成
			orderID := rOrderID + "@" + ipAddr
			od := cm.NewOrder(orderID, ep.CmdOperateDetect(rCmdType, rcvID, 1, detID))
			addToTCPSendQueue(od)
		} else {
			logs.LogMain.Error("错误：没有目标设备编码或者无法获取相关设备编码！")
		}
	}
}

// func AddToWSDataQueue(rID, rOrder []byte) {

// }

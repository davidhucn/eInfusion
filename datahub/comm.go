//Package datahub :此包主要用于并行的数据交互(tcp,websocket)
package datahub

import (
	cm "eInfusion/comm"
	wk "eInfusion/dbwork"
	ep "eInfusion/protocol"
)

// func AddToWSDataQueue(rID, rOrder []byte) {

// }

func AddToTCPQueue(rOrderID string, rDetID string, rCmdType uint8, rArgs string) {
	addDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.AddDetect)
	delDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.DelDetect)
	if rCmdType == addDet || rCmdType == delDet {
		// 根据DetID获取RcvID,如果rcvID不为空
		if wk.GetRcvID(rDetID) != "" {
			rcvID := cm.ConvertStrToBytesByPerTwoChar(wk.GetRcvID(rDetID))
			detID := cm.ConvertStrToBytesByPerTwoChar(rDetID)
			// 获取rcv相关的IP
			ipAddr := wk.GetRcvIP(wk.GetRcvID(rDetID))
			// 重组指令标识:由时间戳+IP地址组成
			orderID := rOrderID + "@" + ipAddr
			addSendQueueMap(orderID, ep.CmdOperateDetect(rCmdType, rcvID, 1, detID))
		} else {
			logs.LogMain.Error("错误：没有目标设备编码或者无法获取相关设备编码！")
		}
	}
}

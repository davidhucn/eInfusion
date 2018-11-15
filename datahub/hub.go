package datahub

import (
	cm "eInfusion/comm"
	ec "eInfusion/dbcomm"
	tsc "eInfusion/trsfscomm"
)

// AddToTCPOrderQueue ：通过TCP协议发送指令至设备
func AddToTCPOrderQueue(rCmd *cm.Cmd) {
	tcpOrderQueue <- rCmd
}

// GetTCPOrderQueue :获取TCP数据
func GetTCPOrderQueue() chan *cm.Cmd {
	return tcpOrderQueue
}

// SendMsgToWeb :回写到web前端
func SendMsgToWeb(rCmd *cm.Cmd) {
	WebMsgQueue <- rCmd
}

// RegisterReqOrdersUnion :登记到请求指令池
func RegisterReqOrdersUnion(rRO *RequestOrder) bool {
	// 指令池里如果有相同操作，终止操作
	for _, v := range ReqOrdersUnion.RequestOrders {
		if v.Args == rRO.Args && v.CmdType == rRO.CmdType && v.TargetID == rRO.TargetID {
			return false
		}
	}
	reqOrderID := NewReqOrdersUnionID(rRO.TargetID)
	ReqOrdersUnion.Lock()
	ReqOrdersUnion.RequestOrders[reqOrderID] = rRO
	ReqOrdersUnion.Unlock()
	return true
}

// UnregisterReqOrdersUnion :注销请求指令池
func UnregisterReqOrdersUnion(rTargetID string, rCmdType uint8) bool {
	for k, v := range ReqOrdersUnion.RequestOrders {
		// 如果被操作设备和操作指令类型相同，则认为已经操作成功，返回结果
		if v.TargetID == rTargetID && rCmdType == v.CmdType {
			// 针对一个设备的操作指令限定只有一个调用者
			ReqOrdersUnion.Lock()
			delete(ReqOrdersUnion.RequestOrders, k)
			ReqOrdersUnion.Unlock()
			return true
		}
	}
	return false
}

// GetReqOrderIDFromUnion :根据设备ID和操作类型获取Ws连接ID
func GetReqOrderIDFromUnion(rTargetID string, rCmdType uint8) string {
	for _, v := range ReqOrdersUnion.RequestOrders {
		if v.TargetID == rTargetID && rCmdType == v.CmdType {
			return v.RequestID
		}
	}
	return ""
}

// SendOrderToDeviceByTCP :添加到TCP指令发送队列
func SendOrderToDeviceByTCP(rRO *RequestOrder) error {
	// 指令池里如果有相同操作，终止操作，返回错误提示
	for _, v := range ReqOrdersUnion.RequestOrders {
		if v.TargetID == rRO.TargetID && v.CmdType == rRO.CmdType && v.Args == rRO.Args {
			return cm.ConvertStrToErr(DataHubMsg.CmdRepeatNotice)
		}
	}
	// 判断是否为检测器
	if ec.IsDetector(rRO.TargetID) {
		rcvID := ec.GetRcvID(rRO.TargetID)
		ipAddr := ec.GetRcvIP(ec.GetRcvID(rRO.TargetID))
		if rcvID == "" || ipAddr == "" {
			// 根据DetID获取RcvID和IP地址失败时，返回错误
			return cm.ConvertStrToErr(DataHubMsg.GetServerDataErr)
		}
		rcvIDbytes := cm.ConvertStrToBytesByPerTwoChar(rcvID)
		detIDbytes := cm.ConvertStrToBytesByPerTwoChar(rRO.TargetID)
		addDet := cm.ConvertHexUnitToDecUnit(tsc.TrsCmdType.AddDetect)
		delDet := cm.ConvertHexUnitToDecUnit(tsc.TrsCmdType.DelDetect)
		if rRO.CmdType == addDet || rRO.CmdType == delDet {
			// TCP指令ID:wsOrderID + 随机字符 + IP地址组成
			tcpOrderID := NewTCPOrderID(NewWSOrderID(rRO.RequestID), ipAddr)
			// FIXME:这里有问题,需重新考虑UnionID
			RegisterReqOrdersUnion(rRO)
			od := cm.NewOrder(tcpOrderID, tsc.CmdOperateDetect(rRO.CmdType, rcvIDbytes, 1, detIDbytes))
			AddToTCPOrderQueue(od)
		} else {
			return cm.ConvertStrToErr(DataHubMsg.CmdInvaildErr)
		}
	}
	return nil
}

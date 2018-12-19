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

// CheckReqOrderUnique 检测指定RequestOrder对象在集合中的唯一性
func CheckReqOrderUnique(rRO *RequestOrder) bool {
	for _, v := range ReqOrdersUnion.RequestOrders {
		if v.CmdType == rRO.CmdType && v.TargetID == rRO.TargetID && v.Args == rRO.Args {
			return false
		}
	}
	return true
}

// RegisterReqOrdersUnion :登记到请求指令池
// 【注意】：参数Args的如果有多项用逗号，分隔
func RegisterReqOrdersUnion(rRO *RequestOrder) bool {
	if !CheckReqOrderUnique(rRO) {
		// 如果该请求动作重复
		return false
	}
	ReqOrdersUnion.Lock()
	ReqOrdersUnion.RequestOrders = append(ReqOrdersUnion.RequestOrders, rRO)
	ReqOrdersUnion.Unlock()
	return true
}

// UnregisterReqOrdersUnion :注销请求指令池
func UnregisterReqOrdersUnion(rTargetID string, rCmdType uint8, rArgs string) bool {
	for i, v := range ReqOrdersUnion.RequestOrders {
		if v.CmdType == rCmdType && v.TargetID == rTargetID && v.Args == rArgs {
			ReqOrdersUnion.Lock()
			ReqOrdersUnion.RequestOrders = append(ReqOrdersUnion.RequestOrders[:i], ReqOrdersUnion.RequestOrders[i+1])
			ReqOrdersUnion.Unlock()
			return true
		}
	}
	return false
}

// GetReqOrderIDFromUnion :根据设备ID和操作类型获取Ws连接ID
func GetReqOrderIDFromUnion(rTargetID string, rCmdType uint8, rArgs string) string {
	for _, v := range ReqOrdersUnion.RequestOrders {
		if v.TargetID == rTargetID && rCmdType == v.CmdType && rArgs == v.Args {
			return v.RequestID
		}
	}
	return ""
}

// SendOrderToDeviceByTCP :添加到TCP指令发送队列
func SendOrderToDeviceByTCP(rRO *RequestOrder) error {
	// 指令池里如果有相同操作，终止操作，返回错误提示
	if !CheckReqOrderUnique(rRO) {
		return cm.ConvertStrToErr(DataHubMsg.CmdRepeatNotice)
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
			od := cm.NewCmd(tcpOrderID, tsc.CmdOperateDetect(rRO.CmdType, rcvIDbytes, 1, detIDbytes))
			AddToTCPOrderQueue(od)
		} else {
			return cm.ConvertStrToErr(DataHubMsg.CmdInvaildErr)
		}
	}
	return nil
}

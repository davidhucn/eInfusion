package datahub

import (
	cm "eInfusion/comm"
	wk "eInfusion/dbwork"
	ep "eInfusion/protocol"
	"strings"
)

// AddToTCPQueue ：通过TCP协议发送指令至设备
func addToTCPSendQueue(rCmd *cm.Cmd) {
	TCPOrderQueue <- rCmd
}

// SendMsgToWeb :回写到web前端
func SendMsgToWeb(rCmd *cm.Cmd) {
	WebMsgQueue <- rCmd
}

// NewTCPOrderID :生成TCP包约定指令序号
func NewTCPOrderID(rStrCnt string, rTCPConnectionID string) string {
	return rStrCnt + "@" + rTCPConnectionID
}

// DecodeToTCPConnID :解析指令ID为TCP连接序号
func DecodeToTCPConnID(rStrCnt string) string {
	return strings.Split(rStrCnt, "@")[1]
}

// RegisterReqOrdersUnion :登记到请求指令池
func RegisterReqOrdersUnion(rRO *RequestOrder) {
	// 指令池里如果有相同操作，终止操作
	for _, v := range ReqOrdersUnion.RequestOrders {
		if v == rRO {
			return
		}
	}
	reqOrderID := rRO.TargetID + "~" + cm.GetRandString(8)
	ReqOrdersUnion.Lock()
	ReqOrdersUnion.RequestOrders[reqOrderID] = rRO
	ReqOrdersUnion.Unlock()
}

// UnregisterReqOrdersUnion :登记到请求指令池
func UnregisterReqOrdersUnion(rReqOrderID string) {
	// TODO:接收到数据后即注销这一操作指令池记录,如果设备长时间无法通讯也注销
	ReqOrdersUnion.Lock()
	delete(ReqOrdersUnion.RequestOrders, rReqOrderID)
	ReqOrdersUnion.Unlock()
}

// SendOrderToDeviceByTCP :添加到TCP指令发送队列
func SendOrderToDeviceByTCP(rRO *RequestOrder) error {
	// 指令池里如果有相同操作，终止操作，返回错误提示
	for _, v := range ReqOrdersUnion.RequestOrders {
		if v == rRO {
			return cm.ConvertStrToErr(DataHubMsg.CmdRepeatNotice)
		}
	}
	// 判断是否为检测器
	if wk.IsDetector(rRO.TargetID) {
		rcvID := wk.GetRcvID(rRO.TargetID)
		ipAddr := wk.GetRcvIP(wk.GetRcvID(rRO.TargetID))
		if rcvID == "" || ipAddr == "" {
			// 根据DetID获取RcvID和IP地址失败时，返回错误
			return cm.ConvertStrToErr(DataHubMsg.GetServerDataErr)
		}
		rcvIDbytes := cm.ConvertStrToBytesByPerTwoChar(rcvID)
		detIDbytes := cm.ConvertStrToBytesByPerTwoChar(rRO.TargetID)
		addDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.AddDetect)
		delDet := cm.ConvertHexUnitToDecUnit(ep.TrsCmdType.DelDetect)
		if rRO.CmdType == addDet || rRO.CmdType == delDet {
			RegisterReqOrdersUnion(rRO)
			// TCP指令标识:wsOrderID + 随机字符 + IP地址组成
			tcpOrderID := NewTCPOrderID(rRO.RequestID, ipAddr)
			od := cm.NewOrder(tcpOrderID, ep.CmdOperateDetect(rRO.CmdType, rcvIDbytes, 1, detIDbytes))
			addToTCPSendQueue(od)
		} else {
			return cm.ConvertStrToErr(DataHubMsg.CmdInvaildErr)
			// logs.LogMain.Error("错误：没有目标设备编码或者无法获取相关设备编码！")
		}
	}
	return nil
}

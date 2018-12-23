package trsfus

import (
	cm "eInfusion/comm"
)

// ReceiverSendOrder :生成接收器发送命令
// SetReceiverReconnectTime ：参数重连时间以秒为单位.
func ReceiverSendOrder(cmdType CmdType, id string, args []string) []byte {
	sd := make([]byte, 0)
	for i := 0; i < len(PacketHeaderContent); i++ {
		sd = append(sd, PacketHeaderContent[i])
	}
	sd = append(sd, byte(0)) //指令总长度，占位
	//	获取接收器ID
	rcvID := cm.ConvertStrToBytesByPerTwoChar(id)
	switch cmdType {
	case GetReceiverState:
		//	获取指令类型
		sd = append(sd, SendCmdMap[GetReceiverState])
		for i := 0; i < len(rcvID); i++ {
			sd = append(sd, rcvID[i])
		}
	case SetReceiverConfig:
		sd = append(sd, SendCmdMap[SetReceiverConfig])
		for i := 0; i < len(rcvID); i++ {
			sd = append(sd, rcvID[i])
		}
		// ip地址
		ip := cm.ConvertStrIPAddrToBytes(args[0])
		for i := 0; i < len(ip); i++ {
			sd = append(sd, ip[i])
		}
		port := cm.ConvertDecToBytes(cm.ConvertBasStrToInt64(10, args[1]))
		for i := 0; i < len(port); i++ {
			sd = append(sd, port[i])
		}
	case SetReceiverReconnectTime: // 重连时间以秒为单位.
		sd = append(sd, SendCmdMap[SetReceiverReconnectTime])
		reconnectTime := cm.ConvertDecToBytes(cm.ConvertBasStrToInt64(10, args[0]))
		for i := 0; i < len(reconnectTime); i++ {
			sd = append(sd, reconnectTime[i])
		}
	}
	sd[1] = uint8(len(sd)) // 重新计算指令总长度
	// 记录到指令集合里面
	// OrdersQueue.LoadOrStore()
	return sd
}

// DetectorSendOrder :生成检测器发送命令
func DetectorSendOrder(cmdType CmdType, rcvid string, detid string) []byte {
	sd := make([]byte, 0)
	for i := 0; i < len(PacketHeaderContent); i++ {
		sd = append(sd, PacketHeaderContent[i])
	}
	sd = append(sd, byte(0)) //指令总长度，占位
	//	获取接收器ID
	rcvID := cm.ConvertStrToBytesByPerTwoChar(rcvid)
	detID := cm.ConvertStrToBytesByPerTwoChar(detid)
	switch cmdType {
	case GetDetectorState:
		sd = append(sd, SendCmdMap[GetDetectorState])

	case DeleteDetector:
		sd = append(sd, SendCmdMap[DeleteDetector])
	}
	for i := 0; i < len(rcvID); i++ {
		sd = append(sd, rcvID[i])
	}
	sd = append(sd, byte(1))
	for i := 0; i < len(detID); i++ {
		sd = append(sd, detID[i])
	}
	sd[1] = uint8(len(sd)) // 重新计算指令总长度
	return sd
}

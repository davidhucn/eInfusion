package trsfus

import (
	cm "eInfusion/comm"
)

// MakeOrderOnReceiver :生成接收器发送命令
// SetReceiverReconnectTime ：参数重连时间以秒为单位.
func MakeOrderOnReceiver(cmdType CmdType, id string, args []string) []byte {
	sd := make([]byte, 0)
	if cmdType == CmdAddDetector || cmdType == CmdDeleteDetector || cmdType == CmdGetDetectorState {
		return sd
	}
	for i := 0; i < len(packetHeaderPrefix); i++ {
		sd = append(sd, packetHeaderPrefix[i])
	}
	sd = append(sd, byte(0)) //指令总长度，占位
	//	获取接收器ID
	rcvID := cm.ConvertStrToBytesByPerTwoChar(id)
	switch cmdType {
	case CmdGetReceiverState:
		//	获取指令类型
		sd = append(sd, SendCmdMap[CmdGetReceiverState])
		for i := 0; i < len(rcvID); i++ {
			sd = append(sd, rcvID[i])
		}
	case CmdSetReceiverConfig:
		sd = append(sd, SendCmdMap[CmdSetReceiverConfig])
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
	case CmdSetReceiverReconnectTime: // 重连时间以秒为单位.
		sd = append(sd, SendCmdMap[CmdSetReceiverReconnectTime])
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

// MakeOrderOnDetector :生成检测器发送命令
func MakeOrderOnDetector(cmdType CmdType, rcvid string, detid string) []byte {
	sd := make([]byte, 0)
	if cmdType == CmdGetReceiverState || cmdType == CmdSetReceiverConfig || cmdType == CmdSetReceiverReconnectTime {
		return sd
	}
	for i := 0; i < len(packetHeaderPrefix); i++ {
		sd = append(sd, packetHeaderPrefix[i])
	}
	sd = append(sd, byte(0)) //指令总长度，占位
	//	获取接收器ID
	rcvID := cm.ConvertStrToBytesByPerTwoChar(rcvid)
	detID := cm.ConvertStrToBytesByPerTwoChar(detid)
	switch cmdType {
	case CmdGetDetectorState:
		sd = append(sd, SendCmdMap[CmdGetDetectorState])

	case CmdDeleteDetector:
		sd = append(sd, SendCmdMap[CmdDeleteDetector])
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

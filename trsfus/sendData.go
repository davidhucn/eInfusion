package trsfus

import (
	cm "eInfusion/comm"
)

// MakeSendOrder :生成接收器相关发送命令
// SetReceiverReconnectTime ：参数重连时间以秒为单位.
func MakeSendOrder(cmdType CmdType, rcvid string, detid string, args []string) []byte {
	sd := make([]byte, 0)
	od := NewOrder(rcvid, detid, cmdType, args)
	if od.matchFromOrderPool() > -1 {
		// 如果指令池内已存在，则返回空指令
		return sd
	}

	for i := 0; i < len(packetHeaderPrefix); i++ {
		sd = append(sd, packetHeaderPrefix[i])
	}
	sd = append(sd, byte(0)) //指令总长度，占位
	//	获取接收器ID
	rcvID := cm.ConvertStrToBytesByPerTwoChar(rcvid)
	detID := cm.ConvertStrToBytesByPerTwoChar(detid)
	sd = append(sd, SendCmdMap[cmdType])
	for i := 0; i < len(rcvID); i++ {
		sd = append(sd, rcvID[i])
	}
	switch cmdType {
	case CmdGetReceiverState:
		// 空
	case CmdSetReceiverConfig:
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
		reconnectTime := cm.ConvertDecToBytes(cm.ConvertBasStrToInt64(10, args[0]))
		for i := 0; i < len(reconnectTime); i++ {
			sd = append(sd, reconnectTime[i])
		}
	case CmdGetDetectorState:
		sd = append(sd, byte(1))
		for i := 0; i < len(detID); i++ {
			sd = append(sd, detID[i])
		}
	case CmdAddDetector:
		sd = append(sd, byte(1))
		for i := 0; i < len(detID); i++ {
			sd = append(sd, detID[i])
		}
	case CmdDeleteDetector:
		sd = append(sd, byte(1))
		for i := 0; i < len(detID); i++ {
			sd = append(sd, detID[i])
		}
	}
	sd[1] = uint8(len(sd)) // 重新计算指令总长度

	// 记录到指令池里面
	od.RegisteToOrderPool()
	return sd
}

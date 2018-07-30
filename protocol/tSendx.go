package protocol

import . "eInfusion/comm"

// 获取指定接收器的状态
func CmdGetRcvStatus(ref_RcvID []byte) []byte {
	var intOrderDataLength = 7
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = G_TsCmd.Header
	sendOrders[1] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, intOrderDataLength))
	//	获取指令类型
	sendOrders[2] = G_TsCmd.GetRcv
	//	获取接收器ID
	for recId := 0; recId < 4; recId++ {
		sendOrders[recId+3] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, ref_RcvID[recId]))
	}
	return sendOrders
}

// 对检测器进行操作（检测、册除、新增）
// 一个设备ID 占4个byte
func CmdOperateDetect(orderType uint8, ref_RcvID []byte, detectAmount int, ref_DetectID []byte) []byte {
	var intOrderDataLength = 7
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = G_TsCmd.Header
	sendOrders[1] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, intOrderDataLength))
	//	获取指令类型
	sendOrders[2] = orderType
	//	获取接收器ID
	for recId := 0; recId < 4; recId++ {
		sendOrders[recId+3] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, ref_RcvID[recId]))
	}
	// 如果与检测器相关的操作(添加、删除、检查)
	if orderType == G_TsCmd.AddDetect || orderType == G_TsCmd.DelDetect || orderType == G_TsCmd.GetDetect {
		if detectAmount > 0 {
			// 检测器数量内容到slice
			sendOrders = append(sendOrders, byte(detectAmount))
			//	添加检测器id到slice
			for devId := 0; devId < 4; devId++ {
				sendOrders = append(sendOrders, ConvertBasStrToUint(16, ConvertBasNumberToStr(16, ref_DetectID[devId])))
			}
		}
	}
	// 根据长度调整第二位
	sendOrders[1] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, len(sendOrders)))
	return sendOrders
}

// 修改接收器网络配置
// FIXME:生成的数据有问题 ，再核对
func CmdSetRcvCfg(ref_RcvID []byte, ref_IP []byte, ref_Port []byte) []byte {
	var intOrderDataLength = 13
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = G_TsCmd.Header
	sendOrders[1] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, intOrderDataLength))
	//	获取指令类型
	sendOrders[2] = G_TsCmd.SetRcvNetCfg
	//	获取接收器ID
	for recId := 0; recId < 4; recId++ {
		sendOrders[recId+3] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, ref_RcvID[recId]))
	}
	// IP地址
	for ipAdd := 0; ipAdd < 4; ipAdd++ {
		sendOrders[ipAdd+7] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, ref_IP[ipAdd]))
	}
	//	端口号
	for portNum := 0; portNum < 2; portNum++ {
		sendOrders[portNum+11] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, ref_Port[portNum]))
	}
	return sendOrders
}

// 设置接收器重连接时间
func CmdSetRcvReconTime(ref_RcvID []byte, ref_IP []byte, ref_ReconTime int) []byte {
	var intOrderDataLength = 9
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = G_TsCmd.Header
	sendOrders[1] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, intOrderDataLength))
	//	获取指令类型
	sendOrders[2] = G_TsCmd.SetReconnTime
	//	获取接收器ID
	for recId := 0; recId < 4; recId++ {
		sendOrders[recId+3] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, ref_RcvID[recId]))
	}
	// 连接时间
	for period := 0; period < 2; period++ {
		sendOrders[period+7] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, ref_IP[period]))
	}
	return sendOrders
}

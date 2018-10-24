package trsfscomm

// CmdGetRcvStatus :获取指定接收器的状态
func CmdGetRcvStatus(rRcvID []byte) []byte {
	var intOrderDataLength = 7
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = TrsDefin.Header
	sendOrders[1] = byte(intOrderDataLength)
	//	获取指令类型
	sendOrders[2] = TrsCmdType.GetRcv
	//	获取接收器ID
	for recID := 0; recID < 4; recID++ {
		sendOrders[recID+3] = rRcvID[recID]
	}
	return sendOrders
}

// CmdOperateDetect :对检测器进行操作（检测、册除、新增）
// 一个设备ID 占4个byte
func CmdOperateDetect(orderType uint8, rRcvID []byte, detectAmount int, rDetID []byte) []byte {
	var intOrderDataLength = 7
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = TrsDefin.Header
	sendOrders[1] = byte(intOrderDataLength)
	//	获取指令类型
	sendOrders[2] = orderType
	//	获取接收器ID
	for recID := 0; recID < 4; recID++ {
		sendOrders[recID+3] = rRcvID[recID]
	}
	// 如果与检测器相关的操作(添加、删除、检查)
	if orderType == TrsCmdType.AddDetect || orderType == TrsCmdType.DelDetect || orderType == TrsCmdType.GetDetect {
		if detectAmount > 0 {
			// 检测器数量内容到slice
			sendOrders = append(sendOrders, byte(detectAmount))
			//	添加检测器id到slice
			for devID := 0; devID < 4; devID++ {
				sendOrders = append(sendOrders, rDetID[devID])
			}
		}
	}
	// 根据长度调整第二位
	sendOrders[1] = byte(len(sendOrders))
	return sendOrders
}

// CmdSetRcvCfg :生成指令，修改接收器网络配置
func CmdSetRcvCfg(rRcvID []byte, rIP []byte, rPort []byte) []byte {
	var intOrderDataLength = 13
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = TrsDefin.Header
	// sendOrders[1] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, intOrderDataLength))
	sendOrders[1] = byte(intOrderDataLength)
	//	获取指令类型
	sendOrders[2] = TrsCmdType.SetRcvNetCfg
	//	获取接收器ID
	for recID := 0; recID < 4; recID++ {
		sendOrders[recID+3] = rRcvID[recID]
	}
	// IP地址
	for ipAdd := 0; ipAdd < 4; ipAdd++ {
		// sendOrders[ipAdd+7] = ConvertBasStrToUint(16, ConvertBasNumberToStr(16, rIP[ipAdd]))
		sendOrders[ipAdd+7] = rIP[ipAdd]
	}
	//	端口号
	for portNum := 0; portNum < 2; portNum++ {
		// sendOrders[portNum+11] = ConvertBasStrToUint(10, ConvertBasNumberToStr(10, rPort[portNum]))
		sendOrders[portNum+11] = rPort[portNum]
	}
	return sendOrders
}

// CmdSetRcvReconTime :设置接收器重连接时间
// FIXME:错误，需要修改，发现遗漏重新连接参数
func CmdSetRcvReconTime(rRcvID []byte, rReconTime []byte) []byte {
	var intOrderDataLength = 9
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = TrsDefin.Header
	sendOrders[1] = byte(intOrderDataLength)
	//	获取指令类型
	sendOrders[2] = TrsCmdType.SetReconnTime
	//	获取接收器ID
	for recID := 0; recID < 4; recID++ {
		sendOrders[recID+3] = rRcvID[recID]
	}
	// 连接时间
	for period := 0; period < 2; period++ {
		sendOrders[period+7] = rReconTime[period]
	}
	return sendOrders
}

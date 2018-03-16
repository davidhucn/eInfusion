package protocol

//import ""
// 对检测器进行操作（检测、册除、新增）
// 一个设备ID 占4个byte
func OperateDetect(orderType uint8, ref_RcvID []byte, detectAmount int, ref_DetectID []byte) []byte {
	var intOrderDataLength = 7
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = c_metaDataHeader
	sendOrders[1] = byte(intOrderDataLength)
	//	获取指令类型
	sendOrders[2] = byte(orderType)
	//	获取接收器ID
	for recId := 0; recId < len(ref_RcvID); recId++ {
		sendOrders[recId+3] = ref_RcvID[recId]
	}
	// 如果与检测器相关的操作(添加、删除、检查)
	if orderType == C_orderType_addDetect || orderType == C_orderType_delDetect || orderType == C_orderType_getDetectStat {
		if detectAmount > 0 {
			// 检测器数量内容到slice
			sendOrders = append(sendOrders, byte(detectAmount))
			//	添加检测器id到slice
			for devId := 0; devId < len(ref_DetectID); devId++ {
				sendOrders = append(sendOrders, ref_DetectID[devId])
			}
		}
	}
	return sendOrders
}

// 修改接收器网络配置
func SetRcvCfg(ref_RcvID []byte, ref_IP []byte, ref_Port []byte) []byte {
	var intOrderDataLength = 13
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = c_metaDataHeader
	sendOrders[1] = byte(intOrderDataLength)
	//	获取指令类型
	sendOrders[2] = C_orderType_setRcvCfg
	//	获取接收器ID
	for recId := 0; recId < 4; recId++ {
		sendOrders[recId+3] = ref_RcvID[recId]
	}
	// IP地址
	for ipAdd := 0; ipAdd < 4; ipAdd++ {
		sendOrders[ipAdd+7] = ref_IP[ipAdd]
	}
	//	端口号
	for portNum := 0; portNum < 2; portNum++ {
		sendOrders[portNum+11] = ref_Port[portNum]
	}
	return sendOrders
}

// 设置接收器重连接时间
func SetRcvReconTime(ref_RcvID []byte, ref_IP []byte, ref_ReconTime int) []byte {
	var intOrderDataLength = 9
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = c_metaDataHeader
	sendOrders[1] = byte(intOrderDataLength)
	//	获取指令类型
	sendOrders[2] = C_orderType_reconnTimePeriod
	//	获取接收器ID
	for recId := 0; recId < 4; recId++ {
		sendOrders[recId+3] = ref_RcvID[recId]
	}
	// 连接时间
	for period := 0; period < 2; period++ {
		sendOrders[period+7] = ref_IP[period]
	}
	return sendOrders
}

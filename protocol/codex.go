package protocol

import (
	"eInfusion/comm"
	"eInfusion/dbOperate"
)

const (
	//	接收数据长度(10进制)
	c_metaDataHeaderLength = 2
	//	包数据中定义长度的帧	(10进制）
	c_metaDataLengthCursor = 1
)

const (
	//	数据协议（报文）头(16进制)
	c_metaDataHeader = 66
	///////////以下为被动接收////////////////////////
	// 接收接收器状态(16进制)
	c_statusValue_recRcvStat = 0
	// 接收检测器状态(16进制)
	c_statusValue_recDetectStat = 1
	//	添加检测器到接收器成功
	c_statusValue_addRcvSuccess = 2
	//	删除检测器成功
	c_statusValue_deleteRcvSuccess = 3
)

///////////以下为主动操作///////////////////////
const (
	//	获取接收器状态(16进制)
	C_orderType_getRcvStat = 10
	//	获取检测器状态(16进制)
	C_orderType_getDetectStat = 11
	//	添加检测器到接受器
	C_orderType_addDetect = 12
	//	删除检测器(16进制)
	C_orderType_delDetect = 13
	//	设置接收器网络配置（IP和port) (16进制)
	C_orderType_setRcvCfg = 14
	//	设备接收器重连接时间
	C_orderType_reconnTimePeriod = 15
)

// 获取包头长度数值
func GetDataHeaderLength() int {
	return c_metaDataHeaderLength
}

//	判断包头是否正确（进制转换）
// 返回：包头是否为真（布尔值），数据包内正文数据包的长度
func DecodeHeader(ref_packHeader []byte, adr_dataLength *int) bool {
	var blnRet bool = false
	var intDataLength int = 0
	// 如果包头长度正确
	if len(ref_packHeader) == c_metaDataHeaderLength {
		//	如果接收的包头内容正确
		if comm.BaseConvert(16, ref_packHeader[0]) == comm.BaseConvert(1, c_metaDataHeader) {
			//	获取包内数据帧的长度,根据协议规定
			intDataLength = int(ref_packHeader[c_metaDataLengthCursor])
			//	包内容帧长 = 包总长度- 包头帧长度
			intDataLength = intDataLength - c_metaDataHeaderLength
			// 函数返回为真
			blnRet = true
		}
	}
	*adr_dataLength = intDataLength
	return blnRet
}

//	处理接收到的包内数据
func DecodeReceiveData(ref_packData []byte) {
	switch ref_packData[0] {
	case c_statusValue_recDetectStat:
		comm.ShowScreen("收到检测器状态..设备编号：", comm.BaseConvert(10, ref_packData[1]), comm.BaseConvert(10, ref_packData[2]),
			comm.BaseConvert(10, ref_packData[3]), comm.BaseConvert(10, ref_packData[4]))
		comm.ShowScreen("其它数据：", ref_packData[5])
	case c_statusValue_recRcvStat:
		comm.ShowScreen("收到接收器状态...，设备编号：", comm.BaseConvert(10, ref_packData[1]), comm.BaseConvert(10, ref_packData[2]),
			comm.BaseConvert(10, ref_packData[3]), comm.BaseConvert(10, ref_packData[4]))
		comm.ShowScreen("其它数据：", ref_packData[5])
	default:
	}
}

// 获取指定接收器的状态
func GetRcvStatus(ref_RcvID []byte) []byte {
	var intOrderDataLength = 7
	//	基本指令内容
	sendOrders := make([]byte, intOrderDataLength)
	sendOrders[0] = c_metaDataHeader
	sendOrders[1] = byte(intOrderDataLength)
	//	sendOrders[1] = comm.ConvertIntToBytes(intOrderDataLength)[0]
	//	获取指令类型
	sendOrders[2] = C_orderType_getRcvStat
	//	获取接收器ID
	for recId := 0; recId < 4; recId++ {
		sendOrders[recId+3] = ref_RcvID[recId]
	}
	return sendOrders
}

// 对检测器进行操作（检测、册除、新增）
// 一个设备ID 占4个byte
func OperateDetect(orderType int, ref_RcvID []byte, detectAmount int, ref_DetectID []byte) []byte {
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

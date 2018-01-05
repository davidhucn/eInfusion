package protocol

import (
	"eInfusion/comm"
)

const (
	//	接收数据长度(10进制)
	c_metaDataHeaderLength = 2
	//	包数据中定义长度的帧	(10进制）
	c_metaDataLengthCursor = 1
	//	数据协议头(16进制)
	c_metaDataHeader = 66
	///////////以下为被动接收////////////
	// 接收接收器状态(16进制)
	c_orderType_recRcvStat = 0
	// 接收检测器状态(16进制)
	c_orderType_recDetectStat = 8
	///////////以下为主动操作//////////
	//	获取接收器状态(16进制)
	c_orderType_getRcvStat = 16
	//	获取检测器状态(16进制)
	c_orderType_getDetectStat = 24
	//	删除检测器(16进制)
	c_orderType_delDetect = 32
	//	设置接收器网络配置（IP和port) (16进制)
	c_orderType_setRcvCfg = 40
)

// 获取包头长度数值
func GetDataHeaderLength() int {
	return c_metaDataHeaderLength
}

//	判断包头是否正确
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

//	处理包内数据
func DecodeData(ref_packData []byte) {
	switch ref_packData[0] {
	case c_orderType_recDetectStat:
		comm.ShowScreen("收到检测器状态..设备编号：", comm.BaseConvert(10, ref_packData[1]), comm.BaseConvert(10, ref_packData[2]),
			comm.BaseConvert(10, ref_packData[3]), comm.BaseConvert(10, ref_packData[4]))
		comm.ShowScreen("其它数据：", ref_packData[5])
	case c_orderType_recRcvStat:
		comm.ShowScreen("收到接收器状态...，设备编号：", comm.BaseConvert(10, ref_packData[1]), comm.BaseConvert(10, ref_packData[2]),
			comm.BaseConvert(10, ref_packData[3]), comm.BaseConvert(10, ref_packData[4]))
		comm.ShowScreen("其它数据：", ref_packData[5])
	default:
	}
}

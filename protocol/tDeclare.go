package protocol

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

//指令集包头信息
//type HeaderPacket struct {
//	Header byte
//	Length byte
//}

////指令集报文内容包
//type ReportPacket struct {
//	FuncKey     byte
//	Receiver    [4]byte
//	DataContent []byte
//}

//检测器对象
type Detector struct {
	Qcode      string
	ReceiverID string
	Disable    bool
}

//接收器对象
type Receiver struct {
	Qcode      string
	DetectorID string
	IPAddr     string
}

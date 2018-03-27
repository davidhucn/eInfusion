package protocol

const (
	//	接收数据长度(10进制)
	c_metaDataHeaderLength int = 2
	//	包数据中定义长度的帧	(10进制）
	c_metaDataLengthCursor int = 1
)

const (
	//	数据协议（报文）头(16进制)
	c_metaDataHeader uint8 = 0x66
	///////////以下为被动接收////////////////////////
	// 接收接收器状态(16进制)
	c_stRcvStat uint8 = 0x00
	// 接收检测器状态(16进制)
	c_stDetectStat uint8 = 0x01
	//	添加检测器到接收器成功
	c_stAddDetectSuccess uint8 = 0x02
	//	删除检测器成功
	c_stDelDetectSuccess uint8 = 0x03
)

///////////以下为主动操作///////////////////////
const (
	//	获取接收器状态(16进制)
	C_orderType_getRcvStat uint8 = 0x10
	//	获取检测器状态(16进制)
	C_orderType_getDetectStat uint8 = 0x11
	//	添加检测器到接受器
	C_orderType_addDetect uint8 = 0x12
	//	删除检测器(16进制)
	C_orderType_delDetect uint8 = 0x13
	//	设置接收器网络配置（IP和port) (16进制)
	C_orderType_setRcvCfg uint8 = 0x14
	//	设备接收器重连接时间
	C_orderType_reconnTimePeriod uint8 = 0x15
)

//检测器对象
type Detector struct {
	ID         string
	ReceiverID string
	Stat       string
	Disable    bool
}

//接收器对象
type Receiver struct {
	Qcode      string
	DetectorID string
	IPAddr     string
}

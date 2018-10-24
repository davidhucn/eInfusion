// trsfscomm :输液警报相关业务协议
package trsfscomm

type transfusionTCPDefination struct {
	Header           uint8
	HeaderLength     int
	PackLengthCursor uint8
}

// TransfusionCmdType 输液报警指令码
type TransfusionCmdType struct {
	RcvState           uint8
	DetectState        uint8
	AddDetectState     uint8
	DelDetectState     uint8
	SetRcvNetCfgState  uint8
	SetReconnTimeState uint8
	GetRcv             uint8
	GetDetect          uint8
	AddDetect          uint8
	DelDetect          uint8
	SetRcvNetCfg       uint8
	SetReconnTime      uint8
}

// TrsDefin :TCP定义
var TrsDefin transfusionTCPDefination

// TrsCmdType :Transfusion指令集
var TrsCmdType TransfusionCmdType

func init() {
	//数据协议（报文）头(16进制)
	TrsDefin.Header = 0x66
	//	接收数据长度(10进制)
	TrsDefin.HeaderLength = 2
	//	包数据中定义长度的帧(10进制）
	TrsDefin.PackLengthCursor = 1
	///////////receive state//////////
	TrsCmdType.RcvState = 0x00
	TrsCmdType.DetectState = 0x01
	TrsCmdType.AddDetectState = 0x02
	TrsCmdType.DelDetectState = 0x03
	TrsCmdType.SetRcvNetCfgState = 0x04
	TrsCmdType.SetReconnTimeState = 0x05
	////////////action cmd////////////
	TrsCmdType.GetRcv = 0x10
	TrsCmdType.GetDetect = 0x11
	TrsCmdType.AddDetect = 0x12
	TrsCmdType.DelDetect = 0x13
	//	设置接收器网络配置（IP和port) (16进制)
	TrsCmdType.SetRcvNetCfg = 0x14
	//	设备接收器重连接时间
	TrsCmdType.SetReconnTime = 0x15

}

//Detector :检测器对象
// Stat:工作状态,0-关机，1-开机
// Alarm: 是否报警，输液条没有液体
type Detector struct {
	QRCode   string
	ID       string
	RcvID    string
	Capacity uint8 //0,1,2,3
	PowerOn  uint8 //工作状态：0-关机，1-开机
	Alarm    uint8 //是否报警，0-正常，1－报警，无药水
}

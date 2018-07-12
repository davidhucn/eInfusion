package protocol

type sTransfusionCode struct {
	Header             uint8
	HeaderLength       int
	PackLengthCursor   uint8
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

//全局输液报警tcp编码
var G_TsCmd sTransfusionCode

func init() {
	//数据协议（报文）头(16进制)
	G_TsCmd.Header = 0x66
	//	接收数据长度(10进制)
	G_TsCmd.HeaderLength = 2
	//	包数据中定义长度的帧(10进制）
	G_TsCmd.PackLengthCursor = 1
	///////////receive state//////////////
	G_TsCmd.RcvState = 0x00
	G_TsCmd.DetectState = 0x01
	G_TsCmd.AddDetectState = 0x02
	G_TsCmd.DelDetectState = 0x03
	G_TsCmd.SetRcvNetCfgState = 0x04
	G_TsCmd.SetReconnTimeState = 0x05
	////////////action cmd///////////////
	G_TsCmd.GetRcv = 0x10
	G_TsCmd.GetDetect = 0x11
	G_TsCmd.AddDetect = 0x12
	G_TsCmd.DelDetect = 0x13
	//	设置接收器网络配置（IP和port) (16进制)
	G_TsCmd.SetRcvNetCfg = 0x14
	//	设备接收器重连接时间
	G_TsCmd.SetReconnTime = 0x15

}

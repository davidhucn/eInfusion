package protocol

import (
	"eInfusion/comm"
	wk "eInfusion/dbworks"
)

//DecodeHeader :判断包头是否正确（进制转换）
// 返回：包头是否为真（布尔值），数据包内正文数据包的长度
func DecodeHeader(ref_packHeader []byte, adr_dataLength *int) bool {
	blnRet := false
	intDataLength := 0
	// 如果包头长度正确
	if len(ref_packHeader) == TrsDefin.HeaderLength {
		//	如果接收的包头内容正确
		if ref_packHeader[0] == TrsDefin.Header {
			//	获取包内数据帧的长度,根据协议规定
			intDataLength = int(ref_packHeader[TrsDefin.PackLengthCursor])
			//	包内数据长度不能为0
			if intDataLength == 0 {
				return false
			}
			//	包内容帧长 = 包总长度 - 包头帧长度
			intDataLength = intDataLength - TrsDefin.HeaderLength
			// 函数返回为真
			blnRet = true
		}
	}
	*adr_dataLength = intDataLength
	return blnRet
}

//	处理接收到的包内数据
func DecodeRcvData(ref_packData []byte, ref_ipAddr string) {
	//	初始化t_device_dict
	//	InitDetInfoToDB(8)

	switch ref_packData[0] {
	//取得接收器状态（得接收器数目）
	case TrsCmdType.RcvState:
		wk.ReceiveRcvStat(ref_packData[1:], ref_ipAddr)
	case TrsCmdType.DetectState:
		wk.ReceiveDetectStat(ref_packData[1:], ref_ipAddr)
	case TrsCmdType.DelDetectState:
		wk.ReceiveDeleteDetect(ref_packData[1:], ref_ipAddr)
	case TrsCmdType.AddDetectState:
		wk.ReceiveAddDetect(ref_packData[1:], ref_ipAddr)
	case TrsCmdType.SetRcvNetCfgState:
		wk.ReceiveSetRcvNetCfgStat(ref_packData[1:], ref_ipAddr)
	case TrsCmdType.SetReconnTimeState:
		wk.ReceiveSetReconnTimeStat(ref_packData[1:], ref_ipAddr)
	default:
		comm.Msg("调试信息，无效数据...")
		//		return
	}
}

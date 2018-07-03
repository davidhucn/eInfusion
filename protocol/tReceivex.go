package protocol

import (
	"eInfusion/comm"
	. "eInfusion/dbWorks"
)

//	判断包头是否正确（进制转换）
// 返回：包头是否为真（布尔值），数据包内正文数据包的长度
func DecodeHeader(ref_packHeader []byte, adr_dataLength *int) bool {
	var blnRet bool = false
	var intDataLength int = 0
	// 如果包头长度正确
	if len(ref_packHeader) == G_TsCmd.HeaderLength {
		//	如果接收的包头内容正确
		if ref_packHeader[0] == G_TsCmd.Header {
			//	获取包内数据帧的长度,根据协议规定
			intDataLength = int(ref_packHeader[G_TsCmd.PackLengthCursor])
			//	包内数据长度不能为0
			if intDataLength == 0 {
				return false
			}
			//	包内容帧长 = 包总长度 - 包头帧长度
			intDataLength = intDataLength - G_TsCmd.HeaderLength
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
	case G_TsCmd.RcvState:
		ReceiveRcvStat(ref_packData[1:], ref_ipAddr)
	case G_TsCmd.DetectState:
		ReceiveDetectStat(ref_packData[1:], ref_ipAddr)
	case G_TsCmd.DelDetectState:
		ReceiveDeleteDetect(ref_packData[1:], ref_ipAddr)
	case G_TsCmd.AddDetectState:
		ReceiveAddDetect(ref_packData[1:], ref_ipAddr)
	case G_TsCmd.SetRcvNetCfgState:
		ReceiveSetRcvNetCfgStat(ref_packData[1:], ref_ipAddr)
	case G_TsCmd.SetReconnTimeState:
	default:
		comm.Msg("调试信息，无效数据...")
		//		return
	}
}

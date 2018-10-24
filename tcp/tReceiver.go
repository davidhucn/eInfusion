package tcp

import tsc "eInfusion/trsfscomm"

//DecodeHeader :判断包头是否正确（进制转换）
// 返回：包头是否为真（布尔值），数据包内正文数据包的长度
func DecodeHeader(rData []byte, rLen *int) bool {
	blnRet := false
	intDataLength := 0
	// 如果包头长度正确
	if len(rData) == tsc.TrsDefin.HeaderLength {
		//	如果接收的包头内容正确
		if rData[0] == tsc.TrsDefin.Header {
			//	获取包内数据帧的长度,根据协议规定
			intDataLength = int(rData[tsc.TrsDefin.PackLengthCursor])
			//	包内数据长度不能为0
			if intDataLength == 0 {
				return false
			}
			//	包内容帧长 = 包总长度 - 包头帧长度
			intDataLength = intDataLength - tsc.TrsDefin.HeaderLength
			// 函数返回为真
			blnRet = true
		}
	}
	*rLen = intDataLength
	return blnRet
}

//DecodeRcvData :处理接收到的包内数据
func DecodeRcvData(rData []byte, rIPAddr string) {
	//	初始化t_device_dict
	//	InitDetInfoToDB(8)

	// 接受指令成功，则返回结果至调用终端
	switch rData[0] {
	case tsc.TrsCmdType.RcvState:
		ReceiveRcvStat(rData[1:], rIPAddr, tsc.TrsCmdType.GetRcv)
	case tsc.TrsCmdType.DetectState:
		ReceiveDetectStat(rData[1:], rIPAddr, tsc.TrsCmdType.GetDetect)
	case tsc.TrsCmdType.DelDetectState:
		ReceiveDeleteDetect(rData[1:], rIPAddr, tsc.TrsCmdType.DelDetect)
	case tsc.TrsCmdType.AddDetectState:
		ReceiveAddDetect(rData[1:], rIPAddr, tsc.TrsCmdType.AddDetect)
	case tsc.TrsCmdType.SetRcvNetCfgState:
		ReceiveSetRcvNetCfgStat(rData[1:], rIPAddr, tsc.TrsCmdType.SetRcvNetCfg)
	case tsc.TrsCmdType.SetReconnTimeState:
		ReceiveSetReconnTimeStat(rData[1:], rIPAddr, tsc.TrsCmdType.SetReconnTime)
	default:
		// comm.Msg("调试信息，无效数据..., 内容：", rData)
		//		return
	}
}

package protocol

import (
	"eInfusion/comm"
	wk "eInfusion/dbwork"
)

//DecodeHeader :判断包头是否正确（进制转换）
// 返回：包头是否为真（布尔值），数据包内正文数据包的长度
func DecodeHeader(rData []byte, rLen *int) bool {
	blnRet := false
	intDataLength := 0
	// 如果包头长度正确
	if len(rData) == TrsDefin.HeaderLength {
		//	如果接收的包头内容正确
		if rData[0] == TrsDefin.Header {
			//	获取包内数据帧的长度,根据协议规定
			intDataLength = int(rData[TrsDefin.PackLengthCursor])
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
	*rLen = intDataLength
	return blnRet
}

//DecodeRcvData :处理接收到的包内数据
func DecodeRcvData(rData []byte, rIPAddr string) {
	//	初始化t_device_dict
	//	InitDetInfoToDB(8)

	// 接受指令成功，则返回结果至调用终端
	switch rData[0] {
	case TrsCmdType.RcvState:
		if wk.ReceiveRcvStat(rData[1:], rIPAddr) {

		}
	case TrsCmdType.DetectState:
		wk.ReceiveDetectStat(rData[1:], rIPAddr)
	case TrsCmdType.DelDetectState:
		wk.ReceiveDeleteDetect(rData[1:], rIPAddr)
	case TrsCmdType.AddDetectState:
		if wk.ReceiveAddDetect(rData[1:], rIPAddr) {
			// 如果添加检测器成功，即开启

		}
	case TrsCmdType.SetRcvNetCfgState:
		wk.ReceiveSetRcvNetCfgStat(rData[1:], rIPAddr)
	case TrsCmdType.SetReconnTimeState:
		wk.ReceiveSetReconnTimeStat(rData[1:], rIPAddr)
	default:
		comm.Msg("调试信息，无效数据..., 内容：", rData)
		//		return
	}
}

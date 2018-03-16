package protocol

//	"eInfusion/comm"
//	ed "eInfusion/dbOperate"
//	"reflect"

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
		if ref_packHeader[0] == c_metaDataHeader {
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
func DecodeRcvData(ref_packData []byte) {

	switch ref_packData[0] {
	//获取接收器状态
	case c_stRcvStat:

	//获取检测器状态
	case c_stDetectStat:

	case c_stDelDetectSuccess:
	case c_stAddDetectSuccess:
		//	default
		//		return
	}
}

package trsfus

// GetReceiveCmdType : 获取接收数据指令类型
// cursor :下标 (0~N)
func GetReceiveCmdType(p []byte, cursor int) {
	if cursor >= 0 && cursor < len(p) {
		ct, ok := ReceiveCmdMap[p[cursor]]
		if ok {
			switch ct {
			case GetReceiverState:

			case GetDetectorState:

			default:
			}
		}
	}
}

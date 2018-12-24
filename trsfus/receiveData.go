package trsfus

// DoReceverData : 获取接收数据指令类型
// cmdTypeCursor :下标 (0 ~ N)
func DoReceverData(p []byte, cmdTypeCursor int) {
	if cmdTypeCursor >= 0 && cmdTypeCursor < len(p) {
		ct, ok := ReceiveCmdMap[p[cmdTypeCursor]]
		if ok {
			switch ct {
			case CmdGetReceiverState:

			case CmdGetDetectorState:

			default:
			}
		}
	}
}

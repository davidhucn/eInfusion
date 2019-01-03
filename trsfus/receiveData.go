package trsfus

import (
	"eInfusion/comm"
	"errors"
)

// DoReceiveData : 获取接收数据指令类型, cmdTypeCursor : 指令相关下标 (0 ~ N)
func DoReceiveData(p []byte, cmdTypeCursor int) error {
	if cmdTypeCursor >= 0 && cmdTypeCursor < len(p) {
		ct, ok := ReceiveCmdMap[p[cmdTypeCursor]]
		if ok {
			switch ct {
			case CmdGetReceiverState:
				if len(p) != 7 {
					return errors.New("错误，数据长度超出定义！")
				}
				return getUploadReceiverState(p[3:])
			case CmdGetDetectorState:

			case CmdAddDetector:

			case CmdSetReceiverConfig:

			case CmdSetReceiverReconnectTime:

			}
		} else {
			return errors.New("错误，指令功能字错误！")
		}
	}
	return nil
}

// getUploadReceiverState :获取上传接收器信息
func getUploadReceiverState(p []byte) error {
	rcvID := comm.ConvertOxBytesToStr(p[0:3])
	// detAmount := comm.ConvertHexUnitToDecUnit(p[4])
	// TODO:匹配指令池内是否有相应发送项
	od := NewOrder(rcvID, "", CmdGetReceiverState, []string{})
	// 如果匹配
	if od.matchFromOrderPool() > -1 {
		// 记录到数据库

		//匹配成功，则注销记录
		od.UnregisterToOrdersPool()
	} else {
		// 不匹配，则错误返回
		return errors.New("未注册操作")
	}

	return nil
}

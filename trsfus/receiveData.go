package trsfus

import (
	"eInfusion/comm"
	"eInfusion/ndb"
	"errors"
	"net"
)

// DoReceiveData : 获取接收数据指令类型, cmdTypeCursor : 指令相关下标 (0 ~ N)
func DoReceiveData(p []byte, cmdTypeCursor int, c *net.TCPConn) error {
	if cmdTypeCursor >= 0 && cmdTypeCursor < len(p) {
		ct, ok := ReceiveCmdMap[p[cmdTypeCursor]]
		if ok {
			switch ct {
			case CmdGetReceiverState:
				if len(p) != 8 {
					l := len(p)
					return errors.New("数据长度超出定义！现长度：" + comm.ConvertBasNumberToStr(10, l))
				}
				return getUploadReceiverState(p[3:], c)
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
func getUploadReceiverState(p []byte, c *net.TCPConn) error {
	rcvID := comm.ConvertOxBytesToStr(p[0:4])
	detAmount := comm.ConvertHexUnitToDecUnit(p[4])
	od := NewOrder(rcvID, "", CmdGetReceiverState, []string{})
	// 如果匹配
	if od.matchFromOrderPool() > -1 {
		// 记录到数据库
		s := "INSERT t_rcv(receiver_id,detector_amount,last_time,ip_addr) VALUES(?,?,?,?);"
		if ndb.DBMain.ExceSQL(s, od.RcvID, detAmount, comm.GetCurrentTime(), comm.GetPureIPAddr(c)) == 0 {
			// 如果记录不成功
			comm.Msg("insert error")
		}
		//匹配成功，则注销记录
		od.UnregisterToOrdersPool()
	} else {
		// 不匹配，则错误返回
		return errors.New("未注册操作")
	}

	return nil
}

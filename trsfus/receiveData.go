package trsfus

import (
	"eInfusion/comm"
	"eInfusion/ndb"
	"eInfusion/tlogs"
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
				return getUpLoadReceiverState(p[3:], c)
			case CmdGetDetectorState:
				return getUpLoadDetectorState(p[3:], c)
			case CmdAddDetector:
				getUpLoadAddDetectorState(p[3:], c)
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
// 未考虑心跳包数据上传，有待后续改进
func getUpLoadReceiverState(p []byte, c *net.TCPConn) error {
	var err error
	rcvID := comm.ConvertOxBytesToStr(p[0:4])
	detAmount := comm.ConvertHexUnitToDecUnit(p[4])
	od := NewOrder(rcvID, "", CmdGetReceiverState, []string{})
	if od.matchFromOrderPool() > -1 {
		// 如果匹配成功
		var id string
		if ndb.DBMain.QueryOneData("SELECT receiver_id FROM t_rcv WHERE receiver_id=?", &id, od.RcvID) {
			if id != "" { // 如果接收器已经登记，则更新操作
				s := "UPDATE t_rcv SET detector_amount=?,last_time=?,ip_addr=?"
				if ndb.DBMain.ExceSQL(s, detAmount, comm.GetCurrentTime(), comm.GetPureIPAddr(c)) > 0 {
					tlogs.DoLog(tlogs.Info, "完成获取接收器:["+od.RcvID+"]信息")
					err = nil
				} else {
					err = errors.New("更新数据库失败，请查看错误日志！")
				}
			} else {
				// 没有登记，则新增检测器信息
				s := "INSERT t_rcv(receiver_id,detector_amount,last_time,ip_addr) VALUES(?,?,?,?);"
				if ndb.DBMain.ExceSQL(s, od.RcvID, detAmount, comm.GetCurrentTime(), comm.GetPureIPAddr(c)) > 0 {
					tlogs.DoLog(tlogs.Info, "完成获取接收器信息")
					err = nil
				} else {
					err = errors.New("更新数据库失败，请查看错误日志！")
				}
			}
		} else {
			err = errors.New("更新数据库失败，请查看错误日志！")
		}
		//匹配成功，则注销记录
		od.UnregisterToOrdersPool()
	} else {
		// 不匹配，则错误返回 -- FIXME: 心跳数据包没有考虑)
		err = errors.New("登记操作错误，平台未发出操作请求，视为无效信息！")
		// 记录日志
		tlogs.DoLog(tlogs.Warn, "接收器:["+rcvID+"] 发来未匹配操作信息")
	}
	return err
}

// getUpLoadDetectorState :获取检测器信息
// 报警会自动发送数据，因此不匹配发送指令池
func getUpLoadDetectorState(p []byte, c *net.TCPConn) error {
	// TODO: Testing
	var err error
	rcvID := comm.ConvertOxBytesToStr(p[0:4])
	// detAmount := comm.ConvertHexUnitToDecUnit(p[4])
	detID := comm.ConvertOxBytesToStr(p[5:9])
	dt := &Detector{}
	BinDetectorStat(p[9], dt)
	dt.ID = detID
	// od := NewOrder(rcvID, detID, CmdGetDetectorState, []string{})
	m := matchRcvToDet(rcvID, detID)
	if comm.CkErr("检测器与接收器匹配校验错误:", tlogs.Error, m) { // 校验验检测器与接收器配对
		return m
	}
	var id string
	ndb.DBMain.QueryOneData("SELECT did FROM t_running WHERE did=?", &id, detID)
	if id != "" {
		s := "UPDATE t_running SET time=?,capacity=?,alarm=?,error=? WHERE did=?"
		if ndb.DBMain.ExceSQL(s, comm.GetCurrentTime(), dt.Capacity, dt.Alarm, 0, detID) == 0 {
			err = errors.New("更新数据库失败，请查看错误日志！")
		} else {
			err = nil
		}
	} else {
		s := "INSERT INTO t_running(did,time,capacity,alarm) VALUES(?,?,?,?)"
		if ndb.DBMain.ExceSQL(s, detID, comm.GetCurrentTime(), dt.Capacity, dt.Alarm) == 0 {
			err = errors.New("更新数据库失败，请查看错误日志！")
		} else {
			err = nil
		}
	}
	return err
}

// getUpLoadAddDetectorState :获取添加/开启检测器结果
func getUpLoadAddDetectorState(p []byte, c *net.TCPConn) error {
	// var err error
	rcvID := comm.ConvertOxBytesToStr(p[0:4])
	// detAmount := comm.ConvertHexUnitToDecUnit(p[4])
	detID := comm.ConvertOxBytesToStr(p[5:9])

	od := NewOrder(rcvID, detID, CmdAddDetector, []string{})
	if od.matchFromOrderPool() > -1 { // 匹配成功
		m := matchRcvToDet(rcvID, detID)
		if comm.CkErr("检测器与接收器匹配校验错误:", tlogs.Error, m) { // 校验验检测器与接收器配对
			return m
		}
		var tas []*ndb.TransacionArgs
		ta := &ndb.TransacionArgs{}
		ta.SQL = "INSERT INTO t_running(did,time) VALUES(?,?)"
		ta.Args = append(ta.Args, detID)
		ta.Args = append(ta.Args, comm.GetCurrentTime())
		tas = append(tas, ta)
		ndb.DBMain.DoTransacion(tas)
	}
	return nil
}

// matchRcvToDet :匹配数据库内检测器与接收器是否对应
// 有错误表示不匹配
func matchRcvToDet(rcvID string, detID string) error {
	var id string
	s := "SELECT detID FROM t_rcv_vs_det WHERE rcvID=? and detID=?"
	if ndb.DBMain.QueryOneData(s, &id, rcvID, detID) {
		if id == "" {
			return errors.New("接收器:[" + rcvID + "]与检测器:[" + detID + "]不匹配！")
		}
	} else {
		return errors.New("获取设备记录表错误，请查看错误日志！")
	}
	return nil
}

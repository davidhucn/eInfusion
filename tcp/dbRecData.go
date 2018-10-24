package tcp

import (
	cm "eInfusion/comm"
	dh "eInfusion/datahub"
	db "eInfusion/tdb"
	logs "eInfusion/tlogs"
	tsc "eInfusion/trsfscomm"
)

//InitDetInfoToDB :初始化生成8个检测器信息到数据库-> t_device_dict
func InitDetInfoToDB(amount int) bool {
	var strSQL string
	var dd []tsc.Detector
	for i := 0; i < amount; i++ {
		var di tsc.Detector
		di.ID = "B000000" + cm.ConvertIntToStr(i)
		// di.Stat = ConvertIntToStr(2)
		di.QRCode = cm.CreateQRID(di.ID)
		dd = append(dd, di)
	}
	for i := 0; i < amount; i++ {
		strSQL = "Insert Into t_device_dict(detector_id,qcode) Values(?,?,?)"
		_, err := db.ExecSQL(strSQL, dd[i].ID, dd[i].QRCode)
		cm.CkErr(db.MsgDB.InsertDataErr, err)
	}
	return true
}

//ReceiveRcvStat :获取接收器状态
func ReceiveRcvStat(packData []byte, ipAddr string, rCmdType uint8) bool {
	var strSQL string
	var err error
	//数据包内数量位置
	intAmountCursor := 4
	//接收器ID
	strRcvID := cm.ConvertOxBytesToStr(packData[:4])
	//检测器数量
	strDetectAmount := cm.ConvertBasNumberToStr(10, packData[intAmountCursor])
	//查询用 接收器ID map
	var mRcvID *map[string]string
	strSQL = "SELECT receiver_id FROM t_rcv WHERE receiver_id=?"
	mRcvID, err = db.QueryOneRow(strSQL, strRcvID)
	if cm.CkErr(db.MsgDB.QueryDataErr, err) {
		return false
	}
	//是否已存在指定的接收器号
	if (*mRcvID)["receiver_id"] == "" {
		//如果没有插入数据
		strSQL = "Insert Into t_rcv(receiver_id,detector_amount,last_time,ip_addr) Values(?,?,?,?)"
		_, err = db.ExecSQL(strSQL, strRcvID, strDetectAmount, cm.GetCurrentTime(), ipAddr)
		if cm.CkErr(db.MsgDB.InsertDataErr, err) {
			return false
		}
	} else {
		//如果已有则更新
		strSQL = "UPDATE t_rcv SET detector_amount=?,last_time=?,ip_addr=? WHERE receiver_id=?"
		_, err = db.ExecSQL(strSQL, strDetectAmount, cm.GetCurrentTime(), strRcvID, ipAddr)
		if cm.CkErr(db.MsgDB.UpdateDataErr, err) {
			return false
		}
	}
	logs.LogMain.Info("接收器[", strRcvID, "]于", cm.GetCurrentDate(), "与服务平台完成通讯!")
	return true
}

//ReceiveDetectStat ：获取检测器状态信息
func ReceiveDetectStat(packData []byte, ipAddr string, rCmdType uint8) bool {
	var dDet []tsc.Detector
	var strSQL string
	var err error
	var mDetID *map[string]string
	//检测器数量所在位置
	intDetAmountCursor := 4
	//	接收器Id
	strRcvID := cm.ConvertOxBytesToStr(packData[0:4])
	//检测器数量
	intDetAmount := cm.ConvertBasStrToInt(10, cm.ConvertBasNumberToStr(10, packData[intDetAmountCursor]))
	//根据数量来存储
	if intDetAmount > 0 {
		begin := 5
		for i := 0; i < intDetAmount; i++ {
			//检测器ID起始位置
			end := begin + 4
			// 在循环内定义检测器对象
			var di tsc.Detector
			di.RcvID = strRcvID
			di.ID = cm.ConvertOxBytesToStr(packData[begin:end])
			// 获取检测器状态
			tsc.BinDetectorStat(packData[end], &di)
			begin = end
			//	判断该检测器是否为device_dict表内已注册设备，如果不是,退出
			// TODO:  目前只核查设备登记表（device_dict）,没有核对配对表（rcv_vs_det），后期考虑更改为配对表
			strSQL = "select did From t_device_dict Where did=?"
			mDetID, err = db.QueryOneRow(strSQL, di.ID)
			if !cm.CkErr(db.MsgDB.QueryDataErr, err) {
				if (*mDetID)["did"] != "" {
					dDet = append(dDet, di)
				} else {
					logs.LogMain.Warn(di.ID, "不是注册设备，来自ip：", ipAddr)
					return false
				}
			}
		}
		// 通过验证后，记录检测器信息(t_running表)
		for i := 0; i < len(dDet); i++ {
			// 判断运行状态表内是否存在
			strSQL = "Select did From t_running Where did=?"
			mDetID, err = db.QueryOneRow(strSQL, dDet[i].ID)
			if cm.CkErr(db.MsgDB.QueryDataErr, err) {
				return false
			}
			// 如果运行状态表内没有相关数据,则插入
			if (*mDetID)["did"] == "" {
				strSQL = `Insert Into t_running(did,time,capacity,alarm)
						 Values(?,?,?,?)`
				_, err = db.ExecSQL(strSQL, dDet[i].ID, cm.GetCurrentTime(), dDet[i].Capacity, dDet[i].Alarm)
				if cm.CkErr(db.MsgDB.InsertDataErr, err) {
					return false
				}
			} else {
				//如果已有则更新
				strSQL = `UPDATE t_running SET time=?,capacity=?,alarm=? WHERE did=?`
				_, err = db.ExecSQL(strSQL, cm.GetCurrentTime(), dDet[i].Capacity, dDet[i].Alarm, dDet[i].ID)
				if cm.CkErr(db.MsgDB.UpdateDataErr, err) {
					return false
				}
			}
			// 记录在日志内
			logs.LogMain.Info("检测器[", dDet[i], "]于", cm.GetCurrentDate(), "与服务平台完成通讯!")
		}
	} else {
		// 返回检测器为零,则退出并报错
		logs.LogMain.Error("错误 ！获取检测器数量为0")
		return false
	}
	return true
}

//ReceiveDeleteDetect :获取删除检测器结果信息
func ReceiveDeleteDetect(packData []byte, ipAddr string, rCmdType uint8) bool {
	var intDetAmount int
	var err error
	var dDet []tsc.Detector
	var strSQL string
	var mDetID *map[string]string
	//检测器数量所在位置
	intDetAmountCursor := 4
	//接收器id
	strRcvID := cm.ConvertOxBytesToStr(packData[:4])
	//检测器数量
	intDetAmount = cm.ConvertBasStrToInt(10, cm.ConvertBasNumberToStr(10, packData[intDetAmountCursor]))
	begin := 5
	for i := 0; i < intDetAmount; i++ {
		end := begin + 4
		var di tsc.Detector
		di.RcvID = strRcvID
		di.ID = cm.ConvertOxBytesToStr(packData[begin:end])
		begin = end
		//	判断该检测器是否为device_dict表内已注册设备，如果不是,则不记录
		strSQL = "Select did FROM t_device_dict Where did=?"
		mDetID, err = db.QueryOneRow(strSQL, di.ID)
		if cm.CkErr(db.MsgDB.QueryDataErr, err) {
			return false
		}
		if (*mDetID)["detector_id"] != "" {
			dDet = append(dDet, di)
		} else {
			logs.LogMain.Warn(di.ID, "不是注册设备，来自ip：", ipAddr)
		}
	}

	//册除：t_device_dict,t_rcv_vs_det,t_rcv
	for i := 0; i < len(dDet); i++ {
		strSQL = "DELETE FROM t_device_dict WHERE did=?"
		_, err = db.ExecSQL(strSQL, dDet[i].ID)
		if cm.CkErr(db.MsgDB.DeleteDataErr, err) {
			return false
		}
		strSQL = "DELETE FROM t_rcv_vs_det WHERE detID=?"
		_, err = db.ExecSQL(strSQL, dDet[i].ID)
		if cm.CkErr(db.MsgDB.DeleteDataErr, err) {
			return false
		}
		//获取现有的检测器数量
		strSQL = "Select detector_amount FROM t_rcv WHERE receiver_id=?"
		mDetID, err = db.QueryOneRow(strSQL, dDet[i].RcvID)
		if cm.CkErr(db.MsgDB.QueryDataErr, err) {
			return false
		}
		if (*mDetID)["detector_amount"] != "" {
			intDetAmount = cm.ConvertBasStrToInt(10, (*mDetID)["detector_amount"])
		} else {
			logs.LogMain.Error(db.MsgDB.QueryDataErr, err)
			return false
		}
		//更新接收器对的检测器数量
		if intDetAmount > 0 {
			intDetAmount = intDetAmount - 1
			strSQL = `UPDATE t_rcv SET detector_amount=?,last_time=?,ip_addr=? WHERE receiver_id=?`
			_, err = db.ExecSQL(strSQL, intDetAmount, cm.GetCurrentTime(), ipAddr, dDet[i].RcvID)
			if cm.CkErr(db.MsgDB.UpdateDataErr, err) {
				return false
			}
			//册除应用信息表
			// strSQL = "DELETE FROM t_apply_dict WHERE qcode=?"
			// _, err = ExecSQL(strSQL, dDet[i].QRCode)
			// if cm.CkErr(MsgDB.DeleteDataErr, err) {
			// 	return false
			// }
			logs.LogMain.Info("成功册除检测器[", dDet[i].ID, "]")
		} else {
			// 接收器附属检测器小于等于0时，提醒出错
			logs.LogMain.Info("册除指定检测器时出错，当前接收器下已无注册检测器！")
			return false
		}
	}
	return true
}

//ReceiveAddDetect :获取添加检测器结果信息
func ReceiveAddDetect(packData []byte, ipAddr string, rCmdType uint8) bool {
	var intDetAmount int
	var err error
	var dDet []tsc.Detector
	var strSQL string
	var mDetID *map[string]string
	//检测器数量所在位置
	intDetAmountCursor := 4
	//接收器id
	strRcvID := cm.ConvertOxBytesToStr(packData[:4])
	//检测器数量
	intDetAmount = cm.ConvertBasStrToInt(10, cm.ConvertBasNumberToStr(10, packData[intDetAmountCursor]))
	begin := 5
	for i := 0; i < intDetAmount; i++ {
		end := begin + 4
		var di tsc.Detector
		di.RcvID = strRcvID
		di.ID = cm.ConvertOxBytesToStr(packData[begin:end])
		// FIXME:检测器QR需要重做接口链接,等蒋少敏的后一步
		di.QRCode = cm.CreateQRID(di.ID)
		begin = end
		dDet = append(dDet, di)
	}
	//插入 t_device_dict,t_rcv_vs_det,t_rcv
	for i := 0; i < len(dDet); i++ {
		strSQL = "Insert Into t_device_dict(did,qcode) Values(?,?)"
		_, err = db.ExecSQL(strSQL, dDet[i].ID, dDet[i].QRCode)
		if cm.CkErr(db.MsgDB.InsertDataErr, err) {
			return false
		}
		strSQL = "Insert Into t_rcv_vs_det(detID,rcvID,time) Values(?,?,?)"
		_, err = db.ExecSQL(strSQL, dDet[i].ID, dDet[i].RcvID, cm.GetCurrentDate())
		if cm.CkErr(db.MsgDB.InsertDataErr, err) {
			return false
		}
		//获取现有的检测器数量
		strSQL = "Select detector_amount FROM t_rcv WHERE receiver_id=?"
		mDetID, err = db.QueryOneRow(strSQL, dDet[i].RcvID)
		if cm.CkErr(db.MsgDB.QueryDataErr, err) {
			return false
		}
		if (*mDetID)["detector_amount"] != "" {
			intDetAmount = cm.ConvertBasStrToInt(10, (*mDetID)["detector_amount"])
		} else {
			logs.LogMain.Error(db.MsgDB.QueryDataErr, err)
			return false
		}
		//更新接收器对应的检测器数量
		intDetAmount = intDetAmount + 1
		strSQL = `update t_rcv set detector_amount=? last_time=? ip_addr=? where receiver_id=?
				Values(?,?,?,?)`
		_, err = db.ExecSQL(strSQL, intDetAmount, cm.GetCurrentDate(), ipAddr, dDet[i].RcvID)
		if cm.CkErr(db.MsgDB.UpdateDataErr, err) {
			return false
		}
		logs.LogMain.Info("成功添加检测器[", dDet[i].ID, "]！")
		// 注销RequestOrderUnions内数据
		dh.UnregisterReqOrdersUnion(dDet[i].ID, rCmdType)
		// 获取web通讯ID
		wsOrderID := dh.GetReqOrderIDFromUnion(dDet[i].ID, rCmdType)
		od := cm.NewOrder(wsOrderID, []byte("开启检测器成功："+dDet[i].ID))
		// 回写到前端
		dh.SendMsgToWeb(od)
	}
	return true
}

//ReceiveSetRcvNetCfgStat ：获取设置网络配置操作结果信息
func ReceiveSetRcvNetCfgStat(packData []byte, ipAddr string, rCmdType uint8) bool {
	var strSQL string
	var err error
	//接收器ID
	strRcvID := cm.ConvertOxBytesToStr(packData[:4])
	//ServerIP地址
	var strServerIP string
	for i := 4; i <= 7; i++ {
		strServerIP += cm.ConvertBasNumberToStr(10, packData[i])
		if i < 7 {
			strServerIP += "."
		}
	}
	//	ServerPort
	intServerPort := cm.ConvertBasStrToInt(16, cm.ConvertOxBytesToStr(packData[8:10]))
	//查询用 接收器ID map
	var mRcvID *map[string]string
	strSQL = "SELECT receiver_id FROM t_rcv WHERE receiver_id=?"
	mRcvID, err = db.QueryOneRow(strSQL, strRcvID)
	if err != nil {
		logs.LogMain.Error(db.MsgDB.InsertDataErr, err)
		return false
	}
	//存在指定的接收器
	if (*mRcvID)["receiver_id"] != "" {
		//更新IP地址和端口设置
		strSQL = "UPDATE t_rcv SET last_time=?,ip_addr=?,target_ip=?,target_port=? WHERE receiver_id=?"
		_, err = db.ExecSQL(strSQL, cm.GetCurrentTime(), ipAddr, strServerIP, intServerPort, strRcvID)
		if err != nil {
			logs.LogMain.Error(db.MsgDB.UpdateDataErr, err)
			return false
		}
	} else {
		//	如果不存在指定的接收器，报错误
		logs.LogMain.Error("更新接收器网络配置出错，不存在指定接收器:", strRcvID)
		return false
	}
	logs.LogMain.Info("成功设置接收器[", strRcvID, "]网络配置！")
	return true
}

//ReceiveSetReconnTimeStat ：获取设置重新连接时间
func ReceiveSetReconnTimeStat(packData []byte, ipAddr string, rCmdType uint8) bool {
	var strSQL string
	var err error
	//接收器ID
	strRcvID := cm.ConvertOxBytesToStr(packData[:4])
	intReconnTime := cm.ConvertBasStrToInt(16, cm.ConvertOxBytesToStr(packData[4:6]))
	//查询用 接收器ID map
	var mRcvID *map[string]string
	strSQL = "SELECT receiver_id FROM t_rcv WHERE receiver_id=?"
	mRcvID, err = db.QueryOneRow(strSQL, strRcvID)
	if err != nil {
		logs.LogMain.Error(db.MsgDB.InsertDataErr, err)
		return false
	}
	//存在指定的接收器
	if (*mRcvID)["receiver_id"] != "" {
		//更新重新连接时间设置
		strSQL = "UPDATE t_rcv SET last_time=?,ip_addr=?,reconn_time=? WHERE receiver_id=?"
		_, err = db.ExecSQL(strSQL, cm.GetCurrentTime(), ipAddr, intReconnTime, strRcvID)
		if err != nil {
			logs.LogMain.Error(db.MsgDB.UpdateDataErr, err)
			return false
		}
	} else {
		//	如果不存在指定的接收器，报错误
		logs.LogMain.Error("更新接收器重连接配置出错，不存在指定接收器:", strRcvID)
		return false
	}
	logs.LogMain.Info("成功设置接受器[", strRcvID, "]重联机时间！")
	return true
}

//此包涉及具体业务的数据库操作
package dbWorks

import (
	cm "eInfusion/comm"
	. "eInfusion/dbOperate"
	"eInfusion/logs"
)

//初始化生成8个检测器信息到数据库-> t_device_dict
func InitDetInfoToDB(amount int) bool {
	var strSQL string
	var dd []Detector
	for i := 0; i < amount; i++ {
		var di Detector
		di.ID = "B000000" + cm.ConvertIntToStr(i)
		// di.Stat = ConvertIntToStr(2)
		di.QRCode = CreateQRID(di.ID)
		dd = append(dd, di)
	}
	for i := 0; i < amount; i++ {
		strSQL = "Insert Into t_device_dict(detector_id,qcode) Values(?,?,?)"
		_, err := ExecSQL(strSQL, dd[i].ID, dd[i].QRCode)
		_ = cm.CkErr(MsgDB.InsertDataErr, err)
	}
	return true
}

//获取接收器状态
func ReceiveRcvStat(packData []byte, ipAddr string) bool {
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
	strSQL = "SELECT receiver_id FROM t_receiver_dict WHERE receiver_id=?"
	mRcvID, err = QueryOneRow(strSQL, strRcvID)
	if err != nil {
		logs.LogMain.Error(MsgDB.InsertDataErr, err)
		return false
	}
	//是否已存在指定的接收器号
	if (*mRcvID)["receiver_id"] == "" {
		//如果没有插入数据
		strSQL = "Insert Into t_receiver_dict(receiver_id,detector_amount,last_time,ip_addr) Values(?,?,?,?)"
		_, err = ExecSQL(strSQL, strRcvID, strDetectAmount, cm.GetCurrentTime(), ipAddr)
		if err != nil {
			logs.LogMain.Error(MsgDB.InsertDataErr, err)
			return false
		}
	} else {
		//如果已有则更新
		strSQL = "UPDATE t_receiver_dict SET detector_amount=?,last_time=?,ip_addr=? WHERE receiver_id=?"
		_, err = ExecSQL(strSQL, strDetectAmount, cm.GetCurrentTime(), strRcvID, ipAddr)
		if err != nil {
			logs.LogMain.Error(MsgDB.UpdateDataErr, err)
			return false
		}
	}
	return true
}

//获取检测器状态信息
func ReceiveDetectStat(packData []byte, ipAddr string) bool {
	var dDet []Detector
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
		for i := 0; i < intDetAmount; i++ {
			//检测器ID起始位置
			begin := 5
			end := begin + 4
			// 在循环内定义检测器对象
			var di Detector
			di.RcvID = strRcvID
			di.ID = cm.ConvertOxBytesToStr(packData[begin:end])
			/*FIXME:接收器状态位，须针对位处理*/ /////////////////////////////////////
			// di.Stat = cm.ConvertBasNumberToStr(10, packData[end])
			cm.Msg("detID:", di.ID)
			cm.Msg("status:", packData[end])
			BinDetectorStat(packData[begin], &di)
			////////////////////////////////////////////////////////////////////////
			begin = end
			//	判断该检测器是否为device_dict表内已注册设备，如果不是,退出
			strSQL = "Select detector_id From t_device_dict Where detector_id=?"
			mDetID, err = QueryOneRow(strSQL, di.ID)
			if !cm.CkErr(MsgDB.QueryDataErr, err) {
				if (*mDetID)["detector_id"] != "" {
					dDet = append(dDet, di)
				} else {
					logs.LogMain.Warn(di.ID, "不是注册设备，来自ip：", ipAddr)
					return false
				}
			}
		}
	}
	//开始存入数据库（t_main）
	//	遍历整个slice
	for i := 0; i < len(dDet); i++ {
		/////////////////////确定 t_receiver_dict
		strSQL = "Select receiver_id From t_receiver_dict Where receiver_id=?"
		mDetID, err = QueryOneRow(strSQL, dDet[i].RcvID)
		if err != nil {
			logs.LogMain.Error(MsgDB.QueryDataErr, err)
			return false
		}
		if (*mDetID)["receiver_id"] == "" {
			strSQL = `Insert Into t_receiver_dict(receiver_id,detector_amount,last_time,ip_addr)
			 		Values(?,?,?,?)`
			_, err = ExecSQL(strSQL, dDet[i].RcvID, len(dDet), cm.GetCurrentTime(), ipAddr)
			if err != nil {
				logs.LogMain.Error(MsgDB.InsertDataErr, err)
				return false
			}
		} else {
			//如果已有则更新
			strSQL = `UPDATE t_receiver_dict SET receiver_id=?,detector_amount=?,last_time=?,
					ip_addr=? WHERE receiver_id=?`
			_, err = ExecSQL(strSQL, dDet[i].RcvID, len(dDet), cm.GetCurrentTime(), ipAddr, dDet[i].RcvID)
			if err != nil {
				logs.LogMain.Error(MsgDB.UpdateDataErr, err)
				return false
			}
		}
	}
	return true
}

//获取删除检测器结果信息
func ReceiveDeleteDetect(packData []byte, ipAddr string) bool {
	var intDetAmount int
	var err error
	var dDet []Detector
	var strSQL string
	var mDetID *map[string]string
	//检测器数量所在位置
	intDetAmountCursor := 4
	//接收器id
	strRcvID := cm.ConvertOxBytesToStr(packData[:4])
	//检测器数量
	intDetAmount = cm.ConvertBasStrToInt(10, cm.ConvertBasNumberToStr(10, packData[intDetAmountCursor]))

	for i := 0; i < intDetAmount; i++ {
		begin := 5
		end := begin + 4
		var di Detector
		di.RcvID = strRcvID
		di.ID = cm.ConvertOxBytesToStr(packData[begin:end])
		begin = end
		//	判断该检测器是否为device_dict表内已注册设备，如果不是,则不记录
		strSQL = "Select detector_id,qcode FROM t_device_dict Where detector_id=?"
		mDetID, err = QueryOneRow(strSQL, di.ID)
		if err != nil {
			logs.LogMain.Error(MsgDB.QueryDataErr, err)
			return false
		}
		if (*mDetID)["detector_id"] != "" {
			di.QRCode = (*mDetID)["qcode"]
			dDet = append(dDet, di)
		} else {
			logs.LogMain.Warn(di.ID, "不是注册设备，来自ip：", ipAddr)
		}
	}

	//册除：t_device_dict,t_match_dict,t_receiver_dict,t_apply_dict
	for i := 0; i < len(dDet); i++ {
		strSQL = "DELETE FROM t_device_dict WHERE detector_id=?"
		_, err = ExecSQL(strSQL, dDet[i].ID)
		if cm.CkErr(MsgDB.DeleteDataErr, err) {
			return false
		}
		strSQL = "DELETE FROM t_match_dict WHERE detector_id=?"
		_, err = ExecSQL(strSQL, dDet[i].ID)
		if cm.CkErr(MsgDB.DeleteDataErr, err) {
			return false
		}
		//获取现有的检测器数量
		strSQL = "Select detector_amount FROM t_receiver_dict WHERE receiver_id=?"
		mDetID, err = QueryOneRow(strSQL, dDet[i].RcvID)
		if cm.CkErr(MsgDB.QueryDataErr, err) {
			return false
		}
		if (*mDetID)["detector_amount"] != "" {
			intDetAmount = cm.ConvertBasStrToInt(10, (*mDetID)["detector_amount"])
		} else {
			logs.LogMain.Error(MsgDB.QueryDataErr, err)
			return false
		}
		//更新接收器对的检测器数量
		intDetAmount = intDetAmount - 1
		strSQL = `UPDATE t_receiver_dict SET detector_amount=?,last_time=?,ip_addr=?
		 		WHERE receiver_id=?`
		_, err = ExecSQL(strSQL, intDetAmount, cm.GetCurrentTime(), ipAddr, dDet[i].RcvID)
		if cm.CkErr(MsgDB.UpdateDataErr, err) {
			return false
		}
		//册除应用信息表
		strSQL = "DELETE FROM t_apply_dict WHERE qcode=?"
		_, err = ExecSQL(strSQL, dDet[i].QRCode)
		if cm.CkErr(MsgDB.DeleteDataErr, err) {
			return false
		}
	}
	return true
}

//获取添加检测器结果信息
func ReceiveAddDetect(packData []byte, ipAddr string) bool {
	var intDetAmount int
	var err error
	var dDet []Detector
	var strSQL string
	var mDetID *map[string]string
	//检测器数量所在位置
	intDetAmountCursor := 4
	//接收器id
	strRcvID := cm.ConvertOxBytesToStr(packData[:4])
	//检测器数量
	intDetAmount = cm.ConvertBasStrToInt(10, cm.ConvertBasNumberToStr(10, packData[intDetAmountCursor]))

	for i := 0; i < intDetAmount; i++ {
		begin := 5
		end := begin + 4
		var di Detector
		di.RcvID = strRcvID
		di.ID = cm.ConvertOxBytesToStr(packData[begin:end])
		di.QRCode = CreateQRID(di.ID)
		begin = end
		dDet = append(dDet, di)
	}
	//插入 t_device_dict,t_match_dict,t_receiver_dict
	for i := 0; i < len(dDet); i++ {
		strSQL = "Insert Into t_device_dict(detector_id,qcode) Values(?,?)"
		_, err = ExecSQL(strSQL, dDet[i].ID, dDet[i].QRCode)
		if cm.CkErr(MsgDB.InsertDataErr, err) {
			return false
		}
		strSQL = "Insert Into t_match_dict(detector_id,receiver_id) Values(?,?)"
		_, err = ExecSQL(strSQL, dDet[i].ID, dDet[i].RcvID)
		if cm.CkErr(MsgDB.InsertDataErr, err) {
			return false
		}
		//获取现有的检测器数量
		strSQL = "Select detector_amount FROM t_receiver_dict WHERE receiver_id=?"
		mDetID, err = QueryOneRow(strSQL, dDet[i].RcvID)
		if cm.CkErr(MsgDB.QueryDataErr, err) {
			return false
		}
		if (*mDetID)["detector_amount"] != "" {
			intDetAmount = cm.ConvertBasStrToInt(10, (*mDetID)["detector_amount"])
		} else {
			logs.LogMain.Error(MsgDB.QueryDataErr, err)
			return false
		}
		//更新接收器对应的检测器数量
		intDetAmount = intDetAmount + 1
		strSQL = `Insert Into t_receiver_dict(detector_id,detector_amount,last_time,ip_addr) 
				Values(?,?,?,?)`
		_, err = ExecSQL(strSQL, dDet[i].ID, intDetAmount, cm.GetCurrentTime(), ipAddr, dDet[i].RcvID)
		if cm.CkErr(MsgDB.InsertDataErr, err) {
			return false
		}
	}
	return true
}

//获取设置网络配置操作结果信息
func ReceiveSetRcvNetCfgStat(packData []byte, ipAddr string) bool {
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
	strSQL = "SELECT receiver_id FROM t_receiver_dict WHERE receiver_id=?"
	mRcvID, err = QueryOneRow(strSQL, strRcvID)
	if err != nil {
		logs.LogMain.Error(MsgDB.InsertDataErr, err)
		return false
	}
	//存在指定的接收器
	if (*mRcvID)["receiver_id"] != "" {
		//更新IP地址和端口设置
		strSQL = "UPDATE t_receiver_dict SET last_time=?,ip_addr=?,target_ip=?,target_port=? WHERE receiver_id=?"
		_, err = ExecSQL(strSQL, cm.GetCurrentTime(), ipAddr, strServerIP, intServerPort, strRcvID)
		if err != nil {
			logs.LogMain.Error(MsgDB.UpdateDataErr, err)
			return false
		}
	} else {
		//	如果不存在指定的接收器，报错误
		logs.LogMain.Error("更新接收器网络配置出错，不存在指定接收器:", strRcvID)
		return false
	}
	return true
}

//获取设置重新连接时间
func ReceiveSetReconnTimeStat(packData []byte, ipAddr string) bool {
	var strSQL string
	var err error
	//接收器ID
	strRcvID := cm.ConvertOxBytesToStr(packData[:4])
	intReconnTime := cm.ConvertBasStrToInt(16, cm.ConvertOxBytesToStr(packData[4:6]))
	//查询用 接收器ID map
	var mRcvID *map[string]string
	strSQL = "SELECT receiver_id FROM t_receiver_dict WHERE receiver_id=?"
	mRcvID, err = QueryOneRow(strSQL, strRcvID)
	if err != nil {
		logs.LogMain.Error(MsgDB.InsertDataErr, err)
		return false
	}
	//存在指定的接收器
	if (*mRcvID)["receiver_id"] != "" {
		//更新重新连接时间设置
		strSQL = "UPDATE t_receiver_dict SET last_time=?,ip_addr=?,reconn_time=? WHERE receiver_id=?"
		_, err = ExecSQL(strSQL, cm.GetCurrentTime(), ipAddr, intReconnTime, strRcvID)
		if err != nil {
			logs.LogMain.Error(MsgDB.UpdateDataErr, err)
			return false
		}
	} else {
		//	如果不存在指定的接收器，报错误
		logs.LogMain.Error("更新接收器重连接配置出错，不存在指定接收器:", strRcvID)
		return false
	}
	return true
}

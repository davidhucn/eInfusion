//此包涉及具体业务的数据库操作
package dbWorks

import (
	. "eInfusion/comm"
	. "eInfusion/dbOperate"
	"eInfusion/logs"
)

//初始化生成8个检测器信息到数据库-> t_device_dict
func InitDetInfoToDB(ref_amount int) bool {
	var strSql string
	var dd []Detector
	for i := 0; i < ref_amount; i++ {
		var di Detector
		di.ID = "B000000" + ConvertIntToStr(i)
		di.Disable = false
		di.Stat = ConvertIntToStr(2)
		di.QRCode = GetQRCodeStr(di.ID)
		dd = append(dd, di)
	}
	for i := 0; i < ref_amount; i++ {
		strSql = "Insert Into t_device_dict(detector_id,qcode,disable) Values(?,?,?)"
		_, err := ExecSQL(strSql, dd[i].ID, dd[i].QRCode, dd[i].Disable)
		if err != nil {
			Msg(C_Msg_DBInsert_Err, err)
		}
	}
	return true
}

//获取接收器状态
func ReceiveRcvStat(packData []byte, ipAddr string) bool {
	var strSql string
	var err error
	//数据包内数量位置
	var intAmountCursor int = 4
	//接收器ID
	strRcvID := ConvertOxBytesToStr(packData[:4])
	//检测器数量
	strDetectAmount := ConvertBasToStr(10, packData[intAmountCursor])
	//查询用 接收器ID map
	var mRcvId *map[string]string
	strSql = "SELECT receiver_id FROM t_receiver_dict WHERE receiver_id=?"
	mRcvId, err = QueryOneRow(strSql, strRcvID)
	if err != nil {
		logs.LogMain.Error(C_Msg_DBInsert_Err, err)
		return false
	}
	//是否已存在指定的接受器号
	if (*mRcvId)["receiver_id"] == "" {
		//如果没有插入数据
		strSql = "Insert Into t_receiver_dict(receiver_id,detector_amount,last_time,ip_addr) Values(?,?,?,?)"
		_, err = ExecSQL(strSql, strRcvID, strDetectAmount, GetCurrentTime(), ipAddr)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBInsert_Err, err)
			return false
		}
	} else {
		//如果已有则更新
		strSql = "UPDATE t_receiver_dict SET detector_amount=?,last_time=?,ip_addr=? WHERE receiver_id=?"
		_, err = ExecSQL(strSql, strDetectAmount, GetCurrentTime(), strRcvID, ipAddr)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBUpdate_Err, err)
			return false
		}
	}
	return true
}

//获取检测器信息
func ReceiveDetectStat(packData []byte, ipAddr string) bool {
	var dDet []Detector
	var strSql string
	var err error
	var mDetId *map[string]string
	//检测器数量所在位置
	var intDetAmountCursor int = 4
	//	接收器Id
	var strRcvID string = ConvertOxBytesToStr(packData[0:4])
	//检测器数量
	intDetAmount := ConvertBasStrToInt(10, ConvertBasToStr(10, packData[intDetAmountCursor]))
	//根据数量来存储
	if intDetAmount > 0 {
		for i := 0; i < intDetAmount; i++ {
			//检测器ID起始位置
			var begin int = 5
			var end int = begin + 4
			var di Detector
			di.RcvID = strRcvID
			di.ID = ConvertOxBytesToStr(packData[begin:end])
			di.Stat = ConvertBasToStr(10, packData[end])
			di.Disable = false
			begin = end
			//	判断该检测器是否为device_dict表内已注册设备，如果不是,则不记录
			strSql = "Select detector_id From t_device_dict Where detector_id=?"
			mDetId, err = QueryOneRow(strSql, di.ID)
			if err != nil {
				logs.LogMain.Error(C_Msg_DBQuery_Err, err)
				return false
			}
			if (*mDetId)["detector_id"] != "" {
				dDet = append(dDet, di)
			} else {
				logs.LogMain.Warn(di.ID, "不是注册设备，来自ip：", ipAddr)
			}
		}
	}
	//开始存入数据库（t_match_dict,t_device_dict,t_receiver_dict）
	//	遍历整个slice
	for i := 0; i < len(dDet); i++ {
		////////////////////确定device_dict表
		strSql = "Select detector_id From t_device_dict Where detector_id=?"
		mDetId, err = QueryOneRow(strSql, dDet[i].ID)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBQuery_Err, err)
			return false
		}
		if (*mDetId)["detector_id"] == "" {
			strSql = "Insert Into t_device_dict(detector_id,disable) Values(?,?)"
			_, err = ExecSQL(strSql, dDet[i].ID, dDet[i].Disable)
			if err != nil {
				logs.LogMain.Error(C_Msg_DBInsert_Err, err)
				return false
			}
		}
		/////////////////////确定 t_match_dict表
		strSql = "Select detector_id From t_match_dict Where detector_id=?"
		mDetId, err = QueryOneRow(strSql, dDet[i].ID)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBQuery_Err, err)
			return false
		}
		if (*mDetId)["detector_id"] == "" {
			strSql = "Insert Into t_match_dict(detector_id,receiver_id) Values(?,?)"
			_, err = ExecSQL(strSql, dDet[i].ID, dDet[i].RcvID)
			if err != nil {
				logs.LogMain.Error(C_Msg_DBInsert_Err, err)
				return false
			}
		} else {
			//如果已有则更新
			strSql = "UPDATE t_match_dict SET detector_id=?,receiver_id=? WHERE receiver_id=?"
			_, err = ExecSQL(strSql, dDet[i].ID, dDet[i].RcvID, dDet[i].RcvID)
			if err != nil {
				logs.LogMain.Error(C_Msg_DBUpdate_Err, err)
				return false
			}
		}
		/////////////////////确定 t_receiver_dict
		strSql = "Select receiver_id From t_receiver_dict Where receiver_id=?"
		mDetId, err = QueryOneRow(strSql, dDet[i].RcvID)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBQuery_Err, err)
			return false
		}
		if (*mDetId)["receiver_id"] == "" {
			strSql = `Insert Into t_receiver_dict(receiver_id,detector_amount,last_time,ip_addr)
			 		Values(?,?,?,?)`
			_, err = ExecSQL(strSql, dDet[i].RcvID, len(dDet), GetCurrentTime(), ipAddr)
			if err != nil {
				logs.LogMain.Error(C_Msg_DBInsert_Err, err)
				return false
			}
		} else {
			//如果已有则更新
			strSql = `UPDATE t_receiver_dict SET receiver_id=?,detector_amount=?,last_time=?,
					ip_addr=? WHERE receiver_id=?`
			_, err = ExecSQL(strSql, dDet[i].RcvID, len(dDet), GetCurrentTime(), ipAddr, dDet[i].RcvID)
			if err != nil {
				logs.LogMain.Error(C_Msg_DBUpdate_Err, err)
				return false
			}
		}
		//////////////////////////////////////////////////////////////////////////////
	}
	return true
}

//获取删除检测器结果信息
func ReceiveDeleteDetect(packData []byte, ipAddr string) bool {
	var intDetAmount int
	var err error
	var dDet []Detector
	var strSql string
	var mDetId *map[string]string
	//检测器数量所在位置
	var intDetAmountCursor int = 4
	//接收器id
	strRcvID := ConvertOxBytesToStr(packData[:4])
	//检测器数量
	intDetAmount = ConvertBasStrToInt(10, ConvertBasToStr(10, packData[intDetAmountCursor]))

	for i := 0; i < intDetAmount; i++ {
		var begin int = 5
		var end int = begin + 4
		var di Detector
		di.RcvID = strRcvID
		di.ID = ConvertOxBytesToStr(packData[begin:end])
		di.Disable = true
		begin = end
		//	判断该检测器是否为device_dict表内已注册设备，如果不是,则不记录
		strSql = "Select detector_id,qcode FROM t_device_dict Where detector_id=?"
		mDetId, err = QueryOneRow(strSql, di.ID)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBQuery_Err, err)
			return false
		}
		if (*mDetId)["detector_id"] != "" {
			di.QRCode = (*mDetId)["qcode"]
			dDet = append(dDet, di)
		} else {
			logs.LogMain.Warn(di.ID, "不是注册设备，来自ip：", ipAddr)
		}
	}

	//册除：t_device_dict,t_match_dict,t_receiver_dict,t_apply_dict
	for i := 0; i < len(dDet); i++ {
		strSql = "DELETE FROM t_device_dict WHERE detector_id=?"
		_, err = ExecSQL(strSql, dDet[i].ID)
		if CkErr(C_Msg_DBDelete_Err, err) {
			return false
		}
		strSql = "DELETE FROM t_match_dict WHERE detector_id=?"
		_, err = ExecSQL(strSql, dDet[i].ID)
		if CkErr(C_Msg_DBDelete_Err, err) {
			return false
		}
		//获取现有的检测器数量
		strSql = "Select detector_amount FROM t_receiver_dict WHERE receiver_id=?"
		mDetId, err = QueryOneRow(strSql, dDet[i].RcvID)
		if CkErr(C_Msg_DBQuery_Err, err) {
			return false
		}
		if (*mDetId)["detector_amount"] != "" {
			intDetAmount = ConvertBasStrToInt(10, (*mDetId)["detector_amount"])
		} else {
			logs.LogMain.Error(C_Msg_DBQuery_Err, err)
			return false
		}
		//更新接收器对的检测器数量
		intDetAmount = intDetAmount - 1
		strSql = `UPDATE t_receiver_dict SET detector_amount=?,last_time=?,ip_addr=?
		 		WHERE receiver_id=?`
		_, err = ExecSQL(strSql, intDetAmount, GetCurrentTime(), ipAddr, dDet[i].RcvID)
		if CkErr(C_Msg_DBDelete_Err, err) {
			return false
		}
		//册除应用信息表
		strSql = "DELETE FROM t_apply_dict WHERE qcode=?"
		_, err = ExecSQL(strSql, dDet[i].QRCode)
		if CkErr(C_Msg_DBDelete_Err, err) {
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
	var strSql string
	var mDetId *map[string]string
	//检测器数量所在位置
	var intDetAmountCursor int = 4
	//接收器id
	strRcvID := ConvertOxBytesToStr(packData[:4])
	//检测器数量
	intDetAmount = ConvertBasStrToInt(10, ConvertBasToStr(10, packData[intDetAmountCursor]))

	for i := 0; i < intDetAmount; i++ {
		var begin int = 5
		var end int = begin + 4
		var di Detector
		di.RcvID = strRcvID
		di.ID = ConvertOxBytesToStr(packData[begin:end])
		di.QRCode = GetQRCodeStr(di.ID)
		di.Disable = false
		begin = end
		dDet = append(dDet, di)
	}
	//插入 t_device_dict,t_match_dict,t_receiver_dict
	for i := 0; i < len(dDet); i++ {
		strSql = "Insert Into t_device_dict(detector_id,qcode,disable) Values(?,?,?)"
		_, err = ExecSQL(strSql, dDet[i].ID, dDet[i].QRCode, dDet[i].Disable)
		if CkErr(C_Msg_DBInsert_Err, err) {
			return false
		}
		strSql = "Insert Into t_match_dict(detector_id,receiver_id) Values(?,?)"
		_, err = ExecSQL(strSql, dDet[i].ID, dDet[i].RcvID)
		if CkErr(C_Msg_DBInsert_Err, err) {
			return false
		}
		//获取现有的检测器数量
		strSql = "Select detector_amount FROM t_receiver_dict WHERE receiver_id=?"
		mDetId, err = QueryOneRow(strSql, dDet[i].RcvID)
		if CkErr(C_Msg_DBQuery_Err, err) {
			return false
		}
		if (*mDetId)["detector_amount"] != "" {
			intDetAmount = ConvertBasStrToInt(10, (*mDetId)["detector_amount"])
		} else {
			logs.LogMain.Error(C_Msg_DBQuery_Err, err)
			return false
		}
		//更新接收器对应的检测器数量
		intDetAmount = intDetAmount + 1
		strSql = `Insert Into t_receiver_dict(detector_id,detector_amount,last_time,ip_addr) 
				Values(?,?,?,?)`
		_, err = ExecSQL(strSql, dDet[i].ID, intDetAmount, GetCurrentTime(), ipAddr, dDet[i].RcvID)
		if CkErr(C_Msg_DBInsert_Err, err) {
			return false
		}
	}
	return true
}

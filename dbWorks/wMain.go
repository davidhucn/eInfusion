//此包涉及具体业务的数据库操作
package dbWorks

import (
	. "eInfusion/comm"
	. "eInfusion/dbOperate"
	"eInfusion/logs"
)

//获取接收器状态
func GetRcvStat(packData []byte) bool {
	var strSql string
	var err error
	//	var afNum int64
	//接收器ID
	strRcvID := ConvertOxBytesToStr(packData[:4])
	//检测器数量
	strDetectAmount := ConvertBasToStr(10, packData[4])
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
		strSql = "Insert Into t_receiver_dict(receiver_id,detector_amount,last_time) Values(?,?,?)"
		_, err = InsertDataFast(strSql, strRcvID, strDetectAmount, GetCurrentTime())
		if err != nil {
			logs.LogMain.Error(C_Msg_DBInsert_Err, err)
			return false
		}
	} else {
		//如果存在则更新
		strSql = "UPDATE t_receiver_dict SET detector_amount=?,last_time=? WHERE receiver_id=?"
		_, err = UpateDataFast(strSql, strDetectAmount, GetCurrentTime(), strRcvID)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBInsert_Err, err)
			return false
		}
	}
	return true
}

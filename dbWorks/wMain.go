//此包涉及具体业务的数据库操作
package dbWorks

import (
	"eInfusion/comm"
	. "eInfusion/dbOperate"
	"eInfusion/logs"
)

//获取接收器状态
func GetRcvStat(packData []byte) {
	var strSql string
	var err error
	var afNum int64
	strRcvID := comm.ConvertOxBytesToStr(packData[:4])
	strDetectAmount := comm.ConvertBasToStr(10, packData[4])
	//	var mRcvId map[string]string
	mRcvId := make(map[string]string)
	strSql = "SELECT Receiver_id FROM t_receiver_dict WHERE receiver_id=?"
	&mRcvId, err = QueryOneRow(strSql, strDetectAmount)
	if err != nil {
		logs.LogMain.Error(C_Msg_DBInsert_Err, err)
		return
	}
	comm.Msg(mRcvId[receiver_id])
	strSql = "Insert Into t_receiver_dict(receiver_id,detector_amount,last_time) Values(?,?,?)"
	afNum, err = InsertDataFast(strSql, strRcvID, strDetectAmount, comm.GetCurrentTime())
	if err != nil {
		logs.LogMain.Error(C_Msg_DBInsert_Err, err)
		return
	}
	if afNum > 0 {
		comm.Msg("get status finished:", afNum)
	}

}

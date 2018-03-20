//此包涉及具体业务的数据库操作
package dbWorks

import (
	"eInfusion/comm"
	. "eInfusion/dbOperate"
	"eInfusion/logs"
)

//获取接收器状态
func GetRcvStat(packData []byte) {

	strRcvID := comm.ConvertOxBytesToStr(packData[:4])
	strSql := "SELECT "
	strSql := "Insert Into t_receiver_dict(receiver_id,detector_amount,last_time) Values(?,?,?)"
	afNum, err := InsertDataFast(strSql, strRcvID, comm.ConvertBasToStr(10, packData[4]), comm.GetCurrentTime())
	if err != nil {
		logs.LogMain.Error(C_Msg_DBInsert_Err, err)
		return
	}
	if afNum > 0 {
		comm.Msg("get status finished:", afNum)
	}

}

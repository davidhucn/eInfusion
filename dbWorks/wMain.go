//此包涉及具体业务的数据库操作
package dbWorks

import (
	. "eInfusion/dbOperate"
	"eInfusion/logs"
)

//获取接收器状态
func GetRcvStat(packData []byte) {

	strSql := "Insert Into t_receiver_dict(receiver_id,detector_amount) Values(?,?)"
	_, err := InsertData(strSql, packData[0], packData[1])

	if err != nil {
		logs.LogMain.Error(C_Msg_DBInsert_Err, err)
		return
	}

}

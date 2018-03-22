//此包涉及具体业务的数据库操作
package dbWorks

import (
	. "eInfusion/comm"
	. "eInfusion/dbOperate"
	"eInfusion/logs"
)

//获取接收器状态
func GetRcvStat(packData []byte, ipAddr string) bool {
	var strSql string
	var err error
	//专用数据内容位置
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
		_, err = InsertDataFast(strSql, strRcvID, strDetectAmount, GetCurrentTime(), ipAddr)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBInsert_Err, err)
			return false
		}
	} else {
		//如果存在则更新
		strSql = "UPDATE t_receiver_dict SET detector_amount=?,last_time=?,ip_addr WHERE receiver_id=?"
		_, err = UpateDataFast(strSql, strDetectAmount, GetCurrentTime(), strRcvID, ipAddr)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBInsert_Err, err)
			return false
		}
	}
	return true
}

//初始化生成8个检测器信息到数据库->t_device_dict
func InitDetInfoToDB() {

}

//获取检测器信息
func GetDetectStat(packData []byte) bool {
	var strSql string
	var err error
	//专用数据内容位置
	var intAmountCursor int = 4
	strDetID := ConvertOxBytesToStr(packData[:4])
	//检测器数量
	intDetAmount := ConvertBasStrToInt(10, ConvertBasToStr(10, packData[intAmountCursor]))
	strDetStat := ConvertBasToStr(10, packData[intAmountCursor])
	//如果检测器为1个
	if intDetAmount == 1 {

	} else {
		//	如果数量大于1个，则需要迭代所有检测器信息
	}
}

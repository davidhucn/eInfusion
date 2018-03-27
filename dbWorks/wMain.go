//此包涉及具体业务的数据库操作
package dbWorks

import (
	. "eInfusion/comm"
	. "eInfusion/dbOperate"
	"eInfusion/logs"
)

//检测器对象
type Detector struct {
	ID         string
	ReceiverID string
	Stat       string
	Disable    bool
}

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
		_, err = ExecSQL(strSql, strRcvID, strDetectAmount, GetCurrentTime(), ipAddr)
		if err != nil {
			logs.LogMain.Error(C_Msg_DBInsert_Err, err)
			return false
		}
	} else {
		//如果已有则更新
		strSql = "UPDATE t_receiver_dict SET detector_amount=?,last_time=?,ip_addr WHERE receiver_id=?"
		_, err = ExecSQL(strSql, strDetectAmount, GetCurrentTime(), strRcvID, ipAddr)
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
func GetDetectStat(packData []byte, ipAddr string) bool {
	var s *Student = new(Detector)
	//	var strSql string
	//	检测器编号数组
	//	var mDetID map[int]string
	mDetID := make(map[string]string)
	//检测器数量所在位置
	var intDetAmountCursor int = 4
	//	strRcvID := ConvertOxBytesToStr(packData[0:4])
	//检测器数量
	intDetAmount := ConvertBasStrToInt(10, ConvertBasToStr(10, packData[intDetAmountCursor]))
	mDetID["Det"] = ConvertOxBytesToStr(packData[5:9])
	//如果检测器为1个
	if intDetAmount > 1 {
		//		for i:=

	}
	Msg("interal data lens:", len(packData))
	return true
}

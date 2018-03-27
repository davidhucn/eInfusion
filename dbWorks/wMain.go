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
	Stat       string //十进制表示
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

	//如果检测器超过1个
	if intDetAmount > 0 {
		for i := 0; i < intDetAmount; i++ {
			//检测器ID起始位置
			var begin int = 5
			var end int = begin + 4
			dDet[i].ReceiverID = strRcvID
			dDet[i].ID = ConvertOxBytesToStr(packData[begin:end])
			dDet[i].Stat = ConvertBasToStr(10, packData[end])
			dDet[i].Disable = false
			begin = end
		}
	}
	//开始存入数据库（t_match_dict,t_device_dict,t_receiver_dict）
	//	确定device_dict表里是否存在，如果存在则更新
	strSql = "Select detector_id From t_device_dict Where detector_id=?"
	mDetId, err = QueryOneRow(strSql, dDet[i])
	if err != nil {
		logs.LogMain.Error(C_Msg_DBInsert_Err, err)
		return false
	}
	//是否已存在指定的接受器号
	if (*mRcvId)["detector_id"] == "" {
		strSql = "Insert Into t_device_dict() "
		return true
	}
}

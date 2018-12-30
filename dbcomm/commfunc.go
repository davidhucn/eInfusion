package dbcomm

import (
	cm "eInfusion/comm"
	db "eInfusion/tdb"
	"eInfusion/tlogs"
)

// GetRcvID :根据DetID获取RcvID
// 表 ：t_rcv_vs_det
func GetRcvID(rDetID string) string {
	var strSQL string
	var mRcvID *map[string]string
	var err error
	strSQL = "SELECT rcvID FROM t_rcv_vs_det WHERE detID=?"
	mRcvID, err = db.QueryOneRow(strSQL, rDetID)
	if cm.CkErr(db.MsgDB.QueryDataErr, tlogs.Error, err) {
		return ""
	}
	return (*mRcvID)["rcvID"]
}

// GetRcvIP :根据RcvID获取其IP地址
func GetRcvIP(rRcvID string) string {
	var strSQL string
	var mRcvIP *map[string]string
	var err error
	strSQL = "SELECT ip_addr FROM t_rcv WHERE receiver_id=?"
	mRcvIP, err = db.QueryOneRow(strSQL, rRcvID)
	if cm.CkErr(db.MsgDB.QueryDataErr, tlogs.Error, err) {
		return ""
	}
	return (*mRcvIP)["ip_addr"]
}

// IsReceiver  :获取设备类型
func IsReceiver(rTargetID string) bool {
	strSQL := "SELECT receiver_id FROM t_rcv WHERE receiver_id=?"
	var mRcvIP *map[string]string
	var err error
	mRcvIP, err = db.QueryOneRow(strSQL, rTargetID)
	if cm.CkErr(db.MsgDB.QueryDataErr, tlogs.Error, err) {
		return false
	}
	if (*mRcvIP)["receiver_id"] == "" {
		return false
	}
	return true
}

// IsDetector  :获取设备类型
func IsDetector(rTargetID string) bool {
	strSQL := "SELECT detID FROM t_rcv_vs_det WHERE detID=?"
	var mRcvIP *map[string]string
	var err error
	mRcvIP, err = db.QueryOneRow(strSQL, rTargetID)
	if cm.CkErr(db.MsgDB.QueryDataErr, tlogs.Error, err) {
		return false
	}
	if (*mRcvIP)["detID"] == "" {
		return false
	}
	return true
}

//此包涉及具体业务的数据库操作
package dbWorks

import (
	"eInfusion/comm"
	. "eInfusion/dbOperate"
	"eInfusion/logs"
)

const (
	c_DB_IPAddr  = "127.0.0.1"
	c_DB_Port    = "3306"
	c_DB_schema  = "transfusion"
	c_DB_UsrName = "root"
	c_DB_Pwd     = "2341656"
)

var G_DB DBConn

func Init() {
	comm.ShowScreen("welcome to the dbworking")
	G_DB.UserName = c_DB_UsrName
	G_DB.Password = c_DB_Pwd
	G_DB.Schema = c_DB_schema
	G_DB.Port = c_DB_Port
	G_DB.IpAddr = c_DB_IPAddr
	err := G_DB.ConnectDB()
	if err != nil {
		logs.LogMain.Critical(C_Msg_DBConnect_Err)
		panic(err)
	}
}

//获取接收器状态
func GetRcvStat(packData []byte) {
	strSql := "INSERT INTO t_receiver_dict(receiver_id,detector_amount) VALUES(?,?)"
	intNum, err := G_DB.InsertData(strSql, packData[0:3], packData[4])
	if err != nil {
		comm.ShowScreen("insert receiver_dict error")
	}
	comm.ShowScreen("affect num:", intNum)
}

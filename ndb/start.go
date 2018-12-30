package ndb

import (
	"eInfusion/comm"
)

// Svr ：全局数据库对象
var Svr *DBx

// InitDB ：初始化数据库
func InitDB() {
	ps := &dbParam{}
	ps.HostNameOrIPAddr = "127.0.0.1"
	ps.Password = "2341656"
	ps.Port = "3306"
	ps.UserName = "root"
	ps.schemaName = "transfusion"
	Svr = NewDBx(ps, DBType.MySQL)
	//////////////////////////////////////
	if !Svr.Connect() {
		comm.Msg("can't connect")
	}
	// type device struct {
	// 	qcode   string
	// 	did     string
	// 	remark  string
	// 	disable int
	// }
	// d := &device{}
	// comm.Msg(Svr.db)
	s := "select * from t_device_dict"
	r := Svr.QueryOneData(s, "B0000000")
	comm.Msg("result:", r)
	/////////////////////////////////////
}

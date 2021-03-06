package ndb

import (
	"eInfusion/tlogs"
)

// DBMain ：全局数据库对象
var DBMain *DBx

// DatabaseInit ：初始化数据库
func DatabaseInit() {

	localDB := NewDBparams(DataBaseType.MySQL, "root", "2341656", "localhost", "3306", "transfusion")
	// sqlite3 := NewDBparams(DataBaseType.Sqlite3, "", "", "./", "", "us")
	DBMain = NewDBx(localDB)

	if DBMain.Connect() {
		tlogs.DoLog(tlogs.Info, DBMsg.DatabaseInitFinish)
	} else {

	}
	// schema := `CREATE TABLE main (
	// 		country text,
	// 		city text NULL,
	// 		telcode integer);`
	// Svr.ExceSQL(schema)

	// var sd sql.NullString
	// Svr.QueryOneData("select disable from t_device_dict limit 1", &sd)
	// a := Svr.ExceSQL("insert into t_rcv_vs_det(detID,rcvID,remark,time) Value(?,?,?,?);", "D0000000", "A0000000", "test", comm.GetCurrentTime())
	// s := "update t_rcv_vs_det set remark=? where detid='D0000000'"
	// a := Svr.ExceSQL(s, "donig")
	// comm.Msg(a)

	// nd, err := sqlx.Connect(DBType.Sqlite3, "./testdb")
	// if err != nil {
	// 	comm.Msg("disconnect")
	// } else {
	// 	var str string
	// 	nd.Get(&str, "select name from main limit 1")
	// 	comm.Msg(str)
	// }

	// 	d, err := sqlx.Connect(DataBaseType.Sqlite3, "./tt")
	// 	comm.Msg(DataBaseType.Sqlite3)
	// 	if err != nil {
	// 		comm.Msg("err:", err)
	// 	} else {
	// 		schema := `CREATE TABLE main (
	// 			country text,
	// 			city text NULL,
	// 			telcode integer);`

	// 		// execute a query on the server
	// 		_, err := d.Exec(schema)
	// 		if err != nil {
	// 			comm.Msg(err)
	// 		} else {

	// 		}
	// 	}
	// }

	//////////////////////////////////////
	// if !Svr.Connect() {

	// }
	// type device struct {
	// 	Qcode   string         `db:"qcode"`
	// 	Did     string         `db:"did"`
	// 	Remark  sql.NullString `db:"remark"`
	// 	Disable int            `db:"disable"`
	// }

	// var d device

	// d := make([]device, 0)
	// comm.Msg(Svr.db)
	// s := "select * from t_device_dict where did=?"

	// Svr.QueryOneData(s, &d, "B0000000")
	// comm.Msg("result:", d)

	/////////////////////////////////////
}

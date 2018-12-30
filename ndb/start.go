package ndb

// Svr ：全局数据库对象
var Svr *DBx

// InitDB ：初始化数据库
func init() {
	ps := NewDBparams("root", "2341656", "localhost", "3306", "transfusion")
	Svr = NewDBx(ps, DBType.MySQL)
	Svr.Connect()
}

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

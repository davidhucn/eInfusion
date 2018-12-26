package ndb

type databaseTypeForDriver struct {
	Oracle  string
	MSSql   string
	Sqlite3 string
	MySQL   string
}

// DBType :数据库类型(驱动名称)
var DBType databaseTypeForDriver

func init() {
	DBType.MySQL = "mysql"
	DBType.Sqlite3 = "sqlite3"
	DBType.Oracle = "oracle"
	DBType.MSSql = "mssql"

	DBMsg.ConnectDBErr = "错误,无法连接到指定数据库！"
	DBMsg.InsertDataErr = "错误,插入数据库操作失败！"
	DBMsg.DeleteDataErr = "错误,册除数据操作失败！"
	DBMsg.QueryDataErr = "错误,查询数据信息失败！"
	DBMsg.UpdateDataErr = "错误,更数据信息失败！"
}

type dbMsg struct {
	ConnectDBErr  string
	InsertDataErr string
	DeleteDataErr string
	QueryDataErr  string
	UpdateDataErr string
}

// DBMsg :数据库消息对象
var DBMsg dbMsg

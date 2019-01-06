package ndb

type databaseTypeForDriver struct {
	Oracle     string
	MSSql      string
	Sqlite3    string
	MySQL      string
	PostgreSQL string
}

// DataBaseType :数据库类型(驱动名称)
var DataBaseType databaseTypeForDriver

func init() {
	DataBaseType.MySQL = "mysql"
	DataBaseType.Sqlite3 = "sqlite3"
	DataBaseType.Oracle = "oracle"
	DataBaseType.MSSql = "mssql"
	DataBaseType.PostgreSQL = "postgres"

	DBMsg.ConnectDBErr = "错误,无法连接到指定数据库！"
	DBMsg.InsertDataErr = "错误,插入数据库操作失败！"
	DBMsg.DeleteDataErr = "错误,册除数据操作失败！"
	DBMsg.QueryDataErr = "错误,查询数据信息失败！"
	DBMsg.UpdateDataErr = "错误,更数据信息失败！"
	DBMsg.DatabaseInitFinish = "提示，主数据库管理模块初始化完成！"
	DBMsg.EnableTransacionFailure = "错误，启用事务失败！"
	DBMsg.TransacionOperateErr = "错误，事务操作失败！"
}

type dbMsg struct {
	ConnectDBErr            string
	InsertDataErr           string
	DeleteDataErr           string
	QueryDataErr            string
	UpdateDataErr           string
	DatabaseInitFinish      string
	TransacionOperateErr    string
	EnableTransacionFailure string
}

// DBMsg :数据库消息对象
var DBMsg dbMsg

// TransacionArgs :事务操作参数
type TransacionArgs struct {
	SQL    string
	Params []interface{}
}

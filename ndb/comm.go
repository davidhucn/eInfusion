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
}

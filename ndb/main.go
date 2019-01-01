package ndb

import (
	cm "eInfusion/comm"
	"eInfusion/tlogs"

	"github.com/jmoiron/sqlx"
	// d "golang.org/x/net/http2/h2demo"

	// mysql 接口
	_ "github.com/go-sql-driver/mysql"
	// sqlite接口
	_ "github.com/mattn/go-sqlite3"
	// oracle接口
	// _ "github.com/mattn/go-oci8"
	// _ "github.com/lib/pq"
)

// var db *sqlx.DB

// DBParam :数据库连接参数
type DBParam struct {
	UserName     string
	Password     string
	Host         string
	Port         string
	SchemaName   string
	DataBaseType string
}

// NewDBparams :新建数据库连接参数(mysql,mssql,oracle)
func NewDBparams(databaseType string, UserName string, Password string, Host string, Port string, SchemaName string) (*DBParam, string) {
	p := &DBParam{
		UserName:     UserName,
		Password:     Password,
		Host:         Host,
		Port:         Port,
		SchemaName:   SchemaName,
		DataBaseType: databaseType,
	}
	var connectStr string
	switch databaseType {
	case DataBaseType.MySQL:
		connectStr = p.UserName + ":" + p.Password + "@tcp(" + p.Host + ":"
		connectStr = connectStr + p.Port + ")/" + p.SchemaName
		connectStr = connectStr + "?charset=utf8" //字符集
	case DataBaseType.Sqlite3:
		connectStr = p.Host
	case DataBaseType.Oracle:
	}
	return p, connectStr
}

// DBx :数据库对象
type DBx struct {
	db            *sqlx.DB
	dbParams      *DBParam
	connectString string
}

// NewDBx :新建数据库对象
func NewDBx(dbParams *DBParam, connectString string) *DBx {
	d := &DBx{
		dbParams:      dbParams,
		connectString: connectString,
	}
	return d
}

// Connect :连接到数据库
func (d *DBx) Connect() bool {
	var err error
	d.db, err = sqlx.Connect(d.dbParams.DataBaseType, d.connectString)
	if cm.CkErr(DBMsg.ConnectDBErr, tlogs.Error, err) {
		// 如果连接失败
		return false
	}
	return true
}

// isConnected :判断是否已连接
func (d *DBx) isConnected() bool {
	if d.db.Ping() == nil {
		return true
	}
	return false
}

// ExceSQL :执行查询语句
func (d *DBx) ExceSQL(s string, args ...interface{}) int64 {
	rs, err := d.db.Exec(s, args...)
	if cm.CkErr(DBMsg.QueryDataErr, tlogs.Error, err) {
		return 0
	}
	inAmount, err := rs.LastInsertId()
	if cm.CkErr(DBMsg.UpdateDataErr, tlogs.Error, err) {
		return 0
	}
	affAmount, err := rs.RowsAffected()
	if cm.CkErr(DBMsg.UpdateDataErr, tlogs.Error, err) {
		return 0
	}
	if affAmount > 0 {
		return affAmount
	}
	if inAmount > 0 {
		return inAmount
	}
	return 0
}

// QueryOneData :查询第一条数据，返回相关struct
func (d *DBx) QueryOneData(sql string, result interface{}, args ...interface{}) bool {
	var err error
	err = d.db.Get(result, sql, args...)
	if cm.CkErr(DBMsg.QueryDataErr, tlogs.Error, err) {
		return false
	}
	return true
}

// QueryDatas :查询批量数据,返回相关[]struct
func (d *DBx) QueryDatas(sql string, results interface{}, args ...interface{}) bool {
	var err error
	err = d.db.Select(&results, sql, args...)
	if cm.CkErr(DBMsg.QueryDataErr, tlogs.Error, err) {
		return false
	}
	return true
}

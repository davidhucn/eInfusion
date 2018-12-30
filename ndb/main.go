package ndb

import (
	cm "eInfusion/comm"
	"eInfusion/tlogs"

	"github.com/jmoiron/sqlx"

	// mysql 接口
	_ "github.com/go-sql-driver/mysql"
)

// var db *sqlx.DB

// 数据库连接参数
type dbParam struct {
	UserName         string
	Password         string
	HostNameOrIPAddr string
	Port             string
	schemaName       string
}

// DBx :数据库对象
type DBx struct {
	db           *sqlx.DB
	params       *dbParam
	databaseType string
}

// NewDBx :新建数据库对象
func NewDBx(param *dbParam, databaseType string) *DBx {
	return &DBx{
		params:       param,
		databaseType: databaseType,
	}
}

// Connect :连接到数据库
func (d *DBx) Connect() bool {
	var err error
	dbsource := d.params.UserName + ":" + d.params.Password + "@tcp(" + d.params.HostNameOrIPAddr + ":"
	dbsource = dbsource + d.params.Port + ")/" + d.params.schemaName
	dbsource = dbsource + "?charset=utf8" //字符集
	d.db, err = sqlx.Connect(d.databaseType, dbsource)
	if cm.CkErr(DBMsg.ConnectDBErr, tlogs.Info, err) {
		// 如果连接失败
		return false
	}
	return true
}

// isConnected :判断是否已连接
func (d *DBx) isConnected() bool {
	if d.db.Stats().OpenConnections > 0 {
		return true
	}
	return false
}

// ExceSQL :执行查询语句
func (d *DBx) ExceSQL(s string, args ...interface{}) bool {
	_, err := d.db.Exec(s, args)
	if cm.CkErr(DBMsg.QueryDataErr, tlogs.Error, err) {
		return false
	}
	// cm.Msg("影响的数据量：", r.RowsAffected())
	return true
}

// QueryOneData :查询第一条数据，返回相关struct
func (d *DBx) QueryOneData(sql string, args ...interface{}) interface{} {
	var rs interface{}
	cm.Msg("sql:", sql)
	if !cm.CkErr(DBMsg.QueryDataErr, tlogs.Error, d.db.Get(rs, sql, args)) {
		return rs
	}
	return nil
}

// QueryDatas :查询批量数据,返回相关[]struct
func (d *DBx) QueryDatas(sql string, args ...interface{}) interface{} {
	var rs interface{}
	if !cm.CkErr(DBMsg.QueryDataErr, tlogs.Error, d.db.Select(rs, sql, args)) {
		return rs
	}
	return nil
}

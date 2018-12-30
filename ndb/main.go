package ndb

import (
	cm "eInfusion/comm"
	"eInfusion/tlogs"

	"github.com/jmoiron/sqlx"

	// mysql 接口
	_ "github.com/go-sql-driver/mysql"
)

// var db *sqlx.DB

// DBParam :数据库连接参数
type DBParam struct {
	UserName         string
	Password         string
	HostNameOrIPAddr string
	Port             string
	SchemaName       string
}

// NewDBparams :新建数据库连接参数
func NewDBparams(UserName string, Password string, HostNameOrIPAddr string, Port string, SchemaName string) *DBParam {
	return &DBParam{
		UserName:         UserName,
		Password:         Password,
		HostNameOrIPAddr: HostNameOrIPAddr,
		Port:             Port,
		SchemaName:       SchemaName,
	}
}

// DBx :数据库对象
type DBx struct {
	db           *sqlx.DB
	params       *DBParam
	databaseType string
}

// NewDBx :新建数据库对象
func NewDBx(param *DBParam, databaseType string) *DBx {
	return &DBx{
		params:       param,
		databaseType: databaseType,
	}
}

// Connect :连接到数据库
func (d *DBx) Connect() bool {
	if d.db.Ping() == nil {
		var err error
		dbsource := d.params.UserName + ":" + d.params.Password + "@tcp(" + d.params.HostNameOrIPAddr + ":"
		dbsource = dbsource + d.params.Port + ")/" + d.params.SchemaName
		dbsource = dbsource + "?charset=utf8" //字符集
		d.db, err = sqlx.Connect(d.databaseType, dbsource)
		if cm.CkErr(DBMsg.ConnectDBErr, tlogs.Error, err) {
			// 如果连接失败
			return false
		}
		return true
	}
}

// isConnected :判断是否已连接
func (d *DBx) isConnected() bool {
	if d.db.Ping() == nil {
		return true
	}
	return false
}

// ExceSQL :执行查询语句
func (d *DBx) ExceSQL(s string, args ...interface{}) bool {
	_, err := d.db.Exec(s, args...)
	if cm.CkErr(DBMsg.QueryDataErr, tlogs.Error, err) {
		return false
	}
	// cm.Msg("影响的数据量：", r.RowsAffected())
	return true
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

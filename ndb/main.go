package ndb

import (
	"github.com/jmoiron/sqlx"
	// mysql 接口
	_ "github.com/go-sql-driver/mysql"
)

// var db *sqlx.DB

// 数据库连接参数
type dbParam struct {
	UserName string
	Pwd      string
	Host     string
	Port     string
}

// DBx :数据库对象
type DBx struct {
	db     *sqlx.DB
	param  dbParam
	dbType databaseTypeForDriver
}

// NewDBx :新建数据库对象
func NewDBx(param dbParam, databaseType databaseTypeForDriver) *DBx {
	return &DBx{
		param:  param,
		dbType: databaseType,
	}
}

// Connect :连接到数据库
func (d *DBx) Connect() {
	var err error
	dbsource := ""

	d.db, err = sqlx.Open(d.dbType, dbsource)
}

package ndb

import (
	"eInfusion/comm"
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
	params       dbParam
	databaseType string
}

// NewDBx :新建数据库对象
func NewDBx(param dbParam, dbType string) *DBx {
	return &DBx{
		params:       param,
		databaseType: dbType,
	}
}

// Connect :连接到数据库
func (d *DBx) Connect() {
	var err error
	// "user:password@tcp(127.0.0.1:3306)/hello")
	dbsource := d.params.UserName + ":" + d.params.Password + "@tcp(" + d.params.HostNameOrIPAddr + "):"
	dbsource = dbsource + d.params.Port + "/" + d.params.schemaName

	d.db, err = sqlx.Open(d.databaseType, dbsource)
	if !comm.CkErr(DBMsg.ConnectDBErr, tlogs.Error, err) {

	}

}

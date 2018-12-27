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
	DB           *sqlx.DB
	Params       *DBParam
	DataBaseType string
}

// NewDBx :新建数据库对象
func NewDBx(params *DBParam, DBType string) *DBx {
	return &DBx{
		Params:       params,
		DataBaseType: DBType,
	}
}

// Connect :连接到数据库
func (d *DBx) Connect() bool {
	var err error
	// "user:password@tcp(127.0.0.1:3306)/hello")
	dbsource := d.Params.UserName + ":" + d.Params.Password + "@tcp(" + d.Params.HostNameOrIPAddr + "):"
	dbsource = dbsource + d.Params.Port + "/" + d.Params.SchemaName

	d.DB, err = sqlx.Open(d.DataBaseType, dbsource)
	cm.Msg("err:", err)
	if !cm.CkErr(DBMsg.ConnectDBErr, tlogs.Error, err) {
		// 连接成功
		// cm.Msg(d.DB.Stats().OpenConnections)
		return true
	}
	return false
}

// QueryData :查询数据
func (d *DBx) QueryData() {

}

// ExecSQL :执行sql语句
func (d *DBx) ExecSQL(s string) bool {

	return true
}

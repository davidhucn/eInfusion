package ndb

import (
	"eInfusion/comm"
	cm "eInfusion/comm"
	"eInfusion/tlogs"
	"errors"

	"github.com/jmoiron/sqlx"
	// d "golang.org/x/net/http2/h2demo"

	// mysql 接口
	_ "github.com/go-sql-driver/mysql"
	// sqlite接口
	_ "github.com/mattn/go-sqlite3"
	// oracle接口
	// _ "github.com/mattn/go-oci8"
	// PostgreSQL 接口:
	_ "github.com/lib/pq"
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
func NewDBparams(databaseType string, UserName string, Password string, Host string, Port string, SchemaName string) *DBParam {
	p := &DBParam{
		UserName:     UserName,
		Password:     Password,
		Host:         Host,
		Port:         Port,
		SchemaName:   SchemaName,
		DataBaseType: databaseType,
	}
	return p
}

// DBx :数据库对象
type DBx struct {
	DB            *sqlx.DB
	dbParams      *DBParam
	connectString string
}

// NewDBx :新建数据库对象
func NewDBx(p *DBParam) *DBx {
	var connectStr string
	switch p.DataBaseType {
	case DataBaseType.MySQL:
		connectStr = p.UserName + ":" + p.Password + "@tcp(" + p.Host + ":"
		connectStr = connectStr + p.Port + ")/" + p.SchemaName
		connectStr = connectStr + "?charset=utf8" //字符集
	case DataBaseType.Sqlite3:
		connectStr = p.Host + p.SchemaName
	case DataBaseType.Oracle:
		connectStr = ""
	case DataBaseType.PostgreSQL:
		connectStr = ""
	}
	d := &DBx{
		dbParams:      p,
		connectString: connectStr,
	}
	return d
}

// Connect :连接到数据库
func (d *DBx) Connect() bool {
	var err error
	d.DB, err = sqlx.Connect(d.dbParams.DataBaseType, d.connectString)
	if cm.CkErr(DBMsg.ConnectDBErr, tlogs.Error, err) {
		// 如果连接失败
		return false
	}
	return true
}

// isConnected :判断是否已连接
func (d *DBx) isConnected() bool {
	if d.DB.Ping() == nil {
		return true
	}
	return false
}

// ExceSQL :执行查询语句,如果是DDL语句则返回都为0
func (d *DBx) ExceSQL(s string, args ...interface{}) int64 {
	rs, err := d.DB.Exec(s, args...)
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
	if affAmount > 0 || inAmount > 0 {
		if affAmount > 0 {
			return affAmount
		}
		return inAmount
	}
	return 0
}

// QueryOneData :查询第一条数据，返回相关struct
func (d *DBx) QueryOneData(sql string, result interface{}, args ...interface{}) bool {
	var err error
	err = d.DB.Get(result, sql, args...)
	if cm.CkErr(DBMsg.QueryDataErr, tlogs.Error, err) {
		return false
	}
	return true
}

// QueryDatas :查询批量数据,返回相关[]struct
func (d *DBx) QueryDatas(sql string, results interface{}, args ...interface{}) bool {
	var err error
	err = d.DB.Select(&results, sql, args...)
	if cm.CkErr(DBMsg.QueryDataErr, tlogs.Error, err) {
		return false
	}
	return true
}

// DoTransacion :进行事务操作，基于TxArgs
func (d *DBx) DoTransacion(t []*TransacionArgs) error {
	tx, err := d.DB.Begin()
	if comm.CkErr(DBMsg.EnableTransacionFailure, tlogs.Error, err) {
		return errors.New(DBMsg.EnableTransacionFailure)
	}
	for _, args := range t {
		_, err := tx.Exec(args.SQL, args.Params...)
		if comm.CkErr(DBMsg.TransacionOperateErr, tlogs.Error, err) {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if comm.CkErr(DBMsg.TransacionOperateErr, tlogs.Error, err) {
		return errors.New(DBMsg.TransacionOperateErr)
	}
	return nil
}

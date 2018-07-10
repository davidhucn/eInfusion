//***********************************************************
//					 Package Name dbOperate				    *
//					  File Name tMysql.go					*
// 						Author:David.Hu						*
//						Date:2018.03.22						*
//				Remark:use the comm sql.Db object			*
//***********************************************************
package dbOperate

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

const (
	c_DB_IPAddr  = "127.0.0.1"
	c_DB_Port    = "3306"
	c_DB_schema  = "transfusion"
	c_DB_UsrName = "root"
	c_DB_Pwd     = "2341656"
)

var G_Db *sql.DB

func init() {
	var err error
	strDataSource := c_DB_UsrName + ":" + c_DB_Pwd + "@tcp(" + c_DB_IPAddr + ":" + c_DB_Port + ")/"
	strDataSource = strDataSource + c_DB_schema + "?charset=utf8"
	G_Db, err = sql.Open("mysql", strDataSource)
	if err != nil {
		panic(err)
	}
}

func IsConnected() bool {
	//	var dbStats sql.DBStats
	if G_Db.Stats().OpenConnections > 0 {
		return true
	}
	return false
}

//快速执行查询到指定数据库内
func ExecSQL(strSql string, args ...interface{}) (affected_Num int64, err error) {
	result, err := G_Db.Exec(strSql, args...)
	if err != nil {
		return 0, err
	}
	affected_Num, _ = result.RowsAffected()
	return affected_Num, err
}

//插入数据到指定数据库内
func InsertData(strSql string, args ...interface{}) (affected_Num int64, err error) {
	var result sql.Result
	stmtIns, err := G_Db.Prepare(strSql)
	if err != nil {
		return 0, err
	}
	defer stmtIns.Close()
	result, err = stmtIns.Exec(args...)
	if err != nil {
		return 0, err
	}
	affected_Num, _ = result.RowsAffected()
	return affected_Num, err
}

//根据条件册除数据库
func DeleteData(strSql string, args ...interface{}) (affected_Num int64, err error) {
	var result sql.Result
	stmtDel, err := G_Db.Prepare(strSql)
	if err != nil {
		return 0, err
	}
	defer stmtDel.Close()
	result, err = stmtDel.Exec(args...)
	affected_Num, _ = result.RowsAffected()
	return affected_Num, err
}

//册除全表
func TruncateTable(strTableName string) (affected_Num int64, err error) {
	var result sql.Result
	result, err = G_Db.Exec("Truncate Table " + strTableName)
	if err != nil {
		return 0, err
	}
	affected_Num, _ = result.RowsAffected()
	return affected_Num, err
}

//查询单条数据,结果皆为string
func QueryOneRow(strSql string, args ...interface{}) (*map[string]string, error) {

	stmtOut, err := G_Db.Prepare(strSql)
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(args...)
	if err != nil {
		return nil, err
	}
	//获取字段对象
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	ret := make(map[string]string, len(scanArgs))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			ret[columns[i]] = value
		}
		break //get the first row only
	}
	return &ret, nil
}

//查询多条数据,结果皆为string
func QueryRows(strSql string, args ...interface{}) (*[]map[string]string, error) {
	stmtOut, err := G_Db.Prepare(strSql)
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(args...)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))

	ret := make([]map[string]string, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		var value string
		vmap := make(map[string]string, len(scanArgs))
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret, vmap)
	}
	return &ret, nil
}

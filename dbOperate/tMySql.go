package dbOperate

import (
	"database/sql"
	//	"eInfusion/comm"

	_ "github.com/Go-SQL-Driver/MySQL"
)

const (
	c_DB_IPAddr  = "127.0.0.1"
	c_DB_Port    = "3306"
	c_DB_schema  = "transfusion"
	c_DB_UsrName = "root"
	c_DB_Pwd     = "2341656"
)

//数据库连接类错误提示信息
const (
	C_Msg_DBConnect_Err  = "错误,无法连接到指定数据库！"
	C_Msg_DBInsert_Err   = "错误,插入数据库操作失败！"
	C_Msg_DBDelete_Err   = "错误,册除数据操作失败！"
	C_Msg_DBTruncate_Err = "错误,册除指定数据表内所有信息失败！"
	C_Msg_DBQuery_Err    = "错误,查询数据信息失败！"
)

var G_Db *sql.DB

func init() {
	//	var err error
	strDataSource := c_DB_UsrName + ":" + c_DB_Pwd + "@tcp(" + c_DB_IPAddr + ":" + c_DB_Port + ")/"
	strDataSource = strDataSource + c_DB_schema + "?charset=utf8"
	G_Db, _ = sql.Open("mysql", strDataSource)

}

//连接数据库
//func ConnectDB() error {
//	//	连接用数据库信息
//	strDataSource := c_DB_UsrName + ":" + c_DB_Pwd + "@tcp(" + c_DB_IPAddr + ":" + c_DB_Port + ")/"
//	strDataSource = strDataSource + c_DB_schema + "?charset=utf8"
//	db, err := sql.Open("mysql", strDataSource)
//	defer db.Close()
//	if err != nil {
//		panic(err.Error())
//		return err
//	}
//	G_Db = db
//	return nil
//}

func IsConnected() bool {
	//	var dbStats sql.DBStats
	if G_Db.Stats().OpenConnections > 0 {
		return true
	}
	return false
}

//快速插入数据到指定数据库内
func InsertDataFast(strSql string, args ...interface{}) (affected_Num int64, err error) {
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

//func QueryFast(strSql string, args ...interface{}) (bool, error) {
//	//	var rows sql.Rows
//	comm.Msg("start...")
//	rows, err := G_Db.Query(strSql, args...)
//	if err != nil {
//		return false, err
//	}
//	defer rows.Close()
//	for rows.Next() {
//		var receiver_id int
//		if err := rows.Scan(&receiver_id); err != nil {
//			comm.Msg("rows err:", err)
//		}
//		//fmt.Printf("name:%s ,id:is %d\n", name, id)
//		comm.Msg(receiver_id)
//	}
//	return true, nil
//}

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

///////////////////////////////////////////////////////////////////////////////////////////////
//func TestDb() {
//	db, err := sql.Open("mysql", "root:2341656@tcp(127.0.0.1:3306)/einfusion?charset=utf8")
//	if err != nil {
//		fmt.Print("database error:")
//		fmt.Println(err)
//		return
//	}
//	defer db.Close()
////////////////insert sample//////////////////////////////////////////////////////
//	var result sql.Result
//	result, err = db.Exec("insert into t_main(ip_address, unit_id,master_id,time) values(?,?,?,?)",
//		"127.0.0.2", "slave_002", "master_001", "2018-01-04")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	lastId, _ := result.LastInsertId()
//	fmt.Println("新插入记录的ID为", lastId)
/////////////////////////////////////////////////////////////////////////////////

//////////////////////////////insert sample by another way///////////////////////////////
//	stmt, err := db.Prepare(`insert into t_main(ip_address, unit_id,master_id,time) values(?,?,?,?)`)
//	checkErr(err)
//	res, err := stmt.Exec("127.0.0.2", "slave_003", "master_002", "2018-01-05")
//	checkErr(err)
//	lastId, _ = res.RowsAffected()
//	fmt.Println("新插入记录的ID为", lastId)
/////////////////////////////////////////////////////////////////////////////////////

/////////////////////select sample by one record/////////////////////////////////////
//	var row *sql.Row
//	row = db.QueryRow("select * from t_main")
//	var ip_address, unit_id, master_id, time string
//	err = row.Scan(&ip_address, &unit_id, &master_id, &time)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(ip_address, "\t", time, "\t", unit_id)
//	fmt.Println(".......................")
/////////////////////////////////////////////////////////////////////////////////

/////////////////////////select sample by records////////////////////////////////////
//	var rows *sql.Rows
//	rows, err = db.Query("select * from t_main")
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	for rows.Next() {
//		rows.Scan(&ip_address, &unit_id, &master_id, &time)
//		fmt.Println(ip_address, "\t", time, "\t", unit_id)
//	}
//	rows.Close()
///////////////////////////////////////////////////////////////////////////////////////

//////////////////////////////////delete sample/////////////////////////////////////
//	db.Exec("truncate table t_test")
//	result, err = db.Exec("delete from t_main where unit_id =?", "slave_002")
//	checkErr(err)

//	delId, _ := result.RowsAffected()
//	fmt.Println("册除的记录数", delId)
///////////////////////////////////////////////////////////////////////////////////

//统一处理错误，待改写
//func checkErr(err error, strMessage string) {
//	if err != nil {
//		fmt.Print(strMessage)
//		//		panic(err)
//	}
//}

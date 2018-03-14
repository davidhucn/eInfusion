package dbOperate

import (
	"database/sql"
	"fmt"

	_ "github.com/Go-SQL-Driver/MySQL"
)

const (
	c_DataBase_IPAddress = "127.0.0.1"
	c_DataBase_Port      = "3306"
	c_DataBase_schema    = "transfusion"
)

//数据库连接类错误提示信息
const (
	c_Msg_DBConnect_Err  = "错误,无法连接到指定数据库！"
	c_Msg_DBInsert_Err   = "错误,插入数据库操作失败！"
	c_Msg_DBDelete_Err   = "错误,册除数据操作失败！"
	c_Msg_DBTruncate_Err = "错误,册除指定数据表内所有信息失败！"
	c_Msg_DBQuery_Err    = "错误,查询数据信息失败！"
)

type DBConn struct {
	UserName    string
	Password    string
	Schema      string
	Port        string
	IpAddr      string
	IsConnected bool
	DbHandler   *sql.DB
}

//连接数据库
func (this *DBConn) ConnectDB() error {
	//	连接用数据库信息
	strDataSource := this.UserName + ":" + this.Password + "@tcp(" + this.IpAddr + ":" + this.Port + ")/"
	strDataSource = strDataSource + this.Schema + "?charset=utf8"
	db, err := sql.Open("mysql", strDataSource)
	defer db.Close()
	if err != nil {
		fmt.Println(c_Msg_DBConnect_Err)
		panic(err.Error())
		return err
	}
	//	checkErr(err, c_Msg_DBConnect_Err)
	this.DbHandler = db
	//	if this.DbHandler.Stats().OpenConnections > 0 {
	//		this.IsConnected = true
	//	} else {
	//		this.IsConnected = false
	//	}
	this.IsConnected = true
	return nil
}

//插入数据到指定数据库内
func (this *DBConn) InsertData(strSql string, args ...interface{}) (affected_Num int64, err error) {
	var result sql.Result
	//	var stmtIns sql.Stmt
	//	如果没有连接数据库则强制连接
	if !this.IsConnected {
		this.ConnectDB()
	}
	stmtIns, err := this.DbHandler.Prepare(strSql)
	if err != nil {
		fmt.Println(c_Msg_DBInsert_Err)
		return
	}
	defer stmtIns.Close()
	result, err = stmtIns.Exec(args...)
	if err != nil {
		fmt.Println(c_Msg_DBInsert_Err)
		return
	}
	affected_Num, _ = result.RowsAffected()
	return affected_Num, err
}

//根据条件册除数据库
func (this *DBConn) DeleteData(strSql string, args ...interface{}) (affected_Num int64, err error) {
	var result sql.Result
	//	var stmtDel sql.Stmt
	//	如果没有连接数据库则强制连接
	if !this.IsConnected {
		this.ConnectDB()
	}
	stmtDel, err := this.DbHandler.Prepare(strSql)
	if err != nil {
		fmt.Println(c_Msg_DBInsert_Err)
		return
	}
	defer stmtDel.Close()
	result, err = stmtDel.Exec(args...)
	affected_Num, _ = result.RowsAffected()
	return affected_Num, err
}

//册除全表
func (this *DBConn) TruncateTable(strTableName string) (affected_Num int64, err error) {
	var result sql.Result
	result, err = this.DbHandler.Exec("Truncate Table " + strTableName)
	if err != nil {
		fmt.Println(c_Msg_DBTruncate_Err)
		return
	}
	affected_Num, _ = result.RowsAffected()
	return affected_Num, err
}

//查询单条数据,结果皆为string
func (this *DBConn) QueryDataOneRow(strSql string, args ...interface{}) (*map[string]string, error) {
	//	如果没有连接数据库则强制连接
	if !this.IsConnected {
		this.ConnectDB()
	}
	stmtOut, err := this.DbHandler.Prepare(strSql)
	if err != nil {
		fmt.Println(c_Msg_DBQuery_Err)
		return nil, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(args...)
	if err != nil {
		fmt.Println(c_Msg_DBQuery_Err)
		return nil, err
	}
	//获取字段对象
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(c_Msg_DBQuery_Err)
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
			fmt.Println(c_Msg_DBQuery_Err)
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
func (this *DBConn) QueryData(strSql string, args ...interface{}) (*[]map[string]string, error) {
	//	var stmtOut sql.Stmt
	//	如果没有连接数据库则强制连接
	if !this.IsConnected {
		this.ConnectDB()
	}
	stmtOut, err := this.DbHandler.Prepare(strSql)
	if err != nil {
		fmt.Println(c_Msg_DBQuery_Err)
		return nil, err
	}
	defer stmtOut.Close()
	rows, err := stmtOut.Query(args...)
	if err != nil {
		fmt.Println(c_Msg_DBQuery_Err)
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(c_Msg_DBQuery_Err)
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
			fmt.Println(c_Msg_DBQuery_Err)
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

//关闭数据库连接
func (this *DBConn) DisconnectDB() error {
	if this.IsConnected {
		err := this.DbHandler.Close()
		if err != nil {
			return err
		}
	}
	return nil
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

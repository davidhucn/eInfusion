package db

import (
	"database/sql"
	"fmt"

	_ "github.com/Go-SQL-Driver/MySQL"
)

//type userinfo struct {
//	name        string
//	description string
//	url         string
//}

func TestDb() {
	db, err := sql.Open("mysql", "root:2341656@tcp(127.0.0.1:3306)/einfusion?charset=utf8")
	if err != nil {
		fmt.Print("database error:")
		fmt.Println(err)
		return
	}
	defer db.Close()
	////////////////insert sample//////////////////////////////////////////////////////
	var result sql.Result
	result, err = db.Exec("insert into t_main(ip_address, unit_id,master_id,time) values(?,?,?,?)",
		"127.0.0.2", "slave_002", "master_001", "2018-01-04")
	if err != nil {
		fmt.Println(err)
		return
	}
	lastId, _ := result.LastInsertId()
	fmt.Println("新插入记录的ID为", lastId)
	/////////////////////////////////////////////////////////////////////////////////

	/////////////////////select sample by one record/////////////////////////////////////
	var row *sql.Row
	row = db.QueryRow("select * from t_main")
	var ip_address, unit_id, master_id, time string
	err = row.Scan(&ip_address, &unit_id, &master_id, &time)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ip_address, "\t", time, "\t", unit_id)
	fmt.Println(".......................")
	/////////////////////////////////////////////////////////////////////////////////

	/////////////////////////select sample by records////////////////////////////////////
	var rows *sql.Rows
	rows, err = db.Query("select * from t_main")
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		rows.Scan(&ip_address, &unit_id, &master_id, &time)
		fmt.Println(ip_address, "\t", time, "\t", unit_id)
	}
	rows.Close()
	///////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////insert sample by another way///////////////////////////////
	stmt, err := db.Prepare(`insert into t_main(ip_address, unit_id,master_id,time) values(?,?,?,?)`)
	checkErr(err)
	res, err := stmt.Exec("127.0.0.2", "slave_003", "master_002", "2018-01-05")
	checkErr(err)
	lastId, _ = res.RowsAffected()
	fmt.Println("新插入记录的ID为", lastId)
	/////////////////////////////////////////////////////////////////////////////////////

	//////////////////////////////////delete sample/////////////////////////////////////
	//	db.Exec("truncate table t_test")
	//	result, err = db.Exec("delete from t_main where unit_id =?", "slave_002")
	//	checkErr(err)

	//	delId, _ := result.RowsAffected()
	//	fmt.Println("册除的记录数", delId)
	///////////////////////////////////////////////////////////////////////////////////
}

func checkErr(err error) {
	if err != nil {
		fmt.Print(err)
		//		panic(err)
	}
}

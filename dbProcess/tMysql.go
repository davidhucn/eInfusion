package dbProcess

//import (
//	"fmt"
//	"runtime"
//	"time"

//	"github.com/jmoiron/sqlx"
//	_ "github.com/lib/pq"
//)

//const (
//	c_DB_IPAddr  = "127.0.0.1"
//	c_DB_Port    = "3306"
//	c_DB_schema  = "transfusion"
//	c_DB_UsrName = "root"
//	c_DB_Pwd     = "2341656"
//)

////用户表结构
//type User struct {
//	Id      int       `db:"id"`
//	Name    string    `db:"name"`
//	Created time.Time `db:"created"`
//	ClassId int       `db:"class_id"`
//}

////Class表结构
//type Class struct {
//	Id   int    `db:"id"`
//	Name string `db:"name"`
//}

//func main() {

//	//1. 连接数据库
//	db, err := sqlx.Open("postgres", "postgres://postgres:123@localhost:5432/test?sslmode=disable")
//	SimplePanic(err)
//	fmt.Println(db)
//	defer db.Close()

//	//2.  crud 增删改查
//	//2.1 增加，删除，修改
//	if false {
//		var crud int = 1 //增改删123
//		switch crud {
//		case 1:
//			_, err = db.Exec("insert into public.user(name) values($1)", "ft")
//			SimplePanic(err)
//		case 2:
//			_, err = db.Exec("update public.user set name ='ftx' where name=$1", "ft")
//			SimplePanic(err)

//		case 3:
//			_, err = db.Exec("delete from public.user where name=$1", "ftx")
//			SimplePanic(err)

//		}
//	}
//	//2.2 查询
//	if true {
//		user := User{}
//		users := make([]User, 0)
//		//查询单个Get
//		if false {
//			err = db.Get(&user, "select * from public.user where id=$1", 1016)
//			SimplePanic(err)
//			fmt.Println(user)
//		}

//		//查询多个Select
//		if false {
//			err = db.Select(&users, "select * from public.user where id>$1", 900)
//			SimplePanic(err)
//		}
//		//2.3使用Queryx和QueryRowx来获取
//		if true {
//			//Queryx()
//			if false {
//				var rows *sqlx.Rows
//				rows, err = db.Queryx("select * from public.user where id>$1", 1011)
//				SimplePanic(err)
//				for rows.Next() {

//					err = rows.Scan(&user.Id, &user.Name, &user.Created, &user.ClassId)
//					SimplePanic(err)
//					users = append(users, user)
//				}
//				SimplePanic(rows.Err())
//				fmt.Println(len(users))
//			}

//			//QueryRowx()
//			if false {
//				var row *sqlx.Row
//				var user = User{}
//				row = db.QueryRowx("select * from public.user where id =$1", 1016)
//				err = row.Scan(&user.Id, &user.Name, &user.Created, &user.ClassId)
//				SimplePanic(row.Err())
//				fmt.Println(user)
//			}
//		}

//	}

//	//3. 事务
//	if false {
//		var tx *sqlx.Tx
//		tx, err = db.Beginx()
//		SimplePanic(err)

//		_, err = tx.Exec("insert into public.user(name) values($1)", "ft3")
//		if err != nil {
//			tx.Rollback()
//			SimplePanic(err)
//		}
//		_, err = tx.Exec("insert into public.user(name) values($1)", "ft4")
//		if err != nil {
//			tx.Rollback()
//			SimplePanic(err)
//		}
//		err = tx.Commit()
//		SimplePanic(err)
//	}

//	//4.特点，命名参数的的sql语句绑定
//	if false {
//		userResult := User{}
//		userQuery := User{}
//		userQuery.Name = "ft"
//		var rows *sqlx.Rows
//		fmt.Println(userQuery.Name)
//		rows, err = db.Queryx("SELECT * FROM public.user WHERE name=$1", userQuery.Name)
//		for rows.Next() {
//			err = rows.Scan(&userResult.Id, &userResult.Name, &userResult.Created, &userResult.ClassId)
//			SimplePanic(err)
//		}
//		SimplePanic(err)
//	}

//}

//func SimplePanic(err error) {
//	if err != nil {
//		_, file, line, _ := runtime.Caller(1)
//		fmt.Println(file, line, err)
//		runtime.Goexit()
//	}

//}

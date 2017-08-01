package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

type User struct {
	id       int
	name     string
	pwd      string
	nickname string
}

/*
func (db *sql.DB) init() {
	db, err := sql.Open("mysql", "root:chenghai3c@tcp(localhost:3306)/student?charset=utf8")
	if err != nil {
		panic(err.Error())
		log.Println(err)
	}
}
*/

func main() {
	//Add()
	//Del()
	//Update()
	//Select()
	//Find()
	Transaction()
}

func Transaction() {
	db, err := sql.Open("mysql", "root:chenghai3c@tcp(localhost:3306)/student?charset=utf8")
	if err != nil {
		panic(err.Error())
		log.Println(err)
	}
	defer db.Close()

	//事务
	tx, err := db.Begin() //声明一个事务的开始
	if err != nil {
		log.Println(err)
		return
	}
	insert_sql := "insert into t_user (name, pwd, nickname) value(?,?,?)"
	insert_stmt, insert_err := tx.Prepare(insert_sql)
	if insert_err != nil {
		log.Println(insert_err)
		return
	}
	insert_res, insert_err := insert_stmt.Exec("Billy", "123abc", "Cat")
	last_insert_id, _ := insert_res.LastInsertId()
	log.Println(last_insert_id)

	//defer tx.Rollback() //回滚之前上面的last_insert_id是有的，但在回滚后该操作没有被提交，被回滚了，所以上面打印的last_insert_id的这条数据是不存在与数据库表中的
	tx.Commit()         //这里提交了上面的操作，所以上面的执行的sql 会在数据库中产生一条数据
}

//find
func Find() {
	db, err := sql.Open("mysql", "root:chenghai3c@tcp(localhost:3306)/student?charset=utf8")
	if err != nil {
		panic(err.Error())
		log.Println(err)
	}
	defer db.Close()

	//查询多条
	select_sql := "select * from t_user where id > ?"
	select_rows, select_err := db.Query(select_sql, 1)
	if select_err != nil {
		log.Println(select_err)
		return
	}
	defer select_rows.Close()

	var user User
	for select_rows.Next() {
		// 直接声明变量
		/*var id int
		var name string
		var pwd string
		var nickname string*/
		// 利用声明好的user来存储数据
		if err := select_rows.Scan(&user.id, &user.name, &user.pwd, &user.nickname); err != nil {
			log.Println(err)
			return
		}
		log.Printf("id=%v, name=%v, pwd=%v, nickname=%v", user.id, user.name, user.pwd, user.nickname)
	}
}

//select
func Select() {
	db, err := sql.Open("mysql", "root:chenghai3c@tcp(localhost:3306)/student?charset=utf8")
	if err != nil {
		panic(err.Error())
		log.Println(err)
	}
	defer db.Close()

	var user User
	select_sql := "select * from t_user where id > ?"
	select_err := db.QueryRow(select_sql, 1).Scan(&user.id, &user.name, &user.pwd, &user.nickname) //查询一条，返回一条结果。并赋值到user这个结构体类型的变量中,就算查询到的是多条，单返回的还是一条
	if select_err != nil {                                                                         //如果没有查询到任何数据就进入if中err：no rows in result set
		log.Println(select_err)
		return
	}
	log.Println(user)
}

//update
func Update() {
	db, err := sql.Open("mysql", "root:chenghai3c@tcp(localhost:3306)/student?charset=utf8")
	if err != nil {
		panic(err.Error())
		log.Println(err)
	}
	defer db.Close()

	update_sql := "update t_user set name=? where id=?"
	update_stmt, update_err := db.Prepare(update_sql)
	if update_err != nil {
		log.Println(update_err)
		return
	}
	update_res, update_err := update_stmt.Exec("Tom", 4)
	if update_err != nil {
		log.Printf("%v", update_err)
		return
	}
	affect_count, _ := update_res.RowsAffected() //返回影响的条数,注意有两个返回值
	log.Printf("%v", affect_count)
}

//delete
func Del() {
	db, err := sql.Open("mysql", "root:chenghai3c@tcp(localhost:3306)/student?charset=utf8")
	if err != nil {
		panic(err.Error())
		log.Println(err)
	}
	defer db.Close()

	del_sql := "delete from t_user where id=?"
	del_stmt, del_err := db.Prepare(del_sql)
	if del_err != nil {
		panic(err.Error())
		log.Println(err)
	}
	del_stmt.Exec(3) //不返回任何结果
}

//insert
func Add() bool {
	name := "Lsyl"
	pwd := "123456"
	nickname := "Jerry"
	b := false
	db, err := sql.Open("mysql", "root:chenghai3c@tcp(localhost:3306)/student?charset=utf8")
	if err != nil {
		panic(err.Error())
		log.Println(err)
		return b
	}
	defer db.Close() //只有在前面用了 panic[抛出异常] 这时defer才能起作用，如果链接数据的时候出问题，他会往err写数据。defer:延迟，这里立刻申请了一个关闭sql 链接的错误，defer 后的方法，或延迟执行。在函数抛出异常一会被执行
	insert_sql := "insert into t_user (name, pwd, nickname) value (?,?,?),(?,?,?),(?,?,?),(?,?,?)"
	stmt, err := db.Prepare(insert_sql) //准备一个sql操作，返回一个*Stmt,用户后面的执行,这个Stmt可以被多次执行，或者并发执行
	/*
	*    这个stmt的主要方法:Exec、Query、QueryRow、Close
	 */
	if err != nil {
		log.Println(err)
		return b
	}
	res, err := stmt.Exec(name, pwd, nickname, name, pwd, nickname, name, pwd, nickname, name, pwd, nickname)
	if err != nil {
		log.Println(err)
		return b
	}
	lastInsertId, err := res.LastInsertId() //批量插入的时候LastInserId返回的是第一条id,单条插入则返回这条的id
	//lastInsertId,err := res.RowsAffected()        //插入的是后RowsAffected 返回的是插入的条数
	if err != nil {
		log.Println(err)
		return b
	}
	//log.Println(reflect.TypeOf(lastInsertId))    //打印变量类型
	last_insert_id_string := strconv.FormatInt(lastInsertId, 10) //int64 转string 需要引入 strconv包
	log.Println("lastInsertId = " + last_insert_id_string)
	return true
}

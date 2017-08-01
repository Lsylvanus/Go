package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:chenghai3c@/squarenum")
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()
	fmt.Println("链接数据库成功")

	fmt.Println(db.Query("SELECT * FROM squareNum WHERE number = 1"))

}
package main

import (
	"io/ioutil"
	"fmt"
)

const sqlLogFileName = "catch.log"

func main() {
	sqlLogName := sqlLogFileName

	sql_b, err := ioutil.ReadFile(sqlLogName)
	//not err == io.EOF
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(sql_b)
	fmt.Println(len(sql_b))
	fmt.Println(1024*1024*2)
	if len(sql_b) > 1024*1024*2 {
		fmt.Printf("read file %s.log's memory ：%d，remove and rebuild it.", sqlLogName, len(sql_b))
	} else {
		fmt.Printf("the file %s is too small to be removed.", sqlLogName)
	}
}
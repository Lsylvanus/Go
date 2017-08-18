package main

import (
	"io/ioutil"
	"fmt"
	"os"
)

const LogFileName = "././catch.log"

func main() {
	sqlLogName := LogFileName

	sql_b, err := ioutil.ReadFile(LogFileName)
	//not err == io.EOF
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(len(sql_b))
	fmt.Println(1024*1024*2)
	if len(sql_b) > 1024*1024*2 {
		fmt.Printf("read file %s.log's memory ：%d，remove and rebuild it.", sqlLogName, len(sql_b))
		err := os.Remove(LogFileName)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Printf("the file %s is too small to be removed.", sqlLogName)
	}
}
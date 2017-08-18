package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"log"
	"os"
	"strconv"
	"time"
)

var e *xorm.Engine
var s *xorm.Session

type Users struct {
	Id        int64     `xorm:"int(11) pk notnull autoincr"`
	Name      string    `xorm:"varchar(16) null"`
	Age       int64     `xorm:"int(11) null"`
	Pwd       string    `xorm:"varchar(20) null"`
	NickName  string    `xorm:"varchar(16) null 'nickname'"`
	CreatedAt time.Time `xorm:"created 'created'"`
	UpdatedAt time.Time `xorm:"updated 'updated'"`
	DeletedAt time.Time `xorm:"deleted 'deleted'"`
	GroupId   int64     `xorm:"index 'index'"`
}

func main() {
	conn()

	if b := create(Users{}); !b {
		return
	}

	users := make([]Users, 5)	//数量，如果超过所需插入数量，则其他字段填默认插入
	users[0].Name = "lsy"
	users[0].Age = 13
	users[1].Name = "esy"
	users[1].Age = 14
	users[2].Name = "what"
	users[2].Age = 15
	users[3].Name = "jim"
	users[3].Age = 16
	users[4].Name = "tom"
	users[4].Age = 17

	aff, err := e.Insert(&users)
	if err != nil {
		log.Println(err.Error())
		return
	}
	affInt := strconv.FormatInt(aff, 10)
	log.Println("insert into t_users :", affInt)
}

func conn() {
	//创建Orm引擎
	var err error
	e, err = xorm.NewEngine("mysql", "root:chenghai3c@/test?charset=utf8&collation=utf8_general_ci")
	if err != nil {
		log.Println(err.Error())
		return
	}
	//测试连接
	errPing := e.Ping()
	if errPing != nil {
		log.Println(errPing.Error())
		return
	}
	log.Println("conn success!")

	//当使用事务处理时，需要创建Session对象。在进行事物处理时，可以混用ORM方法和RAW方法
	s = e.NewSession()
	defer s.Close()
	// add Begin() before any action
	errBegin := s.Begin()
	if errBegin != nil {
		log.Println(errBegin.Error())
		return
	}

	//日志是一个接口，通过设置日志，可以显示SQL，警告以及错误等，默认的显示级别为INFO
	e.ShowSQL(true)
	e.Logger().SetLevel(core.LOG_DEBUG)
	//日志信息保存
	f, err := os.Create("sql_test.log")
	if err != nil {
		log.Println(err.Error())
		return
	}
	e.SetLogger(xorm.NewSimpleLogger(f))

	// 映射规则
	// 默认对应的表名就变成了 t_user 了，而之前默认的是 user
	tMapper := core.NewPrefixMapper(core.SnakeMapper{}, "t_")
	e.SetTableMapper(tMapper)

	//最后一个值得注意的是时区问题，默认xorm采用Local时区，所以默认调用的time.Now()会先被转换成对应的时区。
	e.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
}

func create(beanOrTableName interface{}) bool {
	has, err := e.IsTableExist(beanOrTableName)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if has {
		log.Println("table is exist.")
		return true
	} else {
		err := e.CreateTables(beanOrTableName)
		if err != nil {
			log.Println("create table failed.")
			return false
		}
		log.Println("create table success.")
	}
	return true
}

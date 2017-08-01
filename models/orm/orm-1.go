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

var id int

type Group struct {
	Id   int64  `xorm:"int(11) pk notnull autoincr"`
	Name string `xorm:"varchar(16) null"`
}

type User struct {
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

type UserGroup struct {
	User  `xorm:"extends"`
	Group `xorm:"extends"`
}

var engine *xorm.Engine

func (UserGroup) TableName() string {
	return "t_user"
}

func main() {
	//创建Orm引擎
	var err error
	engine, err = xorm.NewEngine("mysql", "root:chenghai3c@/student?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	//测试连接
	errPing := engine.Ping()
	if errPing != nil {
		panic(errPing.Error())
	}
	log.Printf("conn success!")
	// defer engine.Close()

	//全局的内存缓存
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	engine.SetDefaultCacher(cacher)

	//针对部分表
	//cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	//engine.MapCacher(&user, cacher)

	//禁用某个表的缓存
	//engine.MapCacher(&user, nil)

	//清理某个表的缓存
	//engine.ClearCache(new(User))

	//当使用事务处理时，需要创建Session对象。在进行事物处理时，可以混用ORM方法和RAW方法
	session := engine.NewSession()
	defer session.Close()
	// add Begin() before any action
	errBegin := session.Begin()
	if errBegin != nil {
		panic(errBegin.Error())
	}

	//日志是一个接口，通过设置日志，可以显示SQL，警告以及错误等，默认的显示级别为INFO
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)
	//日志信息保存
	f, err := os.Create("sql.log")
	if err != nil {
		println(err.Error())
		return
	}
	engine.SetLogger(xorm.NewSimpleLogger(f))

	//连接池
	//最大打开连接数
	engine.SetMaxOpenConns(100)
	//连接池的空闲数大小
	engine.SetMaxIdleConns(20)

	// 映射规则
	// 默认对应的表名就变成了 t_user 了，而之前默认的是 user
	tMapper := core.NewPrefixMapper(core.SnakeMapper{}, "t_")
	engine.SetTableMapper(tMapper)

	//最后一个值得注意的是时区问题，默认xorm采用Local时区，所以默认调用的time.Now()会先被转换成对应的时区。
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	//获取表结构信息
	var db []*core.Table
	db, err = engine.DBMetas()
	if err != nil {
		panic(err.Error())
	}
	log.Println(db)

	//插入数据
	user := new(User)
	user.Name = "Randy"
	user.Age = 22
	user.Pwd = "root0011"
	user.NickName = "Ray"
	user.CreatedAt = time.Now()
	user.GroupId = 1
	affected, err := engine.Insert(user)
	if err != nil {
		panic(err.Error())
	}

	// engine.Id(1).Get(&user)
	// SELECT * FROM user Where id = 1
	// engine.Id(core.PK{1, "name"}).Get(&user)
	// SELECT * FROM user Where id =1 AND name= 'name'

	// engine.Cols("age", "name").Get(&user)
	// SELECT age, name FROM user limit 1
	// engine.Cols("age", "name").Find(&users)
	// SELECT age, name FROM user
	// engine.Cols("age", "name").Update(&user)
	// UPDATE user SET age=? AND name=?

	// engine.AllCols().Id(1).Update(&user)
	// UPDATE user SET name = ?, age =?, gender =? WHERE id = 1

	// 和cols相反，此函数指定排除某些指定的字段。注意：此方法和Cols方法不可同时使用。
	// 例1：
	// engine.Omit("age", "gender").Update(&user)
	// UPDATE user SET name = ? AND department = ?
	// 例2：
	// engine.Omit("age, gender").Insert(&user)
	// INSERT INTO user (name) values (?) // 这样的话age和gender会给默认值
	// 例3：
	// engine.Omit("age", "gender").Find(&users)
	// SELECT name FROM user //只select除age和gender字段的其它字段

	// 按照参数中指定的字段归类结果。
	// engine.Distinct("age", "department").Find(&users)
	// SELECT DISTINCT age, department FROM user

	// 禁用自动根据结构体中的值来生成条件
	// engine.Where("name = ?", "lunny").Get(&User{Id:1})
	// SELECT * FROM user where name='lunny' AND id = 1 LIMIT 1
	// engine.Where("name = ?", "lunny").NoAutoCondition().Get(&User{Id:1})
	// SELECT * FROM user where name='lunny' LIMIT 1

	affecttedStr := strconv.FormatInt(affected, 10)
	log.Printf("成功插入" + affecttedStr + "条记录~！")
	// INSERT INTO user (name, pwd, nickname, created) values (?, ?, ?, ?)

	// 单行记录查询
	user1 := new(User)
	user1.Id = 5
	//has, err := engine.Where("name=?", "Billy").Get(user1)
	has, err := engine.Id(user1.Id).Get(user1)
	// 复合主键的获取方法
	// has, errr := engine.Id(xorm.PK{1,2}).Get(user)
	if err != nil {
		panic(err.Error())
	}
	if !has {
		log.Printf("Find nothing ...")
		return
	}
	log.Println(user1)

	// 多行记录查询
	users := make([]User, 0)
	errFind := engine.Where("age > ? or name = ?", 30, "Tom").Limit(20, 10).Find(&users)
	if errFind != nil {
		panic(errFind.Error())
	}
	log.Println(users)

	users1 := make(map[int64]User)
	errId := engine.Find(&users1)
	if errId != nil {
		panic(errId.Error())
	}
	log.Println(users1)

	//Join的使用
	users2 := make([]UserGroup, 0)
	engine.Join("INNER", "t_group", "t_group.id = t_user.index").Find(&users2)
	log.Println(users2)

	//Iterate方法提供逐条执行查询到的记录的方法，他所能使用的条件和Find方法完全相同
	errIter := engine.Where("age > ? or name = ?", 30, "Tom").Iterate(new(User), func(i int, bean interface{}) error {
		//user2 := bean.(*User)
		users := make(map[int64]User)
		err := engine.Find(&users)
		if err != nil {
			panic(err.Error())
		}
		log.Println(users)
		return err
		//do somthing use i and users
	})
	if errIter != nil {
		panic(err.Error())
	}

	//Rows方法和Iterate方法类似，提供逐条执行查询到的记录的方法，不过Rows更加灵活好用。
	user2 := new(User)
	rows, err := engine.Where("id > ?", 30).Rows(user2)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(user2)
		//...
		if err != nil {
			panic(err.Error())
		}
		log.Println(user2)
	}

	//正序排序
	userDescs := make([]User, 0)
	errAsc := engine.Where("id > ?", 50).Desc("id").Find(&userDescs)
	if errAsc != nil {
		panic(errAsc.Error())
	}
	log.Println(userDescs)

	//统计数据使用Count方法，Count方法的参数为struct的指针并且成为查询条件。
	user3 := new(User)
	total, err := engine.Where("id > ?", 30).Count(user3)
	if err != nil {
		panic(err.Error())
	}
	log.Println("满足id > 30的人数为：", total)

	//update更新数据
	user4 := new(User)
	user4.Name = "dyt"
	user4.UpdatedAt = time.Now()
	affectedUpdate, err := engine.Id(45).Cols("name, updated").Update(user4)
	if err != nil {
		panic(err.Error())
	}
	affectedUpdateStr := strconv.FormatInt(affectedUpdate, 10)
	log.Println("成功更新" + affectedUpdateStr + "条记录~！")

	//delete删除数据
	user5 := new(User)
	affectedDelete, err := engine.Id(44).Delete(user5)
	if err != nil {
		panic(err.Error())
	}
	affectedDeleteStr := strconv.FormatInt(affectedDelete, 10)
	log.Println("成功删除" + affectedDeleteStr + "条记录~！")

	/*
	在Delete()时，deleted标记的字段将会被自动更新为当前时间而不是去删除该条记录
	var user6 User
	engine.Id(1).Get(&user6)
	// SELECT * FROM user WHERE id = ?
	engine.Id(1).Delete(&user6)
	// UPDATE user SET ..., deleted_at = ? WHERE id = ?
	engine.Id(1).Get(&user6)
	// 再次调用Get，此时将返回false, nil，即记录不存在
	engine.Id(1).Delete(&user6)
	// 再次调用删除会返回0, nil，即记录不存在

	那么如果记录已经被标记为删除后，要真正的获得该条记录或者真正的删除该条记录，需要启用Unscoped，如下所示：
	var user7 User
	engine.Id(1).Unscoped().Get(&user7)
	// 此时将可以获得记录
	engine.Id(1).Unscoped().Delete(&user7)
	// 此时将可以真正的删除记录
	*/

	// 也可以直接执行一个SQL查询，即Select命令。当调用Query时，第一个返回值results为[]map[string][]byte的形式。
	sql := "SELECT * FROM t_user"
	results, err := engine.Query(sql)
	if err != nil {
		panic(err.Error())
	}
	log.Println(results)

	//也可以直接执行一个SQL命令，即执行Insert， Update， Delete 等操作。
	sql = "UPDATE `t_user` SET NAME = ? WHERE ID = ?"
	res, err := engine.Exec(sql, "Lsyl", 52)
	if err != nil {
		panic(err.Error())
	}
	affect_account, err := res.RowsAffected()	//返回两个结果int64, err
	if err != nil {
		panic(err.Error())
	}
	affect_account_str := strconv.FormatInt(affect_account, 10)
	log.Println("成功更新" + affect_account_str + "条记录~！")

	// add Commit() after all actions
	err = session.Commit()
	if err != nil {
		return
	}
}

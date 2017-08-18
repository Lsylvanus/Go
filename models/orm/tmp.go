package main

import (
	"ERP_1/log"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func getOrdersFromTmp(engine *xorm.Engine, request_id int64) {
	sql := "SELECT REQUEST_ID, SORT_ID, N01, /*status_code*/ N02, /*days_to_fulfill*/ N03, /*hours_to_fulfill*/ F01, /*order_total*/ F02, /*price*/ F03, /*cost*/ F04, /*shipping*/ F05, /*shipping_cost*/ T01, /*order_time*/ T02, /*created*/ C01, /*order_id*/ C02, /*transaction_id*/ C03, /*product_id*/ C04, /*variant_id*/ C05, /*buyer_id*/ C06, /*sku*/ C07, /*size*/ C08, /*color*/ C09, /*state*/ C10, /*user*/ C11, /*user_token*/ C12, /*product_name*/ C13, /*product_image_url*/ C14, /*name*/ C15, /*street_address1*/ C16, /*street_address2*/ C17, /*city*/ C18, /*detail_state*/ C19, /*country*/ C20, /*zipcode*/ C21, /*phone_number*/ C22 /*quantity*/ FROM tmp_wish_orders t WHERE t.REQUEST_ID = ?"
	results, err := engine.Query(sql, request_id)
	if err != nil {
		log.Log.Error("query sql :", err.Error())
		return
	}
	for k, v := range results {
		log.Log.Infof("=============== start : %v ===============", k)
		for ks, vs := range v {
			log.Log.Infof("field is : %v, value is : %v", ks, string(vs))
		}
		log.Log.Infof("================ end : %v ================\n", k)
	}
}

func insertOrdersToTmp(engine *xorm.Engine, request_id int64, sort_id int, status_code int) {
	sql := "INSERT INTO `tmp_wish_orders` ( REQUEST_ID, SORT_ID, N01, /*status_code*/ N02, /*days_to_fulfill*/ N03, /*hours_to_fulfill*/ F01, /*order_total*/ F02, /*price*/ F03, /*cost*/ F04, /*shipping*/ F05, /*shipping_cost*/ T01, /*order_time*/ T02, /*created*/ C01, /*order_id*/ C02, /*transaction_id*/ C03, /*product_id*/ C04, /*variant_id*/ C05, /*buyer_id*/ C06, /*sku*/ C07, /*size*/ C08, /*color*/ C09, /*state*/ C10, /*user*/ C11, /*user_token*/ C12, /*product_name*/ C13, /*product_image_url*/ C14, /*name*/ C15, /*street_address1*/ C16, /*street_address2*/ C17, /*city*/ C18, /*detail_state*/ C19, /*country*/ C20, /*zipcode*/ C21, /*phone_number*/ C22 /*quantity*/ ) SELECT ?, ?, w.status_code, w.days_to_fulfill, w.hours_to_fulfill, w.order_total, w.price, w.cost, w.shipping, w.shipping_cost, w.order_time, w.created, w.order_id, w.transaction_id, w.product_id, w.variant_id, w.buyer_id, w.sku, w.size, w.color, w.state, w.`user`, w.user_token, w.product_name, w.product_image_url, w.`name`, w.street_address1, w.street_address2, w.city, w.detail_state, w.country, w.zipcode, w.phone_number, w.quantity FROM `wish_orders` W WHERE W.status_code = ?"
	res, err := engine.Exec(sql, request_id, sort_id, status_code)
	if err != nil {
		log.Log.Error("exec sql :", err.Error())
		return
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		log.Log.Error("res lastInsertId :", err.Error())
		return
	}
	affStr := strconv.FormatInt(lastInsertId, 10)
	log.Log.Infof("insert into tmp, row begin : %s .", affStr)
}

func main() {
	log.LogInit()
	dbConn := flag.String("db", "root:chenghai3c@/wish_catch?charset=utf8&collation=utf8_general_ci&parseTime=true", "MySQL DB path") //zhusq
	flag.Parse()
	// 连接数据库
	var err error
	engine, err := xorm.NewEngine("mysql", *dbConn)
	if err != nil {
		log.Log.Panic(err)
	}
	err = engine.Ping()
	if err != nil {
		log.Log.Panic(err)
	}
	log.Log.Info("success to connect db")

	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Log.Error(err)
	}
	engine.DatabaseTZ = cst //设置数据库的系统时区
	time.Local = cst        //设置time的时区

	sql := "SELECT REQUEST_ID, SORT_ID, N01, /*status_code*/ N02, /*days_to_fulfill*/ N03, /*hours_to_fulfill*/ F01, /*order_total*/ F02, /*price*/ F03, /*cost*/ F04, /*shipping*/ F05, /*shipping_cost*/ T01, /*order_time*/ T02, /*created*/ C01, /*order_id*/ C02, /*transaction_id*/ C03, /*product_id*/ C04, /*variant_id*/ C05, /*buyer_id*/ C06, /*sku*/ C07, /*size*/ C08, /*color*/ C09, /*state*/ C10, /*user*/ C11, /*user_token*/ C12, /*product_name*/ C13, /*product_image_url*/ C14, /*name*/ C15, /*street_address1*/ C16, /*street_address2*/ C17, /*city*/ C18, /*detail_state*/ C19, /*country*/ C20, /*zipcode*/ C21, /*phone_number*/ C22 /*quantity*/ FROM tmp_wish_orders t WHERE t.REQUEST_ID = ?"
	resultMap, _ := QuerySql(engine, sql, 300)
	for k, re := range resultMap {
		log.Log.Infof("=============== start : %v ===============", k)
		for filed, value := range re {
			log.Log.Infof("filed is : %s, \t\t\tvalue is : %s.", filed, value)
		}
		log.Log.Infof("================ end : %v ================\n", k)
	}
}

func GetReqId(engine *xorm.Engine) (int64, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	reqId := fmt.Sprint(r.Intn(1000))

	sql := "SELECT DISTINCT t.REQUEST_ID FROM tmp_wish_orders t"
	results, err := engine.Query(sql)
	if err != nil {
		log.Log.Error("select request_id :", err.Error())
		return 0, err
	}

	var reqMap map[int]string
	reqMap = make(map[int]string)
	for k, v := range results {
		for _, vs := range v {
			reqMap[k] = string(vs)
		}
	}

	var str string
	for _, v := range reqMap {
		vs := v + " "
		str += vs
	}

	if !strings.Contains(str, reqId) {
		reqInt, err := strconv.ParseInt(reqId, 10, 0)
		if err != nil {
			log.Log.Error("parse int :", err.Error())
			return 0, err
		}
		return reqInt, err
	}

	return 0, err
}

//查找
func QuerySql(engine *xorm.Engine, sql string, args ...interface{}) ([]map[string]string, error) {
	results, err := engine.Query(sql, args...)
	if err != nil {
		log.Log.Error("query sql :", err.Error())
		return nil, err
	}
	ret := make([]map[string]string, 0)
	resultMap := make(map[string]string)
	for _, reMap := range results {
		for filed, value := range reMap {
			resultMap[filed] = string(value)
		}
		ret = append(ret, resultMap)
	}
	return ret, err
}

//插入
func InsSql(engine *xorm.Engine, sql string, args ...interface{}) (string, error) {
	res, err := engine.Exec(sql, args...)
	if err != nil {
		log.Log.Error("exec sql :", err.Error())
		return "0", err
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		log.Log.Error("insert last id :", err.Error())
		return "0", err
	}
	lastInsertIdStr := strconv.FormatInt(lastInsertId, 10)
	return lastInsertIdStr, err
}

//修改和删除
func ExecSql(engine *xorm.Engine, sql string, args ...interface{}) (string, error) {
	res, err := engine.Exec(sql, args...)
	if err != nil {
		log.Log.Error("exec sql :", err.Error())
		return "0", err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Log.Error("rows affected :", err.Error())
		return "0", err
	}
	rowsAffectedStr := strconv.FormatInt(rowsAffected, 10)
	return rowsAffectedStr, err
}
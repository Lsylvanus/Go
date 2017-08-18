package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/thinkboy/log4go"
	//iconv "src/github.com/djimenez/iconv-go"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var gLogger log4go.Logger
var engine *xorm.Engine
var session *xorm.Session

//多token情况
//tokens := map[string]string{"1": "60e8add79d20420c88f251846b229e5c", "2": "8ddac2fa42794124893947b9ff02e857"}
var tokens map[string]string
var resp *http.Response

const sqlLogFileName = "catch_sql.log"

//返回Json
type Retrieve struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    []Data `json:"data"` //cannot unmarshal array into Go value of type main.Order
	Paging  struct {
		Next     string `json:"next"`
		Previous string `json:"previous"`
	} `json:"paging"`
}

type Data struct {
	OrderWithDetail `json:"Order"` //cannot unmarshal object into Go value of type string
}

type OrderWithDetail struct {
	Orders
	//因json中Order与ShippingDetail有相同字段state，如该结构体中存在ShippingDetail，则Order中的state获取不了
	//ShippingDetail
}

//订单信息
type Orders struct {
	Id                           int            `xorm:"autoincr int(11) pk notnull"`
	OrderId                      string         `json:"order_id" xorm:"varchar(64) unique notnull"`
	TransactionId                string         `json:"transaction_id" xorm:"varchar(64) default ''"`
	ProductId                    string         `json:"product_id" xorm:"varchar(64) default ''"`
	VariantId                    string         `json:"variant_id" xorm:"varchar(64) default ''"`
	BuyerId                      string         `json:"buyer_id" xorm:"varchar(64) default ''"`
	Quantity                     string         `json:"quantity" xorm:"varchar(1)"`
	Sku                          string         `json:"sku" xorm:"varchar(64) default ''"`
	Size                         string         `json:"size" xorm:"varchar(32) default ''"`
	Color                        string         `json:"color" xorm:"varchar(16) default ''"`
	State                        string         `json:"state" xorm:"varchar(32) default ''"`
	ShippingProvider             string         `json:"shipping_provider" xorm:"varchar(64) default ''"`
	TrackingNumber               string         `json:"tracking_number" xorm:"varchar(64) default ''"`
	ShippedDateStr               string         `json:"shipped_date" xorm:"-"`
	ShippedDate                  time.Time      `xorm:"timestamp null"`
	ShipNote                     string         `json:"ship_note" xorm:"varchar(200) default ''"`
	LastUpdatedStr               string         `json:"last_updated" xorm:"-"`
	LastUpdated                  time.Time      `xorm:"timestamp null"`
	OrderTotal                   string         `json:"order_total" xorm:"float"`
	DaysToFulfill                string         `json:"days_to_fulfill" xorm:"int(4)"`
	HoursToFulfill               string         `json:"hours_to_fulfill" xorm:"int(6)"`
	ExpectedShipDateStr          string         `json:"expected_ship_date" xorm:"-"`
	ExpectedShipDate             time.Time      `xorm:"timestamp null"`
	Price                        string         `json:"price" xorm:"float"`
	StatusCode                   int            `xorm:"int(2) default 0 notnull"` //0:抓单;1:导出;2:物流号回流;3:物流号推送;4.导出异常;5.无邮寄代号;6.cancelled
	Cost                         string         `json:"cost" xorm:"float"`
	Shipping                     string         `json:"shipping" xorm:"float"`
	ShippingCost                 string         `json:"shipping_cost" xorm:"float"`
	ProductName                  string         `json:"product_name" xorm:"varchar(200) default ''"`
	ProductImageUrl              string         `json:"product_image_url" xorm:"varchar(500) default ''"`
	OrderTimeStr                 string         `json:"order_time" xorm:"-"`
	OrderTime                    time.Time      `xorm:"timestamp null"`
	ShippingDetail               ShippingDetail `json:"ShippingDetail"`
	RefundedBy                   string         `json:"refunded_by" xorm:"varchar(32) default ''"`
	RefundedTimeStr              string         `json:"refunded_time" xorm:"-"`
	RefundedTime                 time.Time      `xorm:"timestamp null"`
	RefundedReason               string         `json:"refunded_reason" xorm:"varchar(200) default ''"`
	IsWishExpress                string         `json:"is_wish_express" xorm:"varchar(16) default ''"`
	RequiresDeliveryConfirmation string         `json:"requires_delivery_confirmation" xorm:"varchar(16) default ''"`
	WeRequiredDeliveryDateStr    string         `json:"we_required_delivery_date" xorm:"-"`
	WeRequiredDeliveryDate       time.Time      `xorm:"timestamp null"`
	TrackingConfirmed            string         `json:"tracking_confirmed" xorm:"varchar(16) default ''"`
	TrackingConfirmedDateStr     string         `json:"tracking_confirmed_date" xorm:"-"`
	TrackingConfirmedDate        time.Time      `xorm:"timestamp null"`
	Name                         string         `xorm:"varchar(64) default ''"`
	StreetAddress1               string         `xorm:"varchar(200) default ''"`
	StreetAddress2               string         `xorm:"varchar(200) default ''"`
	City                         string         `xorm:"varchar(64) default ''"`
	DetailState                  string         `xorm:"varchar(32) default ''"`
	Country                      string         `xorm:"varchar(32) default ''"`
	Zipcode                      string         `xorm:"varchar(16) default ''"`
	PhoneNumber                  string         `xorm:"varchar(16) default ''"`
	User                         string         `xorm:"varchar(16) notnull"`
	UserToken                    string         `xorm:"varchar(64) notnull"`
	Created                      time.Time      `xorm:"timestamp created"`
	Updated                      time.Time      `xorm:"timestamp updated"`
	Version                      int            `xorm:"int(11) version"`
}

//发货详细
type ShippingDetail struct {
	//Vid            int    `xorm:"autoincr int(11) pk notnull"`
	Name           string `json:"name" xorm:"varchar(32) default ''"`
	StreetAddress1 string `json:"street_address1" xorm:"varchar(200) default ''"`
	StreetAddress2 string `json:"street_address2" xorm:"varchar(200) default ''"`
	City           string `json:"city" xorm:"varchar(32) default ''"`
	State          string `json:"state" xorm:"varchar(32) default ''"`
	Country        string `json:"country" xorm:"varchar(32) default ''"`
	Zipcode        string `json:"zipcode" xorm:"varchar(8) pk notnull"`
	PhoneNumber    string `json:"phone_number" xorm:"varchar(16) default ''"`
}

//接口地址参数
type Uri struct {
	Access_token    string
	Since           string
	Limit           string
	Start           string
	WishExpressOnly string
}

//用户token
type Users struct {
	Id       int    `xorm:"autoincr int(11) pk notnull"`
	Src      string `xorm:"varchar(5) default ''"`
	UserCode int    `xorm:"int(11)"`
	UserNick string `xorm:"varchar(20) default ''"`
	Email    string `xorm:"text"`
	Token    string `xorm:"text"`
	Mobile   string `xorm:"text"`
}

//错误信息返回
type ErrorMsg struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    int    `json:"data"`
}

//获取toten
func getToken() map[string]string {
	//多token情况
	//tokens := map[string]string{"1": "60e8add79d20420c88f251846b229e5c", "2": "8ddac2fa42794124893947b9ff02e857"}
	tokens = make(map[string]string)
	users := new(Users)
	total, err := engine.Where("src = ?", "wish").Count(users)
	if err != nil {
		fmt.Println(err.Error())
	}
	numStr := strconv.FormatInt(total, 10)
	gLogger.Info("there are " + numStr + " Wish users in database.")

	users1 := make([]Users, 0)
	err1 := engine.Where("src = ?", "wish").Find(&users1)
	if err1 != nil {
		fmt.Println(err1.Error())
	}
	//gLogger.Info(users2)
	for _, v := range users1 {
		tokens[v.UserNick] = v.Token
	}
	gLogger.Info("current user ：", tokens)
	return tokens
}

//获取起始时间
func getSinceDate() string {
	var date string

	//日期，从数据库中订单日期最新值为标准，拿出来的日期年月日为起始时间，包含有当日可能漏掉的数据
	sql := "SELECT DATE_FORMAT(t.last, '%Y-%m-%d') FROM ( SELECT DATE_SUB( MAX(o.order_time), INTERVAL 1 DAY ) last, MAX(o.order_time) later FROM erp_order o ) t"
	result, _ := engine.Query(sql)
	for _, v := range result {
		for _, vs := range v {
			gLogger.Info("since time ：", string(vs))
			//output, _ := iconv.ConvertString(vs, "unicode", "utf-8")
			date = string(vs)
		}
	}
	return date
}

// 从 API 中获取订单信息(JSON)，并解析成 struct
func fetchOrder(retrieve *Retrieve, start string, limit string, count int, wish_express_only string) {
	var order_b []byte
	//var error_b []byte
	var url string
	//var resp *http.Response

	uri := new(Uri)

	//获取起始时间
	date := getSinceDate()
	if date != "" {
		uri.Since = getSinceDate()
	} else {
		uri.Since = "2017-06-01"
	}

	//翻页，返回json数组数量限制大小为50
	const page int = 50

	//create table
	t_order := new(Orders)
	//查找表是否存在
	b, err := engine.IsTableExist(t_order)
	gLogger.Info("table is exist or not.")
	if !b {
		//创建表
		gLogger.Info("table is not exist, creating...")
		err := engine.CreateTables(t_order)
		if err != nil {
			gLogger.Debug("create table failed.\n", err.Error())
		}
		gLogger.Info("create table successful.")
	} else if err != nil {
		gLogger.Debug("err info : \n", err.Error())
	} else {
		// do nothing //
		gLogger.Info("table is exist.")
		// return
	}

	tokens := getToken()

	//从数据库中获取tokens
	for nick, token := range tokens {
		uri.Access_token = token
		url = "https://merchant.wish.com/api/v2/order/get-fulfill?start=" + start + "&limit=" + limit + "&since=" + uri.Since + "&wish_express_only=" + wish_express_only + "&access_token=" + uri.Access_token

		gLogger.Info("Api url is : \n", url)

		response, err := http.Get(url)
		if err != nil {
			gLogger.Debug(err.Error())
		}
		//defer resp.Body.Close()

		order_b, err = ioutil.ReadAll(response.Body)
		if err != nil {
			gLogger.Debug(err.Error())
		}
		gLogger.Info("get response url's body...")

		retrieve = &Retrieve{}
		//errorMsg = &ErrorMsg{}
		err = json.Unmarshal(order_b, retrieve) // JSON to Struct
		if err != nil {
			gLogger.Debug(err.Error())
		}
		/*错误信息返回
		err = json.Unmarshal(error_b, errorMsg)
		if err != nil {
			gLogger.Debug(err.Error())
		}*/

		data := retrieve.Data
		//没有数据直接返回
		/*避免多用户编历时某个用户没数据而不需要返回
		if len(data) == 0 {
			//关掉资源
			deinitAll()
			return
		}*/
		//fmt.Println(data)

		var order Orders
		var ship_detail ShippingDetail

		//遍历数组的data
		for i, d := range data {
			//fmt.Println(data[i])
			order = d.OrderWithDetail.Orders
			ship_detail = d.OrderWithDetail.Orders.ShippingDetail

			//gLogger.Info("order data : \n", order)
			//gLogger.Info("ship details : \n", ship_detail)
			gLogger.Debug("unmarshal json, current array len : %d ", i)
			//fmt.Printf("================\torder.State：%s\t ship.State：%s\t================\n", order.State, order.ShippingDetail.State)

			//为防止地址信息存在Unicode，特将Unicode转成中文
			//Unicode转码
			//output,_ := iconv.ConvertString(ship_detail.StreetAddress1, "unicode", "utf-8")
			//converter, _ := iconv.NewConverter("unicode", "utf-8")

			//保存json数据
			//logJSONdata(order, ship_detail, i)

			//避免order_id的重复，先查再判断
			var o Orders
			has, err := engine.Where("order_id = ?", order.OrderId).Get(&o)
			//err := engine.In("order_id", []string{order.OrderId}).Find(&order1)
			if err != nil {
				gLogger.Debug("error info : \n", err.Error())
			}

			//orderid不重复，插入，否则不插入
			gLogger.Info("current order is exist in database or not : ", has)
			if has {
				//gLogger.Info("there is a order named : " + order.OrderId + " in database. \n")
			} else {
				//两个结构体存进一个表wish_order中
				order.PhoneNumber = ship_detail.PhoneNumber
				order.City = ship_detail.City
				order.DetailState = ship_detail.State
				order.Name = ship_detail.Name
				order.Country = ship_detail.Country
				order.Zipcode = ship_detail.Zipcode
				order.StreetAddress1 = ship_detail.StreetAddress1
				order.StreetAddress2 = ship_detail.StreetAddress2

				//新增用户信息，用户名和token、状态码(默认为0)
				order.User = nick
				order.UserToken = token
				order.StatusCode = 0

				//日期处理
				timeLayout := "2006-01-02 15:04:05"
				jsonTimeLayout := "2006-01-02T15:04:05"
				loc, _ := time.LoadLocation("Local")

				//字段无返回，设置格式1970-01-01 08:00:01
				tm := time.Unix(1, 0)
				var timeStr string
				timeStr = tm.Format(timeLayout)

				//如果从接口拿到的json格式期间不为空，则需转换成time.Time格式存入数据库
				//否则把时间设置为timeStr
				if order.LastUpdatedStr != "" {
					lastUpdated, _ := time.ParseInLocation(jsonTimeLayout, order.LastUpdatedStr, loc)
					order.LastUpdated = lastUpdated
				} else {
					order.LastUpdated, _ = time.ParseInLocation(timeLayout, timeStr, loc)
				}

				if order.ShippedDateStr != "" {
					shippedDate, _ := time.ParseInLocation(jsonTimeLayout, order.ShippedDateStr, loc)
					order.ShippedDate = shippedDate
				} else {
					order.ShippedDate, _ = time.ParseInLocation(timeLayout, timeStr, loc)
				}

				if order.ExpectedShipDateStr != "" {
					expectedShipDate, _ := time.ParseInLocation(jsonTimeLayout, order.ExpectedShipDateStr, loc)
					order.ExpectedShipDate = expectedShipDate
				} else {
					order.ExpectedShipDate, _ = time.ParseInLocation(timeLayout, timeStr, loc)
				}

				if order.TrackingConfirmedDateStr != "" {
					trackingConfirmedDate, _ := time.ParseInLocation(jsonTimeLayout, order.TrackingConfirmedDateStr, loc)
					order.TrackingConfirmedDate = trackingConfirmedDate
				} else {
					order.TrackingConfirmedDate, _ = time.ParseInLocation(timeLayout, timeStr, loc)
				}
				if order.WeRequiredDeliveryDateStr != "" {
					weRequiredDeliveryDate, _ := time.ParseInLocation(jsonTimeLayout, order.WeRequiredDeliveryDateStr, loc)
					order.WeRequiredDeliveryDate = weRequiredDeliveryDate
				} else {
					order.WeRequiredDeliveryDate, _ = time.ParseInLocation(timeLayout, timeStr, loc)
				}

				if order.OrderTimeStr != "" {
					orderTime, _ := time.ParseInLocation(jsonTimeLayout, order.OrderTimeStr, loc)
					order.OrderTime = orderTime
				} else {
					order.OrderTime, _ = time.ParseInLocation(timeLayout, timeStr, loc)
				}

				if order.RefundedTimeStr != "" {
					refundedTime, _ := time.ParseInLocation(jsonTimeLayout, order.RefundedTimeStr, loc)
					order.RefundedTime = refundedTime
				} else {
					order.RefundedTime, _ = time.ParseInLocation(timeLayout, timeStr, loc)
				}

				//插入表
				_, err := engine.Insert(order)
				if err != nil {
					gLogger.Debug("error info : \n", err.Error())
				}
				gLogger.Info("insert into table success!")
				//affecttedStr := strconv.FormatInt(affected, 10)
				//gLogger.Info("成功，插入表wish_order " + affecttedStr + " 条记录~！\n")

			}
		}

		gLogger.Debug("current url's array len is : %d . compare to default len( %s ), be equivalent to %t", len(data), limit, len(data) == page)

		//超过默认的50条数据，可能还有数据，则翻页
		if len(data) == page {
			gLogger.Debug("change the start page...\n")
			start = strconv.Itoa(page*count + 1)
			limit = strconv.Itoa(page)

			count += 1 //翻页次数
			if start == strconv.Itoa(500+1) {
				//API提供的接口中最大的限制为500	Limit can range from 1 to 500 items and the default is 50
				gLogger.Debug("The maximum LIMIT(500) has been collected.")
				//超过限制返回
				return
			}
			//递归调用，起始位置数据不同
			fetchOrder(retrieve, start, limit, count, wish_express_only)
		}
	}

	//获取完成
}

//初始化
func initAll() {
	gLogger = nil
	gLogger.Info("init logger and database connection...")

	initLogger()
	conn()
	initSqlLogFile()
}

//关闭资源
func deinitAll() {
	err := session.Commit()
	if err != nil {
		return
	}
	defer session.Close()

	if gLogger != nil {
		return
	} else {
		gLogger = nil
		gLogger.Close()
	}
	gLogger.Info("deinit all resources...")
}

//初始化logger
func initLogger() {
	filenameOnly := getCurFileName()
	now := time.Now().Format("2006-01-02")
	var last string
	if strings.Contains(now, "-") {
		last = strings.Replace(now, "-", "", -1)
	}
	logFilename := filenameOnly + last + ".log"

	//控制台输出
	//gLogger.AddFilter("stdout", log4go.INFO, log4go.NewConsoleLogWriter())
	gLogger = log4go.NewDefaultLogger(log4go.INFO)

	//log文件输出
	if _, err := os.Stat(logFilename); err == nil {
		gLogger.Debug("There is an old file named : %s, rebuild at once...\n", logFilename)
		err := os.Remove(logFilename)
		if err != nil {
			gLogger.Debug(err.Error())
		}
	}
	gLogger.AddFilter("log", log4go.FINEST, log4go.NewFileLogWriter(logFilename, false))

	gLogger.Debug("current time is : %s", time.Now().Format("2006-01-02 15:04:05"))

	//return
}

//日志是一个接口，通过设置日志，可以显示SQL，警告以及错误等，默认的显示级别为INFO
func initSqlLogFile() {
	engine.ShowSQL(true)
	engine.Logger().SetLevel(core.LOG_DEBUG)

	//日志信息保存
	if _, err := os.Stat(sqlLogFileName); err == nil {
		gLogger.Debug("There is an old file named : %s, rebuild at once...\n", sqlLogFileName)
		os.Remove(sqlLogFileName)
	}
	f, err := os.Create(sqlLogFileName)
	if err != nil {
		println(err.Error())
	}
	engine.SetLogger(xorm.NewSimpleLogger(f))
	gLogger.Info("create sql log file to save sql exec statements...")
}

// 获取当前文件名称
func getCurFileName() string {
	_, fullFilename, _, _ := runtime.Caller(0)
	filenameWithSuffix := path.Base(fullFilename)
	fileSuffix := path.Ext(filenameWithSuffix)
	filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return filenameOnly
}

//log文件如果过大，清除
func rebuildDataLog() {
	dataLogName := "././" + getCurFileName() + ".log"

	data_b, err := ioutil.ReadFile(dataLogName)
	//not err == io.EOF
	if err != nil {
		fmt.Println(err.Error())
	}
	gLogger.Info(len(data_b))
	if len(data_b) > 1024*1024*2 {
		gLogger.Info("read file %s.log's memory ：%d，remove and rebuild it.", dataLogName, len(data_b))
		initLogger()
	} else {
		gLogger.Info("the file %s is too small to be removed.", dataLogName)
	}
}

func rebuildSqlLog() {
	sqlLogName := "././" + sqlLogFileName

	sql_b, err := ioutil.ReadFile(sqlLogName)
	//not err == io.EOF
	if err != nil {
		fmt.Println(err.Error())
	}
	gLogger.Info(len(sql_b))
	if len(sql_b) > 1024*1024*2 {
		gLogger.Info("read file %s.log's memory ：%d，remove and rebuild it.", sqlLogName, len(sql_b))
		initSqlLogFile()
	} else {
		gLogger.Info("the file %s is too small to be removed.", sqlLogName)
	}
}

//数据库初始化
func conn() {
	//创建Orm引擎
	var err error
	engine, err = xorm.NewEngine("mysql", "root:chenghai3c@/wish_catch?charset=utf8&collation=utf8_general_ci&parseTime=true")
	if err != nil {
		println(err.Error())
	}
	gLogger.Info("create orm engine, connecting database...")
	//测试连接
	errPing := engine.Ping()
	if errPing != nil {
		println(errPing.Error())
	}
	gLogger.Info("conn success!")
	// defer engine.Close()

	//当使用事务处理时，需要创建Session对象。在进行事物处理时，可以混用ORM方法和RAW方法
	session = engine.NewSession()

	// add Begin() before any action
	err1 := session.Begin()
	if err1 != nil {
		println(err1.Error())
	}
	gLogger.Info("session begins...")

	//连接池
	//最大打开连接数
	engine.SetMaxOpenConns(100)
	//连接池的空闲数大小
	engine.SetMaxIdleConns(20)

	// 映射规则
	// 默认对应的表名就变成了 wish_order 了，而之前默认的是 order
	tMapper := core.NewPrefixMapper(core.SnakeMapper{}, "wish_")
	engine.SetTableMapper(tMapper)

	//最后一个值得注意的是时区问题，默认xorm采用Local时区，所以默认调用的time.Now()会先被转换成对应的时区。
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

}

//保存json数据存进log
func logJSONdata(order Orders, detail ShippingDetail, i int) {
	gLogger.Debug("=============================== Data Order [%d] ==============================", i)
	gLogger.Debug("Order.order_id\t\t\t\t\t\t\t=%s", order.OrderId)
	gLogger.Debug("Order.transaction_id\t\t\t\t\t=%s", order.TransactionId)
	gLogger.Debug("Order.product_id\t\t\t\t\t\t=%s", order.ProductId)
	gLogger.Debug("Order.variant_id\t\t\t\t\t\t=%s", order.VariantId)
	gLogger.Debug("Order.buyer_id\t\t\t\t\t\t\t=%s", order.BuyerId)
	gLogger.Debug("Order.quantity\t\t\t\t\t\t\t=%s", order.Quantity)
	gLogger.Debug("Order.sku\t\t\t\t\t\t\t\t=%s", order.Sku)
	gLogger.Debug("Order.size\t\t\t\t\t\t\t\t=%s", order.Size)
	gLogger.Debug("Order.color\t\t\t\t\t\t\t\t=%s", order.Color)
	gLogger.Debug("Order.state\t\t\t\t\t\t\t\t=%s", order.State)
	gLogger.Debug("Order.shipping_provider\t\t\t\t\t=%s", order.ShippingProvider)
	gLogger.Debug("Order.tracking_number\t\t\t\t\t=%s", order.TrackingNumber)
	gLogger.Debug("Order.shipped_date\t\t\t\t\t\t=%s", order.ShippedDate)
	gLogger.Debug("Order.ship_note\t\t\t\t\t\t\t=%s", order.ShipNote)
	gLogger.Debug("Order.last_updated\t\t\t\t\t\t=%s", order.LastUpdated)
	gLogger.Debug("Order.order_total\t\t\t\t\t\t=%s", order.OrderTotal)
	gLogger.Debug("Order.days_to_fulfill\t\t\t\t\t=%s", order.DaysToFulfill)
	gLogger.Debug("Order.hours_to_fulfill\t\t\t\t\t=%s", order.HoursToFulfill)
	gLogger.Debug("Order.expected_ship_date\t\t\t\t=%s", order.ExpectedShipDate)
	gLogger.Debug("Order.price\t\t\t\t\t\t\t\t=%s", order.Price)
	gLogger.Debug("Order.cost\t\t\t\t\t\t\t\t=%s", order.Cost)
	gLogger.Debug("Order.shipping\t\t\t\t\t\t\t=%s", order.Shipping)
	gLogger.Debug("Order.shipping_cost\t\t\t\t\t\t=%s", order.ShippingCost)
	gLogger.Debug("Order.product_name\t\t\t\t\t\t=%s", order.ProductName)
	gLogger.Debug("Order.product_image_url\t\t\t\t\t=%s", order.ProductImageUrl)
	gLogger.Debug("Order.order_time\t\t\t\t\t\t=%s", order.OrderTime)
	gLogger.Debug("Order.shipping_detail\t\t\t\t\t=%s", order.ShippingDetail)
	gLogger.Debug("Order.refunded_by\t\t\t\t\t\t=%s", order.RefundedBy)
	gLogger.Debug("Order.refunded_time\t\t\t\t\t\t=%s", order.RefundedTime)
	gLogger.Debug("Order.refunded_reason\t\t\t\t\t=%s", order.RefundedReason)
	gLogger.Debug("Order.is_wish_express\t\t\t\t\t=%s", order.IsWishExpress)
	gLogger.Debug("Order.requires_delivery_confirmation\t=%s", order.RequiresDeliveryConfirmation)
	gLogger.Debug("Order.we_required_delivery_date\t\t\t=%s", order.WeRequiredDeliveryDate)
	gLogger.Debug("Order.tracking_confirmed\t\t\t\t=%s", order.TrackingConfirmed)
	gLogger.Debug("Order.tracking_confirmed_date\t\t\t=%s", order.TrackingConfirmedDate)
	gLogger.Debug("=============================== Data End [%d] ===============================\n", i)

	gLogger.Debug("=============================== Data Ship [%d] ==============================", i)
	gLogger.Debug("Order.name\t\t\t\t\t\t\t\t=%s", detail.Name)
	gLogger.Debug("Order.street_address1\t\t\t\t\t=%s", detail.StreetAddress1)
	gLogger.Debug("Order.street_address2\t\t\t\t\t=%s", detail.StreetAddress2)
	gLogger.Debug("Order.city\t\t\t\t\t\t\t\t=%s", detail.City)
	gLogger.Debug("Order.detail_state\t\t\t\t\t\t=%s", detail.State)
	gLogger.Debug("Order.country\t\t\t\t\t\t\t=%s", detail.Country)
	gLogger.Debug("Order.zipcode\t\t\t\t\t\t\t=%s", detail.Zipcode)
	gLogger.Debug("Order.phone_number\t\t\t\t\t\t=%s", detail.PhoneNumber)
	gLogger.Debug("=============================== Data End [%d] ===============================\n", i)
}

//main method
func main() {
	initAll()

	/*t := time.Tick(15 * time.Minute)
	i := 0
	for now := range t {
		gLogger.Info("程序执行时间：", now.Format("2006-01-02 15:04:05"))
		fetchOrder(new(Retrieve), "0", "50", 1, "")
		i++
		//执行5个小时
		if i > 19 {
			break
		}
	}*/

	//getToken()

	//rebuildDataLog()
	//rebuildSqlLog()
	//fetchOrder(new(Retrieve), "0", "50", 1, "")
	onTimer()

	deinitAll()

	//go语言中用log4go输出信息时有bug：只输出部分信息，甚至是无任何输出
	time.Sleep(100 * time.Millisecond)
}

func onTime(c <-chan time.Time) {
	for now := range c {
		// now := <- c
		fmt.Println("onTime", now)
		fetchOrder(new(Retrieve), "0", "50", 1, "")
	}
}

//定时器，每15分钟跑一次
func onTimer() {
	timer := time.NewTicker(15 * time.Minute)
	for {
		select {
		case <-timer.C:
			rebuildDataLog()
			rebuildSqlLog()
			fetchOrder(new(Retrieve), "0", "50", 1, "")
		}
	}
}

func startTimer(f func()) {
	go func() {
		for {
			f()
			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}

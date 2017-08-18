package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"net/http"
	"time"
)

type Refresh struct {
	Message string       `json:"message"`
	Code    int          `json:"code"`
	data    RefreshToken `json:"data"`
}

type RefreshToken struct {
	Id             int    `xorm:"autoincr int(11) pk notnull"`
	ExpiryTime     int64  `json:"expiry_time" xorm:"int(11)"`
	TokenType      string `json:"token_type" xorm:"varchar(16)"`
	Token          string `json:"access_token" xorm:"varchar(64)"`
	UserId         int    `xorm:"-"`
	ExpiresIn      int64  `json:"expires_in" xorm:"int(11)"`
	Reason         string `json:"reason" xorm:"-"`
	RefreshToken   string `json:"refresh_token" xorm:"varchar(64)"`
	ClientId       string `json:"client_id" xorm:"varchar(64)"`
	TokenId        string `json:"token_id" xorm:"varchar(64)"`
	MerchantId     string `xorm:"varchar(64)"`
	MerchantUserId string `json:"merchant_user_id" xorm:"varchar(64)"`
	ClientSecret   string `xorm:"varchar(64)"`
	GrantType      string `xorm:"varchar(16)"`
}

func JSONtoStruct(re *RefreshToken) *RefreshToken {
	var expiry []byte
	url := "https://merchant.wish.com/api/v2/oauth/refresh_token?client_id=" + re.ClientId + "&client_secret=" + re.ClientSecret + "&refresh_token=" + re.RefreshToken + "&grant_type=" + re.GrantType

	fmt.Println("Api url is : \n", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	expiry, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}
	resp.Body.Close()
	fmt.Println("get response url's body...")
	fmt.Println(expiry)

	re = &RefreshToken{}
	err = json.Unmarshal(expiry, re) // JSON to Struct
	if err != nil {
		fmt.Println(err.Error())
	}
	return re
}

func wishToken(id int64, engine *xorm.Engine){
	engine.Id(1).Get(&user)
}

func RefreshAccessToken(e *xorm.Engine) {
	re := new(RefreshToken)
	res := make([]RefreshToken, 0)

	err := e.Cols("expiry_time", "client_id", "client_secret", "grant_type", "refresh_token").Find(&res)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, v := range res {
		/*re.ClientId = v.ClientId
		re.ClientSecret = v.ClientSecret
		re.RefreshToken = v.RefreshToken
		re.GrantType = v.GrantType*/
		befE := v.ExpiryTime
		befT := time.Unix(befE, 0)

		nowT := time.Now()
		//判断刷新token的时间
		if befT.Year() == nowT.Year() {
			if befT.Month() == nowT.Month() {

			} else if befT.Month() < nowT.Month() {
				//之前
			}
		} else if befT.Year() < nowT.Year() {
		}

		//刷新后的时间戳

		//更新数据库中的token

	}
}

func main() {
	//创建Orm引擎
	engine, err := xorm.NewEngine("mysql", "root:chenghai3c@/wish_catch?charset=utf8&collation=utf8_general_ci&parseTime=true")
	if err != nil {
		println(err.Error())
	}
	fmt.Println("create orm engine, connecting database...")
	//测试连接
	errPing := engine.Ping()
	if errPing != nil {
		println(errPing.Error())
	}
	fmt.Println("conn success!")
	// defer engine.Close()

	//当使用事务处理时，需要创建Session对象。在进行事物处理时，可以混用ORM方法和RAW方法
	s := engine.NewSession()

	// add Begin() before any action
	err1 := s.Begin()
	if err1 != nil {
		println(err1.Error())
	}
	fmt.Println("session begins...")

	// 默认对应的表名就变成了 wish_order 了，而之前默认的是 order
	tMapper := core.NewPrefixMapper(core.SnakeMapper{}, "erp_")
	engine.SetTableMapper(tMapper)

	//最后一个值得注意的是时区问题，默认xorm采用Local时区，所以默认调用的time.Now()会先被转换成对应的时区。
	engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")

	RefreshAccessToken(engine)
}

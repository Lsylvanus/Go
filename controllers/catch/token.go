package main

import (
	"github.com/go-xorm/xorm"
	"time"
	"ERP_1/datas/dbo"
	"ERP_1/log"
	"encoding/json"
	"net/http"
	"io/ioutil"
)

//字节转string
func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

//刷新token
func WishToken(id int64, engine *xorm.Engine, cst *time.Location) {
	erpUser := new(dbo.ErpUsers)
	//获取erp_users中已知id的所有信息
	has, err := engine.Id(id).Get(erpUser)
	if err != nil {
		log.Log.Error(err.Error())
		return
	}
	if has {
		ref := new(dbo.Refresh)
		secretMap := make(map[int]string)
		var info []byte
		//数据库没有，定义map存储
		secretMap[0] = "6ea9f1b77b984ce9b3909705812f5466"
		secretMap[1] = "9d10867cd71a4878894f6dfab7e9874f"

		//固定参数refresh_token
		ref.Data.GrantType = "refresh_token"
		//读取返回json信息
		log.Log.Info(erpUser.Info)

		info = []byte(erpUser.Info)
		err := json.Unmarshal(info, ref)
		if err != nil {
			log.Log.Error(err.Error())
			return
		}
		if ref.Data.MerchantUserId == "574ff8333d174d5cbddd012b" {
			ref.Data.ClientSecret = secretMap[0]
		} else {
			ref.Data.ClientSecret = secretMap[1]
		}

		var expiry []byte
		url := "https://merchant.wish.com/api/v2/oauth/refresh_token?client_id=" + ref.Data.ClientId + "&client_secret=" + ref.Data.ClientSecret + "&refresh_token=" + ref.Data.RefreshToken + "&grant_type=" + ref.Data.GrantType
		log.Log.Info("refresh token url :", url)
		resp, err := http.Get(url)
		if err != nil {
			log.Log.Error(err.Error())
			return
		}
		defer resp.Body.Close()

		expiry, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Log.Error(err.Error())
			return
		}
		log.Log.Info("get response url's body...")

		// JSON to Struct
		err1 := json.Unmarshal(expiry, ref)
		if err1 != nil {
			log.Log.Error(err1.Error())
			return
		}

		//更新erp_users表的token信息
		users := new(dbo.ErpUsers)
		users.Token = ref.Data.Token
		users.Expire = time.Unix(ref.Data.ExpiryTime, 0)
		users.Info = byteString(expiry)

		//更新指定列
		_, err2 := engine.Id(id).Cols("token", "expire", "info").Update(users)
		if err2 != nil {
			log.Log.Error(err2.Error())
			return
		}
		log.Log.Info("update success!")
	} else {
		log.Log.Error("cannot found any rows in table.")
		return
	}
}
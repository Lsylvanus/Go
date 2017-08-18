package dbo

import "time"

//数据表user
type ErpUsers struct {
	Id            int64
	Src           string
	UserCode      int
	UserNick      string
	Email         string
	Token         string
	ClientSecret  string
	Mobile        string
	Expire        time.Time
	RefreshExpire time.Time
	Info          string
	Created       time.Time `xorm:"created"`
	Updated       time.Time `xorm:"updated"`
}

// 数据表erp_orders
type ErpOrders struct {
	Id             int64
	UserNick       string
	OrderId        int
	Src            string
	PaidAt         time.Time
	Status         string
	StatusCode     string
	TrackingNumber string
	Address1       string
	Address2       string
	City           string
	State          string
	Country        string
	TotalPrice     float64
	Created        time.Time `xorm:"created"`
	Updated        time.Time `xorm:"updated"`
}

// 数据表erp_order_details
type ErpOrderDetails struct {
	Id         int64
	Src        string
	StatusCode string
	OrderId    int
	Sku        string
	Count      int
	Created    time.Time `xorm:"created"`
}

// 数据表erp_sku_nick
type ErpSkuNicks struct {
	Id      int64
	Sku     string
	SkuNick string
	Created time.Time `xorm:"created"`
}

// 数据表erp_skus
type ErpSkus struct {
	Id      int64
	Sku     string
	Created time.Time `xorm:"created"`
}

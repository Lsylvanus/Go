package routers

import (
	"github.com/astaxie/beego"
	"sitepointgoapp/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/hello-world", &controllers.MainController{}, "get:HelloSitepoint")
	beego.Router("/hello-world/:id([0-9]+)", &controllers.MainController{}, "get,post:Get")
}
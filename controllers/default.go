package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (main *MainController) HelloSitepoint() {
	main.Data["Website"] = "Beego Test My site"
	main.Data["Email"] = "lightsylvanus@foxmail.com"
	main.Data["EmailName"] = "Lsylvanus"
	main.Data["Id"] = main.Ctx.Input.Param(":id")
	main.TplName = "default/hello-sitepoint.tpl"
}
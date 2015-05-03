package controllers

import (
	"github.com/astaxie/beego"
	//"merger_point/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
//	sql := "select t.record_time,sum(t.flow_rate) value from r_sector_flow_5 t where 1=1"+
//	" and record_time>= '2015-04-20 00:00:00' and record_time<'2015-04-23 00:00:00' and sp_code='cztv' group by record_time"
//	c.Data["data"] = models.Search(sql,1429459200000,1429718400000)
	c.TplNames = "index.tpl"
}

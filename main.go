package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/lisijie/goblog/models"
	_ "github.com/lisijie/goblog/routers"
	"github.com/lisijie/goblog/util"
)

func init() {

	//设置name=cache 的单例，这里的匿名函数为构造函数（传入到Factory），设置缓存
	util.Factory.Set("cache", func() (interface{}, error) {
		mc := util.NewLruCache(1000)
		return mc, nil
	})

	models.Init()
}

func main() {
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}

	beego.Run()
}

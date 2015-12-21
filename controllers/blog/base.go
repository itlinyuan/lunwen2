package blog

import (
	"github.com/astaxie/beego"
	"github.com/lisijie/goblog/models/option"
	"github.com/lisijie/goblog/util"
	"os"
	"strings"
)

type baseController struct {
	beego.Controller
	moduleName     string
	controllerName string
	actionName     string
	options        map[string]string
	cache          *util.LruCache
}

//在调用controller的方法时，如下面的display、get等，会先调用Prepare方法
func (this *baseController) Prepare() {
	//得到控制器和action名称
	controllerName, actionName := this.GetControllerAndAction()
	this.moduleName = "blog"
	this.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	this.actionName = strings.ToLower(actionName)
	this.options = option.GetOptions()
	this.Data["options"] = this.options
	cache, _ := util.Factory.Get("cache")
	this.cache = cache.(*util.LruCache)
}

//选择显示的主题是前台（default）还是后台（admin）
func (this *baseController) display(tpl string) {
	var theme string
	if v := this.getOption("theme"); v != "" {
		theme = v
	} else {
		theme = "default"
	}
	if _, err := os.Stat(beego.ViewsPath + "/" + theme + "/layout.html"); err == nil {
		this.Layout = theme + "/layout.html"
	}
	this.TplNames = theme + "/" + tpl + ".html"
}

func (this *baseController) getOption(name string) string {
	if v, ok := this.options[name]; ok {
		return v
	} else {
		return ""
	}
}

//系统设置信息
func (this *baseController) setHeadMetas(params ...string) {
	title_buf := make([]string, 0, 3)
	if len(params) == 0 && this.getOption("subtitle") != "" {
		title_buf = append(title_buf, this.getOption("subtitle"))
	}
	if len(params) > 0 {
		title_buf = append(title_buf, params[0])
	}
	title_buf = append(title_buf, this.getOption("sitename"))
	this.Data["title"] = strings.Join(title_buf, " - ")

	if len(params) > 1 {
		this.Data["keywords"] = params[1]
	} else {
		this.Data["keywords"] = this.getOption("keywords")
	}

	if len(params) > 2 {
		this.Data["description"] = params[2]
	} else {
		this.Data["description"] = this.getOption("description")
	}

}

package admin

import (
	"github.com/lisijie/goblog/models"
	"github.com/lisijie/goblog/util"
)

type MoodlistController struct {
	baseController
}

//管理
func (this *MoodlistController) List() {
	var (
		page     int
		pagesize int = 10
		offset   int //偏移量
		list     []*models.Moodlist
		moodlist models.Moodlist
	)

	if page, _ = this.GetInt("page"); page < 1 {
		page = 1
	}
	offset = (page - 1) * pagesize

	count, _ := moodlist.Query().Count()
	if count > 0 {
		moodlist.Query().OrderBy("-is_top", "-time").Limit(pagesize, offset).All(&list)
	}

	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/admin/moodlist/list", true).ToString()
	this.display()
}

//添加
func (this *MoodlistController) Add() {
	var (
		//id       int    = 0
		value    string = this.GetString("saytext")
		likes    int    = 0
		shits    int    = 0
		istop    int8   = 0
		moodlist models.Moodlist
	)

	//id, _ = this.GetInt("id")

	if this.GetString("istop") == "1" {
		istop = 1
	}

	moodlist.IsTop = istop
	moodlist.Value = value
	moodlist.Likes = likes
	moodlist.Shits = shits
	moodlist.Time = this.getTime()
	moodlist.Insert()

	this.Redirect("/admin/moodlist/list", 302)
}

//删除
func (this *MoodlistController) Delete() {
	id, _ := this.GetInt("id")
	moodlist := models.Moodlist{Id: id}
	if moodlist.Read() == nil {
		moodlist.Delete()
	}
	this.Redirect("/admin/moodlist/list", 302)
}

func (this *MoodlistController) UpdateIstop() {
	var moodlist models.Moodlist
	moodlist.Id, _ = this.GetInt("id")
	if this.GetString("istop") == "1" {
		moodlist.IsTop = 1
	} else {
		moodlist.IsTop = 0
	}

	moodlist.Update("is_top")
	this.Redirect("/admin/moodlist/list", 302)
}

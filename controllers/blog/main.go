package blog

import (
	"github.com/lisijie/goblog/models"
	"github.com/lisijie/goblog/util"
	"strconv"
	"strings"
)

type MainController struct {
	baseController
}

//首页, 只显示前N条
func (this *MainController) Index() {
	var (
		list     []*models.Post
		pagesize int
		err      error
		page     int
	)
	//Atoi:字符串->数字，默认为第一页
	if page, err = strconv.Atoi(this.Ctx.Input.Param(":page")); err != nil || page < 1 {
		page = 1
	}
	//默认每页有10篇文章
	if pagesize, err = strconv.Atoi(this.getOption("pagesize")); err != nil || pagesize < 1 {
		pagesize = 10
	}
	//过滤除这两个条件外的所有记录，status：0是已发布，1是草稿箱，2是回收站，
	//Filter 作为orm 的主要作用是过滤查询条件， 最终将会转换为sql 语句中的where 条件语句。
	//urltype:是不是文章的意思（跟心情和留言板做对比）
	query := new(models.Post).Query().Filter("status", 0).Filter("urltype", 0)
	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-istop", "-posttime").Limit(pagesize, (page-1)*pagesize).All(&list) //为list赋值
	}

	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "").ToString() //设置分页
	this.setHeadMetas()
	this.display("index")
}

//文章显示
func (this *MainController) Show() {
	var (
		post models.Post
		err  error
	)

	urlname := this.Ctx.Input.Param(":urlname")
	if urlname != "" {
		post.UrlName = urlname
		err = post.Read("urlname")
	} else {
		id, _ := strconv.Atoi(this.Ctx.Input.Param(":id"))
		post.Id = id
		err = post.Read()
	}
	if err != nil || post.Status != 0 {
		this.Abort("404") //内置404页面
		//this.Abort("500")
	}

	post.Views++
	post.Update("Views")
	//把post.content中出现的_ueditor_page_break_tag_替换为空，-1代表全部
	post.Content = strings.Replace(post.Content, "_ueditor_page_break_tag_", "", -1)

	this.Data["post"] = post
	this.setHeadMetas(post.Title, strings.Trim(post.Tags, ","), post.Title)
	this.display("article")
}

//历史归档
func (this *MainController) Archives() {
	var (
		page     int
		pagesize int
		err      error
		count    int64
		result   map[string][]*models.Post
	)

	if page, err = strconv.Atoi(this.Ctx.Input.Param(":page")); err != nil || page < 1 {
		page = 1
	}

	if pagesize, err = strconv.Atoi(this.getOption("pagesize")); err != nil || pagesize < 1 {
		pagesize = 20 //默认是文章的两倍，因为归档只显示标题
	} else {
		pagesize *= 2
	}

	query := new(models.Post).Query().Filter("status", 0).Filter("urltype", 0)

	count, _ = query.Count()
	result = make(map[string][]*models.Post) //map--key：年份,value:博文
	if count > 0 {
		var list []*models.Post
		query.OrderBy("-post_time").Limit(pagesize, (page-1)*pagesize).All(&list)
		//按年划分，从中遍历取出
		for _, v := range list {
			year := v.PostTime.Format("2006")
			if _, ok := result[year]; !ok {
				result[year] = make([]*models.Post, 0)
			}
			result[year] = append(result[year], v)
		}
	}

	this.Data["page"] = page
	this.Data["count"] = count
	this.Data["pagesize"] = pagesize
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/archives").ToString()
	this.Data["result"] = result
	this.setHeadMetas("归档")
	this.display("archives")
}

//心情列表
func (this *MainController) MoodList() {
	var (
		list     []*models.Moodlist
		pagesize int
		err      error
		page     int
	)

	if page, err = strconv.Atoi(this.Ctx.Input.Param(":page")); err != nil || page < 1 {
		page = 1
	}

	if pagesize, err = strconv.Atoi(this.getOption("pagesize")); err != nil || pagesize < 1 {
		pagesize = 20
	} else {
		pagesize *= 2
	}

	count, _ := new(models.Moodlist).Query().Count()
	if count > 0 {
		new(models.Moodlist).Query().OrderBy("-istop", "-time").Limit(pagesize, (page-1)*pagesize).All(&list) //为list赋值
	}

	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/moodlist").ToString() //设置分页
	this.setHeadMetas("心情列表")
	this.display("moodlist")
}

func (this *MainController) UpdateLikeAndShit() {

	var moodlist models.Moodlist
	id, _ := this.GetInt("id")
	status, _ := this.GetInt("status")
	moodlist.Id = id
	err1 := moodlist.Read()

	if err1 != nil {
		this.Abort("404") //内置404页面
	}
	//	id := 34
	//	moodlist := models.Moodlist{Id: id}
	if status == 0 {
		moodlist.Likes++
		moodlist.Update("Likes")
		//return moodlist.Likes
	} else if status == 1 {
		moodlist.Shits++
		moodlist.Update("Shits")
		//return moodlist.Shits
	}

	var (
		list     []*models.Moodlist
		pagesize int
		err      error
		page     int
	)

	if page, err = strconv.Atoi(this.Ctx.Input.Param(":page")); err != nil || page < 1 {
		page = 1
	}

	if pagesize, err = strconv.Atoi(this.getOption("pagesize")); err != nil || pagesize < 1 {
		pagesize = 20
	} else {
		pagesize *= 2
	}

	count, _ := new(models.Moodlist).Query().Count()
	if count > 0 {
		new(models.Moodlist).Query().OrderBy("-istop", "-time").Limit(pagesize, (page-1)*pagesize).All(&list) //为list赋值
	}

	this.Data["count"] = count
	this.Data["list"] = list
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, "/moodlist").ToString() //设置分页
	this.setHeadMetas("心情列表")
	this.display("moodlist")

	//返回json数据
	this.Data["json"] = moodlist
	this.ServeJson()
}

//分标签查看
func (this *MainController) Category() {
	var (
		page     int
		pagesize int
		name     string
		err      error
		count    int64
		result   map[string][]*models.Post
	)
	name = this.Ctx.Input.Param(":name")
	if page, err = strconv.Atoi(this.Ctx.Input.Param(":page")); err != nil || page < 1 {
		page = 1
	}
	if pagesize, err = strconv.Atoi(this.getOption("pagesize")); err != nil || pagesize < 1 {
		pagesize = 20
	} else {
		pagesize *= 2
	}

	tagpost := new(models.TagPost)
	tag := new(models.Tag)
	tag.Name = name

	if tag.Read("Name") != nil {
		this.Abort("404")
	}

	query := tagpost.Query().Filter("tag_id", tag.Id).Filter("post_status", 0)
	count, _ = query.Count()
	result = make(map[string][]*models.Post)
	if count > 0 {
		var tp []*models.TagPost
		var list []*models.Post
		var pids []int = make([]int, 0)

		query.OrderBy("-post_time").Limit(pagesize, (page-1)*pagesize).All(&tp)
		for _, v := range tp {
			pids = append(pids, v.PostId)
		}

		new(models.Post).Query().Filter("id__in", pids).All(&list)

		for _, v := range list {
			year := v.PostTime.Format("2006")
			if _, ok := result[year]; !ok {
				result[year] = make([]*models.Post, 0)
			}
			result[year] = append(result[year], v)
		}
	}

	this.Data["tag"] = tag
	this.Data["page"] = page
	this.Data["pagesize"] = pagesize
	this.Data["count"] = count
	this.Data["result"] = result
	this.Data["pagebar"] = util.NewPager(page, int(count), pagesize, tag.Link()).ToString()

	this.setHeadMetas(tag.Name, tag.Name, tag.Name)
	this.display("category")
}

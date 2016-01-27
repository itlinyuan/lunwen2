package models

import (
	"github.com/astaxie/beego/orm"
	"regexp"

	"time"
)

type Moodlist struct {
	Id     int
	UserId int       `orm:"index"`
	Value  string    `orm:"type(text)"`
	Time   time.Time `orm:"type(datetime);index"`
	Likes  int
	Shits  int
	IsTop  int8
}

func (m *Moodlist) TableName() string {
	return TableName("moodlist")
}

func (m *Moodlist) Insert() error {
	if _, err := orm.NewOrm().Insert(m); err != nil {
		return err
	}
	return nil
}

func (m *Moodlist) Read(fields ...string) error {
	if err := orm.NewOrm().Read(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Moodlist) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(m, fields...); err != nil {
		return err
	}
	return nil
}

func (m *Moodlist) Delete() error {

	return nil
}

func (m *Moodlist) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(m)
}

//把心情里出现表情文本的替换成表情
func (m *Moodlist) ReplaceValue() string {
	reg := regexp.MustCompile(`\[em_(\d+)\]`)
	regValue := reg.ReplaceAllString(m.Value, "<img src='/static/themes/admin/img/arclist/${1}.gif' border='0' />")
	//m.Value = strings.Replace(m.Value, "", "<img src='/static/themes/admin/img/arclist/1.gif' border='0' />", -1)

	return regValue
}

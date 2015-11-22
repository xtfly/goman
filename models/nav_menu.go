package models

import "github.com/astaxie/beego/orm"

// 导航菜单
type NavMenu struct {
	Id     int64  `json:"id" orm:"pk;auto"`          //ID
	Title  string `json:"tile" orm:"null;size(128)"` //
	Desc   string `json:"description" orm:"null"`    //
	Type   string `json:"type" orm:"null"`           //
	TypeId int64  `json:"type_id" orm:"default(0)"`  //
	Link   string `json:"link" orm:"null;index"`     //链接
	Icon   string `json:"icon" orm:"null"`           //图标
	Sort   int16  `json:"sort" orm:"default(0)"`     //排序
}

func init() {
	orm.RegisterModel(new(NavMenu))
}

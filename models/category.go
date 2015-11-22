package models

import "github.com/astaxie/beego/orm"

// 统一分类
type Category struct {
	Id       int64  `json:"id" orm:"pk;auto"`
	Pid      int64  `json:"parent_id" orm:"index;default(-1)"` // 父一级
	Title    string `json:"title" orm:"index"`                 // 分类名称
	Type     string `json:"type" orm:"null"`                   // 类型
	Icon     string `json:"icon" orm:"null"`                   // 图标
	UrlToken string `json:"url_token" orm:"null;index"`        // URL
	Sort     int16  `json:"sort" `                             // 排序
}

func init() {
	orm.RegisterModel(new(Category))
}

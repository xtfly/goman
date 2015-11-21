package models

import "github.com/astaxie/beego/orm"

//
type Category struct {
	Id       int64     `json:"id" orm:"pk;auto"`
	Title    string    `json:"title" orm:"index"` // 名
	Type     string    `json:"type" orm:"null"`   //
	Icon     string    `json:"icon" orm:"null"`   //
	Parent   *Category `json:"parent_id" orm:"rel(one);index;column(parent_id);null"`
	Sort     int16     `json:"sort" `
	UrlToken string    `json:"url_token" orm:"null;index"` //
}

func init() {
	orm.RegisterModel(new(Category))
}

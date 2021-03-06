package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// 统一草稿记录表
type Draft struct {
	Id       int64     `json:"id" orm:"pk;auto"`
	Uid      int64     `json:"uid" orm:"index"`                              //发布用户ID
	ItemId   int64     `json:"item_id" orm:"index"`                          //关联ID
	ItemType string    `json:"item_type" orm:"null;size(32)"`                //关联类型
	Content  string    `json:"content" orm:"type(text);null"`                //内容
	Time     time.Time `json:"time" orm:"auto_now_add;type(datetime);index"` //添加时间
}

func init() {
	orm.RegisterModel(new(Draft))
}

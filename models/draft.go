package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//
type Draft struct {
	Id      int64     `json:"id" orm:"pk;auto"`
	Uid     int64     `json:"uid" orm:"index"`                              //发布用户ID
	ItemId  int64     `json:"item_id" orm:"index"`                          //发布用户ID
	Type    string    `json:"type" orm:"null;size(32)"`                     //
	Content string    `json:"type" orm:"type(text);null"`                   //内容
	Time    time.Time `json:"time" orm:"auto_now_add;type(datetime);index"` //添加时间
}

func init() {
	orm.RegisterModel(new(Draft))
}

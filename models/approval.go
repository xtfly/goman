package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//
type Approval struct {
	Id  int64 `json:"id" orm:"pk;auto"`
	Uid int64 `json:"uid" orm:"index"` // 用户ID

	Type string    `json:"type" orm:"size(16);null"`
	Data string    `json:"type" orm:"type(text);null"`                   //内容
	Time time.Time `json:"time" orm:"auto_now_add;type(datetime);index"` //添加时间
}

func init() {
	orm.RegisterModel(new(Approval))
}

package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//
type ActiveData struct {
	Id             int64     `json:"id" orm:"pk;auto"`
	User           *Users    `json:"uid" orm:"rel(one);index;column(uid);null"`
	ExpireTime     string    `json:"expire_time" orm:"type(datetime);null"` //
	ActiveCode     string    `json:"active_code" orm:"null;size(32)"`       //
	AcitveTypeCode string    `json:"active_type_code" orm:"null;size(16)"`  //
	AddTime        time.Time `json:"add_time" orm:"type(datetime);auto_now_add"`
	AddIp          string    `json:"add_ip" orm:"null"` //
	ActiveTime     time.Time `json:"active_time" orm:"type(datetime);null"`
	ActiveIp       string    `json:"active_ip" orm:"null"` //
}

func init() {
	orm.RegisterModel(new(ActiveData))
}

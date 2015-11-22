package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//用户邀请记录
type Invitation struct {
	Id           int64     `json:"id" orm:"pk;auto"`                            //ID
	Uid          int64     `json:"uid" orm:"default(0);index"`                  //用户ID
	Code         string    `json:"invitation_code" orm:"size(32);null;index"`   //激活码
	Email        string    `json:"invitation_email" orm:"null;index"`           //激活email
	AddIp        string    `json:"add_ip" orm:"size(32);null"`                  //添加IP
	AddTime      time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"`  //添加时间
	ActiveExpire int16     `json:"active_expire" orm:"default(0)"`              //激活过期
	ActiveIp     string    `json:"active_ip" orm:"size(32);null;index"`         //激活IP
	ActiveTime   time.Time `json:"active_time" orm:"null;type(datetime);index"` //激活时间
	ActiveStatus int8      `json:"active_status" orm:"default(0);index"`        //1已使用0未使用-1已删除
	ActiveUid    int64     `json:"active_uid" orm:"default(0);index"`           //用户ID
}

func init() {
	orm.RegisterModel(new(Invitation))
}

package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// 用户关注表
type UserFollow struct {
	Id     int64  `json:"id" orm:"pk;auto"`
	Fans   *Users `json:"fans_uid" orm:"null;rel(one)"`   // 关注人的UID
	Friend *Users `json:"friend_uid" orm:"null;rel(one)"` // 被关注人的uid

	AddTime time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` // 添加时间
	Ip      string    `json:"ip" valid:"IP"`                              // 客户端ip

	ActiveUrl string `json:"active_url"` // 停留页面
	UserAgent string `json:"user_agent"` // 用户客户端信息
}

func init() {
	orm.RegisterModel(new(UserFollow))
}

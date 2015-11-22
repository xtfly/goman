package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// 在线用户列表
type UserOnline struct {
	Id         int64     `json:"id" orm:"pk;auto"`
	Uid        int64     `json:"uid" orm:"default(0);index"`                          // 用户 ID
	LastActive time.Time `json:"last_active" orm:"auto_now_add;type(datetime);index"` // 上次活动时间
	Ip         string    `json:"ip" valid:"IP"`                                       // 客户端ip
	ActiveUrl  string    `json:"active_url"`                                          // 停留页面
	UserAgent  string    `json:"user_agent"`                                          // 用户客户端信息
}

func init() {
	orm.RegisterModel(new(UserOnline))
}

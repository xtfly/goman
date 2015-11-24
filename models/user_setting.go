package models

import "github.com/astaxie/beego/orm"

// 配置
type UserSetting struct {
	Uid                 int64  `json:"uid" orm:"default(0);pk"`
	DefTimeZone         string `json:"default_timezone" orm:"null"`
	EmailSetting        string `json:"email_settings" orm:"null"`
	NotificationSetting string `json:"Notification_settings" orm:"null"`
}

var usersetting = new(UserSetting)

func init() {
	orm.RegisterModel(usersetting)
}

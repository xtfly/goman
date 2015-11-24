package models

import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/xtfly/goman/kits"
)

//
type ActiveData struct {
	Id             int64     `json:"id" orm:"pk;auto"`                           //
	User           *Users    `json:"uid" orm:"rel(one);index;column(uid);null"`  //
	ExpireTime     int64     `json:"expire_time" orm:"default(0)"`               //
	ActiveCode     string    `json:"active_code" orm:"null;size(32)"`            //
	AcitveTypeCode string    `json:"active_type_code" orm:"null;size(16)"`       //
	AddTime        time.Time `json:"add_time" orm:"type(datetime);auto_now_add"` //
	AddIp          string    `json:"add_ip" orm:"null"`                          //
	ActiveTime     time.Time `json:"active_time" orm:"type(datetime);null"`      //
	ActiveIp       string    `json:"active_ip" orm:"null"`                       //
}

func init() {
	orm.RegisterModel(new(ActiveData))
}

func NewValidByEmail(t *Transaction, uid int64, email string) bool {
	if uid <= 0 {
		return false
	}

	m := &ActiveData{
		User:           &Users{Id: uid},
		ExpireTime:     time.Now().Unix() + 60*60*24,
		ActiveCode:     kits.GenHashStr(email, 4),
		AcitveTypeCode: "VALID_EMAIL",
	}

	_, ok := t.Insert(m)
	// TODO 发送验证邮件
	return ok
}

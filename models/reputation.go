package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Reputation struct {
	Id            int64     `json:"id" orm:"pk;auto"`                      //ID
	User          *Users    `json:"uid" orm:"null;rel(one)"`               //用户UID
	Updated       time.Time `json:"update_time" orm:"type(datetime);null"` //更新
	AgreeCount    int       `json:"agree_count"`                           //赞同数量
	ThanksCount   int       `json:"thanks_count"`                          //感谢数量
	QuestionCount int       `json:"question_count"`                        //问题数量
}

// 多字段唯一键
func (m *Reputation) TableUnique() [][]string {
	return [][]string{
		[]string{"Id", "User"},
	}
}

func init() {
	orm.RegisterModel(new(Reputation))
}

package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type ReputationCatagory struct {
	Id            int64     `json:"id" orm:"pk;auto"`                      //ID
	Uid           int64     `json:"uid" orm:"default(0)"`                  //用户UID
	CatId         int64     `json:"category_id" orm:"default(0)"`          //category_id
	Updated       time.Time `json:"update_time" orm:"type(datetime);null"` //更新
	Reputation    int       `json:"reputation" orm:"default(0);index"`     //
	AgreeCount    int       `json:"agree_count" orm:"default(0)"`          //赞同数量
	ThanksCount   int       `json:"thanks_count" orm:"default(0)"`         //感谢数量
	QuestionCount int       `json:"question_count" orm:"default(0)"`       //问题数量
}

// 多字段唯一键
func (m *ReputationCatagory) TableUnique() [][]string {
	return [][]string{
		[]string{"CatId", "Uid"},
	}
}

//
type ReputationTopic struct {
	Id          int64     `json:"id" orm:"pk;auto"`                      //ID
	Uid         int64     `json:"uid" orm:"default(0);index"`            //用户UID
	TopicId     int64     `json:"topic_id" orm:"default(0);index"`       //用户UID
	Updated     time.Time `json:"update_time" orm:"type(datetime);null"` //更新
	Reputation  int       `json:"reputation" orm:"default(0);index"`     //
	TopicCount  int       `json:"topic_count" orm:"default(0)"`          //威望问题话题计数
	AgreeCount  int       `json:"agree_count" orm:"default(0)"`          //赞同数量
	ThanksCount int       `json:"thanks_count" orm:"default(0)"`         //感谢数量
}

func init() {
	orm.RegisterModel(new(ReputationCatagory), new(ReputationTopic))
}

package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	ILA_NewQuestion int8 = iota
	ILA_AnswerQuestion
	ILA_QuestionAnswer
	ILA_InviteAnswer
	ILA_AnswerInvite
	ILA_ThanksQuestion
	ILA_QuestionThanks
	ILA_ThanksAnswer
	ILA_AnswerThanks
	ILA_AnswerFold
	ILA_BestAnswer
	ILA_Invite
)

// 积分记录表
type IntegralLog struct {
	Id         int64     `json:"id" orm:"pk;auto"`                                    //ID
	Uid        int64     `json:"uid" orm:"default(0);index"`                          //用户ID
	Action     int8      `json:"action" orm:"default(0);index"`                       //
	Integral   int       `json:"integral" orm:"default(0);index"`                     //
	Note       string    `json:"note" orm:"size(128);null"`                           //
	Balance    int       `json:"balance" orm:"default(0)"`                            //
	ItemId     int64     `json:"item_id" orm:"default(0);index"`                      //
	ActiveTime time.Time `json:"active_time" orm:"auto_now_add;type(datetime);index"` //时间

}

func init() {
	orm.RegisterModel(new(IntegralLog))
}

// 增加一条积分日志
func (m *IntegralLog) Add(t *Transaction) bool {
	if m.Integral == 0 {
		return true
	}

	if !UserExistedById(m.Uid) {
		return false
	}

	if _, ok := t.Insert(m); !ok {
		return false
	}

	// 异步处理
	go func() {
		SumUserIntegral(NewTr(), m.Uid)
	}()

	return true
}

// 统计用户的积分
func SumUserIntegral(t *Transaction, uid int64) bool {
	integral, ok := t.Sum("integral_log", "integral", "uid", uid)
	if !ok {
		return false
	}
	return t.UpdateById("Users", uid, orm.Params{"Integral": integral})
}

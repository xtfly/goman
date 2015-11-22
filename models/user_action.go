package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

const (
	_ int8 = iota
	AssTypeQuestion
	AssTypeAnswer
	AssTypeComment
	AssTypeTopic
)

// 用户操作记录
type UserActionHistory struct {
	Id          int64     `json:"id" orm:"pk;auto"`
	Uid         int64     `json:"uid" orm:"default(0);index"`                 //UID
	AssType     int8      `json:"associate_type" orm:"default(0)"`            //关联类型: 1 问题 2 回答 3 评论 4 话题
	AssAction   int16     `json:"associate_action" orm:"default(0)"`          //操作类型
	AssId       int64     `json:"associate_id" orm:"default(0);index"`        //关联ID
	AssAttached int64     `json:"associate_attached" orm:"default(0);index"`  //关联
	AddTime     time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` //添加时间
	Anonymous   bool      `json:"anonymous" orm:"default(0)"`                 //是否匿名
	FoldStatus  bool      `json:"fold_status" orm:"default(0)"`               //
}

// 多字段索引
func (u *UserActionHistory) TableIndex() [][]string {
	return [][]string{
		[]string{"AssType", "AssAction"},
		[]string{"AssId", "AssType", "AssAction"},
		[]string{"Uid", "AssType", "AssAction"},
	}
}

type UserActionHistoryData struct {
	Id      int64              `json:"id" orm:"pk;auto"`
	History *UserActionHistory `json:"hid" orm:"null;rel(one)"` // OneToOne relation

	AssContext  string `json:"associate_content" orm:"type(text);null"`  // 关联
	AssAttached int64  `json:"associate_attached" orm:"type(text);null"` //
	AddonData   int64  `json:"associate_attached" orm:"type(text);null"` // 附加数据
}

type UserActionnHistoryFresh struct {
	Id            int64              `json:"id" orm:"pk;auto"`
	ActionHistory *UserActionHistory `json:"hid" orm:"null;rel(one)"` // OneToOne relation

	AssId     int64 `json:"associate_id" orm:"default(0);index"` // 关联ID
	AssType   int8  `json:"associate_type" orm:"default(0)"`     // 关联类型: 1 问题 2 回答 3 评论 4 话题
	AssAction int16 `json:"associate_action" orm:"default(0)"`   // 操作类型

	AddTime   time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` // 添加时间
	Uid       int64     `json:"uid" orm:"default(0);index"`                 //UID
	Anonymous bool      `json:"anonymous"`                                  // 是否匿名
}

// 多字段索引
func (u *UserActionnHistoryFresh) TableIndex() [][]string {
	return [][]string{
		[]string{"AssType", "AssAction"},
		[]string{"AssId", "AssType", "AssAction"},
		[]string{"Uid", "AssType", "AssAction"},
	}
}

func init() {
	orm.RegisterModel(new(UserActionHistory), new(UserActionHistoryData), new(UserActionnHistoryFresh))
}

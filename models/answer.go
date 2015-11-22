package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//回答列表
type Answer struct {
	Id       int64     `json:"id" orm:"pk;auto"`
	Qid      int64     `json:"qid" orm:"index"`                                      // 问题id
	Uid      int64     `json:"uid" orm:"index"`                                      // 发布问题用户ID
	Category *Category `json:"category_id" orm:"rel(one);index;column(category_id)"` // 分类id

	AgainstCount      int    `json:"against_count" orm:"default(0);index"`      //反对人数
	AgreeCount        int    `json:"agree_count" orm:"default(0);index"`        //支持人数
	CommentCount      int    `json:"comment_count" orm:"default(0)"`            //评论总数
	UnInterestedCount int    `json:"uninterested_count" orm:"default(0);index"` //不感兴趣
	ThanksCount       int    `json:"thanks_count" orm:"default(0);index"`       //感谢数量
	Ip                string `json:"ip" orm:"null;size(32)"`                    //
	HasAttach         bool   `json:"has_attach" orm:"default(0)"`               //是否存在附件
	ForceFold         bool   `json:"force_fold" orm:"default(0)"`               //强制折叠
	Anonymous         bool   `json:"anonymous" orm:"default(0);index"`          //

	PublishSource string    `json:"publish_source" orm:"size(16);null;index"`         //回答内容
	Content       string    `json:"type" orm:"type(text);null"`                       //回答内容
	AddTime       time.Time `json:"add_time" orm:"auto_now_add;type(datetime);index"` //添加时间
}

//回答评论列表
type AnswerComments struct {
	Id     int64   `json:"id" orm:"pk;auto"`
	Uid    int64   `json:"uid" orm:"index"`                // 用户ID
	AtUid  int64   `json:"at_uid" orm:"index"`             // At用户ID
	Answer *Answer `json:"answer_id" orm:"index;rel(one)"` //

	Message string    `json:"message" orm:"type(text);null"`                    //内容
	AddTime time.Time `json:"add_time" orm:"auto_now_add;type(datetime);index"` //添加时间
}

//回答不感兴趣表
type AnswerUninterested struct {
	Id     int64   `json:"id" orm:"pk;auto"`
	Uid    int64   `json:"uid" orm:"index"`                     // 用户ID
	Answer *Answer `json:"answer_id" orm:"index;rel(one);null"` //

	UserName string    `json:"user_name" orm:"size(255);null"`         //
	Time     time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间
}

//回答感谢列表
type AnswerThanks struct {
	Id     int64   `json:"id" orm:"pk;auto"`
	Uid    int64   `json:"uid" orm:"index"`                     // 用户ID
	Answer *Answer `json:"answer_id" orm:"index;rel(one);null"` //

	UserName string    `json:"user_name" orm:"size(255);null"`         //
	Time     time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间
}

//回答赞成列表
type AnswerVote struct {
	Id         int64 `json:"id" orm:"pk;auto"`
	VoteUid    int64 `json:"vote_uid" orm:"index"`         // 用户ID
	AnswerId   int64 `json:"answer_id" orm:"index"`        //
	AnswerUid  int64 `json:"answer_uid" orm:"index"`       //
	RepuFactor int   `json:"repu_factor" orm:"default(0)"` //
	VoteValue  int16 `json:"vote_value" orm:"default(0)"`  //'-1反对 1 支持

	AddTime time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` //添加时间
}

func init() {
	orm.RegisterModel(new(Answer), new(AnswerComments),
		new(AnswerUninterested), new(AnswerThanks),
		new(AnswerVote))
}

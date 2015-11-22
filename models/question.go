package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//问题列表
type Question struct {
	Id  int64 `json:"id" orm:"pk;auto"` //ID
	Uid int64 `json:"uid" orm:"index"`  //发布用户UID

	AddTime    time.Time `json:"add_time" orm:"type(datetime);auto_now_add;index"` //
	UpdateTime time.Time `json:"update_time" orm:"type(datetime);null;index"`      //

	AnswerCount  int `json:"answer_count" orm:"default(0);index"`  //回答计数
	AnswerUsers  int `json:"answer_users" orm:"default(0);index"`  //回答人数
	ViewCount    int `json:"view_count" orm:"default(0)"`          //浏览次数
	FocusCount   int `json:"focus_count" orm:"default(0)"`         //关注数
	CommentCount int `json:"comment_count" orm:"default(0)"`       //评论数
	AgreeCount   int `json:"agree_count" orm:"default(0);index"`   //回复赞同数总和
	AgainstCount int `json:"against_count" orm:"default(0);index"` //回复反对数总和
	ThanksCount  int `json:"thanks_count" orm:"default(0);index"`  //感谢数量

	ActionHistoryId int64     `json:"best_answer" orm:"defualt(0)"`                         //动作的记录表的关连id
	BestAnswer      int64     `json:"best_answer" orm:"defualt(0);index"`                   //最佳回复ID
	LastAnswer      int64     `json:"last_answer" orm:"defualt(0);index"`                   //最后回答ID
	Category        *Category `json:"category_id" orm:"rel(one);column(category_id);index"` //
	RecvEmailId     int64     `json:"received_email_id" orm:"defualt(0);index"`

	PopularValue       int64 `json:"popular_value" orm:"defualt(0);index"`                 //
	PopularValueUpdate int64 `json:"popular_value_update" orm:"type(datetime);null;index"` //

	HasAttach   bool  `json:"has_attach" orm:"default(0);index"` // 是否存在附件
	Lock        bool  `json:"lock" orm:"default(0);index"`       // 是否锁定
	IsRecommend bool  `json:"is_recommend" orm:"default(0)"`     //
	Anonymous   bool  `json:"anonymous" orm:"default(0);index"`  //
	ChapterId   int64 `json:"chapter_id" orm:"index"`            //
	Sort        int16 `json:"sort" orm:"index"`                  //

	Ip              string `json:"ip" orm:"null;size(32)"`                 //
	Content         string `json:"context" orm:"null"`                     //问题内容
	Detail          string `json:"detail" orm:"type(text);null"`           //问题说明
	ContentFullText string `json:"context_fulltext" orm:"type(text);null"` //

	UnverifiedModify      string `json:"unverified_modify" orm:"type(text);null"`        //
	UnverifiedModifyCount int64  `json:"unverified_modify_count" orm:"defualt(0);index"` //
}

//问题评论列表
type QuestionComments struct {
	Id       int64     `json:"id" orm:"pk;auto"`
	Uid      int64     `json:"uid" orm:"index"`                                  // 用户ID
	Question *Question `json:"-" orm:"index;rel(one);column(question_id)"`       //
	Message  string    `json:"message" orm:"type(text);null"`                    //内容
	AddTime  time.Time `json:"add_time" orm:"auto_now_add;type(datetime);index"` //添加时间
}

//邀请问答
type QuestionInvite struct {
	Id           int64     `json:"id" orm:"pk;auto"`
	Question     *Question `json:"-" orm:"index;rel(one);column(question_id)"` //
	SenderUid    int64     `json:"sender_uid" orm:"index"`                     //
	RecipientUid int64     `json:"recipient_uid" orm:"index"`                  //
	Email        string    `json:"email" orm:"size(128);null;index"`           //受邀Email

	AddTime       time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` //添加时间
	AvailableTime time.Time `json:"available_time" orm:"null;type(datetime)"`   //生效时间
}

//问题感谢列表
type QuestionThanks struct {
	Id       int64     `json:"id" orm:"pk;auto"`
	Uid      int64     `json:"uid" orm:"index"`                            // 用户ID
	Question *Question `json:"-" orm:"index;rel(one);column(question_id)"` //
	UserName string    `json:"user_name" orm:"size(255);null"`             //
	Time     time.Time `json:"time" orm:"auto_now_add;type(datetime)"`     //添加时间
}

//问题关注表
type QuestionFocus struct {
	Id       int64     `json:"id" orm:"pk;auto"`
	Uid      int64     `json:"uid" orm:"index"`                            // 用户ID
	Question *Question `json:"-" orm:"index;rel(one);column(question_id)"` //
	Time     time.Time `json:"time" orm:"auto_now_add;type(datetime)"`     //添加时间
}

//问题不感兴趣表
type QuestionUninterested struct {
	Id       int64     `json:"id" orm:"pk;auto"`
	Uid      int64     `json:"uid" orm:"index"`                            // 用户ID
	Question *Question `json:"-" orm:"index;rel(one);column(question_id)"` //
	UserName string    `json:"user_name" orm:"size(255);null"`             //
	Time     time.Time `json:"time" orm:"auto_now_add;type(datetime)"`     //添加时间
}

func init() {
	orm.RegisterModel(new(Question), new(QuestionThanks), new(QuestionComments),
		new(QuestionInvite), new(QuestionFocus), new(QuestionUninterested))
}

package models

//
import (
	"time"

	"github.com/astaxie/beego/orm"
)

//话题表
type Topic struct {
	Id      int64     `json:"id" orm:"pk;auto"`                       //ID
	Pid     int64     `json:"parent_id" orm:"default(-1);index"`      //
	Title   string    `json:"title" orm:"size(128);unique"`           //话题标题
	Desc    string    `json:"description" orm:"size(255);null"`       //话题描述
	Picture string    `json:"picture" orm:"size(255);null"`           //话题图片
	AddTime time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间

	Lock                 bool      `json:"lock" orm:"default(0)"`                //话题是否锁定
	DisussCount          int       `json:"discuss_count" orm:"default(0);index"` //讨论计数
	DisussCountLastWeek  int       `json:"discuss_count" orm:"default(0);index"` //讨论计数
	DisussCountLastMonth int       `json:"discuss_count" orm:"default(0);index"` //讨论计数
	DisussCountUpdate    time.Time `json:"time" orm:"type(datetime);null"`       //
	FocusCount           int       `json:"focus_count" orm:"default(0);index"`   //关注计数
	UserRelated          bool      `json:"user_related" orm:"default(0);index"`  //是否被用户关联
	IsParent             bool      `json:"is_parent" orm:"default(0)"`           //是否

	UrlToken string `json:"url_token" orm:"size(255);null"` //
	SeoTitle string `json:"seo_title" orm:"size(255);null"` //
	MergedId int64  `json:"merged_id" orm:"default(0)"`     //

}

//话题关注表
type TopicFocus struct {
	Id  int64 `json:"id" orm:"pk;auto"`           //ID
	Tid int64 `json:"topic_id" orm:"index"`       //
	Uid int64 `json:"uid" orm:"default(0);index"` //话题标题

	AddTime time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间
}

//话题合并表
type TopicMerge struct {
	Id       int64 `json:"id" orm:"pk;auto"`                 //ID
	SourceId int64 `json:"source_id" orm:"default(0);index"` //
	TargetId int64 `json:"target_id" orm:"default(0);index"` //
	Uid      int64 `json:"uid" orm:"default(0);index"`       //

	AddTime time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间
}

//话题关联表
type TopicRelation struct {
	Id     int64  `json:"id" orm:"pk;auto"`               //ID
	Tid    int64  `json:"topic_id" orm:"index"`           //
	ItemId int64  `json:"item_id" orm:"default(0);index"` //
	Uid    int64  `json:"uid" orm:"default(0);index"`     //
	Type   string `json:"type" orm:"size(16);null"`

	AddTime time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间
}

func init() {
	orm.RegisterModel(new(Topic), new(TopicFocus), new(TopicMerge), new(TopicRelation))
}

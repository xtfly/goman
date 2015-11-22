package models

//
import (
	"time"

	"github.com/astaxie/beego/orm"
)

//
type PostIndex struct {
	Id           int64     `json:"id" orm:"pk;auto"`                               //ID
	Uid          int64     `json:"uid" orm:"default(0);index"`                     //
	PostId       int64     `json:"post_id" orm:"default(0);index"`                 //
	PostType     string    `json:"post_type" orm:"size(16);null"`                  //
	AddTime      time.Time `json:"time" orm:"auto_now_add;type(datetime)"`         // 添加时间
	UpdateTime   time.Time `json:"update_time" orm:"auto_now;type(datetime);null"` // 添加时间
	CateId       int64     `json:"category_id" orm:"default(0);index"`             //
	IsRecommend  bool      `json:"is_recommend" orm:"default(0);index"`            //
	Viewcount    int       `json:"view_count" orm:"default(0);index"`              //
	Agreecount   int       `json:"agree_count" orm:"default(0);index"`             //
	Answercount  int       `json:"answer_count" orm:"default(0);index"`            //
	PopularValue int       `json:"popular_value" orm:"default(0)"`                 //
	Anonymous    bool      `json:"anonymous" orm:"default(0)"`                     //
	Lock         bool      `json:"lock" orm:"default(0)"`                          //
}

func init() {
	orm.RegisterModel(new(PostIndex))
}

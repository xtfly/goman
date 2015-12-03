package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Article struct {
	Id          int64     `json:"id" orm:"pk;auto"`                           //ID
	Uid         int64     `json:"uid" orm:"index"`                            //
	Title       string    `json:"title" orm:""`                               //
	CatId       int64     `json:"category_id" orm:"null"`                     //
	Comments    int64     `json:"comments" orm:"default(0);index"`            //
	Views       int64     `json:"views" orm:"default(0);index"`               //
	Votes       int64     `json:"votes" orm:"default(0)"`                     //
	AddTime     time.Time `json:"add_time" orm:"type(datetime);auto_now_add"` //
	HasAttach   bool      `json:"has_attach" orm:"default(0);index"`          //
	Lock        bool      `json:"lock" orm:"default(0);index"`                //
	IsRecommend bool      `json:"is_recommend" orm:"default(0)"`              //
	ChapterId   int64     `json:"chapter_id" orm:"index"`                     //
	Sort        int16     `json:"sort" orm:"index"`                           //

	Message       string `json:"message" orm:"type(text);null"`        //
	TitleFullText string `json:"title_fulltext" orm:"type(text);null"` //
}

//
type ArticleComments struct {
	Id    int64 `json:"id" orm:"pk;auto"`
	Uid   int64 `json:"uid" orm:"index"`        // 用户ID
	AtUid int64 `json:"at_uid" orm:"index"`     // At用户ID
	Aid   int64 `json:"article_id" orm:"index"` //

	Votes   int       `json:"votes" orm:"default(0)"`
	Message string    `json:"message" orm:"type(text);null"`                    //内容
	AddTime time.Time `json:"add_time" orm:"auto_now_add;type(datetime);index"` //添加时间
}

//
type ArticleVote struct {
	Id   int64 `json:"id" orm:"pk;auto"`
	Uid  int64 `json:"uid" orm:"index"`         // 用户ID
	Aid  int64 `json:"article_id" orm:"index"`  //
	Auid int64 `json:"article_uid" orm:"index"` //

	RepuFactor int   `json:"repu_factor" orm:"default(0)"` //
	Rating     int16 `json:"rating" orm:"default(0)"`      //

	Type    string    `json:"type" orm:"size(16);null;index"`
	AddTime time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` //添加时间
}

func init() {
	orm.RegisterModel(new(Article), new(ArticleComments), new(ArticleVote))
}

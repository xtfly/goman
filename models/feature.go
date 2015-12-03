package models

import "github.com/astaxie/beego/orm"

//专题
type Feature struct {
	Id         int64  `json:"id" orm:"pk;auto"`                     //ID
	Title      string `json:"title" orm:"size(200);null;index"`     //专题标题
	Desc       string `json:"description" orm:"size(255);null"`     //专题标题
	Icon       string `json:"icon" orm:"size(255);null"`            //专题标题
	TopicCount int    `json:"topic_count" orm:"default(0)"`         //话题计数
	Enable     bool   `json:"enabled" orm:"default(0)"`             //
	UrlToken   string `json:"url_token" orm:"size(128);null;index"` //
	SeoTitle   string `json:"seo_title" orm:"size(128);null"`       //
	Css        string `json:"css" orm:"type(text);null"`            //自定义CSS
}

//专题与话题关联
type FeatureTopic struct {
	Id  int64 `json:"id" orm:"pk;auto"`       //ID
	Fid int64 `json:"feature_id" orm:"index"` //专题ID
	Tid int64 `json:"topic_id" orm:"index"`   //话题ID
}

func init() {
	orm.RegisterModel(new(Feature), new(FeatureTopic))
}

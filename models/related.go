package models

//
import (
	"time"

	"github.com/astaxie/beego/orm"
)

//
type RelatedTopic struct {
	Id        int64 `json:"id" orm:"pk;auto"`                  //ID
	TopicId   int64 `json:"topic_id" orm:"default(0);index"`   //话题
	RelatedId int64 `json:"related_id" orm:"default(0);index"` //相关话题
}

type RelatedLink struct {
	Id       int64     `json:"id" orm:"pk;auto"`                                 //ID
	Uid      int64     `json:"uid" orm:"default(0);index"`                       //
	ItemId   int64     `json:"item_id" orm:"default(0);index"`                   //
	ItemType string    `json:"item_type" orm:"size(32);null"`                    //
	Link     string    `json:"link" orm:"size(255);null"`                        //
	AddTime  time.Time `json:"add_time" orm:"type(datetime);auto_now_add;index"` //
}

func init() {
	orm.RegisterModel(new(RelatedTopic), new(RelatedLink))
}

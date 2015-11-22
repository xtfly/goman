package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//
type Report struct {
	Id       int64     `json:"id" orm:"pk;auto"`                                    //ID
	Uid      int64     `json:"uid" orm:"default(0)"`                                //举报用户id
	Type     string    `json:"type" orm:"size(32);null;index"`                      //类别
	TargetId int64     `json:"target_id" orm:"default(0)"`                          //
	Reason   string    `json:"reason" orm:"size(255);null"`                         //举报理由
	Url      string    `json:"url" orm:"size(255);null"`                            //
	Status   int       `json:"status" orm:"default(0);index"`                       //是否处理
	AddTime  time.Time `json:"active_time" orm:"auto_now_add;type(datetime);index"` //举报时间

}

func init() {
	orm.RegisterModel(new(Report))
}

package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// 积分记录表
type IntegralLog struct {
	Id         int64     `json:"id" orm:"pk;auto"`                                    //ID
	Uid        int64     `json:"uid" orm:"default(0);index"`                          //用户ID
	Action     string    `json:"action" orm:"size(32);null;index"`                    //
	Integral   int       `json:"integral" orm:"default(0);index"`                     //
	Note       string    `json:"note" orm:"size(128);null"`                           //
	Balance    int       `json:"balance" orm:"default(0)"`                            //
	ItemId     int64     `json:"item_id" orm:"default(0);index"`                      //
	ActiveTime time.Time `json:"active_time" orm:"auto_now_add;type(datetime);index"` //时间

}

func init() {
	orm.RegisterModel(new(IntegralLog))
}

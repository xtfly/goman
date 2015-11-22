package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//
type SearchCache struct {
	Id      int64     `json:"id" orm:"pk;auto"`                                    //ID
	Hash    int64     `json:"hash" orm:"size(32);index"`                           //
	Data    string    `json:"data" orm:"type(text);null"`                          //
	AddTime time.Time `json:"active_time" orm:"auto_now_add;type(datetime);index"` //时间

}

func init() {
	orm.RegisterModel(new(SearchCache))
}

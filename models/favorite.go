package models

//
import (
	"time"

	"github.com/astaxie/beego/orm"
)

// 收藏表
type Favorite struct {
	Id       int64  `json:"id" orm:"pk;auto"`               //ID
	Uid      int64  `json:"uid" orm:"default(0);index"`     //
	ItemId   int64  `json:"item_id" orm:"default(0);index"` //
	ItemType string `json:"item_type" orm:"size(32);null"`  //

	Time time.Time `json:"time" orm:"auto_now_add;type(datetime)"` // 添加时间
}

//
type FavoriteTag struct {
	Id       int64  `json:"id" orm:"pk;auto"`               //ID
	Uid      int64  `json:"uid" orm:"default(0);index"`     //
	ItemId   int64  `json:"item_id" orm:"default(0);index"` //
	ItemType string `json:"item_type" orm:"size(32);null"`  //
	Title    string `json:"title" orm:"size(128);null"`     //
}

func init() {
	orm.RegisterModel(new(Favorite), new(FavoriteTag))
}

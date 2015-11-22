package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// 附件记录
type Attach struct {
	Id           int64     `json:"id" orm:"pk;auto"`
	FileName     string    `json:"file_name" orm:"null"`                 //附件名称
	AccessKey    string    `json:"access_key" orm:"null;index"`          //接入Key
	FileLocation string    `json:"file_location" orm:"null"`             //文件位置
	IsImage      bool      `json:"is_image" orm:"default(0);index"`      //是否为图片
	ItemType     string    `json:"item_type" orm:"null;size(32)"`        //关联类型
	ItemId       int64     `json:"item_id" orm:"default(0)"`             //关联ID
	WaitApproval bool      `json:"wait_approval" orm:"default(0);index"` //
	AddTime      time.Time `json:"add_time" orm:"type(datetime);auto_now_add"`
}

func init() {
	orm.RegisterModel(new(Attach))
}

// 多字段唯一键
func (m *Attach) TableUnique() [][]string {
	return [][]string{
		[]string{"ItemId", "ItemType"},
	}
}

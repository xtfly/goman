package models

//
import (
	"time"

	"github.com/astaxie/beego/orm"
)

// 用户通知
type Notification struct {
	Id           int64     `json:"id" orm:"pk;auto"`                       //ID
	SenderUid    int64     `json:"sender_uid" orm:"default(0);index"`      //发送者ID
	RecipientUid int64     `json:"recipient_uid" orm:"default(0)"`         //接收者ID
	ActionType   int8      `json:"action_type" orm:"default(0);index"`     //操作类型
	ModelType    int16     `json:"model_type" orm:"default(0);index"`      //
	SourceId     int64     `json:"source_id" orm:"default(0);index"`       //关联 ID
	ReadFlag     bool      `json:"read_flag" orm:"default(0)"`             //阅读状态
	Time         time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间
}

// 多字段索引
func (m *Notification) TableIndex() [][]string {
	return [][]string{
		[]string{"RecipientUid", "ReadFlag"},
	}
}

//
type NotificationData struct {
	Id   int64 `json:"id" orm:"pk;auto"`      //ID
	Data int64 `json:"data" orm:"type(text)"` //

}

func init() {
	orm.RegisterModel(new(Notification), new(NotificationData))
}

package models

//
import (
	"time"

	"github.com/astaxie/beego/orm"
)

type InboxMsg struct {
	Id              int64        `json:"id" orm:"pk;auto"`                                 //ID
	Uid             int64        `json:"uid" orm:"default(0);index"`                       //发送者 ID
	InboxDialog     *InboxDialog `json:"dialog_id" orm:"rel(one);index;column(dialog_id)"` //对话 ID
	SenderRemove    bool         `json:"sender_remove" orm:"default(0);index"`             //
	RecipientRemove bool         `json:"recipient_remove" orm:"default(0);index"`          //
	Receipt         int64        `json:"receipt" orm:"default(0);index"`                   //

	Message string    `json:"message" orm:"type(text);null"`          //内容
	Time    time.Time `json:"time" orm:"auto_now_add;type(datetime)"` // 添加时间
}

//
type InboxDialog struct {
	Id              int64     `json:"id" orm:"pk;auto"`                           //ID
	SenderUid       int64     `json:"sender_uid" orm:"default(0);index"`          //发送者UID
	SenderUnread    int       `json:"sender_unread" orm:"default(0)"`             //发送者未读
	SenderCount     int       `json:"sender_count" orm:"default(0)"`              //发送者显示对话条数
	RecipientUid    int64     `json:"recipient_uid" orm:"default(0);index"`       //接收者UID
	RecipientUnread int       `json:"recipient_unread" orm:"default(0)"`          //接收者未读
	RecipientCount  int       `json:"recipient_count" orm:"default(0)"`           //接收者显示对话条数
	AddTime         time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` // 添加时间
	UpdateTime      time.Time `json:"update_time" orm:"auto_now;type(datetime)"`  // 最后更新时间
}

func init() {
	orm.RegisterModel(new(InboxMsg), new(InboxDialog))
}

package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// 用户关注表
type UserFollow struct {
	Id        int64 `json:"id" orm:"pk;auto"`
	FansUid   int64 `json:"fans_uid" orm:"null;index"`   // 关注人的UID
	FriendUid int64 `json:"friend_uid" orm:"null;index"` // 被关注人的uid

	AddTime time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` // 添加时间
	Ip      string    `json:"ip" valid:"IP"`                              // 客户端ip

	ActiveUrl string `json:"active_url"` // 停留页面
	UserAgent string `json:"user_agent"` // 用户客户端信息
}

var ufollow = new(UserFollow)

func init() {
	orm.RegisterModel(ufollow)
}

func AddUserFollow(t *Transaction, fans, friend int64) bool {
	if fans == friend {
		return true
	}

	// 用户不存在
	if !UserExistedById(fans) || !UserExistedById(friend) {
		return false
	}

	// 关系已存在
	if UFollowExistedById(fans, friend) {
		return true
	}

	m := &UserFollow{
		FansUid:   fans,
		FriendUid: friend,
	}
	_, ok := t.Insert(m)
	if !ok {
		return false
	}

	// 更新各自的统计，需要性能优化
	go func() {
		nt := NewTr()
		updateUFollowRelation(nt, fans)
		updateUFollowRelation(nt, friend)
	}()

	return true
}

func UFollowExistedById(fans, friend int64) bool {
	return NewTr().Query("UserFollow").Filter("FansUId", fans).
		Filter("FriendUId", friend).Exist()
}

func updateUFollowRelation(t *Transaction, uid int64) bool {
	fansc, ok1 := t.Count("UserFollow", "FriendUId", uid)
	if !ok1 {
		return false
	}

	friendc, ok2 := t.Count("UserFollow", "FansUId", uid)
	if !ok2 {
		return false
	}

	return t.UpdateById("Users", uid,
		orm.Params{"FansCount": fansc, "FriendCount": friendc})
}

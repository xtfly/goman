package models

import (
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/orm"
	"github.com/xtfly/goman/kits"
)

const (
	GenderUnknown = iota
	GenderMale
	GenderFemale
)

type Users struct {
	Id        int64  `json:"uid" orm:"pk;auto"`                 //用户的 UID
	UserName  string `json:"user_name" orm:"index;unique"`      //用户名
	Email     string `json:"email" orm:"index"`                 //EMAIL
	Mobile    string `json:"mobile" orm:"index"`                //用户手机
	Salt      string `json:"-"`                                 //加密的盐值
	Password  string `json:"-"`                                 //加密之后的密文
	GroupId   int64  `json:"group_id" orm:"index"`              //用户组
	RepuGroup int64  `json:"reputation_group" orm:"default(0)"` //威望对应组

	// 基本信息
	Avatar     string    `json:"avatar" orm:"null;index;"`                   //头像文件
	Gender     int8      `json:"gender" orm:"default(0)"`                    //0：Unknown, 1: Male， 2：Female
	Birthday   time.Time `json:"birthday" orm:"null;type(date)"`             //生日
	Province   string    `json:"province" orm:"null"`                        //省
	City       string    `json:"city" orm:"null"`                            //市
	Intro      string    `json:"introduction" orm:"null"`                    //个人简介
	Signature  string    `json:"signature" orm:"null"`                       //个人签名
	JobId      int64     `json:"job_id" orm:"default(0)"`                    //职业ID
	RegTime    time.Time `json:"reg_time" orm:"auto_now_add;type(datetime)"` //注册时间
	Updated    time.Time `json:"updated" orm:"type(datetime);null"`          //资料更新时间
	RegIp      string    `json:"reg_ip" orm:"null"`                          //注册IP
	LastLogin  time.Time `json:"last_login" orm:"type(datetime);null"`       //最后登录时间
	LastIp     string    `json:"last_ip" orm:"null"`                         //最后登录的IP
	LoginCount int64     `json:"login_times" orm:"default(0)"`               //登录次数
	OnlineTime int64     `json:"online_time" orm:"default(0)"`               //在线时间
	LastActive time.Time `json:"last_active" orm:"type(datetime);null"`      //最后登录时间

	// 统计
	NotificationUnread  int64     `json:"notification_unread" orm:"default(0)"`             //未读系统通知
	InboxUnread         int64     `json:"inbox_unread" orm:"default(0)"`                    //未读短信息
	InboxRecv           int8      `json:"inbox_recv" orm:"default(0)"`                      //0-所有人可以发给我,1-我关注的人
	FansCount           int64     `json:"fans_count" orm:"default(0)"`                      //粉丝数
	FriendCount         int64     `json:"friend_count" orm:"default(0)"`                    //观众数
	InviteCount         int64     `json:"invite_count" orm:"default(0)"`                    //邀请我回答数量
	ArticleCount        int64     `json:"article_count" orm:"default(0)"`                   //文章数量
	QuestionCount       int64     `json:"question_count" orm:"default(0)"`                  //问题数量
	AnswerCount         int64     `json:"answer_count" orm:"default(0)"`                    //回答数量
	TopicFocusCount     int64     `json:"topic_focus_count" orm:"default(0)"`               //关注话题数量
	InvitationAvailable int64     `json:"invitation_available" orm:"default(0)"`            //邀请数量
	AgreeCount          int64     `json:"agree_count" orm:"default(0)"`                     //赞同数量
	ThanksCount         int64     `json:"thanks_count" orm:"default(0)"`                    //感谢数量
	ViewsCount          int64     `json:"views_count" orm:"default(0)"`                     //个人主页查看数量
	Reputation          int64     `json:"reputation" orm:"index;default(0)"`                //威望
	ReputationUpdated   time.Time `json:"reputation_update_time" orm:"type(datetime);null"` //威望更新
	Integral            int64     `json:"integral" orm:"default(0)"`                        //积分
	DraftCount          int64     `json:"draft_count" orm:"default(0)"`                     //

	// 安全
	Forbidden       bool      `json:"forbidden" orm:"default(0)"`                 //是否禁止用户
	FirstLogin      bool      `json:"first_login" orm:"default(1)"`               //首次登录标记
	ValidEmail      bool      `json:"valid_email" orm:"default(0)"`               //邮箱验证
	Verified        string    `json:"-" orm:"null"`                               //验证码
	WeiboVist       bool      `json:"weibo_visit" orm:"default(1)"`               //微博允许访问
	CommonEmail     string    `json:"common_email" orm:"null"`                    //常用邮箱
	UrlToken        string    `json:"url_token" orm:"index;null"`                 //个性网址
	UrlTokenUpdated time.Time `json:"url_token_update" orm:"type(datetime);null"` //个性网址更新

	//RecentTopics string `json:"recent_topics" orm:"size(1024)"`

	ExtAttrs map[string]interface{} `json:"-" orm:"-"`
}

var users = new(Users)

func init() {
	orm.RegisterModel(users)
}

//----------------------------------------------------------
//检查用户名,电子邮件地址是否已经存在
func UserExistedByName(name string) bool {
	return NewTr().Query("Users").Filter("UserName", name).Exist()
}

func UserExistedByEmail(email string) bool {
	return NewTr().Query("Users").Filter("Email", email).Exist()
}

func UserExistedById(id int64) bool {
	return NewTr().Query("Users").Filter("Id", id).Exist()
}

//----------------------------------------------------------
// 增加一条记录
func (m *Users) Add(t *Transaction) (int64, bool) {
	m.LastIp = m.RegIp

	//
	m.Password, m.Salt = kits.GenPasswd(m.Password, 8)

	// 生成用户的URL
	m.UrlToken = kits.NewRandWithPrefix(m.Email, 8)
	m.UrlTokenUpdated = time.Now()

	m.InvitationAvailable = syscfg.Ra.NewerInviteNum
	m.RepuGroup = 5 // TODO
	m.FirstLogin = true

	// 前面要求对邀请码验证
	if syscfg.Ra.RegisterValidType == RegValidNone ||
		(syscfg.Ra.RegisterValidType == RegValidEmail && syscfg.Ra.RegisterType == RegTypeInvite) {
		m.GroupId = GroupNormal
	}

	//
	id, ok := t.Insert(m)
	if !ok {
		return id, ok
	}

	// 增加关注关系表
	fuids := strings.Split(syscfg.Ra.DefFocusUids, ",")
	for _, fuid := range fuids {
		if !AddUserFollow(t, id, com.StrTo(fuid).MustInt64()) {
			return id, false
		}
	}

	// 初始化积分
	if !AddIntegralLog(t, id, IntegralRegister, syscfg.Ir.DefautlIntegral, "初始资本") {
		return id, false
	}

	// 默认的配置
	us := &UserSetting{
		Uid:                 id,
		EmailSetting:        syscfg.Ra.EmailSettings,
		NotificationSetting: syscfg.Ra.NotificationSettings,
	}
	if _, ok := t.Insert(us); !ok {
		return id, false
	}

	return id, ok
}

//----------------------------------------------------------
// 检查用户，邮件与密码是否正确
func (m *Users) CheckSignin(input string, password string) bool {
	if input == "" || password == "" {
		return false
	}

	f := "Email"
	if strings.Contains(input, "@") {
		m.Email = input
	} else {
		m.UserName = input
		f = "UserName"
	}

	t := NewTr()
	if !t.Read(m, f) {
		return false
	}

	if !kits.CmpPasswd(password, m.Salt, m.Password) {
		return false
	}
	return true
}

//----------------------------------------------------------
// 获取活跃用户 (非垃圾用户)
// 好友 & 粉丝 > 5, 回复 > 5, 根据登陆时间, 倒序
func GetActivityUsers(limit int64, uid int64) ([]*Users, bool) {
	var users []*Users
	_, err := NewTr().Query("Users").Filter("FansCount__gt", 5).Filter("FriendCount__gt", 5).Filter("AnswerCount__gt", 1).
		Exclude("Id", uid).Limit(limit).OrderBy("-LastLogin").All(&users)
	return users, err != nil
}

// 获取推荐用户
func GetRecommendRandUser(limit int64, uid int64, recommends []string) ([]*Users, bool) {
	var users []*Users

	uns := make([]interface{}, len(recommends))
	for _, v := range recommends {
		uns = append(uns, v)
	}
	_, err := NewTr().Query("Users").Filter("UserName__gt", uns).
		Exclude("Id", uid).Limit(limit).All(&users) // TODO by rand
	return users, err != nil
}

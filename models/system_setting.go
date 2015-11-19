package models

import (
	"reflect"

	"github.com/astaxie/beego/orm"
	"github.com/xtfly/goman/kits/m2s"

	log "github.com/Sirupsen/logrus"
)

// site global information
type SiteInfo struct {
	SiteName  string `map:"site_name"`
	SiteDesc  string `map:"description"`
	Keywords  string `map:"keywords"`
	UploadUrl string `map:"upload_url"`
	UploadDir string `map:"upload_dir"`
	IcpBeian  string `map:"icp_beian"`
}

// enum of register type
const (
	RegTypeOpen   = iota // 开放注册
	RegTypeInvite        // 邀请注册
	RegTypeClose         // 关闭注册
)

// enum of user register rule
const (
	UserRuleNotLimit               = iota // 不限制
	UserRuleChineseLetterNumUnline        // 汉字/字母/数字/下划线 Chinese characters / letter / number / underline
	UserRuleLetterNumUnline               // 字母/数字/下划线
	UserRuleChinese                       // 汉字
)

// enum of user send email condition
const (
	UserEmailFollowme       = iota // 当有人关注我
	UserEmailQuestionInvite        // 有人邀请我回答问题
	UserEmailNewAnswer             // 我关注的问题有了新回复
	UserEmailNewMessage            // 有人向我发送私信
	UserEmailQuestionModify        // 我的问题被编辑
)

// enum of user send notification condition
const (
	UserNotifyFocusMe            = iota //有人关注了我
	UserNotifyNewReply                  //我关注的问题有了新的回复
	UserNotifyInivteMe                  //有人邀请我回复问题
	UserNotifyCommented                 //我的问题被评论
	UserNotifyCommentReplied            //我的问题评论被回复
	UserNotifyQuestCommentAtMe          //有问题评论提到我
	UserNotifyAnswerCommentAtme         //有回答评论提到我
	UserNotifyReplyAtme                 //有回答提到我
	UserNotifyReplyUp                   //我的回复收到赞同
	UserNotifyReplyTks                  //我的回复收到感謝
	UserNotifyQuestTks                  //我发布的问题收到感謝
	UserNotifyQuestModify               //我的问题被编辑
	UserNotityQuestRedirect             //我发布的问题被重定向
	UserNotifyArticleComment            //我的文章被评论
	UserNotifyArticleCommentAtMe        //有文章评论提到我
)

// the setting of user register access information
type RegisterAccess struct {
	DefaultTimezone      string `map:"default_timezone"`
	SiteClose            bool   `map:"site_close,string"`
	SiteNotice           string `map:"close_notice"`
	RegisterCaptcha      bool   `map:"register_seccode,string"`
	RegisterValidType    string `map:"register_valid_type"`
	RegisterType         int    `map:"register_type,string"`
	UsernameRule         int    `map:"username_rule,string"`
	UsernameLenMinx      int    `map:"username_length_min,string"`
	UsernameLenMax       int    `map:"username_length_max,string"`
	CensorUser           string `map:"censoruser"`
	DefFocusUids         string `map:"def_focus_uids"`
	WelcomeRecmdusers    string `map:"welcome_recommend_users"`
	NewerInviteNum       int    `map:"newer_invitation_num,string"`
	WelcomeMsgPm         string `map:"welcome_message_pm"`
	EmailSettings        string `map:"set_email_settings"`
	NotificationSettings string `map:"set_notification_settings"`
}

func (m *RegisterAccess) EmailSettingExist(i int) bool {
	return existInStr(m.EmailSettings, i)
}

func (m *RegisterAccess) NotificationSettingExist(i int) bool {
	return existInStr(m.NotificationSettings, i)
}

func existInStr(s string, i int) bool {
	if len(s) > i {
		if s[i] == '1' {
			return true
		}
	}
	return false
}

// System setting table
type SystemSetting struct {
	Id    int64  `orm:"pk;auto"`
	Name  string `orm:"size(64)"`   //字段名
	Value string `orm:"size(1024)"` //变量值
}

func init() {
	orm.RegisterModel(new(SystemSetting))
}

//
type GlobalSetting struct {
	Si *SiteInfo
	Ra *RegisterAccess
}

func NewGlobalSetting() *GlobalSetting {
	return &GlobalSetting{
		Si: &SiteInfo{},
		Ra: &RegisterAccess{},
	}
}

//-----------------------------------------------------------------------------
// Load all setting from db
func (gss *GlobalSetting) LoadAll() bool {
	var lists []orm.ParamsList
	num, err := NewTr().Query("SystemSetting").ValuesList(&lists, "Name", "Value")
	if err != nil {
		log.Errorf("Load data from SystemSetting failed. error = %s", err.Error())
		return false
	}

	kvs := make(map[string]interface{}, num)
	log.Printf("Result Nums in SystemSetting : %d", num)
	for _, row := range lists {
		kvs[row[0].(string)] = row[1]
	}

	// 给每个成员变量赋值
	v := reflect.ValueOf(gss).Elem()
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		if err := m2s.Map2Struct(kvs, fv.Interface()); err != nil {
			log.Errorf("Assign value to  failed. error = %s", err.Error())
			return false
		}
	}

	return true
}

// Update site info to db
func (gss *GlobalSetting) UpdateSiteInfo() bool {
	return gss.update(m2s.Struct2Map(gss.Si))
}

func (gss *GlobalSetting) update(kvs map[string]interface{}) bool {
	if len(kvs) <= 0 {
		log.Errorf("Save data to SystemSetting failed, map size is zero.")
		return false
	}

	t := NewTr()
	t.Begin()
	defer t.End()

	for k, v := range kvs {
		if _, t.Err = t.Query("SystemSetting").Filter("Name", k).Update(orm.Params{"Value": v.(string)}); t.Err != nil {
			log.Errorf("Save data to SystemSetting failed. error = %s", t.Err.Error())
			return false
		}
	}
	return true
}

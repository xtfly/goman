package models

import (
	"encoding/json"
	"reflect"

	"github.com/astaxie/beego/orm"
	"github.com/xtfly/goman/kits/m2s"

	log "github.com/Sirupsen/logrus"
)

// site global information
type SiteInfo struct {
	SiteName  string `map:"site_name"`   //网站名称
	SiteDesc  string `map:"description"` //网站简介
	Keywords  string `map:"keywords"`    //网站关键词
	UploadUrl string `map:"upload_url"`  //上传目录外部访问 URL 地址
	UploadDir string `map:"upload_dir"`  //上传文件存放绝对路径
	IcpBeian  string `map:"icp_beian"`   //网站 ICP 备案号
	Static    string `map:"static_url"`  //static 目录资源 URL 地址
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
	DefaultTimezone      string `map:"default_timezone"`            //默认时区
	SiteClose            bool   `map:"site_close,string"`           //站点关闭
	SiteNotice           string `map:"close_notice"`                //站点关闭的提示
	RegisterCaptcha      bool   `map:"register_seccode,string"`     //新用户注册显示验证码
	RegisterValidType    bool   `map:"register_valid_type,string"`  //新用户注册验证类型
	RegisterType         int    `map:"register_type,string"`        //注册类型
	UsernameRule         int    `map:"username_rule,string"`        //用户名规则
	UsernameLenMin       int    `map:"username_length_min,string"`  //用户名最少字符数
	UsernameLenMax       int    `map:"username_length_max,string"`  //用户名最多字符数:
	CensorUser           string `map:"censoruser"`                  //用户注册名不允许出现以下关键字
	DefFocusUids         string `map:"def_focus_uids"`              //用户注册后默认关注的用户 ID
	WelcomeRecmdusers    string `map:"welcome_recommend_users"`     //首次登录推荐用户列表
	NewerInviteNum       int    `map:"newer_invitation_num,string"` //新用户注册获得邀请数量
	WelcomeMsgPm         string `map:"welcome_message_pm"`          //新用户注册系统发送的欢迎内容
	EmailSettings        string `map:"set_email_settings"`          //新用户默认邮件提醒设置
	NotificationSettings string `map:"set_notification_settings"`   //新用户默认通知设置
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

// 站点功能
type SiteCapibility struct {
	SiteAnnounce        string `map:"site_announce"`                 //网站公告
	UnreadFlushInterval int    `map:"unread_flush_interval,string"`  //通知未读数刷新间隔时间
	AutoQuestLockDay    int    `map:"auto_question_lock_day,string"` //问题自动锁定时间
	StatisticCode       string `map:"statistic_code"`                //网站统计代码
	ReportReason        string `map:"report_reason"`                 //问题举报理由选项
	ReportMsgUid        string `map:"report_message_uid"`            //有新举报与认证申请私信提醒用户 ID
	TimeStyle           int    `map:"time_style,string"`             //系统时间显示格式:
	AdminLoginSecCode   bool   `map:"admin_login_seccode,string"`    //管理员后台登录是否需要验证码:
	ReportDiagnostics   bool   `map:"report_diagnostics,string"`     //发送诊断数据帮助改善产品
}

//用户权限
type UserAuthority struct {
	AnswerUnique     bool `map:"answer_unique,string"`        //用户对每个问题的回复限制
	AnswerSelfQuest  bool `map:"answer_self_question,string"` //允许用户回复自己发起的问题
	AnonymouseEnable bool `map:"anonymous_enable,string"`     //允许匿名发起或回复
}

// 积分与威望
type IntegralReputation struct {
	Enabled            bool   `map:"integral_system_enabled,string"`                     //使用积分系统
	Unit               string `map:"integral_unit"`                                      //积分单位
	DefautlIntegral    int64  `map:"integral_system_config_register,string"`             //新用户注册默认拥有积分
	FinishProfile      int64  `map:"integral_system_config_profile,string"`              //用户完善资料获得积分
	InviteUser         int64  `map:"integral_system_config_invite,string"`               //用户邀请他人注册且被邀请人成功注册
	NewQuestion        int64  `map:"integral_system_config_new_question,string"`         //发起问题
	NewAnswer          int64  `map:"integral_system_config_new_answer,string"`           //回复问题
	AnswerChangeSource bool   `map:"integral_system_config_answer_change_source,string"` //问题被回复时增加发起者积分
	BestAnswer         int64  `map:"integral_system_config_best_answer,string"`          //回复被评为最佳回复
	Thanks             int64  `map:"integral_system_config_thanks,string"`               //感谢回复
	AnswerFold         int64  `map:"integral_system_config_answer_fold,string"`          //回复被折叠
	InviteAnswer       int64  `map:"integral_system_config_invite_answer,string"`        //发起者邀请用户回答问题且收到答案
	PubrepFactor       int64  `map:"publisher_reputation_factor,string"`                 //发起者赞同反对威望系数
	BestrepFactor      int64  `map:"best_answer_reput,string"`                           //最佳回复威望系数
	ReplogFactor       int64  `map:"reputation_log_factor,string"`                       //对数底值
}

//内容设置
type ContextSetting struct {
	QuickPublish          bool   `map:"quick_publish,string"`                //快捷模式
	AdvEditor             bool   `map:"advanced_editor_enable,string"`       //编辑器设置
	CategoryEnable        bool   `map:"category_enable,string"`              //启用分类功能
	UploadEnable          bool   `map:"upload_enable,string"`                //允许上传附件
	AllowedUploadType     string `map:"allowed_upload_types"`                //允许的附件文件类型
	UploadSizeLimit       int    `map:"upload_size_limit,string"`            //允许上传最大附件大小
	AnswerLengthLow       int    `map:"answer_length_lower,string"`          //回复内容最小字符数限制
	QuestTitleLimit       int    `map:"question_title_limit,string"`         //问题标题最大字符数限制
	ForceAddTopic         bool   `map:"new_question_force_add_topic,string"` //新问题强制要求添加话题
	QuestTopicLimit       int    `map:"question_topics_limit,string"`        //问题话题数量限制
	UnfoldComments        bool   `map:"unfold_question_comments,string"`     //自动展开评论
	CommentLimit          int    `map:"comment_limit,string"`                //评论内容最大字符数限制
	TopicTitleLimit       int    `map:"topic_title_limit,string"`            //话题标题最大字符数限制
	AnswerEditTime        int    `map:"answer_edit_time,string"`             //回复编辑/删除有效时间
	UninterestedFold      int    `map:"uninterested_fold,string"`            //“没有帮助”数量达到多少个时折叠回复
	BestAnswerDay         int    `map:"best_answer_day,string"`              //系统自动选出最佳回复时间
	BestAnswerMinCount    int    `map:"best_answer_min_count,string"`        //参与自动选出最佳回复的问题最小回复数
	BestAgreeMinCount     int    `map:"best_agree_min_count,string"`         //参与自动选出最佳回复的最小赞同数
	ReadQuestLastDays     int    `map:"reader_questions_last_days,string"`   //阅读器获取最近多少天的热门问题
	ReadQuestAgreeCount   int    `map:"reader_questions_agree_count,string"` //阅读器热门问题赞同数需大于或等于
	AutoCreateSocialTopic bool   `map:"auto_create_social_topics,string"`    //自动建立地区/职业/公司话题
	SensitiveWords        string `map:"sensitive_words"`                     //敏感词列表
	EnableHelp            bool   `map:"enable_help_center,string"`           //启用帮助中心
}

//界面设置
type PageSetting struct {
	UiStyle        string `map:"ui_style"`                      //用户界面风格
	IndexPerPage   int    `map:"index_per_page,string"`         //最新动态显示条数
	NotifyPerPage  int    `map:"notifications_per_page,string"` //最新动态显示条数
	ContentPerPage int    `map:"contents_per_page,string"`      //内容列表页显示条数
	RecomdUsersNum int    `map:"recommend_users_number,string"` //首页推荐用户和话题数量
}

// System setting table
type SystemSetting struct {
	Id    int64  `orm:"pk;auto"`
	Name  string `orm:"size(128)"`  //字段名
	Value string `orm:"size(2048)"` //变量值
}

func init() {
	orm.RegisterModel(new(SystemSetting))
}

//
type GlobalSetting struct {
	Si *SiteInfo
	Ra *RegisterAccess
	Sc *SiteCapibility
	Ua *UserAuthority
	Ir *IntegralReputation
	Cs *ContextSetting
	Ps *PageSetting
}

func NewGlobalSetting() *GlobalSetting {
	return &GlobalSetting{
		Si: &SiteInfo{},
		Ra: &RegisterAccess{},
		Sc: &SiteCapibility{},
		Ua: &UserAuthority{},
		Ir: &IntegralReputation{},
		Cs: &ContextSetting{},
		Ps: &PageSetting{},
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

	s, _ := json.Marshal(gss)
	log.Printf("System setting=%s", string(s))
	return true
}

// Update site info to db
func (gss *GlobalSetting) UpdateSiteInfo() bool {
	return gss.update(m2s.Struct2Map(gss.Si))
}

func (gss *GlobalSetting) UpdateRegisterAccess() bool {
	return gss.update(m2s.Struct2Map(gss.Ra))
}

func (gss *GlobalSetting) UpdateSiteCapibility() bool {
	return gss.update(m2s.Struct2Map(gss.Sc))
}

func (gss *GlobalSetting) UpdateUserAuthority() bool {
	return gss.update(m2s.Struct2Map(gss.Ua))
}

func (gss *GlobalSetting) UpdateIntegralReputation() bool {
	return gss.update(m2s.Struct2Map(gss.Ir))
}

func (gss *GlobalSetting) UpdateContextSetting() bool {
	return gss.update(m2s.Struct2Map(gss.Cs))
}

func (gss *GlobalSetting) UpdatePageSetting() bool {
	return gss.update(m2s.Struct2Map(gss.Ps))
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

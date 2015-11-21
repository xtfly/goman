package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//用户组
type UsersGroup struct {
	Id         int64  `json:"gid" orm:"pk;auto"`
	Type       int8   `json:"type" orm:"index"`            // 0-会员组 1-系统组
	Custom     bool   `json:"custom" orm:"index"`          // 是否自定义
	Name       string `json:"name"`                        //
	Permission string `json:"permission" orm:"size(1024)"` // 权限设置

	RepuLower  string `json:"repu_lower"`  //
	RepuHiger  string `json:"repu_higer"`  //
	RepuFactor string `json:"repu_factor"` // 威望系数
}

func init() {
	orm.RegisterModel(new(UsersGroup))
}

type UserPermission struct {
	IsAdmin          int `json:"is_admin"`
	IsModerator      int `json:"is_moderator"`
	EditQuestion     int `json:"edit_question"`
	RedirectQuestion int `json:"redirect_question"`
	EditTopic        int `json:"edit_topic"`
	ManageTopic      int `json:"manage_topic"`
	CreateTopic      int `json:"create_topic"`
	UploadAttach     int `json:"upload_attach"`
	PublishUrl       int `json:"publish_url"`
	PublishArticle   int `json:"publish_article"`
	EditArticle      int `json:"edit_article"`
	EditQuestTopic   int `json:"edit_question_topic"`
	PublishComment   int `json:"publish_comment"`

	HumanValid      int          `json:"human_valid"`
	QuestValidHour  int          `json:"question_valid_hour"`
	AnswerValidHour int          `json:"answer_valid_hour"`
	ApprovalTime    ApprovalTime `json:"publish_approval_time"`

	VisitSite     int `json:"visit_site"`
	VisitExplore  int `json:"visit_explore"`
	SearchAvail   int `json:"search_avail"`
	VisitQuestion int `json:"visit_question"`
	VisitTopic    int `json:"visit_topic"`
	VisitFeature  int `json:"visit_feature"`
	VisitPeople   int `json:"visit_people"`
	VisitChapter  int `json:"visit_chapter"`
	AnswerShow    int `json:"answer_show"`
}

const UPermissionOn = 1

type ApprovalTime struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

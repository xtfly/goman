package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

//工作经历
type WorkExperience struct {
	Id  int64 `json:"id" orm:"pk;auto"`           //
	Uid int64 `json:"uid" orm:"default(0);index"` //

	StartYear   time.Time `json:"start_year" orm:"type(date)"` // 开始年份
	EndYear     time.Time `json:"end_year" orm:"type(date)"`   // 开始年份
	CompanyName string    `json:"company_name" orm:"null"`     // 公司名
	JobId       int64     `json:"job_id" orm:"null"`           // OneToOne relation

	AddTime time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` // 添加时间
}

//教育经历
type EduExperience struct {
	Id  int64 `json:"id" orm:"pk;auto"`           //
	Uid int64 `json:"uid" orm:"default(0);index"` //

	EduYears   time.Time `json:"education_years" orm:"type(date)"` // 入学年份
	SchoolName string    `json:"school_name" orm:"null"`           // 学校名
	SchoolType string    `json:"school_type" orm:"null"`           // 学校类别
	Department string    `json:"departments" orm:"null"`           // 院系

	AddTime time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` // 添加时间
}

func init() {
	orm.RegisterModel(new(WorkExperience), new(EduExperience))
}

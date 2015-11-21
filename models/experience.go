package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

// 工作经历
type WorkExperience struct {
	Id   int64  `json:"id" orm:"pk;auto"`
	User *Users `json:"uid" orm:"null;rel(one)"` // OneToOne relation

	StartYear   time.Time `json:"start_year" orm:"type(date)"` // 开始年份
	EndYear     time.Time `json:"end_year" orm:"type(date)"`   // 开始年份
	CompanyName string    `json:"company_name"`                // 公司名
	Job         *Jobs     `json:"job_id" orm:"null;rel(one)"`  // OneToOne relation

	AddTime time.Time `json:"add_time" orm:"auto_now_add;type(datetime)"` // 添加时间
}

func init() {
	orm.RegisterModel(new(WorkExperience))
}

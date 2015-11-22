package models

import "github.com/astaxie/beego/orm"

// 工作类型
type Jobs struct {
	Id   int64  `json:"id" orm:"pk;auto"`
	Name string `json:"job_name"` // 名
}

func init() {
	orm.RegisterModel(new(Jobs))
}

// 查询所有的工作类型
func AllJobs() []Jobs {
	var jobs []Jobs
	NewTr().Query("Jobs").OrderBy("Id").All(&jobs)
	return jobs
}

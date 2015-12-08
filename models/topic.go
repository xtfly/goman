package models

//
import (
	"time"

	"github.com/astaxie/beego/orm"
	"github.com/xtfly/gokits"
)

//话题表
type Topic struct {
	Id      int64     `json:"id" orm:"pk;auto"`                       //ID
	Pid     int64     `json:"parent_id" orm:"default(-1);index"`      //
	Title   string    `json:"title" orm:"size(128);unique"`           //话题标题
	Desc    string    `json:"description" orm:"size(255);null"`       //话题描述
	Picture string    `json:"picture" orm:"size(255);null"`           //话题图片
	AddTime time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间

	Lock                 bool      `json:"lock" orm:"default(0)"`                //话题是否锁定
	DisussCount          int       `json:"discuss_count" orm:"default(0);index"` //讨论计数
	DisussCountLastWeek  int       `json:"discuss_count" orm:"default(0);index"` //讨论计数
	DisussCountLastMonth int       `json:"discuss_count" orm:"default(0);index"` //讨论计数
	DisussCountUpdate    time.Time `json:"time" orm:"type(datetime);null"`       //
	FocusCount           int       `json:"focus_count" orm:"default(0);index"`   //关注计数
	UserRelated          bool      `json:"user_related" orm:"default(0);index"`  //是否被用户关联
	IsParent             bool      `json:"is_parent" orm:"default(0)"`           //是否

	UrlToken string `json:"url_token" orm:"size(255);null"` //
	SeoTitle string `json:"seo_title" orm:"size(255);null"` //
	MergedId int64  `json:"merged_id" orm:"default(0)"`     //

	ExtAttrs map[string]interface{} `json:"-" orm:"-"`
}

type Topics []*Topic

//话题关注表
type TopicFocus struct {
	Id  int64 `json:"id" orm:"pk;auto"`           //ID
	Tid int64 `json:"topic_id" orm:"index"`       //
	Uid int64 `json:"uid" orm:"default(0);index"` //话题标题

	AddTime time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间
}

//话题合并表
type TopicMerge struct {
	Id       int64 `json:"id" orm:"pk;auto"`                 //ID
	SourceId int64 `json:"source_id" orm:"default(0);index"` //
	TargetId int64 `json:"target_id" orm:"default(0);index"` //
	Uid      int64 `json:"uid" orm:"default(0);index"`       //

	AddTime time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间
}

//话题关联表
type TopicRelation struct {
	Id     int64  `json:"id" orm:"pk;auto"`               //ID
	Tid    int64  `json:"topic_id" orm:"index"`           //
	ItemId int64  `json:"item_id" orm:"default(0);index"` //
	Uid    int64  `json:"uid" orm:"default(0);index"`     //
	Type   string `json:"type" orm:"size(16);null"`

	AddTime time.Time `json:"time" orm:"auto_now_add;type(datetime)"` //添加时间
}

func init() {
	orm.RegisterModel(new(Topic), new(TopicFocus), new(TopicMerge), new(TopicRelation))
}

//----------------------------------------------------------
// 随机选出几个Topic
func GetTopicExts(limit int64, uid int64) ([]*Topic, bool) {
	var topics []*Topic
	t := NewTr()
	_, err := t.Query("Topic").Limit(limit).All(&topics) // TODO by rand
	if err != nil {
		return topics, false
	}

	for _, v := range topics {
		v.ExtAttrs = map[string]interface{}{"HasFocus": t.ExistedV2("TopicFocus", orm.Params{"Uid": uid, "Tid": v.Id})}
	}
	return topics, true
}

//----------------------------------------------------------
func AddTopic(t *Transaction, name string) bool {
	_, ok := t.Insert(&Topic{Title: name, Desc: name})
	return ok
}

//----------------------------------------------------------
func GetFoucsTopicIdsByUid(uid int64) ([]int64, bool) {
	var tfs []*TopicFocus
	if _, err := NewTr().Query("Topicfocus").Filter("Uid", uid).All(&tfs); err != nil {
		return nil, false
	}

	ids := make([]int64, len(tfs))
	for _, v := range tfs {
		ids = append(ids, v.Tid)
	}

	return ids, true
}

//----------------------------------------------------------
func GetItemIdsByTopicIds(tids []int64, stype string, limit int) ([]int64, bool) {
	var trs []*TopicRelation
	if _, err := NewTr().Query("TopicRelation").Filter("Tid__in", gokits.SliceInt64To(tids)).
		Filter("Type", stype).Limit(limit).All(&trs); err != nil {
		return nil, false
	}

	ids := make([]int64, len(trs))
	for _, v := range trs {
		ids = append(ids, v.ItemId)
	}

	return ids, true
}

//----------------------------------------------------------
func GetTopicsByItemIds(iids []int64) (map[int64]Topics, bool) {
	t := NewTr()
	var trs []*TopicRelation
	if _, err := t.Query("TopicRelation").Filter("ItemId__in", gokits.SliceInt64To(iids)).All(&trs); err != nil {
		return nil, false
	}

	tids := make([]int64, len(trs))
	for _, v := range trs {
		tids = append(tids, v.Tid)
	}

	var ts []*Topic
	if _, err := t.Query("Topic").Filter("Id__in", gokits.SliceInt64To(tids)).All(&ts); err != nil {
		return nil, false
	}

	tsmap := make(map[int64]*Topic)
	for _, v := range ts {
		tsmap[v.Id] = v
	}

	item2tsmap := make(map[int64]Topics)
	for _, v := range trs {
		vts, existed := item2tsmap[v.ItemId]
		if !existed {
			vts = make(Topics, 0)
		}
		if t, e := tsmap[v.Tid]; e {
			vts = append(vts, t)
		}
		item2tsmap[v.ItemId] = vts
	}

	return item2tsmap, true
}

//----------------------------------------------------------

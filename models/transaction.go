package models

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/astaxie/beego/orm"
)

// --------------------------------------------------------
// a simple transaction
type Transaction struct {
	o   orm.Ormer
	Err error
}

func NewTr() *Transaction { return &Transaction{o: orm.NewOrm(), Err: nil} }

func (t *Transaction) Begin() *Transaction {
	t.o.Begin()
	return t
}

func (t *Transaction) End() {
	if err := recover(); err != nil {
		log.Errorf("Operate db failed, %v", err)
		t.o.Rollback()
		return
	}

	if t.Err == nil {
		t.o.Commit()
	} else {
		log.Errorf("Operate db failed, %v", t.Err.Error())
		t.o.Rollback()
	}
}

// --------------------------------------------------------
func (t *Transaction) Read(m interface{}, fields ...string) bool {
	if t.Err = t.o.Read(m, fields...); t.Err != nil {
		log.Errorf("Read data from db failed. error = %s", t.Err.Error())
		return false
	}
	return true
}

func (t *Transaction) Insert(m interface{}) (int64, bool) {
	if id, err := orm.NewOrm().Insert(m); err != nil {
		log.Errorf("Insert data to db failed. error = %s", err.Error())
		t.Err = err
		return id, false
	} else {
		return id, true
	}
}

func (t *Transaction) Update(m interface{}, fields ...string) (int64, bool) {
	if rownum, err := t.o.Update(m, fields...); err != nil {
		log.Errorf("Update data to db failed. error = %s", err.Error())
		t.Err = err
		return rownum, false
	} else {
		return rownum, true
	}
}

func (t *Transaction) Delete(m interface{}) bool {
	if _, t.Err = t.o.Delete(m); t.Err != nil {
		log.Errorf("Delete data from db failed. error = %s", t.Err.Error())
		return false
	} else {
		return true
	}
}

//----------------------------------------------------------
func (t *Transaction) Query(table string) orm.QuerySeter {
	return t.o.QueryTable(table)
}

//----------------------------------------------------------
// 判断行是否存在
func (t *Transaction) Existed(table string, field string, value interface{}) bool {
	return t.Query(table).Filter(field, value).Exist()
}

func (t *Transaction) ExistedV2(table string, params orm.Params) bool {
	return t.filters(table, params).Exist()
}

func (t *Transaction) filters(table string, params orm.Params) orm.QuerySeter {
	qs := t.Query(table)
	for k, v := range params {
		qs.Filter(k, v)
	}
	return qs
}

// 根据某个字查询统计行数
func (t *Transaction) Count(table string, field string, value interface{}) (int64, bool) {
	if c, err := t.Query(table).Filter(field, value).Count(); err != nil {
		return c, false
	} else {
		return c, true
	}
}

func (t *Transaction) CountV2(table string, params orm.Params) (int64, bool) {
	if c, err := t.filters(table, params).Count(); err != nil {
		return c, false
	} else {
		return c, true
	}
}

// 根据某个字查询统计字段值的求和
func (t *Transaction) Sum(table string, field string, where string, value interface{}) (int64, bool) {
	qb, _ := orm.NewQueryBuilder("mysql")
	qb.Select(fmt.Sprintf("SUM(%s)", field)).From(table).Where(where + " = ?")
	var sum int64 = 0
	if err := t.o.Raw(qb.String(), value).QueryRow(&sum); err != nil {
		return 0, false
	}
	return sum, true
}

//----------------------------------------------------------
// 根据Id来查询并更新数据
func (t *Transaction) UpdateById(table string, id int64, params orm.Params) bool {
	return t.UpdateByField(table, "Id", id, params)
}

// 根据某个字段查询并更新数据
func (t *Transaction) UpdateByField(table string, field string, value interface{}, params orm.Params) bool {
	if _, err := t.Query(table).Filter(field, value).Update(params); err != nil {
		return false
	}
	return true
}

func (t *Transaction) UpdateByFieldV2(table string, fileds orm.Params, params orm.Params) bool {
	if _, err := t.filters(table, fileds).Update(params); err != nil {
		return false
	}
	return true
}

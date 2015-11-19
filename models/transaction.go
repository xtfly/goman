package models

import (
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

func (t *Transaction) Begin() {
	t.o.Begin()
}

func (t *Transaction) End() {
	if t.Err == nil {
		t.o.Commit()
	} else {
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

func (t *Transaction) Query(table string) orm.QuerySeter {
	return t.o.QueryTable(table)
}

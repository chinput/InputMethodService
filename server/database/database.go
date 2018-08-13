package database

import (
	"code.aliyun.com/JRY/mtquery/module/merror"
	"code.aliyun.com/JRY/mtquery/module/mtype"
)

type OrderBy map[string]int

type DataBaseQuery struct {
	Id        string
	TableName string
	Condition mtype.IM
	Limit     int
	Skip      int
	OrderBy   OrderBy
	Data      interface{}
	Update    mtype.IM
	Quto      int
}

type DataBaser interface {
	//please use defer Close() when you use Open()
	//such as
	//defer databaser.Close()
	Open(dbhost string, dbname string) error
	Close() error
	QueryInit(query *DataBaseQuery) error
	Find(result []mtype.IM) error
	FindLike(result []mtype.IM) error
	Insert() (string, error)
	Update() error
	UpdateAll() error
	Delete() error
	InsertIfNotExist() (string, error)
}

func New(dbtype string, dbhost string, dbname string) (DataBaser, error) {
	var (
		conn DataBaser
		err  error
	)
	err = nil
	if dbtype == "mongo" {
		var newCon Mmongo
		conn = &newCon
	} else if dbtype == "mysql" {
		var newCon Mmysql
		conn = &newCon
	} else {
		err = merror.New("dbtypeerr",
			"database type invalid,need mongo/mysql, given is %s",
			dbtype)
		conn = nil
	}

	if err == nil {
		err = conn.Open(dbhost, dbname)
	}
	return conn, err
}

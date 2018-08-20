package database

import "gopkg.in/mgo.v2/bson"

type OrderBy map[string]int

type DataBaseQuery struct {
	Id        string
	TableName string
	Condition bson.M
	Limit     int
	Skip      int
	OrderBy   OrderBy
	Data      interface{}
	Update    bson.M
	Quto      int
}

type DataBaser interface {
	//please use defer Close() when you use Open()
	//such as
	//defer databaser.Close()
	Open(dbhost string, dbname string) error
	Close() error
	QueryInit(query *DataBaseQuery) error
	Find(result []bson.M) error
	FindLike(result []bson.M) error
	Insert() (string, error)
	Update() error
	UpdateAll() error
	Delete() error
	InsertIfNotExist() (string, error)
}

func New(dbhost string, dbname string) (DataBaser, error) {
	var (
		conn DataBaser
	)
	conn = new(Mmongo)
	return conn, conn.Open(dbhost, dbname)
}

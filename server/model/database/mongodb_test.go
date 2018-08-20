package database

import (
	"fmt"
	"log"
	"testing"

	"gopkg.in/mgo.v2/bson"
)

type person struct {
	Name  string
	Phone string
}

func TestMongodb(t *testing.T) {
	s, err := New("localhost:27017", "mydb")
	if err != nil {
		log.Fatal(err)
	}
	query := DataBaseQuery{
		Id:        "",
		TableName: "test",
		Condition: nil,
		Limit:     1,
		Skip:      0,
		OrderBy:   OrderBy{"_id": -1},
		Data:      person{"Bob", "192178238"},
		Update:    nil,
		Quto:      0,
	}
	result := make([]bson.M, 10)
	s.QueryInit(&query)
	s.Insert()
	s.QueryInit(&query)
	s.Find(result)
	fmt.Println(result)
	s.Close()
}

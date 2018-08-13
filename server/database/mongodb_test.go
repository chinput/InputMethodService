package database

import (
	"fmt"
	"log"
	"testing"

	"code.aliyun.com/JRY/mtquery/module/mtype"
)

type person struct {
	Name  string
	Phone string
}

func TestMongodb(t *testing.T) {
	s, err := New("mongo", "localhost:27017", "mydb")
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
	result := make([]mtype.IM, 10)
	s.QueryInit(&query)
	s.Insert()
	s.QueryInit(&query)
	s.Find(result)
	fmt.Println(result)
	s.Close()
}

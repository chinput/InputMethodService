package mmodel

import (
	"fmt"
	"log"
	"testing"

	"code.aliyun.com/JRY/mtquery/module/mtype"
)

func TestModel(t *testing.T) {
	conf := InitConf{
		Dbtype:    "mongo",
		Dbhost:    "localhost:27017",
		Dbname:    "mydb",
		Findlimit: 10,
	}
	InitByHand(conf)

	m := New()
	u, err := m.Copy("user")
	if err != nil {
		log.Fatal(err)
	}
	u.Add(mtype.IM{"name": "Bob", "phone": "123456"})
	data := u.FindOne(nil)
	fmt.Println(data)
}

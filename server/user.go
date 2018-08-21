package main

import (
	"log"

	"github.com/chinput/InputMethodService/server/config"
	"github.com/chinput/InputMethodService/server/model"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

type User struct {
	*model.Model
}

type DataUser struct {
	Username      string `json:"username" bson:"username"`
	Password      string `json:"password" bson:"password"`
	Nicknime      string `json:"nicknime" bson:"nicknime"`
	DataPath      string `json:"data_path" bson:"data_path"`
	Email         string `json:"email" bson:"email"`
	Register_Time int64  `json:"register_time" bson:"register_time"`
}

func newUser(m model.Modeler) *User {
	mm := newConnection(TABLE_User, m)
	return &User{mm}
}

func ConnectDB() {
	conf := model.InitConf{
		Dbhost:    config.DBUrl(),
		Dbname:    config.DBName(),
		Dbtype:    "mongo",
		Findlimit: 10,
	}
	//	c, error := database.New("mongo", "localhost:27017", "test")
	model.InitByHand(conf)
	m := model.New()
	defer m.End()

	user, err := m.Copy("user")
	if err != nil {
		log.Fatal(err)
	}

	insertData := DataUser{
		Username:      "hello world",
		Password:      "123456",
		Register_Time: 123456,
	}
	newUid, err := user.Add(insertData)
	log.Println(newUid, insertData)

	allData := user.FindMany(nil)

	user.DeleteById(newUid)

	log.Println(allData)
}

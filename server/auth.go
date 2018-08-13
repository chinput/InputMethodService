package main

import (
	"log"

	"code.aliyun.com/JRY/mtquery/module/mmodel"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

type Auth struct {
	*mmodel.Model
}

type DataAuth struct {
	Username      string `json:"username" bson:"username"`
	Password      string `json:"password" bson:"password"`
	Nicknime      string `json:"nicknime" bson:"nicknime"`
	DataPath      string `json:"data_path" bson:"data_path"`
	Email         string `json:"email" bson:"email"`
	Register_Time uint64 `json:"register_time" bson:"register_time"`
}

func newAuth(m mmodel.Modeler) *Auth {
	mm := newConnetction(TABLE_Auth, m)
	return &Auth{mm}
}

func ConnectDB() {
	conf := mmodel.InitConf{
		Dbhost:    "localhost",
		Dbname:    "tttdd",
		Dbtype:    "mongo",
		Findlimit: 10,
	}
	//	c, error := database.New("mongo", "localhost:27017", "test")
	mmodel.InitByHand(conf)
	m := mmodel.New()
	defer m.End()

	user, err := m.Copy("user")
	if err != nil {
		log.Fatal(err)
	}

	insertData := DataAuth{
		Username:      "hello world",
		Password:      "123456",
		Register_Time: 123456,
	}

	newUid := user.Add(insertData)
	log.Println(newUid, insertData)

	allData := user.FindMany(nil)

	user.DeleteById(newUid)

	log.Println(allData)
}

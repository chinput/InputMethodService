package database

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MmongoQuery struct {
	Condition bson.M
	Data      bson.M
	Update    bson.M
	Id        bson.ObjectId
	Limit     int
	Skip      int
	OrderBy   string
}

type Mmongo struct {
	Execed       int
	Session      *mgo.Session
	Db           *mgo.Database
	DbCollection *mgo.Collection
	Query        *DataBaseQuery
	MQuery       *MmongoQuery
}

/*
   var (
       Conf = config.D
   )

   session, err := mgo.Dial(Conf.DbHost) //连接数据库
   if err != nil {
       panic(err)
   }
   defer session.Close()
   session.SetMode(mgo.Monotonic, true)

   db := session.DB(Conf.DbName) //数据库名称
   collection := db.C("person2") //如果该集合已经存在的话，则直接返回
   temp := &Person{ID: bson.NewObjectId(), NAME: "admin", PHONE: "139xxxxx"}

   err = collection.Insert(temp)
   if err != nil {
       panic(err)
   }

    err = collection.FindId(Obid(5)).One(&result)
    collection.Find(bson.M{"_id": bson.NewObjectId()}).Apply(change, &baseId)



*/

/*
	Please defer Close When you use it:
	defer Mmongo.Close()
*/

// The mgo package has it's own connection pool.
// We do not need to write it by ourselves anymore.

const (
	//max connect num of mongodb
	connNum = 2
)

var (
	connIndex = 0
	//connections = make([]*mgo.Session, connNum)
)

var (
	oriSession *mgo.Session = nil
)

func (m *Mmongo) Open(dbhost string, dbname string) error {
	var (
		mquery MmongoQuery
		err    error
	)
	/*
		if connections[connIndex] == nil {
			debug.PrintMsg(connIndex, "create a new Database sesseion")
			newcon, err := mgo.Dial(dbhost)
			//connections[connIndex], err = mgo.Dial(dbhost)
			if err != nil {
				return err
			}
			connections[connIndex] = newcon
		}
		debug.PrintMsg("use Database session ", connIndex)
		connIndex = (connIndex + 1) & (connNum - 1)
	*/
	if oriSession == nil {
		oriSession, err = mgo.Dial(dbhost)
		if err != nil {
			return err
		}
		//oriSession = newcon
	}
	m.Session = oriSession.Copy() //connections[connIndex]
	m.Session.SetMode(mgo.Monotonic, true)
	m.Db = m.Session.DB(dbname)
	m.MQuery = &mquery
	return err
}

func (m *Mmongo) Close() (err error) {
	m.Session.Close()
	//m.Session = nil
	//m = nil
	//err = nil
	return nil
}

func (m *Mmongo) QueryInit(query *DataBaseQuery) (err error) {
	if query.TableName == "" {
		err = fmt.Errorf("table name is blank")
		return err
	} else {
		err = nil
	}
	m.DbCollection = m.Db.C(query.TableName)
	m.Query = query
	if query.Id != "" {
		//	m.MQuery.Id, err = m.string2Id(query.Id)
		m.MQuery.Id = bson.ObjectIdHex(query.Id)
		if err != nil {
			return err
		}
	}
	err = m.initCondition()
	if err != nil {
		return err
	}
	m.initQueryInt()
	m.MQuery.Data = nil
	m.Execed = 0
	return err
}

func (m *Mmongo) Find(result []bson.M) (err error) {
	if m.Execed == 0 {
		err = m.DbCollection.
			Find(m.MQuery.Condition).
			Limit(m.MQuery.Limit).
			Skip(m.MQuery.Skip).
			Sort(m.MQuery.OrderBy).
			All(&result)
		if err != nil {
			return err
		}

		for i := 0; i < len(result); i++ {
			result[i]["_id"] = m.id2String(result[i]["_id"].(bson.ObjectId))
		}
		m.Execed = 1
	} else {
		err = fmt.Errorf("need init")
	}
	return err
}

func (m *Mmongo) FindLike(result []bson.M) (err error) {
	var (
		newcond bson.M
		realnum int
	)
	if m.Execed == 0 {
		if m.MQuery.Condition != nil {
			for key, data := range m.MQuery.Condition {
				datastr := data.(string)
				conds := strings.Split(datastr, " ")
				num := len(conds)
				for i := 0; i < num; i++ {
					if conds[i] != "" {
						realnum++
					}
				}
				condsarr := make([]interface{}, realnum)
				realnum = 0
				for i := 0; i < num; i++ {
					if conds[i] != "" {
						condsarr[realnum] = bson.M{key: bson.M{"$regex": conds[i], "$options": "i"}}
						realnum++
					}
				}

				newcond = bson.M{"$or": condsarr}
				break
			}
			m.MQuery.Condition = newcond
		}

		return m.Find(result)

	} else {
		return fmt.Errorf("need init")
	}
}

func Struct2Map(obj interface{}) bson.M {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	var data = make(bson.M)
	for i := 0; i < t.NumField(); i++ {
		key := strings.ToLower(t.Field(i).Name)
		data[key] = v.Field(i).Interface()
	}
	return data
}

func (m *Mmongo) Insert() (id string, err error) {
	if m.Execed == 0 {

		if len(m.MQuery.Id) == 0 {
			m.MQuery.Id = bson.NewObjectId()
			//log.Println(m.Query.Data)
			log.Println(false, m.MQuery.Id)
		} else {
			log.Println(true, m.MQuery.Id)
			//	m.MQuery.Id = bson.ObjectIdHex(m.Query.Id)
		}

		t := reflect.TypeOf(m.Query.Data)
		if reflect.TypeOf(m.MQuery.Data) != t {
			m.MQuery.Data = Struct2Map(m.Query.Data)
		} else {
			m.MQuery.Data = m.Query.Data.(bson.M)
		}
		m.MQuery.Data["_id"] = m.MQuery.Id
		err = m.DbCollection.Insert(m.MQuery.Data)
		if err != nil {
			return "", err
		}
		id = m.id2String(m.MQuery.Id)
		m.Execed = 1
	} else {
		id = ""
		err = fmt.Errorf("need init")
	}

	return id, err
}

func (m *Mmongo) Update() (err error) {
	if m.Execed == 0 {
		m.initUpdate()
		err = m.DbCollection.Update(m.MQuery.Condition, m.MQuery.Update)
		if err != nil {
			return err
		}
		m.Execed = 1
	} else {
		err = fmt.Errorf("need init")
	}
	return err
}

func (m *Mmongo) UpdateAll() (err error) {
	if m.Execed == 0 {
		m.initUpdate()
		_, err = m.DbCollection.UpdateAll(m.MQuery.Condition, m.MQuery.Update)
		if err != nil {
			return err
		}
		m.Execed = 1
	} else {
		err = fmt.Errorf("need init")
	}
	return err
}

func (m *Mmongo) Delete() (err error) {
	if m.Execed == 0 {
		err = m.DbCollection.Remove(m.MQuery.Condition)
		m.Execed = 1
		if err != nil {
			return err
		}
	} else {
		err = fmt.Errorf("need init")
	}

	return err
}

func (m *Mmongo) InsertIfNotExist() (Id string, err error) {
	var (
		result = make([]bson.M, 1)
	)
	if m.Execed == 0 {
		m.MQuery.Limit = 1
		err = m.Find(result)
		if err != nil {
			return "", err
		}
		if result[0] == nil {
			m.MQuery.Data = m.Query.Condition
			m.MQuery.Id = bson.NewObjectId()
			m.MQuery.Data["_id"] = m.MQuery.Id
			err = m.DbCollection.Insert(m.MQuery.Data)
			if err != nil {
				return "", err
			}
			Id = m.id2String(m.MQuery.Id)
			m.Execed = 1
		} else {
			Id = result[0]["_id"].(string)
		}
	} else {
		Id = "0"
		err = fmt.Errorf("need init")
	}
	return Id, err
}

func (m *Mmongo) initCondition() error {
	var (
		err error
	)
	if m.Query.Condition == nil {
		m.MQuery.Condition = nil
	} else if len(m.Query.Condition) > 0 {
		m.MQuery.Condition = m.Query.Condition
	} else {
		m.MQuery.Condition = nil
	}
	if Id, Found := m.Query.Condition["_id"]; Found {
		m.MQuery.Condition["_id"], err = m.string2Id(Id)
	}
	return err
}

func (m *Mmongo) initUpdate() error {
	var (
		err error
	)
	if m.Query.Quto == 2 {
		m.MQuery.Update = bson.M{"$inc": m.Query.Update}
	} else {
		m.MQuery.Update = bson.M{"$set": m.Query.Update}
	}
	if Id, Found := m.Query.Update["_id"]; Found {
		m.MQuery.Update["_id"], err = m.string2Id(Id)
	}
	return err
}

func (m *Mmongo) initQueryInt() {
	if m.Query.Limit != 0 {
		m.MQuery.Limit = m.Query.Limit
	} else {
		m.MQuery.Limit = 1
	}

	if m.Query.OrderBy != nil {
		for field, num := range m.Query.OrderBy {
			if num == 1 {
				m.MQuery.OrderBy = field
			} else {
				m.MQuery.OrderBy = "-" + field
			}
		}
	} else {
		m.MQuery.OrderBy = "-_id"
	}

	m.MQuery.Skip = m.Query.Skip
}

func (m *Mmongo) id2String(Id bson.ObjectId) string {
	result := fmt.Sprintf("%x", string(Id))
	return result
}

func (m *Mmongo) string2Id(in interface{}) (bson.ObjectId, error) {
	var (
		Id  bson.ObjectId
		err error
	)
	t := reflect.TypeOf(in)
	if t == reflect.TypeOf("string") {
		inStr := in.(string)
		if inStr != "" {
			if len(inStr) != 24 {
				err = fmt.Errorf("id invalid")
			} else {
				Id = bson.ObjectIdHex(inStr)
				err = nil
			}
		} else {
			err = fmt.Errorf("blank string for objectid")
		}
	} else if t == reflect.TypeOf(bson.ObjectId("")) {
		Id = in.(bson.ObjectId)
		err = nil
	} else {
		err = fmt.Errorf("type not support")
	}
	return Id, err
}

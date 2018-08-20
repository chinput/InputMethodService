package model

import (
	"errors"

	"github.com/chinput/InputMethodService/server/model/database"
	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	TbName string
	Db     database.DataBaser
}

type Modeler interface {
	Copy(string) (*Model, error)
}

/*
type DataBaseQuery struct {
    Id        string
    TableName string
    Condition bson.M
    Limit     int
    Skip      int
    OrderBy   OrderBy
    Data      interface{}
    Update    bson.M
    Quto      int //2 for add ,others for set
}
type DataBaser interface {
    //please use defer Close() when you use Open()
    //such as
    //defer databaser.Close()
    Open(dbhost string, dbname string) error
    Close() error
    QueryInit(query *DataBaseQuery) error
    Find(result []bson.M) error
    Insert() (string, error)
    Update() error
    UpdateAll() error
    Delete() error
    InsertIfNotExist() (string, error)
}

*/

type InitConf struct {
	Dbtype    string
	Dbhost    string
	Dbname    string
	Findlimit int
}

var (
	confDb InitConf
	inited bool = false
)

func InitByHand(data InitConf) {
	confDb = data
	inited = true
}

func New() *Model {
	var (
		model Model
		err   error
	)
	if !inited {
		return nil
	}
	model.Db, err = database.New(confDb.Dbhost, confDb.Dbname)
	if err != nil {
		return nil
	}
	return &model
}

func DataRealLen(in *[]bson.M) *[]bson.M {
	var (
		res = make([]bson.M, 0, confDb.Findlimit)
		num = 0
	)
	for _, val := range *in {
		if val != nil {
			num++
		}
	}
	res = append(res, (*in)[0:num]...)
	return (&res)
}

func (m *Model) Copy(tb string) (*Model, error) {
	var (
		model Model
		err   error
	)
	if tb != "" {
		if m.Db != nil {
			model.TbName = tb
			model.Db = m.Db
		} else {
			err = errors.New("model_not_init" + "old model not init")
		}
	} else {
		err = errors.New("tbname_empty" + "table name is empty")
	}
	return &model, err
}

func (m *Model) End() {
	if m.Db != nil {
		m.Db.Close()
		m.Db = nil
		m = nil
	}
}

func (m *Model) newQuery() *database.DataBaseQuery {
	query := &database.DataBaseQuery{}
	query.TableName = m.TbName
	return query
}

func (m *Model) FindMany(cond bson.M, skip ...int) *[]bson.M {
	var (
		limit  = confDb.Findlimit
		result = make([]bson.M, limit)
		query  = m.newQuery()
		err    error
	)
	if len(skip) == 0 {
		query.Skip = 0
	} else {
		query.Skip = skip[0]
	}
	query.Limit = limit
	if cond != nil {
		query.Condition = cond
	}

	err = m.Db.QueryInit(query)
	if err != nil {
		return nil
	}
	err = m.Db.Find(result)
	return DataRealLen(&result)
}

func (m *Model) FindOne(cond bson.M) *bson.M {
	var (
		limit  = 1
		result = make([]bson.M, limit)
		query  = m.newQuery()
		err    error
	)
	query.Skip = 0
	if cond != nil {
		query.Condition = cond
	}

	err = m.Db.QueryInit(query)
	if err != nil {
		return nil
	}
	err = m.Db.Find(result)
	if len(result[0]) == 0 {
		return nil
	} else {
		return &result[0]
	}
}

func (m *Model) FindLike(cond bson.M, skip ...int) *[]bson.M {
	var (
		limit  = confDb.Findlimit
		result = make([]bson.M, limit)
		query  = m.newQuery()
		err    error
	)

	if len(skip) == 0 {
		query.Skip = 0
	} else {
		query.Skip = skip[0]
	}

	query.Limit = limit

	query.Condition = cond

	err = m.Db.QueryInit(query)
	if err != nil {
		return nil
	}
	if cond == nil {
		err = m.Db.Find(result)
	} else {
		err = m.Db.FindLike(result)
	}
	return DataRealLen(&result)
}

func (m *Model) FindManyFromOtherTable(table string, cond bson.M, skip ...int) (*[]bson.M, error) {
	var (
		limit  = confDb.Findlimit
		result = make([]bson.M, limit)
		query  = m.newQuery()
		err    error
	)
	if table == "" {
		return nil, errors.New("tbname_empty" + "table name is empty")
	} else {
		query.TableName = table
	}

	if len(skip) == 0 {
		query.Skip = 0
	} else {
		query.Skip = skip[0]
	}
	query.Limit = limit
	if cond != nil {
		query.Condition = cond
	}
	err = m.Db.QueryInit(query)
	if err != nil {
		return nil, err
	}
	err = m.Db.Find(result)
	if err != nil {
		return nil, err
	}
	return DataRealLen(&result), nil
}

func (m *Model) FindOneFromOtherTable(table string, cond bson.M) (*bson.M, error) {
	var (
		limit  = 1
		result = make([]bson.M, limit)
		query  = m.newQuery()
		err    error
	)
	if table != "" {
		query.TableName = table
	} else {
		return nil, errors.New("tbname_empty" + "table name is empty")
	}
	query.Skip = 0
	if cond != nil {
		query.Condition = cond
	}
	err = m.Db.QueryInit(query)
	if err != nil {
		return nil, err
	}
	err = m.Db.Find(result)
	if err != nil {
		return nil, err
	}
	if len(result[0]) == 0 {
		return nil, nil
	} else {
		return &result[0], nil
	}
}

func (m *Model) FindViaPage(page ...int) *[]bson.M {
	var (
		pg   int
		skip int
	)
	if len(page) == 0 {
		pg = 1
	} else {
		if page[0] <= 0 {
			pg = 1
		} else {
			pg = page[0]
		}
	}
	skip = (pg - 1) * confDb.Findlimit
	return m.FindMany(nil, skip)
}

func (m *Model) Add(data bson.M, id2 ...string) string {
	var (
		query = m.newQuery()
		id    string
		err   error
	)
	/*
		t := reflect.TypeOf(data)
		if reflect.TypeOf(insData) != t {
			insDataAddr := mtype.Struct2Map(data)
			insData = *insDataAddr
		} else {
			insData = data.(bson.M)
		}
	*/
	query.Data = data

	if len(id2) == 1 {
		query.Id = id2[0]
	}

	err = m.Db.QueryInit(query)
	if err != nil {
		id = ""
		return id
	}
	id, err = m.Db.Insert()
	return id
}

func (m *Model) AddToOtherTable(table string, data bson.M) (string, error) {
	var (
		query = m.newQuery()
		id    string
		err   error
	)
	query.Data = data
	if table != "" {
		query.TableName = table
	} else {
		return "", errors.New("tbname_empty" + "table name is empty")
	}
	err = m.Db.QueryInit(query)
	if err != nil {
		return "", err
	}
	id, err = m.Db.Insert()
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *Model) UpdateById(id string, data bson.M, muti ...bool) error {
	var (
		query = m.newQuery()
		err   error
	)
	query.Update = data
	query.Condition = bson.M{
		"_id": id,
	}
	if len(muti) == 0 {
		query.Quto = 1
	} else {
		if muti[0] == true {
			query.Quto = 2
		} else {
			query.Quto = 1
		}
	}
	err = m.Db.QueryInit(query)
	if err != nil {
		return err
	}
	err = m.Db.Update()
	return err
}

func (m *Model) DeleteById(id string) error {
	var (
		query = m.newQuery()
		err   error
	)
	query.Condition = bson.M{
		"_id": id,
	}
	err = m.Db.QueryInit(query)
	if err != nil {
		return err
	}
	err = m.Db.Delete()
	return err
}

func (m *Model) FindById(id string) *bson.M {
	cond := bson.M{
		"_id": id,
	}
	return m.FindOne(cond)
}

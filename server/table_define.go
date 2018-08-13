package main

import (
	"code.aliyun.com/JRY/mtquery/module/mmodel"
)

type TableName string

const (
	TABLE_Auth TableName = "auth"
)

func newConnetction(name TableName, m mmodel.Modeler) *mmodel.Model {
	if m != nil {
		m2, err := m.Copy(string(name))
		if err != nil {
			return newConnetction(name, nil)
		}

		return m2
	}

	m0 := mmodel.New()
	return m0
}

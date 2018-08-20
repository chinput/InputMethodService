package main

import "github.com/chinput/InputMethodService/server/model"

type TableName string

const (
	TABLE_Auth TableName = "auth"
)

func newConnetction(name TableName, m model.Modeler) *model.Model {
	if m != nil {
		m2, err := m.Copy(string(name))
		if err != nil {
			return newConnetction(name, nil)
		}

		return m2
	}

	m0 := model.New()
	return m0
}

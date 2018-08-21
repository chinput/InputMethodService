package main

import "github.com/chinput/InputMethodService/server/model"

type TableName string

const (
	// 登录状态，存储用户的登录 token，输入法启动时验证
	TABLE_Auth TableName = "auth"

	// 存储用户数据地址，包括用户名，密码，以及数据文件位置
	TABLE_User TableName = "user"

// 操作记录不存数据库，存进文本里
)

func newConnection(name TableName, m model.Modeler) *model.Model {
	if m != nil {
		m2, err := m.Copy(string(name))
		if err != nil {
			return newConnection(name, nil)
		}

		return m2
	}

	m0 := model.New()
	return m0
}

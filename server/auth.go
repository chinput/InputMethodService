package main

import (
	"time"

	"github.com/chinput/InputMethodService/server/model"
)

type Auth struct {
	*model.Model
}

type DataAuth struct {
	Uid        string `json:"uid" bson:"uid"`
	Token      string `json:"token" bson:"token"`
	Time       int64  `json:"time" bson:"time"`
	LogoutTime int64  `json:"logout_time" bson:"logout_time"`
}

func NewAuth(m model.Modeler) *Auth {
	return &Auth{
		newConnection(TABLE_Auth, m),
	}
}

func (a *Auth) AddAnAuthToken(uid string) (string, error) {
	data := DataAuth{
		Uid:        uid,
		Token:      newRandString(32),
		Time:       time.Now().Unix(),
		LogoutTime: 0,
	}

	_, err := a.Add(data)
	if err != nil {
		return "", err
	}

	return data.Token, nil
}

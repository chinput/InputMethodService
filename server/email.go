package main

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/chinput/InputMethodService/server/email"
)

type registeCode struct {
	email string
	code  string
	time  time.Time
}

type registerCodeGroup struct {
	group map[string]*registeCode
	lock  sync.Mutex
}

var regpool *registerCodeGroup
var randR *rand.Rand
var sBase = "QWERTYUIOPASDFGHJKLZXCVBNM1234567890"

func init() {
	regpool = new(registerCodeGroup)
	regpool.group = make(map[string]*registeCode)
	randR = rand.New(rand.NewSource(time.Now().Unix()))
}

func newRandString(length int) string {

	res := ""
	index := 0
	for i := 0; i < length; i++ {
		index = randR.Int() & 31
		res += sBase[index : index+1]
	}

	return res
}

func newRegisterCode() string {
	return newRandString(5)
}

func SendRegisterCode(email_addr string) error {
	regpool.lock.Lock()
	defer regpool.lock.Unlock()

	now := time.Now()
	exist := regpool.group[email_addr]
	if exist != nil {
		if now.Before(exist.time.Add(time.Minute)) {
			return errors.New("Please wait for at least one minute")
		}
	}

	code := newRegisterCode()

	log.Println("register code:", code)

	regpool.group[email_addr] = &registeCode{
		email: email_addr,
		code:  code,
		time:  now,
	}

	email.SendRegisterCode(email_addr, code)
	return nil
}

package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Link struct {
	expire time.Time
	path   string
	lock   sync.RWMutex
}

func (l *Link) RemoveIfExpire(now time.Time) bool {
	if l.expire.After(now) {
		return false
	}

	go func() {
		l.lock.Lock()
		defer l.lock.Unlock()
		log.Println("Remove expire file:", l.path)
		err := os.Remove(l.path)
		if err != nil {
			log.Println(err)
		}
	}()

	return true
}

func (l *Link) Reader() io.ReadCloser {
	l.lock.RLock()
	defer l.lock.RUnlock()

	r, err := os.Open(l.path)
	if err != nil {
		return nil
	}
	return r
}

type LinkGroup struct {
	lock   sync.RWMutex
	data   map[string]*Link
	during time.Duration

	end     bool
	endChan chan bool
}

func NewLinkGroup(duration time.Duration) *LinkGroup {
	lg := new(LinkGroup)
	lg.during = duration
	lg.data = make(map[string]*Link)
	lg.end = false
	lg.endChan = make(chan bool, 1)
	lg.Setup()
	return lg
}

func (lg *LinkGroup) Check() {
	//	shouldDelete := make([]string, 0, len(lg.data))
	now := time.Now()
	for k, v := range lg.data {
		if v == nil {
			delete(lg.data, k)
		}

		if v.RemoveIfExpire(now) {
			delete(lg.data, k)
		}
	}
}

func (lg *LinkGroup) CheckExpire() {
	for (!lg.end) || len(lg.data) > 0 {
		<-time.After(lg.during)
		log.Println("LinkGroup check expire")
		lg.Check()
	}

	lg.endChan <- true
}

func (lg *LinkGroup) Setup() {
	go lg.CheckExpire()
}

func (lg *LinkGroup) WaitForStop() {
	lg.end = true
	log.Println("Send the end signal")
	<-lg.endChan
	log.Println("Finished to check expire")
}

const randKeyStr = "qwertyuiopasdfghjklzxcvbnm012345"

var krd = rand.New(rand.NewSource(time.Now().Unix()))
var rdNum = krd.Int() & 65535

func newRandKey() string {
	key := ""

	for i := 0; i < 16; i++ {
		num := krd.Int() & 31
		key += randKeyStr[num : num+1]
	}

	rdNum++

	key += fmt.Sprintf("%x", rdNum)

	return key
}

/*
expire time.Time
path   string
lock   sync.RWMutex
*/

func (lg *LinkGroup) AddOne(path string) string {
	newId := newRandKey()
	lg.lock.Lock()
	defer lg.lock.Unlock()
	l := new(Link)
	l.expire = time.Now().Add(time.Second * 5)
	l.path = path
	lg.data[newId] = l
	return newId
}

func (lg *LinkGroup) ReaderOf(key string) io.ReadCloser {
	lg.lock.RLock()
	defer lg.lock.RUnlock()
	l := lg.data[key]
	if l == nil {
		return nil
	}

	return l.Reader()
}

var kLinkGroup *LinkGroup

func init() {
	kLinkGroup = NewLinkGroup(time.Second * 3)
}

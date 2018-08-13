package main

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Test_ConnectDB(t *testing.T) {
	ConnectDB()
}

func Test_DeleteMap(t *testing.T) {
	map1 := make(map[string]int)

	map1["one"] = 1
	map1["two"] = 2
	map1["three"] = 3
	map1["four"] = 4

	fmt.Println(map1)

	for k, v := range map1 {
		if v < 3 {
			delete(map1, k)
		}
	}

	fmt.Println(map1)
}

func Test_Rand(t *testing.T) {
	log.Println(newRandKey(), newRandKey(), newRandKey())
}

func Test_LinkGroup(t *testing.T) {
	num := 0
	for num < 100 {
		num++
		path := "/data/" + fmt.Sprint(num)
		log.Println(num, ": prepare to add path", path)
		go func() {
			<-time.After(time.Microsecond * 10000 * time.Duration(krd.Uint64()&511))
			log.Println(path, ": Add path", path)
			kLinkGroup.AddOne(path)
		}()
	}

	<-time.After(time.Second)

	kLinkGroup.WaitForStop()
}

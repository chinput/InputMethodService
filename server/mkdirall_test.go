package main

import (
	"log"
	"os"
	"testing"
)

func Test_mkdir(t *testing.T) {
	err := os.MkdirAll("./tmp/user/hhhh", 0755)
	log.Println(err)
}

package main

import (
	"github.com/doowsing/rediuse"
	"log"
)

var rdbHandler = rediuse.NewRdbHandler(GetR)

type User struct {
	Id   int
	Name string
}

func main() {
	rdbHandler.Set("id", 1)

	id, err := rdbHandler.Get("id").Int()
	if err != nil {
		log.Printf("get id err:%s\n", err)
	}

	idFlag, err := rdbHandler.Hget("flag", id).Bool()
	if err != nil {
		log.Printf("idFlag 1 is %s!", idFlag)
	}

	user := &User{
		Id:   1,
		Name: "张三",
	}
	rdbHandler.Set("user", user)
	user1 := &User{}
	err = rdbHandler.Get("user").Struct(user1)
	log.Printf("user name:%s\n", user1.Name)
}

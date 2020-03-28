package main

import (
	"github.com/doowsing/rediuse"
	"log"
)

var rdbHandler = rediuse.NewRdbHandler(GetRedisConn)

type User struct {
	Id   int
	Name string
}

func main() {
	rdbHandler.Set("id", 1)

	id, err := rdbHandler.Get("id").Int()
	if err != nil {
		log.Printf("get id err:%s\n", err)
	} else {
		log.Printf("get id :%d\n", id)
	}

	idFlag, err := rdbHandler.Hget("flag", id).Bool()
	if err != nil {
		log.Printf("idFlag 1 is %v!", idFlag)
	}

	user := &User{
		Id:   1,
		Name: "张三",
	}
	rdbHandler.Set("user", user)
	user1 := &User{}
	err = rdbHandler.Get("user").Struct(user1)
	if err != nil {
		log.Printf("get user err:%s\n", err)
	} else {
		log.Printf("user name:%s\n", user1.Name)
	}

}

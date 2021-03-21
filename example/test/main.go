package main

import (
	"fmt"
	"github.com/doowsing/rediuse"
)

var rdbHandler = rediuse.NewRdbHandler(GetRedisConn)

func hashOperation() {

	type Petime struct {
		Name string
	}
	rdbHandler.Delete("test_id_petime")
	id2petime := make(map[int]*Petime)
	id2petime[1] = &Petime{"mask1"}
	id2petime[2] = &Petime{"mask2"}
	id2petime[3] = &Petime{"mask3"}
	var err error
	err = rdbHandler.Hmset("test_id_petime", id2petime).Error()
	//getResult := rdbHandler.HgetAll("test_id_id")
	//id2petimeget := make(map[int]int)
	//err = getResult.ScanMap(&id2petimeget)
	//fmt.Printf("id2petimeget:%v\n",id2petimeget)

	//fmt.Printf("err:%v\n",err)

	id2petimeget := make(map[int]*Petime)
	getResult := rdbHandler.HgetAll("test_id_petime")
	getResult.This()
	err = getResult.ScanMap(&id2petimeget)
	_ = err
	//fmt.Printf("ifce:%v\n",ifce)
	fmt.Printf("id2petime:%s\n", id2petime)
	fmt.Printf("id2petimeget:%s\n", id2petimeget)
	fmt.Printf("err:%v\n", err)
}

func main() {
	hashOperation()
}

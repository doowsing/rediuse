package main

import "testing"

func BenchmarkHashOperation(b *testing.B) {
	type Petime struct {
		Name string
	}
	id2petime := make(map[int]Petime)
	var err error
	err = rdbHandler.Hmset("test_id_petime1", id2petime).Error()
	id2petimeget := make(map[float64]Petime)
	getResult := rdbHandler.Hmget("test_id_petime1", 5, 2)
	for i := 0; i < b.N; i++ {

		//getResult := rdbHandler.HgetAll("test_id_id")
		//id2petimeget := make(map[int]int)
		//err = getResult.ScanMap(&id2petimeget)
		//fmt.Printf("id2petimeget:%v\n",id2petimeget)

		//fmt.Printf("err:%v\n",err)

		err = getResult.ScanMapByMergeArgs(&id2petimeget)
		_ = err
	}
}

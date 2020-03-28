# rediuse

简化对redigo包的操作，将redigo包操作后的结果进行封装，以便更加方便的获取想要的数据格式。
可转换的数据格式包括：int,bool,float64,string,[]byte,以及各种常用类型的slice,结构体。
对传入的结构体、切片、字典会自动做进行json编码。
可使用单机与集群，详情见example文件夹。
## Example
```go
package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
    "github.com/doowsing/rediuse"
	"time"
)

// redis 连接池
var pool *redis.Pool

//根据配置初始化打开redis连接
func init() {
	pool = &redis.Pool{
		MaxIdle:     20,
		MaxActive:   30,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				return nil, err
			}
			//TODO 不加有时候池子链接失败
			// 线上环境redis配置密码, 则需要加上这句AUTH
			//_,err = c.Do("AUTH","24245@163.com")
			return c, err
		},
		//testOnBorrow 向资源池借用连接时是否做连接有效性检测(ping)，无效连接会被移除 默认值 false 业务量很大时候建议设置为false(多一次ping的开销)。
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	fmt.Printf("redis start on 6379\n")
}
func GetRedisPool() *redis.Pool {
	return pool
}

// 获取redis全局实例
func GetR() redis.Conn {
	return pool.Get()
}

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
```

使用获得redis.Conn的函数进行初始化。

[redis单机连接库](https://github.com/gomodule/redigo/redis)

[redis集群连接库](https://github.com/mna/redisc)

两个库的接口是一样的，便于开发和线上切换

目前加入的redis使用函数有限，可自行添加